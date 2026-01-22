package routes

import (
	"net/http"
	"projet1/handlers"
	"projet1/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	//Route publique
	r.POST("/login", handlers.LoginHandler)

	//Routes protégées par middleware
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware()) //Syntaxe de l'appel du middleware
	{
		users := protected.Group("/users")
		{
			users.POST("/", handlers.CreateUser)
			users.GET("/", handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)

			//Routes Fonctionnalités
			users.GET("/paginated_users", handlers.GetPaginatedUser)
			users.PATCH("/:id", handlers.UpdateUserPartial)
			users.GET("/user_by_email", handlers.FindUserByEmail)
			users.POST("/upload_file/:user_id", handlers.UploadFile) //Route pour importer un fichier
			users.StaticFS("/files", http.Dir("./upload"))
			users.GET("/get_file/:file_id", handlers.ServeFile) //Route pour récuperer un fichier de la base
			users.GET("/user_files/:user_id", handlers.GetUserFiles)
			users.GET("/paginated_files", handlers.PaginatedFiles)
			users.GET("/activity_overview_anonyme", handlers.GetAllUsersActivity_anonyme)
			users.GET("/activity_overview", handlers.GetAllUsersActivity)
			users.GET("/user_overview", handlers.GetUsersActivity)
			users.GET("/global_stat", handlers.GlobalStats)
			users.GET("/global_stat_channel", handlers.GlobalStats_channel)

		}

		tasks := protected.Group("/tasks")
		{
			tasks.POST("/", handlers.CreateTask)
			tasks.GET("/", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)

			//Fonctionnalitées
			tasks.GET("/paginated", handlers.GetPaginatedTasks)
			tasks.GET("/filtrer", handlers.FiltrerTask)
			tasks.GET("/rate/:user_id", handlers.CompletionRate)
			tasks.GET("/filtre_date", handlers.GetTasksByDate)
		}
	}
}
