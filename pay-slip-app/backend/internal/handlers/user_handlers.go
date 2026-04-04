package handlers

import (
	"encoding/json"
	"net/http"
	"pay-slip-app/internal/models"
	"pay-slip-app/internal/utils"
)

// ── User handlers ─────────────────────────────────────────────────────────────

// GetCurrentUser handles GET /api/me
func (h *PaySlipHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := mustGetUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	jsonResponse(w, http.StatusOK, user)
}

// GetUsers handles GET /api/users  [admin only]
func (h *PaySlipHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	currentUser := mustGetUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if currentUser.Role != models.UserRoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, users)
}

// GetUsersV2 handles GET /api/v2/users  [admin only]
func (h *PaySlipHandler) GetUsersV2(w http.ResponseWriter, r *http.Request) {
	currentUser := mustGetUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if currentUser.Role != models.UserRoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	limit, afterID, afterCreatedAt, err := utils.ParsePagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, total, err := h.UserService.GetUsers(limit, afterID, afterCreatedAt)
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	data := users
	var nextCursor *string
	if limit > 0 && len(users) > limit {
		data = users[:limit]
		last := data[limit-1]
		nextCursor = utils.EncodeCursor(last.CreatedAt, last.ID)
	}

	jsonResponse(w, http.StatusOK, models.UsersResponse{
		Data:       data,
		Total:      total,
		NextCursor: nextCursor,
	})
}

// UpdateUserRole handles PUT /api/users/{id}/role  [admin only]
func (h *PaySlipHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	currentUser := mustGetUser(r)
	if currentUser == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if currentUser.Role != models.UserRoleAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := r.PathValue("id")
	var req models.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.UserService.UpdateUserRole(userID, req.Role); err != nil {
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "User role updated successfully"})
}
