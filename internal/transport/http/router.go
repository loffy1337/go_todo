package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	var router *chi.Mux = chi.NewRouter()

	// Установка стандартных middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X_CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Установка путь health для того, чтобы понять, что http-сервер "жив"
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	// Установка путей основной API
	router.Route("/api/v1", func(api chi.Router) {
		// TODO: Вызывать контроллеры, когда сделаю
	})

	return router
}
