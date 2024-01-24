package enrich

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (e Enrich) getCountry(surname string) (string, error) {
	type Resp struct {
		Country []struct {
			Country     string  `json:"country_id"`
			Probability float32 `json:"probability"`
		} `json:"country"`
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.nationalize.io/", nil)
	if err != nil {
		return "", fmt.Errorf("getAge : http.NewRequest err %w", err)
	}
	q := req.URL.Query()
	q.Add("name", surname)
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

	countrys := Resp{}
	err = render.DecodeJSON(res.Body, &countrys)
	if err != nil {
		return "", fmt.Errorf("getAge : render.DecodeJSON : %w", err)
	}

	if len(countrys.Country) == 0 {
		return "", ErrInvalidParameter
	}

	var country string
	var maxProbability = float32(0)

	for _, v := range countrys.Country {
		if v.Probability > maxProbability {
			maxProbability = v.Probability
			country = v.Country
		}
	}

	return country, nil
}
