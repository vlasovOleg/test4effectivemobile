package getbyid_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"test4effectivemobile/internal/http-rest/handlers/person/getbyid"
	"test4effectivemobile/internal/http-rest/handlers/person/getbyid/mocks"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/lib/logger/handlers/slogdiscard"
	"test4effectivemobile/internal/model"
	"test4effectivemobile/internal/storage"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
)

func TestGetOk(t *testing.T) {
	person := model.Person{
		ID:         1,
		Name:       "Test",
		Surname:    "Test",
		Patronymic: "Test",
		Age:        1,
		Gender:     "Test",
		Country:    "Test",
	}

	geter := mocks.NewGeterById(t)
	geter.On("GetByID", int64(1)).Return(person, nil).Once()

	r := httptest.NewRequest(http.MethodGet, "/person/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	getbyid.New(slogdiscard.NewDiscardLogger(), geter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := getbyid.Resp{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, person, data.Person)
}

func TestGetBadID(t *testing.T) {
	geter := mocks.NewGeterById(t)

	r := httptest.NewRequest(http.MethodGet, "/person/g", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "g")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	getbyid.New(slogdiscard.NewDiscardLogger(), geter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := getbyid.Resp{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.StatusError, data.Status)
	require.Equal(t, response.ErrBadReq, data.Error)
}

func TestGetNotFound(t *testing.T) {
	geter := mocks.NewGeterById(t)
	geter.On(
		"GetByID", int64(1)).
		Return(
			model.Person{},
			storage.ErrNotFound).
		Once()

	r := httptest.NewRequest(http.MethodGet, "/person/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	getbyid.New(slogdiscard.NewDiscardLogger(), geter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := getbyid.Resp{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.StatusError, data.Status)
	require.Equal(t, response.ErrNotFound, data.Error)
}
