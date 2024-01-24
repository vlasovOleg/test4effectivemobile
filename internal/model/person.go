package model

type E string

func (e E) Error() string {
	return string(e)
}

const (
	ErrMissValue = E("miss value")
)

type Person struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        uint   `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}

func (p Person) Validate() error {
	// ID for Validate before BD

	if len(p.Name) == 0 {
		return ErrMissValue
	}

	if len(p.Surname) == 0 {
		return ErrMissValue
	}

	// Patronymic omitempty

	if p.Age == uint(0) {
		return ErrMissValue
	}

	if len(p.Gender) == 0 {
		return ErrMissValue
	}

	if len(p.Country) == 0 {
		return ErrMissValue
	}
	return nil
}
