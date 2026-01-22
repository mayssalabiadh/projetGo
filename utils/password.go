package utils

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(c *gin.Context, password string) (string, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur lors de la génération du password"})
		JSONAppError(c, ErrBadRequest, err)
		return "", err
	}
	return string(passwordHashed), err
}
