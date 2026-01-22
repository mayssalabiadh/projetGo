package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"projet1/database"
	"projet1/models"
	"projet1/response"
	"projet1/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Créer un utilisateur
// @Description Création de l'utilisateur avec les champs JSON fournis
// @Tags Utilisateur
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param				user		body	models.User		true		"Les données de l'utilisateur à créer"
// @Success	201			{object}	utils.AppSuccessCRUD
// @Failure	400			{object}	utils.AppError 				"Requête invalide"
// @Failure	500			{object}	utils.AppError 				"Erreur interne"
// @Router /api/users/ [post]
func CreateUser(c *gin.Context) {
	var user models.User

	//Récuperation du JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		//utils.JSONError(c, http.StatusBadRequest, err, "Erreur de validation des données")
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Hashage du mot de passe
	hashedPassword, err := utils.HashPassword(c, user.Password)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de hachage du mot de passe"})
		//utils.JSONError(c, http.StatusUnauthorized, err, "Erreur de hachage du mot de passe")
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//L'aasertion du mot de passe
	user.Password = hashedPassword

	//Génération du UUID pour le nouvel utilisateur
	UserUUID := uuid.New()
	user.ID = UserUUID

	for i := range user.Tasks {
		user.Tasks[i].ID = uuid.New() //Création de l'ID de la tâche
		user.Tasks[i].UserID = UserUUID
	}

	//Création dans la base de données
	if err := database.DB.Create(&user).Error; err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de création"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	//c.JSON(http.StatusCreated, user)
	//utils.JSONResponseSuccess(c, http.StatusCreated, user, "Utilisateur créé avec succès")
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordCreated, user)
}

// @Summary Extraire les utilisateurs
// @Description Extraire les utilisateurs avec tous les tâches
// @Tags Utilisateur
// @Security BearerAuth
// @Produce json
// @Success 200 			{object} 	utils.AppSuccessCRUD
// @Failure 400				{object}	utils.AppError 				"Requête invalide"
// @Router /api/users/ [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	//L'utilisation de Preload
	if err := database.DB.Preload("Tasks").Find(&users).Error; err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de l'extraction des utilisateurs"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}
	//c.JSON(http.StatusOK, users)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, users)
}

// @Summary Extraire un utilisateur
// @Description Extraire un utilisatuer avec son ID
// @Tags Utilisateur
// @Security BearerAuth
// @Produce json
// @Param		id 			path		string			true		"ID de l'itulisateur (UUID)"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Utilisateur introuvable"
// @Router /api/users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	//Récuperation du param ID
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors du parse de l'uuid"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//charger l'utilisateur de la base de données
	if err := database.DB.Preload("Tasks").First(&user, "id = ?", userID).Error; err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur non trouvé"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, user)
}

// @Summary Mettre à jour un utilisateur
// @Description	Mettre à jour les informations d'un utilisateur avec son ID
// @Tags Utilisateur
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param		id 			path		string 			true 			"L'ID de l'utilisateur"
// @Param		user		body		models.User		true			"Nouvelles données utilisateur"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Utilisateur introuvable"
// @Router  /api/users/{id} [put]
func UpdateUser(c *gin.Context) {
	//Récuperation du param ID
	userID, err := uuid.Parse(c.Param("id"))
	log.Println(userID)
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors du parse de l'uuid"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	var user models.User
	//Vérification de l'existance de l'utilisateur et récuperation des information de l'utilisateur
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}

	//Liaison de la structure avec le JSON
	if err := c.ShouldBindJSON(&user); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Mise à jour de l'utilisateur
	res := database.DB.Save(&user)
	if res.Error != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de l'utilisateur"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordUpdated, user)
}

// @Summary Supprimer un utilisateur
// @Description Supprimer un utilisateur par son ID
// @Tags Utilisateur
// @Security BearerAuth
// @Produce json
// @Param		id 			path		string 			true 			"L'ID de l'utilisateur"
// @Success		200 		{object}	utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Utilisateur introuvable"
// @Router  /api/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	//Récuperation du param ID
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors du parse de l'uuid"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	var user models.User
	if err := database.DB.Delete(&user, "id = ?", userID).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Erreur lors de la suppression de l'utilisateur"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordDelete, nil)
}

