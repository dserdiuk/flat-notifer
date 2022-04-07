package source

import (
	"context"
	"fmt"
	"github.com/dserdiuk/flat-notifier/internal/model"
	"testing"
)

func TestGetNewFlats(t *testing.T) {
	source := NewSsSource("Sort.SortExpression=\"OrderDate\"%20DESC&RealEstateTypeId=5&RealEstateDealTypeId=1&CityIdList=96&StatusField.FieldId=34&StatusField.Type=SingleSelect&StatusField.StandardField=Status&Rooms=3&Rooms=4&Rooms=5&Rooms=6&PriceType=false&CurrencyId=1")
	ch := make(chan []*model.Flat)
	ctx := context.TODO()
	go func() {
		err := source.GetNewFlats(ctx, ch)
		if err != nil {
			t.Error(err)
		}
	}()

	for {
		flats, ok := <-ch
		if ok == false {
			break
		} else {
			if len(flats) > 0 {
				for _, flat := range flats {
					fmt.Println(flat)
				}
			}
		}
	}
}
