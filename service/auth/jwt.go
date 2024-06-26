package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gitKashish/ecommerce-api-go/config"
	"github.com/gitKashish/ecommerce-api-go/types"
	"github.com/gitKashish/ecommerce-api-go/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	// Create a new JWT Token establishing its signing method (not yet signed).
	// Along with mapped claims (using `jwt.MapClaims()` method).
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	// Signing the token with pre-determined signing method.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// JWT Authorization middleware.
// TODO : Change application of middleware from..
// Traditional closures -> Middleware Chaining.
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from the user request
		tokenString := getTokenFromRequest(r)

		// validate JWT token
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		// if it is we need to fetch the userID from the DB (id from the token)
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, _ := strconv.Atoi(str)

		user, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// set context "userID" to the userID.
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		// Execute wrapped HandlerFunc. It will now execute with...
		// updated context.
		handlerFunc(w, r)
	}
}

// Check if a token string received in a request is valid or not.
func validateToken(tokenString string) (*jwt.Token, error) {
	// `jwt.Parse` takes in the token string and the JWT secret to...
	// extract theh fields of the Token.
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// jwt.KeyFunc(this) returns the secrete key (loaded from environment)...
		// for parsing if the token signing etc. is proper.
		// If it returns `nil` then the `jwt.Parse()` method would return error...
		// indicating an invalid token.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})
}

// Util function to get token string from HTTP "Authorization" header.
func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

// Util function (to make code cleaner).
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

// Get authorized userID from current context.
// Should be used once context has been updated.
func GetUseIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