// @Summary 	Extraire les utilisateurs avec pagination
// @Description	Extraire les utilisateurs avec pagination, en fonction du page et limit
// @Tags 	Utilisateur
// @Security	BearerAuth
// @Produce json
// @Param   	page 		query		string 		false			"Les pages"
// @Param		limit		query		string		false			"La limite des elements"
// @Success		200 		{object}		utils.AppSuccessCRUD
// @Failure		400			{object}	utils.AppError 				"Requête invalide"
// @Failure		404			{object}	utils.AppError 				"Utilisateur introuvable"
// @Router  /api/users/paginated_users [get]
func GetPaginatedUser(c *gin.Context) {

	//Récuperation du page
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de la récupération de la numéro du page"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	//Récuperation du limite
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de la récupération de la limite"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}
	//Calcul de l'offset
	offset := (page - 1) * limit

	var users []models.User
	if err := database.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateurs introuvables"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}

	//Retourner la résultat
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, users)
}

// @Summary Mise à jour partielle de l'utilisateur
// @Description Mise à jour partielle de l'utilisateur par son ID
// @Tags Utilisateur
// @Security 	BearerAuth
// @Accept json
// @Produce json
// @Param		id 				path			string						true		"L'ID de l'utilisateur"
// @Param		updateUser		body			response.UpdateUser			true		"Les nouvelles données de l'utilisateur"
// @Success		200 			{object}		utils.AppSuccessCRUD
// @Failure		400				{object}		utils.AppError 							"Requête invalide"
// @Failure		404				{object}		utils.AppError 							"Utilisateur introuvable"
// @Router  /api/users/{id} [patch]
func UpdateUserPartial(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur lors du parse de l'UUID"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "erreur lors de la récupération"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	var updateUser response.UpdateUser
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "requete invalide"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Mise à jour Partielle
	// if updateUser.Nom != "" && updateUser.Prenom != "" {
	// 	if err := database.DB.Model(&user).Updates(map[string]interface{}{
	// 		"Nom":    updateUser.Nom,
	// 		"Prenom": updateUser.Prenom,
	// 	}).Error; err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "mise à jour échouée"})
	// 		return
	// 	}
	// }

	if err := database.DB.Model(&user).Updates(models.User{
		Nom:    updateUser.Nom,
		Prenom: updateUser.Prenom,
	}).Error; err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "mise à jour échouée"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordUpdated, user)

}

// @Summary Chercher un utilisateur
// @Description Chercher un utilisateur par son Email
// @Tags Utilisateur
// @Security BearerAuth
// @Produce json
// @Param		email			query			string					true		"Email de l'utilisateur"
// @Success		200 			{object}		utils.AppSuccessCRUD
// @Failure		404				{object}		utils.AppError 							"Utilisateur introuvable"
// @Router /api/users/user_by_email [get]
func FindUserByEmail(c *gin.Context) {
	email := c.Query("email")

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		utils.JSONAppError(c, utils.ErrUserNotFound, err)
		return
	}
	//c.JSON(http.StatusOK, user)
	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, user)
}

// @Summary Connexion de l'utilisateur
// @Description Connexion de l'utilisateur avec email et mot de passe
// @Tags Utilisateur
// @Accept json
// @Produce json
// @Param		loginRequest	body			response.LoginRequest			true	"les coordonnées de l'utilisateurs"
// @Success		200 			{object}		utils.AppSuccessCRUD
// @Failure		400				{object}		utils.AppError 							"Requête invalide"
// @Failure		404				{object}		utils.AppError 							"Utilisateur introuvable"
// @Failure		401				{object}		utils.AppError 							"Les coordonnées invalides"
// @Router  	/login [post]
func LoginHandler(c *gin.Context) {

	// declarer la structure a recevoir (email, password)
	var loginRequest response.LoginRequest

	//Vérification du JSON avec shouldbindJSON
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"error": "Email ou mt de passe invalide"})
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Recherche par email
	var user models.User
	if err := database.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur introuvable"})
		utils.JSONAppError(c, utils.ErrInvalidCrendentials, err)
		return
	}

	//verifier le mot de passe CompreHashAndPassword
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		//c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		utils.JSONAppError(c, utils.ErrInvalidCrendentials, err)
		return
	}

	//generation du token JWT
	//Besoin: GenerateToken
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	//retourner le token à l'utilisateur
	//c.JSON(http.StatusOK, gin.H{"token": token})
	utils.JSONAppSuccessCRUD(c, utils.SuccessLogin, token)
}

