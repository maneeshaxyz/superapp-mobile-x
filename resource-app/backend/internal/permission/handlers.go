package permission

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleCreatePermission(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreatePermissionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		permission := ResourcePermission{
			ResourceID:     req.ResourceID,
			GroupID:        req.GroupID,
			PermissionType: req.PermissionType,
		}

		if err := svc.CreatePermission(&permission); err != nil {
			statusCode, responseBody := mapPermissionErrorToResponse(err, "Failed to create permission")
			c.JSON(statusCode, responseBody)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"success": true, "data": permission})
	}
}

func HandleUpdatePermissionType(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req UpdatePermissionTypeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updated, err := svc.UpdatePermissionType(id, req.PermissionType)
		if err != nil {
			statusCode, responseBody := mapPermissionErrorToResponse(err, "Failed to update permission")
			c.JSON(statusCode, responseBody)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
	}
}

func HandleDeletePermission(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := svc.DeletePermission(id); err != nil {
			statusCode, responseBody := mapPermissionErrorToResponse(err, "Failed to delete permission")
			c.JSON(statusCode, responseBody)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": true})
	}
}

func HandleGetGroupPermissions(svc *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("id")

		permissions, err := svc.GetPermissionsByGroupID(groupID)
		if err != nil {
			statusCode, responseBody := mapPermissionErrorToResponse(err, "Failed to fetch group permissions")
			c.JSON(statusCode, responseBody)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": permissions})
	}
}
