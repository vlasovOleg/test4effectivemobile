package storage

import (
	"test4effectivemobile/internal/model"
)

type E string

func (e E) Error() string {
	return string(e)
}

const (
	ErrInternal = E("internal database error")
	ErrNotFound = E("not found")
)

type Person interface {
	Save(model.Person) (int64, error)
	GetByID(int64) (model.Person, error)
	Update(model.Person) error
	Delete(int64) error
	Get(f model.Filter) ([]model.Person, error)
}
