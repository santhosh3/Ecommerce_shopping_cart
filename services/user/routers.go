package user

import (
	"fmt"
	"net/http"
	"unsafe"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/santhosh3/ECOM/types"
	"github.com/santhosh3/ECOM/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST");
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	//var register types.RegisterUserPayload
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request)  {
	var creds types.LoginUser
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	
	if unsafe.Sizeof(creds) == 0 {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"msg" : "Request body is empty"})
	}

	if err := utils.Validate.Struct(creds); err != nil {
		errors := err.(validator.ValidationErrors);
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v",errors))
		return
	}

	fmt.Println(creds.Email);
	user, err := h.store.GetUserByEmail(creds.Email);
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"msg" : err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}