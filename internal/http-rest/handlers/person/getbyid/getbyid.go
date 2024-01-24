package getbyid

import (
	"errors"
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

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=GeterById
type GeterByID interface {
	GetByID(int64) (model.Person, error)
}

type Resp struct {
	response.Response
	Person model.Person `json:"person"`
}

func New(log *slog.Logger, geter GeterByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.person.getbyid"

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

		person, err := geter.GetByID(id)
		if errors.Is(storage.ErrNotFound, err) {
			log.Debug("person not found", "err", err.Error())
			response.SendErrorNotFound(r, w)
			return
		}
		if err != nil {
			log.Error("get person error", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Resp{
			Response: response.Ok(),
			Person:   person,
		})
	}
}
