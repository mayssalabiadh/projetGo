package handlers

import (
	"log"
	"net/http"
	"projet1/database"
	"projet1/models"
	"projet1/response"
	"projet1/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Créer une tâche
// @Description Création d'une tâche avec les champs JSON fournis
// @Tags Tâche
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param				task		body						models.Task		true		"Les données de tache à créer"
// @Success	201			{object}	utils.AppSuccessCRUD
// @Failure	400			{object}	utils.AppError 				"Requête invalide"
// @Failure	500			{object}	utils.AppError 				"Erreur interne"
// @Router /api/tasks/ [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	if err := database.DB.Create(&task).Error; err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de création"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}
	//c.JSON(http.StatusCreated, task)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordCreated, task)
}

// @Summary Extraire les tâches
// @Description Extraire les tâches avec de tout les utilisateurs
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Success 200 			{object} 	utils.AppSuccessCRUD
// @Failure 400				{object}	utils.AppError 				"Requête invalide"
// @Router /api/tasks/ [get]
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	//c.JSON(http.StatusOK, tasks)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, tasks)
}

// @Summary Extraire une tâche
// @Description Extraire une tâche avec son ID
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Param		id 			path		string			true		"ID de la tâche (UUID)"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Tâche introuvable"
// @Router /api/tasks/{id} [get]
func GetTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var task models.Task
	if err := database.DB.First(&task, "id = ?", id).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}
	//c.JSON(http.StatusOK, task)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, task)
}

// @Summary Mettre à jour une tâche
// @Description	Mettre à jour les informations d'une tâche avec son ID
// @Tags Tâche
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param		id 			path		string 			true 			"L'ID du tâche"
// @Param		task		body		models.Task		true			"Nouvelles données du Tâche"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Tâche introuvable"
// @Router  /api/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var task models.Task
	if err := database.DB.First(&task, "id = ?", id).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Tâche non trouvée"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}
	if err := c.ShouldBindJSON(&task); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	database.DB.Save(&task)
	//c.JSON(http.StatusOK, task)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordUpdated, task)
}

// @Summary Supprimer une tâche
// @Description Supprimer une tâche par son ID
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Param		id 			path		string 			true 			"L'ID du tâche"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Tâche introuvable"
// @Router  /api/users/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	database.DB.Delete(&models.Task{}, "id = ?", id)

	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordDelete, nil)
}

// @Summary 	Extraire les tâches avec pagination
// @Description	Extraire les tâches avec pagination, en fonction du page et limit
// @Tags 	Tâche
// @Security	BearerAuth
// @Produce json
// @Param   	page 		query		string 		false			"Les pages"
// @Param		limit		query		string		false			"La limite des elements"
// @Success		200 		{object}		utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Tâche introuvable"
// @Router  /api/tasks/paginated [get]
func GetPaginatedTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var tasks []models.Task
	database.DB.Limit(limit).Offset(offset).Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// @Summary filterer les tâches
// @Description filterer les tâches par l'ID de l'utilisateur et la status
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Param		user_id 			query		string 			false 			"L'ID de l'utilisateur"
// @Param		completed 			query		string 			false 			"Status du tâche"
// @Success		200 				{object}	utils.AppSuccessCRUD
// @Failure		400					{object}	utils.AppError 				"Requête invalide"
// @Router  /api/tasks/filtrer [get]
func FiltrerTask(c *gin.Context) {
	//Récuperation du UUID
	userID, err := uuid.Parse(c.DefaultQuery("user_id", ""))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Error in parsing uuid"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	//Récuperation du champ completed
	completed := (c.DefaultQuery("completed", ""))

	//initialisation de la requête
	query := database.DB.Model(&models.Task{})

	if userID != uuid.Nil {
		log.Println(userID)
		query = query.Where("user_id = ?", userID)
	}

	if completed != "" {
		tacheStatus, err := strconv.ParseBool(completed)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "error lors de la conversion du status de la tache"})
			utils.JSONAppError(c, utils.ErrBadRequest, err)
			return
		}
		query = query.Where("completed = ?", tacheStatus)
	}

	//Déclarer le slice des taches

	var tasks []models.Task
	if err := query.Find(&tasks).Error; err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "error lors de la récupération des taches"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//c.JSON(http.StatusOK, tasks)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, tasks)
}

