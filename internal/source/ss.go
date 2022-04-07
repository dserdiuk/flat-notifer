package source

import (
	"context"
	"github.com/dserdiuk/flat-notifier/internal/model"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

const BaseSsUrl = "https://ss.ge/en/real-estate/realestateapplicationlist?"

type MySsSource struct {
	Options       string
	LastCheckTime time.Time
	blocksRe      *regexp.Regexp
	createdAtRe   *regexp.Regexp
	urlRe         *regexp.Regexp
}

func NewSsSource(options string) *MySsSource {
	return &MySsSource{
		Options:       options,
		LastCheckTime: time.Now(),
		blocksRe:      regexp.MustCompile(`(?m)<div class="latest_article_each ".*>`),
		createdAtRe:   regexp.MustCompile(`(?ms)(\d+.\d+.\d+ \/ \d+\:\d+)`),
		urlRe:         regexp.MustCompile(`(?ms)href="(.*?)"`),
	}
}

func (s *MySsSource) getUnixTime(date string) int64 {
	layout := "02.01.2006 / 15:04"
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Now().Unix()
	}

	t = t.Add(-(time.Hour * Timezone))
	return t.Unix()
}

func (s *MySsSource) GetNewFlats(context context.Context, ch chan []*model.Flat) error {
	resp, err := http.Get(BaseSsUrl + s.Options)
	if err != nil {
		return err
	}

	defer func() {
		resp.Body.Close()
		s.LastCheckTime = time.Now()
	}()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := s.blocksRe.Split(string(response), -1)

	var flats []*model.Flat
	for _, flatBlock := range result {

		createdAt := s.createdAtRe.FindString(flatBlock)
		if createdAt == "" {
			continue
		}

		if s.getUnixTime(createdAt) < s.LastCheckTime.Unix() {
			continue
		}

		url := s.urlRe.FindStringSubmatch(flatBlock)
		if len(url) < 2 {
			continue
		}
		modelFlat := model.Flat{Url: "https://ss.ge" + url[1]}
		flats = append(flats, &modelFlat)
	}

	ch <- flats
	return nil
}
