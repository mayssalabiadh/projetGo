package utils

import "github.com/gin-gonic/gin"

//la méthode classique pour renvoyer une erreur
func JSONError(c *gin.Context, status int, err error, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"error":   err.Error(),
		"message": message,
	})
}

//La méthode redcommandée
func JSONAppError(c *gin.Context, appErr AppError, err error) {
	c.JSON(appErr.Status, gin.H{
		"status":  appErr.Status,
		"message": appErr.Message,
		"code":    appErr.Code,
		"error":   err,
		"success": false,
	})
}
