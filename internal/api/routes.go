// Package api provides the web interface
package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Routes is the api entrypoint
func (app *Config) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.prometheusMiddleware)
	mux.Mount("/api", app.apiRouter())

	return mux
}

func (app *Config) apiRouter() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/log", app.WriteLog)
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

// prometheusMiddleware ticks up prometheus counters
func (app *Config) prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// rctx := chi.RouteContext(r.Context())
		// path := rctx.RoutePattern()

		path := currentRoutePattern(r)

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

func currentRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available
		return pattern
	}

	routePath := r.URL.Path
	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()
	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// No matching pattern, so just return the request path.
		// Depending on your use case, it might make sense to
		// return an empty string or error here instead
		return routePath
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()
}
