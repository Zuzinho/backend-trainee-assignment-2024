package main

import (
	"avito_hr/pkg/banner"
	"avito_hr/pkg/env"
	"avito_hr/pkg/handler"
	"avito_hr/pkg/middleware"
	"avito_hr/pkg/session"
	"avito_hr/pkg/user"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// docker network create avito_network
// docker build -t my_app .
// docker run --network avito_network -d --name avito_container -p 8080:8080 my_app
// docker compose up

func main() {
	log.SetLevel(log.DebugLevel)

	log.Debug("Initializing env vars...")
	env.MustInit()

	log.Debug("Getting env vars...")
	jwtConfig := env.MustJWTConfig()
	dbConnString := env.MustDBConnectionString()
	maxConnCount := env.MustMaxConnectionCount()
	port := env.MustPort()

	log.Debug("Connecting to db...")
	db, err := sql.Open("postgres", dbConnString)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Pinging db...")
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Debug("Trying up migration...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./tmp/migrate",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Debug("Setting max open conns for db...")
	db.SetMaxOpenConns(maxConnCount)

	log.Debug("Initializing repositories...")
	bannersRepo := banner.NewBannersDBRepository(db)
	bannersTempRepo := banner.NewBannersTempMemoryRepository(bannersRepo)
	usersRepo := user.NewUsersDBRepository(db)
	log.Debug("Initializing session manager...")
	sessionManager := session.NewSessionsManager(jwtConfig)
	log.Debug("Initializing middleware...")
	middle := middleware.NewMiddleware(sessionManager)

	log.Debug("Initializing handlers...")
	authHandler := handler.NewAuthHandler(usersRepo, sessionManager)
	appHandler := handler.NewAppHandler(bannersRepo, bannersTempRepo)

	log.Debug("Initializing router...")
	router := mux.NewRouter()

	router.HandleFunc("/user_banner", appHandler.GetUserBanner).Methods(http.MethodGet)
	router.HandleFunc("/banner", appHandler.GetBanners).Methods(http.MethodGet)
	router.HandleFunc("/banner", appHandler.CreateBanner).Methods(http.MethodPost)
	router.HandleFunc("/banner/{banner_id:[0-9]+}", appHandler.UpdateBanner).Methods(http.MethodPatch)
	router.HandleFunc("/banner/{banner_id:[0-9]+}", appHandler.DeleteBanner).Methods(http.MethodDelete)

	router.HandleFunc("/login", authHandler.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/register", authHandler.SignUp).Methods(http.MethodPost)
	middle.AddUnmonitoredQuery(http.MethodPost, "/login")
	middle.AddUnmonitoredQuery(http.MethodPost, "/register")

	log.Debug("Exec gorutine for temp memory repository...")
	ticker := time.NewTicker(5 * time.Minute)
	go bannersTempRepo.UpdateBanners(ticker)

	api := middle.PackInMiddleware(router)

	log.Debugf("Started on port :%s", port)

	log.Println(http.ListenAndServe(":"+port, api))
}
