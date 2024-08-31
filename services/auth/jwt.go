package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/types"
	"github.com/santhosh3/ECOM/utils"
)
type contextKey string
var UserKey contextKey = "user"

func CreateJWT(secret []byte, userId int) (string, error)  {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId" : strconv.Itoa(int(userId)),
		"expiresAt" : time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return string(tokenString), err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		tokenString := utils.GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)

		//Any Error From Validation
		if err != nil {
			permissionDenied(w)
			return
		}

		//If Token is not Valid
		if !token.Valid {
			permissionDenied(w)
			return
		}
		
		claims := token.Claims.(jwt.MapClaims)
		str := claims["UserId"].(string)

		userId, err := strconv.Atoi(str)
		if err != nil {
			permissionDenied(w)
			return
		}

		user, err := store.GetUserById(int16(userId))
		if err != nil {
			permissionDenied(w)
			return
		}
		
		// Add a user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		// Call the func if token is Valid
		handlerFunc(w,r)
	}	
}