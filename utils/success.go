package utils

import "net/http"

type AppSuccessCRUD struct {
	Code    string
	Message string
	Status  int
}

var (
	SuccessRecordCreated = AppSuccessCRUD{
		Code:    "RECORD_CREATED",
		Message: "Enregistrement créé avec succès",
		Status:  http.StatusCreated,
	}

	SuccessRecordUpdated = AppSuccessCRUD{
		Code:    "RECORD_UPDATED",
		Message: "Enregistrement mis à jour avec succès",
		Status:  http.StatusOK,
	}

	SuccessRecordDelete = AppSuccessCRUD{
		Code:    "RECORD_DELETED",
		Message: "Enregistrement supprimé avec succès",
		Status:  http.StatusOK,
	}

	SuccessRecordFetched = AppSuccessCRUD{
		Code:    "RECORD_FETCHED",
		Message: "Enregistrement récupéré avec succès",
		Status:  http.StatusOK,
	}

	SuccessLogin = AppSuccessCRUD{
		Code:    "LOGIN_SUCCESSFUL",
		Message: "Connexion réussie",
		Status:  http.StatusOK,
	}
)
