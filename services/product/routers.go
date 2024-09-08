package product

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	config "github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/models"
	"github.com/santhosh3/ECOM/types"
	"github.com/santhosh3/ECOM/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) ProductRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.GetAllProducts).Methods(http.MethodGet)
	router.HandleFunc("/product", h.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/product/{productId}", h.GetProductById).Methods(http.MethodGet)
	router.HandleFunc("/product/{productId}", h.DeleteProductById).Methods(http.MethodDelete)
	router.HandleFunc("/product/{productId}", h.UpdateProductById).Methods(http.MethodPut)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	// Extracting query parameters
	size := queryParams.Get("size")
	name := queryParams.Get("name")
	priceGreaterThan := queryParams.Get("priceGreaterThan")
	priceLessThan := queryParams.Get("priceLessThan")
	priceSort := queryParams.Get("priceSort")

	// Validate priceSort
	if priceSort != "" {
		if _, err := strconv.Atoi(priceSort); err != nil || (priceSort != "1" && priceSort != "-1") {
			// http.Error(w, `{"status":false, "message":"please provide 1 or -1"}`, http.StatusBadRequest)
			// return
			utils.WriteError(w, http.StatusBadRequest, err)
		}
	}

	// Call the store function to fetch products
	products, err := h.store.GetFilteredProducts( size, name, priceGreaterThan, priceLessThan, priceSort)
	if err != nil || len(products) == 0 {
		http.Error(w, `{"status":false, "message":"No product found"}`, http.StatusBadRequest)
		return
	}
}

func (h *Handler) UpdateProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["productId"]
	var product models.Product
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	product = models.Product{
		Title:          r.FormValue("title"),
		Description:    r.FormValue("description"),
		Price:          utils.ConvertStringToFloat((r.FormValue("price"))),
		IsFreeShipping: utils.ConvertStringToBool(r.FormValue("is_free_shipping")),
		CurrencyId:     r.FormValue("currency_id"),
		Installments:   int(utils.ConvertStringToFloat(r.FormValue("installments"))),
	}
	if len(r.FormValue("available_size")) != 0 {
		product.AvailableSize = utils.ConvertStringToArray(r.FormValue("available_size"))
	}

	file, handler, err := r.FormFile("product_image")
	if err != nil {
		product.ProductImage = ""
	} else {
		defer file.Close()
		product.ProductImage = handler.Filename
	}
	productInfo, err := h.store.UpdateProductById(uint64(utils.ConvertStringToFloat(productId)), product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to update product with ID %s", productId))
		return
	}
	utils.WriteJSON(w, http.StatusOK, productInfo)
}

func (h *Handler) DeleteProductById(w http.ResponseWriter, r *http.Request) {
	//take id from params
	vars := mux.Vars(r)
	productId := vars["productId"]
	_, err := h.store.DeleteProductById(int16(utils.ConvertStringToFloat(productId)))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"error": "product got deleted"})
}

func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	//take id from params
	vars := mux.Vars(r)
	productId := vars["productId"]
	product, err := h.store.GetProductById(int16(utils.ConvertStringToFloat(productId)))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	product.ProductImage = fmt.Sprintf("%s:%s%s", config.Envs.PublicHost, config.Envs.Port, product.ProductImage)
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusOK, err)
		return
	}

	product = models.Product{
		Title:          r.FormValue("title"),
		Description:    r.FormValue("description"),
		Price:          utils.ConvertStringToFloat((r.FormValue("price"))),
		IsFreeShipping: utils.ConvertStringToBool(r.FormValue("is_free_shipping")),
		AvailableSize:  utils.ConvertStringToArray(r.FormValue("available_size")),
		CurrencyId:     r.FormValue("currency_id"),
		Installments:   int(utils.ConvertStringToFloat(r.FormValue("installments"))),
	}

	file, handler, err := r.FormFile("product_image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product_image is missing: %v", err))
		return
	}
	defer file.Close()

	// Create a unique file name and save the file
	folderPath := "./uploads/products"

	//generate unique file name
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
	filePath := filepath.Join(folderPath, fileName)

	//Ensure folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, os.ModePerm)
	}

	//Creating a file in path of folder
	out, err := os.Create(filePath)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not able to create a file"})
		return
	}
	defer out.Close()

	//write a file content to a new file
	_, err = io.Copy(out, file)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not able to create a file"})
		return
	}

	product.ProductImage = fmt.Sprintf("/api/v1/productImage/%s", fileName)

	//Validating payload of
	if err := utils.Validate.Struct(product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		os.Remove(filePath)
		return
	}

	//Create Product
	prod, err := h.store.CreateProduct(product)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error occured: %v", err))
		os.Remove(filePath)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, prod)
}
