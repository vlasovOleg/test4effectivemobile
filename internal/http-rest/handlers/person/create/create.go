package create

import (
	"errors"
	"log/slog"
	"net/http"
	"test4effectivemobile/internal/enrich"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/model"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type E string

func (e E) Error() string {
	return string(e)
}

const (
	ErrHaveNoName    = E("have no name")
	ErrHaveNoSurname = E("have no surname")
)

type Req struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Resp struct {
	response.Response
	ID int64 `json:"id,omitempty"`
}

func (r Req) Validate() error {
	if len(r.Name) == 0 {
		return ErrHaveNoName
	}
	if len(r.Surname) == 0 {
		return ErrHaveNoSurname
	}
	return nil
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Saver
type Saver interface {
	Save(model.Person) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Enricher
type Enricher interface {
	Enrich(p model.Person) (model.Person, error)
}

func New(log *slog.Logger, saver Saver, enricher Enricher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.person.create"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Req
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Debug("request prepare", "err", err.Error())
			response.SendErrorBadReq(r, w, "bad json")
			return
		}

		log.Debug("request", "body", req)

		if err := req.Validate(); err != nil {
			log.Debug("request validate error", "err", err.Error())
			response.SendErrorBadReq(r, w, err.Error())
			return
		}

		person := model.Person{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
		}

		var err error
		person, err = enricher.Enrich(person)
		if errors.Is(err, enrich.ErrInvalidParameter) {
			log.Debug("enricher validate error", "err", err.Error())
			response.SendErrorBadReq(r, w, "")
			return
		}
		if err != nil {
			log.Error("enrich person error", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		person.ID, err = saver.Save(person)
		if err != nil {
			log.Error("save person error", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		log.Debug("response prepare completed", "person", person)

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, Resp{
			Response: response.Ok(),
			ID:       person.ID,
		})
	}
}
