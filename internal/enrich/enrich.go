package enrich

import (
	"errors"
	"net/http"
	"test4effectivemobile/internal/model"
)

type E string

func (e E) Error() string {
	return string(e)
}

const (
	ErrCanNotGetAge     = E("can not get age")
	ErrCanNotGetCountry = E("can not get country")
	ErrCanNotGetGender  = E("can not get gender")
	ErrInvalidParameter = E("invalid parameter")
)

type Enrich struct {
	httpClient http.Client
}

func New(h http.Client) *Enrich {
	return &Enrich{
		httpClient: h,
	}
}

func (e Enrich) Enrich(p model.Person) (model.Person, error) {
	var err error

	p.Age, err = e.getAge(p.Name)
	if err != nil {
		return model.Person{}, errors.Join(ErrCanNotGetAge, err)
	}

	p.Country, err = e.getCountry(p.Surname)
	if err != nil {
		return model.Person{}, errors.Join(ErrCanNotGetCountry, err)
	}

	p.Gender, err = e.getGender(p.Name, p.Country)
	if err != nil {
		return model.Person{}, errors.Join(ErrCanNotGetGender, err)
	}

	return p, nil
}
