package utils

import (
	"projet1/database"
	"projet1/models"
	"projet1/response"
	"sync"
)

func WaitAndClose(ch chan response.UserResume, taskCh chan response.TaskResult, wg *sync.WaitGroup) {
	wg.Wait()
	close(ch)
	close(taskCh)
}

func UserOverviewFunc(ch chan response.UserResume, taskCh chan response.TaskResult, wg *sync.WaitGroup, user models.User) {
	//L'interieur de la fonction
	defer wg.Done()
	var tasks []models.Task
	//L'extraction des tâches
	err := database.DB.Where("user_id = ?", user.ID).Find(&tasks).Error
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
		ID:               user.ID,
		Nom:              user.Nom,
		TotalTasks:       total,
		CompletedTasks:   completed,
		CompletedPercent: rate,
	}
}
