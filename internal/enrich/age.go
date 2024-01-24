package enrich

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (e Enrich) getAge(name string) (uint, error) {
	type Resp struct {
		Age uint `json:"age"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.agify.io/", nil)
	if err != nil {
		return 0, fmt.Errorf("getAge : http.NewRequest err %w", err)
	}
	q := req.URL.Query()
	q.Add("name", name)
	req.URL.RawQuery = q.Encode()

	res, err := e.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("getAge : e.httpClient.Do(req) %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode == http.StatusUnprocessableEntity {
		return 0, ErrInvalidParameter
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("getAge : response StatusCode %d", res.StatusCode)
	}

	age := Resp{}
	err = render.DecodeJSON(res.Body, &age)
	if err != nil {
		return 0, fmt.Errorf("getAge : render.DecodeJSO %w", err)
	}

	if age.Age == 0 {
		return 0, ErrInvalidParameter
	}

	return age.Age, nil
}
