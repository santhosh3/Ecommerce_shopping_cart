package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	config "github.com/santhosh3/ECOM/Config"
	"github.com/santhosh3/ECOM/types"
	"github.com/santhosh3/ECOM/utils"
)

type contextKey string

var UserKey contextKey = "user"

func GenerateAccessToken(userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.AccessJWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":    strconv.Itoa(int(userId)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Envs.AccessJWTSecret))
	if err != nil {
		return "", err
	}

	return string(tokenString), err
}

func GenerateJWT(secret []byte, userId uint64, expiration int64) (string, error) {
	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"UserId":    strconv.FormatUint(userId, 10),
		"expiresAt": time.Now().Add(time.Duration(expiration) * time.Second).Unix(),
	})

	// Sign the token with the given secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err // Return an empty string and the error if signing fails
	}
	return string(tokenString), nil // Return the signed token and nil for no error
}

// generateTokens creates both access and refresh tokens concurrently.
func GenerateTokens(userId uint64, store types.UserStore) (string, string, error) {
	var wg sync.WaitGroup
	var errChan = make(chan error, 2)

	var accessToken, refreshToken string

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		accessToken, err = GenerateJWT([]byte(config.Envs.AccessJWTSecret), uint64(userId), config.Envs.AccessJWTExpirationInSeconds)
		if err != nil {
			errChan <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		refreshToken, err = GenerateJWT([]byte(config.Envs.RefreshJWTSecret), uint64(userId), config.Envs.RefreshJWTExpirationInSeconds)
		if err != nil {
			errChan <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := store.LoggingUser(uint64(userId));
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}

func ValidateJWT(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)
		token, err := ValidateJWT(tokenString, []byte(config.Envs.AccessJWTSecret))

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

		if !user.Status {
			permissionDenied(w)
			return
		}

		// Add a user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		// Call the func if token is Valid
		handlerFunc(w, r)
	}
}


// VerifyRefreshToken checks the validity of the refresh token.
func VerifyRefreshToken(tokenString string, store types.UserStore) (int16, error) {
	// Validate the JWT token
	token, err := ValidateJWT(tokenString, []byte(config.Envs.RefreshJWTSecret))
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, fmt.Errorf("invalid refresh token: token is not valid")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid refresh token: failed to parse claims")
	}

	// Get UserId from claims
	userIdStr, ok := claims["UserId"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid refresh token: missing UserId in claims")
	}

	// Convert UserId to int
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token: UserId is not a valid integer")
	}

	// Retrieve user by ID
	user, err := store.GetUserById(int16(userId))
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token: user not found")
	}

	// Check user status
	if !user.Status {
		return 0, fmt.Errorf("invalid refresh token: user is inactive")
	}
	return int16(user.ID), nil
}