package source

import (
	"context"
	"github.com/dserdiuk/flat-notifier/internal/model"
	"time"
)

type Mock struct {
	sleepTime time.Duration
}

func NewMock(time time.Duration) Mock {
	return Mock{sleepTime: time}
}

func (m Mock) GetNewFlats(context context.Context, ch chan []*model.Flat) {
	var flats []*model.Flat
	flats = append(flats, &model.Flat{Url: "https://localhost/test.php"})
	time.Sleep(time.Second * m.sleepTime)
	ch <- flats
}
