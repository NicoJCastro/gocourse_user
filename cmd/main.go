package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NicoJCastro/gocourse_user/internal/user"
	"github.com/NicoJCastro/gocourse_user/pkg/bootstrap"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	//router
	router := mux.NewRouter()

	_ = godotenv.Load("../.env")
	//logger
	logger := bootstrap.InitLogger()

	//db
	db, err := bootstrap.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	pagLimitDef := os.Getenv("PAGINATION_LIMIT_DEFAUL")
	if pagLimitDef == "" {
		logger.Fatal("PAGINATION_LIMIT_DEFAUL is not set")
	}

	//Antes de usar el servicio, se debe crear el repositorio
	userRepo := user.NewRepository(logger, db)
	userService := user.NewService(logger, userRepo)
	userEndpoints := user.MakeEndpoints(userService, user.Config{LimPageDef: pagLimitDef})

	//user endpoints

	router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	adress := "localhost:" + port

	srv := &http.Server{
		Handler:      router,
		Addr:         adress,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
