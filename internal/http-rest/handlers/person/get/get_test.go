package get_test

import (
	"net/http"
	"net/http/httptest"
	"test4effectivemobile/internal/http-rest/handlers/person/get"
	"test4effectivemobile/internal/http-rest/handlers/person/get/mocks"
	"test4effectivemobile/internal/lib/logger/handlers/slogdiscard"
	"test4effectivemobile/internal/model"
	"testing"
)

func TestGetOk(t *testing.T) {
	filter := model.Filter{
		AgeMin:  1,
		AgeMax:  2,
		Gender:  "3",
		Country: "4",
		Offset:  5,
		Limit:   6,
	}

	geter := mocks.NewGeter(t)
	geter.On("Get", filter).Return([]model.Person{}, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "http://domain.com/person/?agemax=2&gender=3&country=4&offset=5&limit=6&agemin=1", nil)
	w := httptest.NewRecorder()

	get.New(slogdiscard.NewDiscardLogger(), geter).ServeHTTP(w, r)
}
