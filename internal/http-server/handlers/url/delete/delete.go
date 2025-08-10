package delete

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"
)

type URLDeleter interface {
	DeleteURL(ctx context.Context, alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(r.Context(), alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Error("url not found", slog.String("alias", alias))

				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, resp.Error("URL not found"))

				return
			}
			log.Error("failed to delete url", slog.String("error", err.Error()))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete url"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		w.WriteHeader(http.StatusNoContent)
	}
}