// @Summary Upload fu fichier
// @Description Upload d'un fichier avec l'ID de li'utilisateur. Les extensions autorisées sont : .pdf, .doc, .docx
// @Tags 		Utilisateur
// @Security	BearerAuth
// @Accept 		multipart/form-data
// @Produce		json
// @Param		user_id			path			string				true				"User ID"
// @Param		file			formData		file				true				"Fichier pour l'upload"
// @Success 	200 			{object}		models.File 							"Fichier importer avec succées"
// @Failure		400				{object}		utils.AppError 							"Requête invalide"
// @Failure		500				{object}		utils.AppError 							"Erreur Interne su serveur"
// @Router    	/api/users/upload_file/{user_id} [post]
func UploadFile(c *gin.Context) {

	//Récupération de l'ID de l'utilisateur
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Récupérer le fichier
	file, err := c.FormFile("file")
	if err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	// Vérification de l'extension
	var errExt error
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
		utils.JSONError(c, http.StatusBadRequest, errExt, "extension non autorisée")
		return
	}

	//Génération du nom unique pour le fichier
	fileID := uuid.New().String()
	filename := fileID + ext

	//Génération du path local
	localPath := fmt.Sprintf("upload/%s", filename)

	//Génération de l'URL
	URLPath := "/files/" + filename

	//Sauvegarde physique et dans la base de données
	//1ere etape
	if err := c.SaveUploadedFile(file, localPath); err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	//2eme etape
	//Création de l'objet file
	newFile := models.File{
		FileName: file.Filename,
		FileType: ext,
		Size:     file.Size,
		Path:     localPath,
		URL:      URLPath,
		UserID:   userID,
	}

	//Enregistrement de l'objet
	if err := database.DB.Create(&newFile).Error; err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	//retourner la réponse
	utils.JSONAppSuccess(c, "File uploaded", newFile)
}

// @Summary Servir un fichier de la base de données
// @Description	Servir un fichier de la base de données avec son ID
// @Tags Utilisateur
// @Security BearerAuth
// @Produce octet-stream
// @Param	file_id				path			string						true		"File ID"
// @Success 200					{file}			string									"File Content"
// @Failure		400				{object}		utils.AppError 							"Requête invalide"
// @Failure		500				{object}		utils.AppError 							"Erreur Interne su serveur"
// @Router 		/api/users/get_file/{file_id}  [get]
func ServeFile(c *gin.Context) {
	//Récuperation de l'ID fu fichier
	fileID, err := uuid.Parse(c.Param("file_id"))
	fmt.Println(fileID)
	if err != nil {
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Chercher le fichier par son ID dans la base de données
	var file models.File
	if err := database.DB.Model(&models.File{}).Where("id = ?", fileID).First(&file).Error; err != nil {
		utils.JSONAppError(c, utils.ErrRecordNotFound, err)
		return
	}

	//Ajouter l'URL dans le contexte
	c.File(file.Path) //Le chemin vers le fichier
}

// @Summary 		Récupérer l'utilisateur et ces fichiers
// @Description 	Récupérer l'utilisateur et ces fichiers avec le user ID
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Param			user_id			path			string									true		"L'ID de l'utilisateur"
// @Success			200			 	{object}		map[string]interface{}
// @Failure			400				{object}		utils.AppError 							"Requête invalide"
// @Failure			404				{object}		utils.AppError 							"Utilisateur introuvable"
// @Failure			500				{object}		utils.AppError 							"Erreur Interne su serveur"
// @Router  /api/users/user_files/{user_id}  [get]
func GetUserFiles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	var user models.User
	if err := database.DB.Preload("Files").First(&user, "id = ?", userID).Error; err != nil {
		utils.JSONAppError(c, utils.ErrRecordNotFound, err)
		return
	}

	utils.JSONAppSuccess(c, "l'utilisateur avec ces documents", user)
}

// @Summary 			Récupération des fichiers
// @Description 		Récupération des fichiers avec pagination
// @Tags				Utilisateur
// @Security			BearerAuth
// @Produce				json
// @Param				page		 	query 				string 				false 				"le numero du page"
// @Param				limit			query				string				false				"la limite des elements"
// @Success				200				{object}			map[string]interface{}
// @Failure				400				{object}			utils.AppError 							"Requête invalide"
// @Failure				404				{object}			utils.AppError 							"Utilisateur introuvable"
// @Failure				500				{object}			utils.AppError 							"Erreur Interne su serveur"
// @Router 				/api/users/paginated_files  [get]
func PaginatedFiles(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		utils.JSONAppError(c, utils.ErrBadRequest, err)
		return
	}

	//Calcul de l'offset
	offset := (page - 1) * limit

	var files []models.File
	if err := database.DB.Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		utils.JSONAppError(c, utils.ErrRecordNotFound, err)
		return
	}

	utils.JSONAppSuccessCRUD(c, utils.SuccessRecordFetched, files)
}

