package get

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"test4effectivemobile/internal/http-rest/response"
	"test4effectivemobile/internal/model"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=Geter
type Geter interface {
	Get(f model.Filter) ([]model.Person, error)
}

type Resp struct {
	response.Response
	Offset  uint64         `json:"offset"`
	Count   int            `json:"count"`
	Persons []model.Person `json:"persons"`
}

func New(log *slog.Logger, geter Geter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.person.get"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			log.Debug("parse query error", "err", err.Error())
			response.SendErrorBadReq(r, w, err.Error())
			return
		}

		filter, err := prepareFilter(values)
		if err != nil {
			log.Debug("prepare filter error", "err", err.Error())
			response.SendErrorBadReq(r, w, err.Error())
			return
		}

		log.Debug("filter is ready", "filter", filter)

		persons, err := geter.Get(filter)
		if err != nil {
			log.Error("get persons error", "err", err.Error())
			response.SendErrorInternal(r, w)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Resp{
			Response: response.Ok(),
			Offset:   filter.Offset,
			Count:    len(persons),
			Persons:  persons,
		})
	}
}

func prepareFilter(values url.Values) (model.Filter, error) {
	filter := model.Filter{}
	filter.Limit = 1000
	var err error

	if param, ok := values["agemax"]; ok {
		var age uint64
		age, err = strconv.ParseUint(param[0], 10, 32)
		filter.AgeMax = uint(age)
	}
	if err != nil {
		return model.Filter{}, fmt.Errorf("URL param parsing error %w", err)
	}

	if param, ok := values["agemin"]; ok {
		var age uint64
		age, err = strconv.ParseUint(param[0], 10, 32)
		filter.AgeMin = uint(age)
	}
	if err != nil {
		return model.Filter{}, fmt.Errorf("URL param parsing error %w", err)
	}

	if param, ok := values["gender"]; ok {
		filter.Gender = param[0]
	}

	if param, ok := values["country"]; ok {
		filter.Country = param[0]
	}

	if param, ok := values["offset"]; ok {
		filter.Offset, err = strconv.ParseUint(param[0], 10, 64)
	}
	if err != nil {
		return model.Filter{}, fmt.Errorf("URL param parsing error %w", err)
	}

	if param, ok := values["limit"]; ok {
		filter.Limit, err = strconv.ParseUint(param[0], 10, 64)
	}
	if err != nil {
		return model.Filter{}, fmt.Errorf("URL param parsing error %w", err)
	}

	return filter, nil
}
