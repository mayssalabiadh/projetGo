package utils

import "net/http"

type AppError struct {
	Code    string
	Message string
	Status  int
}

var (
	ErrInvalidCrendentials = AppError{
		Code:    "INVALID_CRENDENTIALS",
		Message: "Email ou mot de passe incorrect",
		Status:  http.StatusUnauthorized,
	}

	ErrUserNotFound = AppError{
		Code:    "USER_NOT_FOUND",
		Message: "Utilisatuer introuvable",
		Status:  http.StatusNotFound,
	}

	ErrRecordNotFound = AppError{
		Code:    "RECORD_NOT_FOUND",
		Message: "Enregistrement introuvable",
		Status:  http.StatusNotFound,
	}

	ErrBadRequest = AppError{
		Code:    "BAD_REQUEST",
		Message: "Requête invalide",
		Status:  http.StatusBadRequest,
	}

	ErrInternal = AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Une erreur interne est survenue",
		Status:  http.StatusInternalServerError,
	}

	ErrTokenMissing = AppError{
		Code:    "TOKEN_MISSING",
		Message: "Token manqant dans l'en-tête Authorization",
		Status:  http.StatusUnauthorized,
	}

	ErrInvalidToken = AppError{
		Code:    "INVALID_TOKEN",
		Message: "Token invalide ou expiré",
		Status:  http.StatusUnauthorized,
	}

	ErrAccessDenied = AppError{
		Code:    "ACCESS_DENIED",
		Message: "Accès refusé",
		Status:  http.StatusForbidden,
	}

	ErrValidationFailed = AppError{
		Code:    "VALIDATION_FAILED",
		Message: "Echec de validation des données envoyées",
		Status:  http.StatusUnprocessableEntity,
	}
)
