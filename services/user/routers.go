package user

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
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
	router.HandleFunc("/address", auth.WithJWTAuth(h.addAddress, h.store)).Methods(http.MethodPost)
	router.HandleFunc("/profile", auth.WithJWTAuth(h.GetUserById, h.store)).Methods(http.MethodGet)
	router.HandleFunc("/remove", auth.WithJWTAuth(h.DeleteUserById, h.store)).Methods(http.MethodDelete)
	router.HandleFunc("/update", auth.WithJWTAuth(h.UpdateUser, h.store)).Methods(http.MethodPut)
	router.HandleFunc("/forgetPassword", h.ForgetUserPassword).Methods(http.MethodPost)
	router.HandleFunc("/updatePassword", h.UpdateUserPassword).Methods(http.MethodPut)
	router.HandleFunc("/checkOTP", h.CheckOTP).Methods(http.MethodPost)
	router.HandleFunc("/generateAccessToken", h.GenerateAccessToken).Methods(http.MethodPost)
	router.HandleFunc("/logout", h.LogoutUser).Methods(http.MethodDelete)
}

func (h *Handler) CheckOTP(w http.ResponseWriter, r *http.Request) {
	var creds types.CheckOTPPayload
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	status, err := h.store.CheckOTPByEmail(creds.Email, creds.OTP)
	if !status {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("%s", err))
		return
	}
	utils.WriteJSON(w, http.StatusBadRequest, map[string]bool{"status": status})
}

func (h *Handler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	var creds types.RefreshTokenPayload
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(creds); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	userId, err := auth.VerifyRefreshToken(creds.Token, h.store)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	err = h.store.LogOutUser(int16(userId))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"error": false, "message": "Logged Out Successfully"})
}

func (h *Handler) GenerateAccessToken(w http.ResponseWriter, r *http.Request) {
	var creds types.RefreshTokenPayload
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(creds); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	userId, err := auth.VerifyRefreshToken(creds.Token, h.store)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	//secret []byte, userId int, expiration int64
	accessToken, err := auth.GenerateJWT([]byte(config.Envs.AccessJWTSecret), uint64(userId), config.Envs.AccessJWTExpirationInSeconds)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Access Token created successfully", "accessToken": accessToken})
}

func (h *Handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var creds types.UpdatePasswordCreds
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(creds); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	if creds.Password != creds.ConfirmPassword {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("confirmed Password is not mathching"))
		return
	}

	message, err := h.store.UpdatePasswordByEmail(creds.Password, creds.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not able to update password %s", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "SUCCESS", "message": message})
}

func (h *Handler) ForgetUserPassword(w http.ResponseWriter, r *http.Request) {
	otp := utils.GenerateOTP()
	var creds types.ForgetUserPassword
	if err := utils.ParseJSON(r, &creds); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Environment configurations
	email := creds.Email
	HostPassword := config.Envs.HostPassword
	SMTPPort := int(config.Envs.SMTPPort)
	SMTPHost := config.Envs.SMTPHost
	HostMail := config.Envs.HostMail

	// Check if the email exists in the DB
	profile, err := h.store.GetUserByEmail(email)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid email, please register"})
		return
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, 2)

	// Goroutine 1: Send OTP via email
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := utils.SendOTP(SMTPPort, otp, SMTPHost, HostMail, email, HostPassword); err != nil {
			errChan <- err
		}
	}()

	// Goroutine 2: Insert OTP into DB
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := h.store.InsertOTP(*profile, strconv.Itoa(otp)); err != nil {
			errChan <- err
		}
	}()

	// Goroutine 3: Background OTP Removal after 5 minutes
	go func() {
		h.store.RemoveOTP(*profile)
	}()

	// Wait for the first two goroutines to finish
	wg.Wait()
	close(errChan)

	// Handle any errors that occurred in the first two operations
	for err := range errChan {
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"email": email, "error": "An error occurred, please try again"})
			return
		}
	}

	// Send success response
	utils.WriteJSON(w, http.StatusOK, map[string]string{"success": "OTP sent successfully"})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) addAddress(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.store.CreateAddress(address)
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
	user, err := h.store.GetUserById(int16(userID))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "user not found"})
		return
	}
	userMap := map[string]interface{}{
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"email":         user.Email,
		"profile_image": user.ProfileImage,
		"phone_number":  user.PhoneNumber,
	}
	utils.WriteJSON(w, http.StatusOK, userMap)
}

func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	// Taking userId from middleware
	userID, ok := r.Context().Value(auth.UserKey).(uint64)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("user ID is missing or of incorrect type"))
		return
	}
	error, err := h.store.DeleteUserById(userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Not able to delete user"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": error})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var filePath string
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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Email already exists please try to Login"})
		return
	}
	var checkFile bool
	//Checking Profile Image exists from request
	file, handler, err := r.FormFile("profile_image")
	if err != nil {
		user = models.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
			ProfileImage: fmt.Sprintf("/api/v1/profile/%s", "abc.jpeg"),
			PhoneNumber:  r.FormValue("phone"),
			Email:        r.FormValue("email"),
		}
		checkFile = false
	} else {
		defer file.Close()
		// Create a unique file name and save the file
		folderPath := "./uploads/profiles"

		//generate unique file name
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
		filePath = filepath.Join(folderPath, fileName)

		//create a user Modal
		user = models.User{
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Password:     r.FormValue("password"),
			ProfileImage: fmt.Sprintf("/api/v1/profile/%s", fileName),
			PhoneNumber:  r.FormValue("phone"),
			Email:        r.FormValue("email"),
		}
		checkFile = true
	}

	// hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		os.Remove(filePath)
		return
	}

	user.Password = hashedPassword

	// Creating the user in the user table
	Userprofile, err := h.store.CreateUser(user)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		os.Remove(filePath)
		return
	}

	if checkFile {
		// Creating a file in path of folder
		out, err := os.Create(filePath)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not able to create a file"})
			return
		}
		defer out.Close()

		// write a file content to a new file
		_, err = io.Copy(out, file)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not able to create a file"})
			return
		}
	}

	accessToken, refreshToken, err := auth.GenerateTokens(Userprofile.ID, h.store)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		os.Remove(filePath)
		return
	}
	Userprofile.AccessToken = accessToken
	Userprofile.RefreshToken = refreshToken
	utils.WriteJSON(w, http.StatusCreated, Userprofile)
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

	accessToken, refreshToken, err := auth.GenerateTokens(user.ID, h.store)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "Logged in successfully", "accessToken": accessToken, "refreshToken": refreshToken})
}
