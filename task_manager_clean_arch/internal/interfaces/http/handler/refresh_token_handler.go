package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yiheyistm/task_manager/internal/domain"
)

type RefreshTokenHandler struct {
	RefreshTokenUsecase domain.RefreshTokenUsecase
}

func (rtc *RefreshTokenHandler) RefreshToken(c *gin.Context) {
	var request domain.RefreshTokenRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := rtc.RefreshTokenUsecase.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	username := claims["username"].(string)
	user, err := rtc.RefreshTokenUsecase.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response, err := rtc.RefreshTokenUsecase.GenerateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
