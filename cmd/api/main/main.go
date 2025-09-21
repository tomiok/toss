package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tomiok/toss/cmd/api"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	run()
}

func run() {
	loggerSetup()
	r := chi.NewRouter()
	deps := api.NewDeps()
	routes(deps, r)

	srv := api.Server{
		Server: &http.Server{
			Addr: ":9999",
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r,
		},
	}

	srv.Start()
}

func routes(deps api.Deps, r chi.Router) {
	r.Use(middleware.Logger, AddCors())

	r.Route("/uploads", func(r chi.Router) {
		r.Get("/", do(deps.UploadHandler.UploadView))
		r.Post("/", do(deps.UploadHandler.Upload))
	})

	fileServer(r)
}

func AddCors() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Content-Type-Options"},
		AllowCredentials: false,
		MaxAge:           500,
	})
}

func fileServer(r chi.Router) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	fs(r, "/static", filesDir)
}

// fs conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("file server does not permit any URL parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		h := http.StripPrefix(pathPrefix, http.FileServer(root))
		h.ServeHTTP(w, r)
	})
}

func loggerSetup() {
	opts := &slog.HandlerOptions{
		// Use the ReplaceAttr function on the handler options
		// to be able to replace any single attribute in the log output
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// check that we are handling the time key
			if a.Key != slog.TimeKey {
				return a
			}

			t := a.Value.Time()

			// change the value from a time.Time to a String
			// where the string has the correct time format.
			a.Value = slog.StringValue(t.Format(time.DateTime))

			return a
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

type WebHandler func(w http.ResponseWriter, r *http.Request) error

func do(f WebHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			requestID := middleware.GetReqID(r.Context())
			slog.Info("cannot process request", slog.String("RequestID", requestID), slog.Any("err", err))

			//todo render err page here
			//web.ReturnErr(w, err)
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
