package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/santhosh3/ECOM/services/user"
)


type APIServer struct {
	addr string
	db *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router :=  mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db);
	userHandler := user.NewHandler(userStore);
	userHandler.RegisterRoutes(subRouter);

	log.Println("App is Listening on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