// @Summary  		Récupération des utilisateurs avec ces résumés
// @Description		Récupération des utilisateurs avec ces résumés et le taux de completion des taches (Avec Goroutines)
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Success			200 				{object}					map[string]interface{}
// @Failure			404					{object}					utils.AppError 							"Utilisateur introuvable"
// @Router			/api/users/activity_overview_anonyme  [get]
func GetAllUsersActivity_anonyme(c *gin.Context) {
	start := time.Now()
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan response.UserResume)
	taskCh := make(chan response.TaskResult, len(users))

	for _, user := range users {
		wg.Add(1)
		//L'appel à la fonction
		go func(u models.User) {
			//L'interieur de la fonction
			defer wg.Done()
			var tasks []models.Task
			//L'extraction des tâches
			err := database.DB.Where("user_id = ?", u.ID).Find(&tasks).Error
			if err != nil {
				//Envoie de l'erreur
				taskCh <- response.TaskResult{
					Tasks: tasks,
					Err:   err,
				}
			}

			//le total des tâches
			total := len(tasks)

			//le nombre des taches compléte
			completed := 0
			for _, t := range tasks {
				if t.Completed {
					completed++
				}
			}

			rate := 0.0
			//vérification si le total est > 0 (éviter la divison par 0)
			if total > 0 {
				rate = (float64(completed) / float64(total)) * 100
			}

			//L'envoie du res avec le channel
			ch <- response.UserResume{
				ID:               u.ID,
				Nom:              u.Nom,
				TotalTasks:       total,
				CompletedTasks:   completed,
				CompletedPercent: rate,
			}
		}(user) // Passer la copie des informations de user
	}

	//Fermeture du channel aprés que toutes les goroutines ont terminé
	go utils.WaitAndClose(ch, taskCh, &wg)

	//La deuxieme methode
	// go func() {
	// 	wg.Wait()
	// 	close(ch)
	// }()

	//La réception des envoies
	var resumes []response.UserResume
	for res := range ch {
		resumes = append(resumes, res)
	}

	//Récupérer les erreurs du channel
	for errRes := range taskCh {
		if errRes.Err != nil {
			utils.JSONAppError(c, utils.ErrRecordNotFound, errRes.Err)
			return
		}
	}

	duration := time.Since(start)

	fmt.Println("Durée : ", duration)
	utils.JSONAppSuccess(c, "C'est le résumé des utiilisateurs", resumes)

}

// @Summary  		Récupération des utilisateurs avec ces résumés
// @Description		Récupération des utilisateurs avec ces résumés et le taux de completion des taches (Avec Goroutines)
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Success			200 				{object}					map[string]interface{}
// @Failure			404					{object}					utils.AppError 							"Utilisateur introuvable"
// @Router			/api/users/activity_overview  [get]
func GetAllUsersActivity(c *gin.Context) {
	start := time.Now()
	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan response.UserResume)
	taskCh := make(chan response.TaskResult, len(users))

	for _, user := range users {
		wg.Add(1)
		//L'appel à la fonction
		go utils.UserOverviewFunc(ch, taskCh, &wg, user)
	}

	//Fermeture du channel aprés que toutes les goroutines ont terminé
	//go utils.WaitAndClose(ch, taskCh, &wg)

	//La deuxieme methode
	go func() {
		wg.Wait()
		close(ch)
		close(taskCh)
	}()

	//La réception des envoies
	var resumes []response.UserResume
	for res := range ch {
		resumes = append(resumes, res)
	}

	//Récupérer les erreurs du channel
	for errRes := range taskCh {
		if errRes.Err != nil {
			utils.JSONAppError(c, utils.ErrRecordNotFound, errRes.Err)
			return
		}
	}

	duration := time.Since(start)

	fmt.Println("Durée : ", duration)
	utils.JSONAppSuccess(c, "C'est le résumé des utiilisateurs", resumes)

}

