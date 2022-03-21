package source

import (
	"context"
	"github.com/dserdiuk/flat-notifier/internal/model"
)

type Source interface {
	GetNewFlats(context context.Context, ch chan []*model.Flat)
}
