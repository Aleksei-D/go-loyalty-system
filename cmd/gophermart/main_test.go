package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/config"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain/mocks"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/router"
	"github.com/Aleksei-D/go-loyalty-system/internal/service"
	crypto2 "github.com/Aleksei-D/go-loyalty-system/internal/utils/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	newLogin      = "newLogin"
	newPassword   = "newPassword"
	existLogin    = "existLogin"
	existPassword = "existPassword"
	validNewOrder = "9278923470"
	existOrder    = "12345678903"
	existOrder2   = "346436439"
	inValidOrder  = "1234567893"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body []byte, metadata map[string]string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	for key, value := range metadata {
		req.Header.Add(key, value)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		code         int
		method       string
		login        string
		password     string
		assertHeader func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:         "positive test for register user handler",
			path:         "/api/user/register",
			method:       http.MethodPost,
			code:         http.StatusOK,
			login:        newLogin,
			password:     newPassword,
			assertHeader: assert.NotEqual,
		},
		{
			name:         "negative test for register user handler",
			path:         "/api/user/register",
			method:       http.MethodPost,
			code:         http.StatusConflict,
			login:        existLogin,
			password:     existPassword,
			assertHeader: assert.Equal,
		},
		{
			name:         "positive test for login user handler",
			path:         "/api/user/login",
			method:       http.MethodPost,
			code:         http.StatusOK,
			login:        existLogin,
			password:     existPassword,
			assertHeader: assert.NotEqual,
		},
		{
			name:         "negative test for login user handler",
			path:         "/api/user/login",
			method:       http.MethodPost,
			code:         http.StatusUnauthorized,
			login:        newLogin,
			password:     newPassword,
			assertHeader: assert.Equal,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMockRepo := mocks.NewMockUserRepository(ctrl)
	balanceMockRepo := mocks.NewMockBalanceRepository(ctrl)
	orderMockRepo := mocks.NewMockOrderRepository(ctrl)
	withdrawalMockRepo := mocks.NewMockWithdrawalRepository(ctrl)
	serviceApp := service.NewService(
		balanceMockRepo,
		orderMockRepo,
		userMockRepo,
		withdrawalMockRepo,
	)
	userMockRepo.EXPECT().IsExist(gomock.Any(), newLogin).Return(false, nil).AnyTimes()
	userMockRepo.EXPECT().Create(gomock.Any(), &models.User{Login: newLogin, Password: newPassword}).Return(&models.User{Login: newLogin, Password: newPassword}, nil).AnyTimes()
	userMockRepo.EXPECT().IsExist(gomock.Any(), existLogin).Return(true, nil).AnyTimes()
	userMockRepo.EXPECT().GetByLogin(gomock.Any(), existLogin).Return(&models.User{Login: existLogin, Password: newPassword}, nil).AnyTimes()
	userMockRepo.EXPECT().GetByLogin(gomock.Any(), newLogin).Return(nil, nil).AnyTimes()
	err := config.InitDefaultEnv()
	require.NoError(t, err)

	cfg, err := config.InitConfig()
	require.NoError(t, err)

	ts := httptest.NewServer(router.NewRouter(serviceApp, *cfg.SecretKey))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			userJSON, err := json.Marshal(models.User{Login: v.login, Password: v.password})
			require.NoError(t, err)

			metadata := map[string]string{}
			resp, _ := testRequest(t, ts, v.method, v.path, userJSON, metadata)
			defer resp.Body.Close()

			assert.Equal(t, v.code, resp.StatusCode)
			token := resp.Header.Get("Authorization")
			v.assertHeader(t, "", token)
		})
	}
}

func TestAddOrderHandler(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		login       string
		orderNumber string
		token       bool
	}{
		{
			name:        "positive test for add order handler",
			code:        http.StatusAccepted,
			login:       existLogin,
			orderNumber: validNewOrder,
			token:       true,
		},
		{
			name:        "positive test for add exist order handler",
			code:        http.StatusOK,
			login:       existLogin,
			orderNumber: existOrder2,
			token:       true,
		},
		{
			name:        "negative test for add order another user handler",
			code:        http.StatusConflict,
			login:       existLogin,
			orderNumber: existOrder,
			token:       true,
		},
		{
			name:        "negative test for add invalid order handler",
			code:        http.StatusUnprocessableEntity,
			login:       existLogin,
			orderNumber: inValidOrder,
			token:       true,
		},
		{
			name:        "negative test for add order handler",
			code:        http.StatusUnauthorized,
			login:       existLogin,
			orderNumber: validNewOrder,
			token:       false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMockRepo := mocks.NewMockUserRepository(ctrl)
	balanceMockRepo := mocks.NewMockBalanceRepository(ctrl)
	orderMockRepo := mocks.NewMockOrderRepository(ctrl)
	withdrawalMockRepo := mocks.NewMockWithdrawalRepository(ctrl)
	serviceApp := service.NewService(
		balanceMockRepo,
		orderMockRepo,
		userMockRepo,
		withdrawalMockRepo,
	)
	nowDatetime := models.CustomTime{Time: time.Now()}
	orderMockRepo.EXPECT().IsExist(gomock.Any(), validNewOrder).Return(false, nil).AnyTimes()
	orderMockRepo.EXPECT().IsExist(gomock.Any(), existOrder2).Return(true, nil).AnyTimes()
	orderMockRepo.EXPECT().IsExist(gomock.Any(), existOrder).Return(true, nil).AnyTimes()
	orderMockRepo.EXPECT().Add(gomock.Any(), existLogin, validNewOrder).Return(&models.Order{Login: existLogin, Number: validNewOrder, Status: models.OrderStatusNew, UploadedAt: nowDatetime}, nil).AnyTimes()
	orderMockRepo.EXPECT().GetOrderByNumber(gomock.Any(), existOrder).Return(&models.Order{Login: newLogin, Number: existOrder, Status: models.OrderStatusNew, UploadedAt: nowDatetime}, nil).AnyTimes()
	orderMockRepo.EXPECT().GetOrderByNumber(gomock.Any(), existOrder2).Return(&models.Order{Login: existLogin, Number: existOrder2, Status: models.OrderStatusNew, UploadedAt: nowDatetime}, nil).AnyTimes()

	err := config.InitDefaultEnv()
	require.NoError(t, err)

	cfg, err := config.InitConfig()
	require.NoError(t, err)

	ts := httptest.NewServer(router.NewRouter(serviceApp, *cfg.SecretKey))
	defer ts.Close()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			var tokenString string
			if v.token {
				token, err := crypto2.CreateToken(v.login, *cfg.SecretKey)
				require.NoError(t, err)
				tokenString = token
			}

			metadata := make(map[string]string)
			metadata["Authorization"] = fmt.Sprintf("Bearer %s", tokenString)

			resp, _ := testRequest(t, ts, http.MethodPost, "/api/user/orders", []byte(v.orderNumber), metadata)
			defer resp.Body.Close()

			assert.Equal(t, v.code, resp.StatusCode)
		})
	}
}
