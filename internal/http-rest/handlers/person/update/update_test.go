package update_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	pupdate "test4effectivemobile/internal/http-rest/handlers/person/update"
	"test4effectivemobile/internal/http-rest/handlers/person/update/mocks"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/lib/logger/handlers/slogdiscard"
	"test4effectivemobile/internal/model"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
)

func TestUpddateOk(t *testing.T) {
	person := model.Person{
		ID:         1,
		Name:       "Dmitriy",
		Surname:    "Ushakov",
		Patronymic: "Vasilevich",
		Age:        1,
		Gender:     "G",
		Country:    "C",
	}
	updater := mocks.NewUpdater(t)
	updater.On("Update", person).Return(nil).Once()

	body := strings.NewReader(
		`{
			"name": "Dmitriy", 
			"surname": "Ushakov", 
			"patronymic": "Vasilevich",
			"age": 1,
			"Gender": "G",
			"Country": "C"
		}`,
	)
	r := httptest.NewRequest(http.MethodPatch, "/person/1", body)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	pupdate.New(slogdiscard.NewDiscardLogger(), updater).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := response.Response{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.Ok(), data)
}

func TestUpddateBadReq(t *testing.T) {
	updater := mocks.NewUpdater(t)

	requests := []string{
		`{
			"surname": "Ushakov", 
			"patronymic": "Vasilevich",
			"age": 1,
			"Gender": "G",
			"Country": "C"
		}`,
		`{
			"name": "Dmitriy",
			"patronymic": "Vasilevich",
			"age": 1,
			"Gender": "G",
			"Country": "C"
		}`,
		`{
			"name": "Dmitriy",
			"surname": "Ushakov",
			"patronymic": "Vasilevich",
			"Gender": "G",
			"Country": "C"
		}`,
		`{
			"name": "Dmitriy",
			"surname": "Ushakov",
			"patronymic": "Vasilevich",
			"age": 1,
			"Country": "C"
		}`,
		`{
			"name": "Dmitriy",
			"surname": "Ushakov",
			"patronymic": "Vasilevich",
			"age": 1,
			"Gender": "G"
		}`,
		"test",
	}

	for _, b := range requests {
		body := strings.NewReader(b)

		r := httptest.NewRequest(http.MethodPatch, "/person/1", body)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("ID", "1")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		pupdate.New(slogdiscard.NewDiscardLogger(), updater).ServeHTTP(w, r)
		res := w.Result()
		defer func() { _ = res.Body.Close() }()
		data := response.Response{}
		err := render.DecodeJSON(res.Body, &data)
		require.NoError(t, err)

		require.Equal(t, response.ErrBadReq, data.Error)
	}
}
