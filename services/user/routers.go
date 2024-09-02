package user

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	config "github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/models"
	"github.com/santhosh3/ECOM/services/auth"
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
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.register).Methods(http.MethodPost)
	router.HandleFunc("/address", auth.WithJWTAuth(h.shippingAddress, h.store)).Methods(http.MethodPost)
	router.HandleFunc("/profile", auth.WithJWTAuth(h.GetUserById, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/remove", auth.WithJWTAuth(h.DeleteUserById, h.store)).Methods(http.MethodDelete)
	router.HandleFunc("/update", auth.WithJWTAuth(h.UpdateUser, h.store)).Methods(http.MethodPut)
}

func (h *Handler) ForgetUserPassword(w http.ResponseWriter, r *http.Request) {
	type ForgetUserPassword struct {
		email string
	}
	var creds ForgetUserPassword
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//Checking mail is present on DB or not
	profile, err := h.store.GetUserByEmail(creds.email)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"msg": "Invalid email please register"})
		return
	}
	otp, err := utils.SendOTP(int(config.Envs.SMTPPort), config.Envs.SMTPHost, config.Envs.HostMail,creds.email,config.Envs.HostPassword)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"msg": "Invalid email please register"})
		return
	}

	err = h.store.InsertOTP(*profile, string(otp))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"msg": "Invalid email please register"})
		return
	}
	fmt.Println(otp);

	//update otp to db
	//remove otp from db after 5 mins
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request)  {
	//Taking userId from middleware 
	userID, ok := r.Context().Value(auth.UserKey).(uint64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("userId is missing or of incorrect type"))
		return
	}
	var creds models.User
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.store.UpdateUserById(userID, creds)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to update user"))
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) shippingAddress(w http.ResponseWriter, r *http.Request) {
	//Taking userId from middleware 
	userID, ok := r.Context().Value(auth.UserKey).(uint64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user ID is missing or of incorrect type"))
		return
	}
	//Assigning to value address
	var address types.Address
	if err := utils.ParseJSON(r, &address); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	//Validating JSON of shipping Address
	if err := utils.Validate.Struct(address); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	address.ShippingAddress.UserID = userID
	address.BillingAddress.UserID = userID
	user, err := h.store.CreateAddress(address);
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	// Taking userId from middleware 
	userID, ok := r.Context().Value(auth.UserKey).(uint64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user ID is missing or of incorrect type"))
		return
	}
	err, user := h.store.GetUserById(int16(userID))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "register user not found"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, user);
}

func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request)  {
	// Taking userId from middleware 
	userID, ok := r.Context().Value(auth.UserKey).(uint64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user ID is missing or of incorrect type"))
		return
	}
	msg, err := h.store.DeleteUserById(userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Not able to delete user"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message":msg});
}



func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	user = models.User{
		FirstName:   r.FormValue("first_name"),
		LastName:    r.FormValue("last_name"),
		Password:    r.FormValue("password"),
		PhoneNumber: r.FormValue("phone"),
		Email:       r.FormValue("email"),
	}

	//Validating request body
	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	//Checking mail is present on DB or not
	_, err = h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"msg": "Email already exists please try to Login"})
		return
	}

	//Checking Profile Image exists from request
	file, handler, err := r.FormFile("profile_image")
	if err != nil {
		user = models.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
			ProfileImage: fmt.Sprintf("%s:%s/api/v1/profile/%s", config.Envs.PublicHost, config.Envs.Port, "abc.jpeg"),
			PhoneNumber:  r.FormValue("phone"),
			Email:        r.FormValue("email"),
		}
	} else {
		defer file.Close()
		// Create a unique file name and save the file
		folderPath := "./uploads/profiles"

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

		//create a user Modal
		user = models.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
			ProfileImage: fmt.Sprintf("%s:%s/api/v1/profile/%s", config.Envs.PublicHost, config.Envs.Port, fileName),
			PhoneNumber:  r.FormValue("phone"),
			Email:        r.FormValue("email"),
		}
	}

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user.Password = hashedPassword

	// Creating the user in the user table
	Userprofile, err := h.store.CreateUser(user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(Userprofile.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"user": Userprofile, "token": token})
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var creds types.LoginUser
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(creds); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	user, err := h.store.GetUserByEmail(creds.Email)
	if err != nil {
		utils.WriteJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	if !auth.ComparePasswords(user.Password, []byte(creds.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, int(user.ID))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"user": user, "token": token})
}
