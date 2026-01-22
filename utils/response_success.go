package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// la méthode classique pour renvoyer la réponse de succées
func JSONResponseSuccess(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, gin.H{
		"status":  status,
		"data":    data,
		"message": message,
	})
}

// La méthode recommandée pour les API CRUD
func JSONAppSuccessCRUD(c *gin.Context, appSuccessCRUD AppSuccessCRUD, data interface{}) {
	c.JSON(appSuccessCRUD.Status, gin.H{
		"status":  appSuccessCRUD.Status,
		"message": appSuccessCRUD.Message,
		"code":    appSuccessCRUD.Code,
		"data":    data,
		"success": true,
	})
}

// La méthode recommandée pour les autres APIs des fonctionnalitées
func JSONAppSuccess(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": message,
		"data":    data,
		"success": true,
	})
}
