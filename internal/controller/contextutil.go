package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
)

const (
	ctxKeyCurrentUserID = "current_user_id"
)

func GetCurrentUserIDFromContext(c *gin.Context) (int, error) {
	currentUserID := c.GetInt(ctxKeyCurrentUserID)
	if currentUserID <= 0 {
		return 0, errors.New("cannot get current user id. check the endpoint is protected")
	}
	return currentUserID, nil
}

func SetCurrentUserIDToContext(userID int, c *gin.Context) {
	c.Set(ctxKeyCurrentUserID, userID)
}