// @Summary  		Récupération des utilisateurs avec ces résumés
// @Description		Récupération des utilisateurs avec ces résumés et le taux de completion des taches (Sans Goroutines)
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Success			200 				{object}					map[string]interface{}
// @Failure			404					{object}					utils.AppError 							"Utilisateur introuvable"
// @Router			/api/users/user_overview  [get]
func GetUsersActivity(c *gin.Context) {
	start := time.Now()

	var users []models.User

	if err := database.DB.Find(&users).Error; err != nil {
		utils.JSONAppError(c, utils.ErrInternal, err)
		return
	}

	var resumes []response.UserResume

	for _, user := range users {
		var tasks []models.Task
		//L'extraction des tâches
		if err := database.DB.Where("user_id = ?", user.ID).Find(&tasks).Error; err != nil {
			utils.JSONAppError(c, utils.ErrRecordNotFound, err)
			return
		}

		//le total des tâches
		total := len(tasks)

		//le nombre des taches compléte
		completed := 0
		for _, t := range tasks {
			if t.Completed {
				completed++
			}
		}

		rate := 0.0
		//vérification si le total est > 0 (éviter la divison par 0)
		if total > 0 {
			rate = (float64(completed) / float64(total)) * 100
		}

		//L'envoie du res avec le channel
		resume := response.UserResume{
			ID:               user.ID,
			Nom:              user.Nom,
			TotalTasks:       total,
			CompletedTasks:   completed,
			CompletedPercent: rate,
		}
		resumes = append(resumes, resume)
	}

	duration := time.Since(start)

	fmt.Println("Durée : ", duration)

	utils.JSONAppSuccess(c, "C'est le résumé des utiilisateurs", resumes)

}

// @Summary  		Global stat
// @Description		Global stat sans channels
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Success			200 				{object}					map[string]interface{}
// @Failure			404					{object}					utils.AppError 							"Utilisateur introuvable"
// @Router			/api/users/global_stat  [get]
func GlobalStats(c *gin.Context) {
	var (
		totalUsers     int64
		totalTasks     int64
		completedTasks int64

		errUsers     error
		errTasks     error
		errCompleted error
	)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		errUsers = database.DB.Model(&models.User{}).Count(&totalUsers).Error
	}()

	go func() {
		defer wg.Done()
		errTasks = database.DB.Model(&models.Task{}).Count(&totalTasks).Error
	}()

	go func() {
		defer wg.Done()
		errCompleted = database.DB.Model(&models.Task{}).Where("completed = ?", true).Count(&completedTasks).Error
	}()

	wg.Wait()

	if errUsers != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du comptage des utilisateurs"})
		return
	}
	if errTasks != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du comptage des tâches"})
		return
	}
	if errCompleted != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du comptage des tâches complétées"})
		return
	}

	completedPercent := 0.0
	if totalTasks > 0 {
		completedPercent = float64(completedTasks) * 100 / float64(totalTasks)
	}

	res := response.UserStat{
		TotalUser:      totalUsers,
		TotalTasks:     totalTasks,
		TotalCompleted: completedTasks,
		Rate:           completedPercent,
	}

	utils.JSONAppSuccess(c, "statistiques globale des utilisateurs", res)
}

// @Summary  		Global stat avec channel
// @Description		Global stat avec channels pour la gestion des erreurs
// @Tags			Utilisateur
// @Security		BearerAuth
// @Produce			json
// @Success			200 				{object}					map[string]interface{}
// @Failure			404					{object}					utils.AppError 							"Utilisateur introuvable"
// @Router			/api/users/global_stat_overview  [get]
func GlobalStats_channel(c *gin.Context) {
	var (
		totalUsers     int64
		totalTasks     int64
		completedTasks int64
	)

	var wg sync.WaitGroup

	statsChError := make(chan error, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		err := database.DB.Model(&models.User{}).Count(&totalUsers).Error
		statsChError <- err
	}()

	go func() {
		defer wg.Done()
		err := database.DB.Model(&models.Task{}).Count(&totalTasks).Error
		statsChError <- err
	}()

	go func() {
		defer wg.Done()
		err := database.DB.Model(&models.Task{}).Where("completed = ?", true).Count(&completedTasks).Error
		statsChError <- err
	}()

	//La fermeture du canal après que toutes les goroutines sont terminées
	go func() {
		wg.Wait()
		close(statsChError)
	}()

	for errCh := range statsChError {
		if errCh != nil {
			utils.JSONAppError(c, utils.ErrRecordNotFound, errCh)
			return
		}
	}

	completedPercent := 0.0
	if totalTasks > 0 {
		completedPercent = float64(completedTasks) * 100 / float64(totalTasks)
	}

	res := response.UserStat{
		TotalUser:      totalUsers,
		TotalTasks:     totalTasks,
		TotalCompleted: completedTasks,
		Rate:           completedPercent,
	}

	utils.JSONAppSuccess(c, "statistiques globale des utilisateurs", res)
}
