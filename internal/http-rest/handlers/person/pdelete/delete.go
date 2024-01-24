package pdelete

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Deleter
type Deleter interface {
	Delete(int64) error
}

func New(log *slog.Logger, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.person.delete"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("New request")

		id, err := strconv.ParseInt(chi.URLParam(r, "ID"), 10, 64)
		if err != nil {
			log.Debug("Parse ID", "err", err.Error())
			response.SendErrorBadReq(r, w, response.ErrWrongID)
			return
		}

		err = deleter.Delete(id)
		if errors.Is(err, storage.ErrNotFound) {
			log.Debug("person not found", "err", err.Error())
			response.SendErrorBadReq(r, w, response.ErrNotFound)
		}
		if err != nil {
			log.Error("can`t delete person", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.Ok())
	}
}
