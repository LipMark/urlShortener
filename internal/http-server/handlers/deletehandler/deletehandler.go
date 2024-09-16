package deletehandler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"urlShortener/internal/http-server/handlers"
	"urlShortener/internal/storage"
	"urlShortener/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewDelete(log *slog.Logger, deleter handlers.URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const errDeleter = "handlers.deletehandler.NewDelete"

		log := log.With(
			slog.String("errDeleter", errDeleter),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("can't parse url id", util.SlogErr(err))
			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "invalid id",
			})
			return
		}

		err = deleter.DeleteURL(strconv.Itoa(id))

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url id not found", slog.Int("id", id))

			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "url id not found",
			})

			return
		}

		if err != nil {
			log.Error("failed to delete url", util.SlogErr(err))

			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "failed to delete url",
			})

			return
		}

		log.Info("url deleted", slog.Int("id", id))

		render.JSON(w, r, handlers.Response{
			Status: "OK",
		})
	}
}
