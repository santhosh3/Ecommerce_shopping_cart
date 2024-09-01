package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/santhosh3/ECOM/services/user"
	"gorm.io/gorm"
)


type APIServer struct {
	addr string
	db *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}

func (s *APIServer) Run() error {
	router :=  mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	//StaticFiles
	router.PathPrefix("/api/v1/profile/").Handler(http.StripPrefix("/api/v1/profile/",http.FileServer(http.Dir("./uploads/profiles"))))

	userStore := user.NewStore(s.db);
	userHandler := user.NewHandler(userStore);
	userHandler.RegisterRoutes(subRouter);

	log.Println("App is Listening on port http://localhost", s.addr)
	return http.ListenAndServe(s.addr, router)
}