// @Summary Taux de completion
// @Description Taux de completion des tâches par l'utilisateur
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Param		user_id 			path		string 			true 			"L'ID de l'utilisateur"
// @Success		200 				{object}	response.CompletionRate
// @Failure		400					{object}	utils.AppError 				"Requête invalide"
// @Router  /api/tasks/rate/{user_id} [get]
func CompletionRate(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur lors de parsing de l'UUID"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	var total, completed int64

	//Calculer le nombre total des taches pour cet utilisateur
	if err := database.DB.Model(&models.Task{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		log.Println(err.Error())
		//c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}

	//Calculer le nombre total des taches completes pour cet utilisateur
	if err := database.DB.Model(&models.Task{}).Where("user_id = ? AND completed = true", userID).Count(&completed).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, gin.H{"completion_rate": "0%"})
		return
	}

	rate := float64(completed) / float64(total) * 100

	CompletionRate := strconv.FormatFloat(rate, 'f', 2, 64) + "%"

	completion_rate := response.CompletionRate{
		Rate: CompletionRate,
	}

	//c.JSON(http.StatusOK, gin.H{"completion_rate": CompletionRate})
	utils.JSONAppSuccess(c, "C'est le taux de complétion des tâches", completion_rate)
}

// @Summary Filtrer les taches par date
// @Description Filtrer les taches par date avec limite
// @Tags Tâche
// @Security BearerAuth
// @Produce json
// @Param		start 				query		string 			false 			"Date de début"
// @Param		end 				query		string 			false 			"Date de fin"
// @Param		jours 				query		string 			false 			"Nombre des jours en arrière"
// @Success		200 				{array}		models.Task
// @Failure		400					{object}	utils.AppError 				"Requête invalide"
// @Router  /api/tasks/filtre_date [get]
func GetTasksByDate(c *gin.Context) {

	start := c.Query("start")
	end := c.Query("end")

	joursStr := c.Query("jours")

	var jours int
	var err error

	if joursStr == "" {
		jours = 30
	} else {
		jours, err = strconv.Atoi(joursStr)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur de la conversion de nb du jours"})
			utils.JSONAppError(c, utils.ErrBadRequest, err)
			return
		}
	}

	var startDate, endDate time.Time

	if start == "" {
		if jours == 0 || jours > 60 {
			startDate = time.Now().AddDate(0, 0, -30)
		} else {
			startDate = time.Now().AddDate(0, 0, -jours)
		}
	} else {
		today := time.Now().AddDate(0, 0, -60).Truncate(24 * time.Hour)
		startDate, err = time.Parse("2006-01-02", start)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur de la conversion de la date de début"})
			utils.JSONAppError(c, utils.ErrBadRequest, err)
			return
		}
		if startDate.Before(today) {
			startDate = time.Now().AddDate(0, 0, -30)
		}
	}

	if end == "" {
		endDate = time.Now()
	} else {
		today := time.Now().Truncate(24 * time.Hour)
		endDate, err = time.Parse("2006-01-02", end)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur de la conversion de la date de début"})
			utils.JSONAppError(c, utils.ErrBadRequest, err)
			return
		}
		if endDate.After(today) {
			endDate = time.Now()
		}
	}

	var tasks []models.Task
	if err := database.DB.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&tasks).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "erreur de la récupération des taches "})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}

	//c.JSON(http.StatusOK, tasks)
	utils.JSONAppSuccess(c, "Les tâches pour l'intervalle du date donnée", tasks)
}
