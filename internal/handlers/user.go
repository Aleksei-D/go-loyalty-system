package handlers

import (
	"encoding/json"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	crypto2 "github.com/Aleksei-D/go-loyalty-system/pkg/utils/crypto"
	"io"
	"net/http"
)

type UserHandlers struct {
	us        *service.UserService
	secretKey string
}

func NewUserHandlers(us *service.UserService, secretKey string) *UserHandlers {
	return &UserHandlers{us, secretKey}
}

func (u *UserHandlers) ApiUserRegisterHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf, &user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := crypto2.HashPassword(*user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, exists := u.us.GetByLogin(r.Context(), *user.Login); exists {
			http.Error(w, "User already Exist", http.StatusConflict)
			return
		}

		user.Password = &hashedPassword
		registeredUser, err := u.us.Create(r.Context(), &user)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		registeredUserJson, err := json.Marshal(registeredUser)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "service/json")
		w.WriteHeader(http.StatusOK)
		w.Write(registeredUserJson)
	}
}

func (u *UserHandlers) ApiUserLoginHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginUser models.User
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err = json.Unmarshal(buf, &loginUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, exists := u.us.GetByLogin(r.Context(), *loginUser.Login)
		if !exists || !crypto2.CheckPasswordHash(*loginUser.Password, *user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := crypto2.CreateToken(*user.Login, u.secretKey)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "invalid marshaling", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
		w.Header().Set("Content-Type", "service/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJson)
	}
}
