package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/santhosh3/ECOM/database"
	globalmiddlewares "github.com/santhosh3/ECOM/services/globalMiddlewares"
	"github.com/santhosh3/ECOM/services/product"
	"github.com/santhosh3/ECOM/services/user"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
	rdb  *redis.Client
}

func NewAPIServer(addr string, db *gorm.DB, rdb *redis.Client) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		rdb:  rdb,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(globalmiddlewares.RateLimitingMiddleware(s.rdb))
	router.Use(globalmiddlewares.CorsMiddleware)
	router.Use(database.DBQueryTimeoutMiddleware(s.db))
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	//StaticFiles for users
	router.PathPrefix("/api/v1/profile/").Handler(http.StripPrefix("/api/v1/profile/", http.FileServer(http.Dir("./uploads/profiles"))))

	//StaticFiles for products
	router.PathPrefix("/api/v1/productImage/").Handler(http.StripPrefix("/api/v1/productImage/", http.FileServer(http.Dir("./uploads/products"))))

	//User Routes

	folderPath := "./uploads/profiles"
	//Ensure folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, os.ModePerm)
	}

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subRouter)

	//Product Routes

	//Ensure folder exists

	folderPath = "./uploads/products"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, os.ModePerm)
	}

	productStore := product.NewStore(s.db)
	ProductHandler := product.NewHandler(productStore)
	ProductHandler.ProductRoutes(subRouter)

	log.Println("App is Listening on port http://localhost", s.addr)
	return http.ListenAndServe(s.addr, router)
}
