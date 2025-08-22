package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
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

func (u *UserHandlers) APIUserRegisterHandler() func(http.ResponseWriter, *http.Request) {
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

		_, err = u.us.CreateUser(r.Context(), &user)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrUserAlreadyExists):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		tokenString, err := u.us.GetToken(user.Login, u.secretKey)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}

func (u *UserHandlers) APIUserLoginHandler() func(http.ResponseWriter, *http.Request) {
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

		err = u.us.Login(r.Context(), &loginUser)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrInvalidCredentials):
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			default:
				http.Error(w, "server error", http.StatusInternalServerError)
			}
			return
		}

		tokenString, err := u.us.GetToken(loginUser.Login, u.secretKey)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}
