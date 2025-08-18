package handlers

import (
	"encoding/json"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	crypto2 "github.com/Aleksei-D/go-loyalty-system/internal/utils/crypto"
	"go.uber.org/zap"
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

		ok, err := u.us.IsExist(r.Context(), user.Login)
		if err != nil {
			logger.Log.Warn(err.Error(), zap.Error(err))
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if ok {
			http.Error(w, "User already Exist", http.StatusConflict)
			return
		}

		_, err = u.us.Create(r.Context(), &user)
		if err != nil {
			logger.Log.Warn("User Create Error", zap.Error(err))
			http.Error(w, "create user error", http.StatusInternalServerError)
			return
		}

		tokenString, err := crypto2.CreateToken(user.Login, u.secretKey)
		if err != nil {
			logger.Log.Warn(err.Error(), zap.Error(err))
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

		existUser, err := u.us.GetByLogin(r.Context(), loginUser.Login)
		if err != nil {
			logger.Log.Warn(err.Error(), zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		if existUser == nil || loginUser.Password != existUser.Password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		tokenString, err := crypto2.CreateToken(existUser.Login, u.secretKey)
		if err != nil {
			logger.Log.Warn(err.Error(), zap.Error(err))
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Authorization", tokenString)
		w.WriteHeader(http.StatusOK)
	}
}
