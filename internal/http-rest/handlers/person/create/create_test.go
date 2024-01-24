package create_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"test4effectivemobile/internal/http-rest/handlers/person/create"
	"test4effectivemobile/internal/http-rest/handlers/person/create/mocks"
	"test4effectivemobile/internal/lib/logger/handlers/slogdiscard"
	"test4effectivemobile/internal/model"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewpersonOk(t *testing.T) {
	saver := mocks.NewSaver(t)
	enricher := mocks.NewEnricher(t)

	body := strings.NewReader(`{
			"name": "Dmitriy",
			"surname": "Ushakov",
			"patronymic": "Vasilevich"
		}`)
	person := model.Person{
		ID:         1,
		Name:       "Dmitriy",
		Surname:    "Ushakov",
		Patronymic: "Vasilevich",
		Age:        10,
		Gender:     "g",
		Country:    "c",
	}

	enricher.On("Enrich", mock.Anything).Return(person, nil)
	saver.On("Save", mock.Anything).Return(int64(1), nil)

	r := httptest.NewRequest(http.MethodPost, "/person", body)
	w := httptest.NewRecorder()

	create.New(slogdiscard.NewDiscardLogger(), saver, enricher).ServeHTTP(w, r)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()
	data := create.Resp{}
	err := render.DecodeJSON(res.Body, &data)

	require.NoError(t, err)
	require.Empty(t, data.Error)
	require.Empty(t, data.Message)
	require.Equal(t, int64(1), data.ID)
}

func TestNewpersonNoValue(t *testing.T) {
	testCases := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "no name",
			body: `{
				"surname": "Ushakov",
		 		"patronymic": "Vasilevich"
			 }`,
			err: "have no name",
		},
		{
			name: "no surname",
			body: `{
				"name": "Dmitriy",
				"patronymic": "Vasilevich"
			 }`,
			err: "have no surname",
		},
	}

	saver := mocks.NewSaver(t)
	enricher := mocks.NewEnricher(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(tc.body)

			r := httptest.NewRequest(http.MethodPost, "/person", body)
			w := httptest.NewRecorder()

			create.New(slogdiscard.NewDiscardLogger(), saver, enricher).ServeHTTP(w, r)
			res := w.Result()
			defer func() { _ = res.Body.Close() }()
			data := create.Resp{}
			err := render.DecodeJSON(res.Body, &data)

			require.NoError(t, err)
			require.Equal(t, tc.err, data.Message)
			require.Empty(t, data.ID)
		})
	}
}
