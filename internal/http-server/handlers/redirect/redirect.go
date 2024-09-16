package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"urlShortener/internal/http-server/handlers"
	"urlShortener/internal/storage"
	"urlShortener/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func NewRedirect(log *slog.Logger, urlGetter handlers.URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const redirecter = "handlers.redirect.NewRedirect"

		log := log.With(
			slog.String("redirecter", redirecter),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "invalid request, missing alias",
			})

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "url not found",
			})

			return
		}
		if err != nil {
			log.Error("failed to get url", util.SlogErr(err))

			render.JSON(w, r, handlers.Response{
				Status: "renderError",
				Error:  "internal error",
			})

			return
		}

		log.Info("got url", slog.String("url", resURL))

		// redirect to found url
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
