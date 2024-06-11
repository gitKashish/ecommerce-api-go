package user

import (
	"fmt"
	"net/http"

	"github.com/gitKashish/ecommerce-api-go/config"
	"github.com/gitKashish/ecommerce-api-go/service/auth"
	"github.com/gitKashish/ecommerce-api-go/types"
	"github.com/gitKashish/ecommerce-api-go/utils"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /login", h.handleLogin)
	router.HandleFunc("POST /register", h.handleRegister)
}

// ---- HandlerFunc for USER LOGIN & JWT GENERATION ----
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// General Flow :
	// 1. Parse request JSON to appropriate payload type.
	// 2. Validate payload structure.
	// 3. Get User by email.
	// 4. Password Verification.
	// 5. Generate and respond with JWT token & http.StatusOK.

	var payload types.LoginUserPayload

	// Decoding http.Request -> LoginUserPayload
	// returns error if request body is empty.

	/* LoginUserPayload Structure:
	Email string (required)
	Password string (required)
	*/
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Validating structure of payload by tags.
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// Getting types.User object from DB using Email from the payload.
	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	// Comparing Hashed password in types.User object & PlainText password in payload.
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	// Generating JWT token for session authentication...
	// if payload password is correct.
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// sending back the JWT authentication token once auth is completed.
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

// ---- HandlerFunc for REGISTERING NEW USER ----
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// General Flow :
	// 1. Parse request JSON to appropriate payload type.
	// 2. Validate payload structure.
	// 3. Checking if user already exists.
	// 4. If not, hashing user password.
	// 4. Create a new entry in the DB.
	// 5. Respond with http.StatusCreated.

	var payload types.RegisterUserPayload

	// Decoding http.Request -> RegisterUserPayload
	// returns error if request body is empty.

	/* RegisterUserPayload Structure :
	FirstName string (required)
	LastName string (required)
	Email string (required, email)
	Password string (required)
	*/
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validating structure of payload by tags.
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// Checking if user exists already...
	// if error == nil it means user already exists.
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	// Hashing the payload.Password.
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Creating a new user entry in the DB.
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Responding with http.StatusCreated.
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}
