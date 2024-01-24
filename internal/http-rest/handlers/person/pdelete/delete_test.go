package pdelete_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"test4effectivemobile/internal/http-rest/handlers/person/pdelete/mocks"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/lib/logger/handlers/slogdiscard"
	"testing"

	"test4effectivemobile/internal/http-rest/handlers/person/pdelete"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
)

func TestDeleteOK(t *testing.T) {
	deleter := mocks.NewDeleter(t)
	deleter.On("Delete", int64(1)).Return(nil).Once()

	r := httptest.NewRequest(http.MethodDelete, "/person/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	pdelete.New(slogdiscard.NewDiscardLogger(), deleter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := response.Response{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.StatusOK, data.Status)
	require.Empty(t, data.Error)
	require.Empty(t, data.Message)
}

func TestDeleteBadID(t *testing.T) {
	deleter := mocks.NewDeleter(t)

	r := httptest.NewRequest(http.MethodGet, "/person/d", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "d")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	pdelete.New(slogdiscard.NewDiscardLogger(), deleter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := response.Response{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.StatusError, data.Status)
	require.Equal(t, response.ErrBadReq, data.Error)
}

func TestDeleteInternalErr(t *testing.T) {
	deleter := mocks.NewDeleter(t)
	deleter.On("Delete", int64(1)).Return(errors.New(response.ErrInternal)).Once()

	r := httptest.NewRequest(http.MethodGet, "/person/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	pdelete.New(slogdiscard.NewDiscardLogger(), deleter).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := response.Response{}
	err := render.DecodeJSON(res.Body, &data)
	require.NoError(t, err)

	require.Equal(t, response.StatusError, data.Status)
}
