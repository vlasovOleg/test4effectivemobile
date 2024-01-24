package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"test4effectivemobile/internal/config"
	"test4effectivemobile/internal/enrich"
	pcreate "test4effectivemobile/internal/http-rest/handlers/person/create"
	pget "test4effectivemobile/internal/http-rest/handlers/person/get"
	pgetbyid "test4effectivemobile/internal/http-rest/handlers/person/getbyid"
	"test4effectivemobile/internal/http-rest/handlers/person/pdelete"
	pupdate "test4effectivemobile/internal/http-rest/handlers/person/update"
	"test4effectivemobile/internal/http-rest/handlers/ping"
	mwlogger "test4effectivemobile/internal/http-rest/middleware"
	"test4effectivemobile/internal/lib/logger/handlers/slogpretty"
	"test4effectivemobile/internal/storage/sqlstore"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	logLvlLocal = "local"
	logLvlDebug = "debug"
	logLvlInfo  = "info"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := loadConfig()
	log := initLogger(cfg.LogLvl)

	log.Info("starting tt4em")
	log.Debug(
		"debug messages are enabled",
		"cfg", cfg,
	)

	strorage := sqlstore.New(mustInitPostgres(cfg.Postgresql.URL))
	enrich := enrich.New(http.Client{Timeout: cfg.Enrich.RequestTimeout})

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(cfg.CtxTimeout))

	router.Get("/ping", ping.New(log))
	router.Route("/person", func(r chi.Router) {
		r.Post("/", pcreate.New(log, strorage.Person(), enrich))
		r.Get("/{ID}", pgetbyid.New(log, strorage.Person()))
		r.Get("/", pget.New(log, strorage.Person()))
		r.Delete("/{ID}", pdelete.New(log, strorage.Person()))
		r.Patch("/{ID}", pupdate.New(log, strorage.Person()))
	})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
	}

	log.Info("starting server", "address", cfg.Address)

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("ListenAndServe.ServerClosed")
			return
		}
		if err != nil {
			log.Error("failed to start server", err)
			done <- os.Kill
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Error("stopping server", "http server shutdown", err)
	}

	if err := strorage.Close(); err != nil {
		log.Error("stopping server error", "error", err)
	} else {
		log.Info("strorage closed")
	}
}

func initLogger(loglvl string) *slog.Logger {
	var log *slog.Logger
	switch loglvl {
	case logLvlLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := opts.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	case logLvlDebug:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case logLvlInfo:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func mustInitPostgres(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	return db
}

func loadConfig() config.Config {
	getEnv := func(key string, defaultVal string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return defaultVal
	}

	getEnvAsDuration := func(key string, defaultVal time.Duration) time.Duration {
		value, exists := os.LookupEnv(key)
		if !exists {
			return defaultVal
		}
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return defaultVal
		}
		return time.Duration(i)
	}

	return config.Config{
		LogLvl: getEnv("LOGLVL", logLvlLocal),
		Postgresql: config.Postgresql{
			URL: getEnv("POSTGRESQL_URL", ""),
		},
		HTTPServer: config.HTTPServer{
			Address:      getEnv("SERVER_ADDR", ":8080"),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 5000000000),  //nolint:gomnd // defaultVal
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 5000000000), //nolint:gomnd // defaultVal
			CtxTimeout:   getEnvAsDuration("SERVER_CTX_TIMEOUT", 6000000000),   //nolint:gomnd // defaultVal
		},
		Enrich: config.Enrich{
			RequestTimeout: getEnvAsDuration("ENRICH_TIMEOUT", 5000000000), //nolint:gomnd // defaultVal
		},
	}
}
