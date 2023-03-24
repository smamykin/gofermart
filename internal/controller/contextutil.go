package controller

import "github.com/gin-gonic/gin"

const (
	ctxKeyCurrentUserId = "current_user_id"
)

func GetCurrentUserIDFromContext(c *gin.Context) int {
	currentUserID := c.GetInt(ctxKeyCurrentUserId)
	if currentUserID <= 0 {
		panic("cannot get current user id. check the endpoint is protected.")
	}
	return currentUserID
}

func SetCurrentUserIDToContext(userID int, c *gin.Context) {
	c.Set(ctxKeyCurrentUserId, userID)
}
