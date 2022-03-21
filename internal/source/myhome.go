package source

import (
	"context"
	"encoding/json"
	"github.com/dserdiuk/flat-notifier/internal/model"
	"io/ioutil"
	"net/http"
	"time"
)

const BaseUrl = "https://www.myhome.ge/en/s/?"
const Timezone = 4

type MyHomeSource struct {
	Options       string
	LastCheckTime time.Time
}

func NewMyHomeSource(options string) *MyHomeSource {
	return &MyHomeSource{
		Options:       options,
		LastCheckTime: time.Now(),
	}
}

type MyHomeResult struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	Data          struct {
		Prs []struct {
			ProductID      string      `json:"product_id"`
			UserID         string      `json:"user_id"`
			ParentID       interface{} `json:"parent_id"`
			MaklerID       interface{} `json:"makler_id"`
			HasLogo        interface{} `json:"has_logo"`
			MaklerName     interface{} `json:"makler_name"`
			LocID          string      `json:"loc_id"`
			StreetAddress  string      `json:"street_address"`
			YardSize       string      `json:"yard_size"`
			YardSizeTypeID string      `json:"yard_size_type_id"`
			SubmissionID   string      `json:"submission_id"`
			AdtypeID       string      `json:"adtype_id"`
			ProductTypeID  string      `json:"product_type_id"`
			Price          string      `json:"price"`
			Photo          string      `json:"photo"`
			PhotoVer       string      `json:"photo_ver"`
			PhotosCount    string      `json:"photos_count"`
			AreaSizeValue  string      `json:"area_size_value"`
			VideoURL       string      `json:"video_url"`
			CurrencyID     string      `json:"currency_id"`
			OrderDate      string      `json:"order_date"`
			PriceTypeID    string      `json:"price_type_id"`
			Vip            string      `json:"vip"`
			Color          string      `json:"color"`
			EstateTypeID   string      `json:"estate_type_id"`
			AreaSize       string      `json:"area_size"`
			AreaSizeTypeID string      `json:"area_size_type_id"`
			Comment        interface{} `json:"comment"`
			MapLat         string      `json:"map_lat"`
			MapLon         string      `json:"map_lon"`
			LLiving        string      `json:"l_living"`
			SpecialPersons string      `json:"special_persons"`
			Rooms          string      `json:"rooms"`
			Bedrooms       string      `json:"bedrooms"`
			Floor          string      `json:"floor"`
			ParkingID      string      `json:"parking_id"`
			Canalization   string      `json:"canalization"`
			Water          string      `json:"water"`
			Road           string      `json:"road"`
			Electricity    string      `json:"electricity"`
			OwnerTypeID    string      `json:"owner_type_id"`
			OsmID          string      `json:"osm_id"`
			NameJSON       string      `json:"name_json"`
			PathwayJSON    string      `json:"pathway_json"`
			Homeselfie     string      `json:"homeselfie"`
			SeoTitleJSON   string      `json:"seo_title_json"`
			SeoNameJSON    interface{} `json:"seo_name_json"`
		} `json:"Prs"`
		MapData struct {
			Points []struct {
				ProductID string `json:"product_id"`
				FoundByID int    `json:"FoundByID"`
				Href      string `json:"href"`
				Price     struct {
					Amount   string `json:"amount"`
					Currency string `json:"currency"`
				} `json:"price"`
				Center struct {
					Lat string `json:"lat"`
					Lng string `json:"lng"`
				} `json:"center"`
				Title     string `json:"title"`
				OwnerType string `json:"owner_type"`
				Img       string `json:"img"`
				Avatar    string `json:"avatar"`
				Letter    string `json:"letter"`
			} `json:"Points"`
		} `json:"MapData"`
	} `json:"Data"`
}

func getUnixTime(date string) int64 {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Now().Unix()
	}

	t = t.Add(-(time.Hour * Timezone))
	return t.Unix()
}

func (s *MyHomeSource) GetNewFlats(context context.Context, ch chan []*model.Flat) {
	resp, err := http.Get(BaseUrl + s.Options)
	if err != nil {
		return
	}

	defer func() {
		resp.Body.Close()
		s.LastCheckTime = time.Now()
	}()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var res MyHomeResult
	err = json.Unmarshal(response, &res)
	if err != nil {
		return
	}

	var flats []*model.Flat
	for _, flat := range res.Data.Prs {
		if getUnixTime(flat.OrderDate) < s.LastCheckTime.Unix() {
			//break
		}

		for _, point := range res.Data.MapData.Points {
			if point.ProductID != flat.ProductID {
				continue
			}
			modelFlat := model.Flat{Url: point.Href}
			flats = append(flats, &modelFlat)
			break
		}

		break
	}
	ch <- flats
}
