package ping

import (
	"log/slog"
	"net/http"
	"test4effectivemobile/internal/http-rest/response"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Rsp struct {
	response.Response
	Message string `json:"message"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.ping"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		log.Info("ping")

		render.JSON(w, r, Rsp{
			Response: response.Ok(),
			Message:  "Pong",
		})
	}
}
