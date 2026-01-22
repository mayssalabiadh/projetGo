package middleware

import (
	"projet1/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var errMiddleware error
		if authHeader == "" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "token manquant"})
			utils.JSONAppError(c, utils.ErrTokenMissing, errMiddleware)
			c.Abort() //pour bloquer
			return
		}

		//vérification que le header commence par "Bearer"
		var errTokenSplit error
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Format de token invalide"})
			utils.JSONAppError(c, utils.ErrInvalidToken, errTokenSplit)
			c.Abort() //pour bloquer
			return
		}

		//Extraire seulement le token
		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			utils.JSONAppError(c, utils.ErrInvalidToken, err)
			c.Abort() //pour bloquer
			return
		}

		//Ajouter les informations de l'utilisateur au contexte
		c.Set("user_id", claims.UserID)
		c.Next()
	}

}
