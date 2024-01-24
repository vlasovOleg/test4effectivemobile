package enrich

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (e Enrich) getGender(name, country string) (string, error) {
	type Resp struct {
		Gender string `json:"gender"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.genderize.io", nil)
	if err != nil {
		return "", fmt.Errorf("getAge : http.NewRequest err %w", err)
	}
	q := req.URL.Query()
	q.Add("name", name)
	q.Add("country_id", country)
	req.URL.RawQuery = q.Encode()

	res, err := e.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("getAge : e.httpClient.Do(req) %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode == http.StatusUnprocessableEntity {
		return "", ErrInvalidParameter
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("getAge : response StatusCode %d", res.StatusCode)
	}

	gender := Resp{}
	err = render.DecodeJSON(res.Body, &gender)
	if err != nil {
		return "", fmt.Errorf("getAge : render.DecodeJSON : %w", err)
	}

	if len(gender.Gender) == 0 {
		return "", ErrInvalidParameter
	}

	return gender.Gender, nil
}
