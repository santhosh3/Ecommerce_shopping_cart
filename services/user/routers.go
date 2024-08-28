package user

import (
	"net/http"
	"unsafe"

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

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request)  {
	var creds types.LoginUser
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if unsafe.Sizeof(creds) == 0 {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"msg" : "Request body is empty"})
	}
}