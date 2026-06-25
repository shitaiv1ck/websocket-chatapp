package users_http_transport

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	mock_users_service "github.com/shitaiv1ck/realtime-chat/internal/features/users/transport/http/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const userID = 1

func TestCreateUserHandler(t *testing.T) {
	type mockBehavior func(s *mock_users_service.MockUsersService, user domains.User)

	testCases := []struct {
		name                 string
		inputBody            string
		inputUser            domains.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Create user",
			inputBody: `{"username":"test", "password":"12345678"}`,
			inputUser: domains.User{
				Username: "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mock_users_service.MockUsersService, user domains.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(domains.User{
					ID:       1,
					Username: "test",
				}, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id":1, "username":"test", "is_online":false}`,
		},
		{
			name:                 "Invalid argument",
			inputBody:            `{"username":"test", "password":"1234567"}`,
			mockBehavior:         func(s *mock_users_service.MockUsersService, user domains.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"failed to validate json: invalid argument", "message":"failed to decode and validate"}`,
		},
		{
			name:      "User already exists",
			inputBody: `{"username":"test","password":"12345678"}`,
			inputUser: domains.User{
				Username: "test",
				Password: "12345678",
			},
			mockBehavior: func(s *mock_users_service.MockUsersService, user domains.User) {
				s.EXPECT().CreateUser(gomock.Any(), user).Return(domains.User{}, core_errors.ErrConflict)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"already exists", "message":"failed to create user"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_users_service.NewMockUsersService(c)
			testCase.mockBehavior(service, testCase.inputUser)

			transport := &UsersHTTPTransport{service: service}

			r := http.NewServeMux()
			r.Handle("POST /users", transport.CreateUserHandler())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestGetMeHandler(t *testing.T) {
	type mockBehavior func(s *mock_users_service.MockUsersService, userID int)

	testCases := []struct {
		name                 string
		inputUserID          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Get current user",
			inputUserID: userID,
			mockBehavior: func(s *mock_users_service.MockUsersService, userID int) {
				s.EXPECT().GetUser(gomock.Any(), userID).Return(domains.User{
					ID:       userID,
					Username: "test",
					IsOnline: true,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1, "username":"test", "is_online":true}`,
		},
		{
			name:        "Failed to authorize current user",
			inputUserID: userID,
			mockBehavior: func(s *mock_users_service.MockUsersService, userID int) {
				s.EXPECT().GetUser(gomock.Any(), userID).Return(domains.User{}, core_errors.ErrUnauthorize)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"error":"invalid user id", "message":"failed to authorize"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_users_service.NewMockUsersService(c)
			testCase.mockBehavior(service, testCase.inputUserID)

			transport := &UsersHTTPTransport{service: service}

			r := http.NewServeMux()
			r.Handle("GET /users/me", auth(transport.GetMeHandler()))

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users/me", nil)

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestGetUsersHandler(t *testing.T) {
	type mockBehavior func(s *mock_users_service.MockUsersService, search *string, limit *int, offset *int)

	testCases := []struct {
		name                 string
		inputSearch          *string
		inputLimit           *int
		inputOffset          *int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Get all users",
			mockBehavior: func(s *mock_users_service.MockUsersService, search *string, limit, offset *int) {
				s.EXPECT().GetUsers(gomock.Any(), search, limit, offset).Return([]domains.User{
					{
						ID:       1,
						Username: "test1",
						IsOnline: true,
					},
					{
						ID:       2,
						Username: "test2",
						IsOnline: false,
					},
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `[{"id":1, "username":"test1", "is_online":true},
			 {"id":2, "username":"test2", "is_online":false}]`,
		},
		{
			name: "Internal server error",
			mockBehavior: func(s *mock_users_service.MockUsersService, search *string, limit, offset *int) {
				s.EXPECT().GetUsers(gomock.Any(), search, limit, offset).Return([]domains.User{}, errors.New("internal server error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"internal server error", "message":"failed to get users"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_users_service.NewMockUsersService(c)
			testCase.mockBehavior(
				service,
				testCase.inputSearch,
				testCase.inputLimit,
				testCase.inputOffset,
			)

			transport := &UsersHTTPTransport{service: service}

			r := http.NewServeMux()
			r.Handle("GET /users", transport.GetUsersHandler())

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users", nil)

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestPatchUserHandler(t *testing.T) {
	type mockBehavior func(s *mock_users_service.MockUsersService, userID int, patch domains.UserPatch)

	testCases := []struct {
		name                 string
		inputUserID          int
		inputBody            string
		inputPatch           domains.UserPatch
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Patch current user",
			inputUserID: userID,
			inputBody:   `{"username":"test2", "old_password":"12345678", "new_password":"87654321"}`,
			inputPatch: domains.UserPatch{
				Username:    domains.Nullable[string]{Value: new("test2"), Set: true},
				OldPassword: domains.Nullable[string]{Value: new("12345678"), Set: true},
				NewPassword: domains.Nullable[string]{Value: new("87654321"), Set: true},
			},
			mockBehavior: func(s *mock_users_service.MockUsersService, userID int, patch domains.UserPatch) {
				s.EXPECT().PatchUser(gomock.Any(), userID, patch).Return(domains.User{
					ID:       userID,
					Username: "test2",
					IsOnline: true,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"id":1, "username":"test2", "is_online":true}`,
		},
		{
			name:        "User with input username already exists",
			inputUserID: userID,
			inputBody:   `{"username":"test2"}`,
			inputPatch: domains.UserPatch{
				Username: domains.Nullable[string]{Value: new("test2"), Set: true},
			},
			mockBehavior: func(s *mock_users_service.MockUsersService, userID int, patch domains.UserPatch) {
				s.EXPECT().PatchUser(gomock.Any(), userID, patch).Return(domains.User{}, core_errors.ErrConflict)
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: `{"error":"already exists", "message":"failed to patch user"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_users_service.NewMockUsersService(c)
			testCase.mockBehavior(service, testCase.inputUserID, testCase.inputPatch)

			transport := &UsersHTTPTransport{service: service}

			r := http.NewServeMux()
			r.Handle("PATCH /users", auth(transport.PatchUserHandler()))

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/users", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user_id", userID)))
	})
}
