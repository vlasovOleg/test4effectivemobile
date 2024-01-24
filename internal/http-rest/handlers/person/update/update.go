package update

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/model"
	"test4effectivemobile/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Req struct {
	model.Person
}

type E string

func (e E) Error() string {
	return string(e)
}

const (
	ErrMissValue = E("miss value")
)

func (r Req) Validate() error {
	if len(r.Name) == 0 {
		return ErrMissValue
	}

	if len(r.Surname) == 0 {
		return ErrMissValue
	}

	// Patronymic omitempty

	if r.Age == uint(0) {
		return ErrMissValue
	}

	if len(r.Gender) == 0 {
		return ErrMissValue
	}

	if len(r.Country) == 0 {
		return ErrMissValue
	}
	return nil
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Updater
type Updater interface {
	Update(model.Person) error
}

func New(log *slog.Logger, updater Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.person.update"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
		if err != nil {
			log.Debug("parse ID error", "err", err.Error())
			response.SendErrorBadReq(r, w, response.ErrWrongID)
			return
		}

		var req Req
		if err = render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("request prepare error", "err", err.Error())
			response.SendErrorBadReq(r, w, "bad json")
			return
		}
		log.Info("request", "body", fmt.Sprintf("%#v", req))

		if err = req.Validate(); err != nil {
			log.Debug("request validate error", "err", err.Error())
			response.SendErrorBadReq(r, w, err.Error())
			return
		}

		req.Person.ID = id
		err = updater.Update(req.Person)
		if errors.Is(err, storage.ErrNotFound) {
			log.Debug("person not found", "err", err.Error())
			response.SendErrorNotFound(r, w)
			return
		}
		if err != nil {
			log.Error("update person error", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Ok())
	}
}
