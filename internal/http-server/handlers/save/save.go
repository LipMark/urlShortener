package save

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"urlShortener/internal/random"
	"urlShortener/internal/storage"
	"urlShortener/util"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const customAliasLength = 6

// Request describes request to the server. validate is used to check if url is present
type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// Response describes server response
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

// ValidationError checks given request and returns it in a more readable format
func ValidationError(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return strings.Join(errMsgs, ", ")
}

func NewSave(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const saver = "handlers.save.NewSave"

		log := log.With(
			slog.String("saver", saver),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, Response{
				Status: "renderError",
				Error:  "empty request",
			})
			return
		}

		if err != nil {
			log.Error("failed to decode request body", util.SlogErr(err))

			render.JSON(w, r, Response{
				Status: "renderError",
				Error:  "failed to decode",
			})

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", util.SlogErr(err))

			render.JSON(w, r, ValidationError(validateErr))

			return
		}

		alias := req.Alias
		// create random alias if alias field is empty
		if alias == "" {
			alias = random.NewRandomAlias(customAliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, ("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", util.SlogErr(err))

			render.JSON(w, r, ("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Status: "OK",
			Alias:  alias,
		})
	}
}
