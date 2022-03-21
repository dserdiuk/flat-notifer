package service

import (
	"context"
	"github.com/dserdiuk/flat-notifier/internal/model"
	"github.com/dserdiuk/flat-notifier/internal/notifier"
	"github.com/dserdiuk/flat-notifier/internal/source"
	"log"
	"sync"
	"time"
)

type Checker struct {
	Sources  []source.Source
	Notifier notifier.Notifier
	Wg       sync.WaitGroup
}

func NewCheckService(sources []source.Source, notifier notifier.Notifier) *Checker {
	return &Checker{
		Sources:  sources,
		Notifier: notifier,
		Wg:       sync.WaitGroup{},
	}
}

func (c *Checker) Start() {
	ch := make(chan []*model.Flat)
	go c.CheckAllSources(ch)
	go c.WaitNewFlats(ch)
}

func (c *Checker) Notify(flat *model.Flat) {
	err := c.Notifier.Notify(flat.Url)
	if err != nil {
		log.Printf("failed to notify: %v\n", err)
	}
}

func (c *Checker) WaitNewFlats(ch chan []*model.Flat) {
	for {
		flats, ok := <-ch
		if ok == false {
			break
		} else {
			c.Wg.Done()
			if len(flats) > 0 {
				for _, flat := range flats {
					c.Notify(flat)
				}
			}
		}
	}
}

func (c *Checker) CheckAllSources(ch chan []*model.Flat) {
	ctx := context.TODO()
	for _, src := range c.Sources {
		go func() {
			err := src.GetNewFlats(ctx, ch)
			if err != nil {
				log.Printf("failed to get flats: %v\n", err)
			}
		}()
		c.Wg.Add(1)
	}
	c.Wg.Wait()
	time.Sleep(time.Second * 60)
	c.CheckAllSources(ch)
}
