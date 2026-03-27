package permission

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mapPermissionErrorToResponse(err error, defaultMsg string) (int, gin.H) {
	switch {
	case errors.Is(err, ErrInvalidPermissionType):
		return http.StatusBadRequest, gin.H{"error": err.Error()}
	case errors.Is(err, ErrResourceNotFound), errors.Is(err, ErrGroupNotFound), errors.Is(err, ErrPermissionNotFound):
		return http.StatusNotFound, gin.H{"error": err.Error()}
	case errors.Is(err, ErrPermissionConflict):
		return http.StatusConflict, gin.H{"error": err.Error()}
	default:
		return http.StatusInternalServerError, gin.H{"error": defaultMsg}
	}
}