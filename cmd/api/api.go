package api

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/santhosh3/ECOM/services/auth"
	"github.com/santhosh3/ECOM/services/product"
	"github.com/santhosh3/ECOM/services/user"
	"gorm.io/gorm"
)


type APIServer struct {
	addr string
	db *gorm.DB
	rdb *redis.Client
}

func NewAPIServer(addr string, db *gorm.DB, rdb *redis.Client) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
		rdb: rdb,
	}
}

func (s *APIServer) Run() error {
	router :=  mux.NewRouter()
	router.Use(auth.LoggingMiddleware(s.rdb))
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	//StaticFiles for users
	router.PathPrefix("/api/v1/profile/").Handler(http.StripPrefix("/api/v1/profile/",http.FileServer(http.Dir("./uploads/profiles"))))

	//StaticFiles for products
	router.PathPrefix("/api/v1/productImage/").Handler(http.StripPrefix("/api/v1/productImage/",http.FileServer(http.Dir("./uploads/products"))))

	//User Routes
	userStore := user.NewStore(s.db);
	userHandler := user.NewHandler(userStore);
	userHandler.RegisterRoutes(subRouter);

	//Product Routes
	productStore := product.NewStore(s.db)
	ProductHandler := product.NewHandler(productStore)
	ProductHandler.ProductRoutes(subRouter)


	log.Println("App is Listening on port http://localhost", s.addr)
	return http.ListenAndServe(s.addr, router)
}
