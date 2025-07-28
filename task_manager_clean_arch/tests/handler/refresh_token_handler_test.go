package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/suite"
	"github.com/yiheyistm/task_manager/internal/domain"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/dto"
	"github.com/yiheyistm/task_manager/internal/interfaces/http/handler"
	"github.com/yiheyistm/task_manager/mocks/mocks_security"
)

// RefreshTokenHandlerSuite defines the test suite for RefreshTokenHandler
type RefreshTokenHandlerSuite struct {
	suite.Suite
	mockRefreshTokenUsecase *mocks_security.IRefreshTokenUsecase
	handler                 *handler.RefreshTokenHandler
}

// SetupTest initializes the mocks and handler before each test
func (s *RefreshTokenHandlerSuite) SetupTest() {
	s.mockRefreshTokenUsecase = mocks_security.NewIRefreshTokenUsecase(s.T())
	s.handler = &handler.RefreshTokenHandler{
		RefreshTokenUsecase: s.mockRefreshTokenUsecase,
	}
}

// TestRefreshTokenHandlerSuite runs the test suite
func TestRefreshTokenHandlerSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenHandlerSuite))
}

// TestRefreshToken tests the RefreshToken method
func (s *RefreshTokenHandlerSuite) TestRefreshToken() {
	s.Run("Success", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{RefreshToken: "valid_refresh_token"}
		claims := jwt.MapClaims{"username": "abebe"}
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		tokens := domain.RefreshToken{AccessToken: "new_access_token", RefreshToken: "new_refresh_token"}
		s.mockRefreshTokenUsecase.On("ValidateRefreshToken", refreshTokenRequest.RefreshToken).Return(claims, nil)
		s.mockRefreshTokenUsecase.On("GetByUsername", "abebe").Return(user, nil)
		s.mockRefreshTokenUsecase.On("GenerateTokens", *user).Return(tokens, nil)

		body, _ := json.Marshal(refreshTokenRequest)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RefreshToken(c)

		s.Equal(http.StatusOK, w.Code)
		var response dto.RefreshTokenResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal(tokens.AccessToken, response.AccessToken)
		s.Equal(tokens.RefreshToken, response.RefreshToken)
	})

	s.Run("InvalidJSON", func() {
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", strings.NewReader("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RefreshToken(c)

		s.Equal(http.StatusBadRequest, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Contains(response["error"], "invalid character")
	})
	s.resetMocks()
	s.Run("InvalidRefreshToken", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{RefreshToken: "invalid_refresh_token"}
		s.mockRefreshTokenUsecase.On("ValidateRefreshToken", refreshTokenRequest.RefreshToken).Return(jwt.MapClaims{}, errors.New("invalid refresh token"))

		body, _ := json.Marshal(refreshTokenRequest)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("invalid refresh token", response["error"])
	})

	s.resetMocks()
	s.Run("UserNotFound", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{RefreshToken: "valid_refresh_token"}
		claims := jwt.MapClaims{"username": "abebe"}
		s.mockRefreshTokenUsecase.On("ValidateRefreshToken", refreshTokenRequest.RefreshToken).Return(claims, nil)
		s.mockRefreshTokenUsecase.On("GetByUsername", "abebe").Return(nil, errors.New("user not found"))

		body, _ := json.Marshal(refreshTokenRequest)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RefreshToken(c)

		s.Equal(http.StatusUnauthorized, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("user not found", response["error"])
	})
	s.resetMocks()

	s.Run("TokenGenerationError", func() {
		refreshTokenRequest := dto.RefreshTokenRequest{RefreshToken: "valid_refresh_token"}
		claims := jwt.MapClaims{"username": "abebe"}
		user := &domain.User{ID: "1", Username: "abebe", Email: "abebe@example.com", Role: "user"}
		s.mockRefreshTokenUsecase.On("ValidateRefreshToken", refreshTokenRequest.RefreshToken).Return(claims, nil)
		s.mockRefreshTokenUsecase.On("GetByUsername", "abebe").Return(user, nil)
		s.mockRefreshTokenUsecase.On("GenerateTokens", *user).Return(domain.RefreshToken{}, errors.New("token generation failed"))

		body, _ := json.Marshal(refreshTokenRequest)
		req := httptest.NewRequest(http.MethodPost, "/refresh-token", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		s.handler.RefreshToken(c)

		s.Equal(http.StatusInternalServerError, w.Code)
		var response gin.H
		json.Unmarshal(w.Body.Bytes(), &response)
		s.Equal("token generation failed", response["error"])
	})
}

func (s *RefreshTokenHandlerSuite) resetMocks() {
	s.mockRefreshTokenUsecase.Calls = nil
	s.mockRefreshTokenUsecase.ExpectedCalls = nil
}
