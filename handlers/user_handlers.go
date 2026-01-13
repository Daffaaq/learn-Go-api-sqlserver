// ---------------- handlers/user_handler.go ----------------
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"go-sqlserver-api/config"
	"go-sqlserver-api/helpers"
	"go-sqlserver-api/models"
	"go-sqlserver-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Repo      repositories.UserStore
	Validator *validator.Validate
}

// ------------------- GET USERS -------------------
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	users, err := h.Repo.GetUsersWithPagination(r.Context(), limit, offset)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users", err.Error())
		return
	}

	total, _ := h.Repo.CountUsers(r.Context()) // optional

	response := map[string]interface{}{
		"data": users,
		"meta": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	}

	helpers.RespondWithJSON(w, http.StatusOK, response)
}

// ------------------- CREATE USER -------------------
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := h.Validator.Struct(u); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	exists, err := h.Repo.IsEmailExists(r.Context(), u.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err.Error())
		return
	}
	if exists {
		helpers.RespondWithError(w, http.StatusBadRequest, "Email already exists", nil)
		return
	}

	if err := h.Repo.CreateUser(r.Context(), u); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created"})
}

// ------------------- UPDATE USER -------------------
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	u.ID = id

	if err := h.Validator.Struct(u); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	if err := h.Repo.UpdateUser(r.Context(), u); err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to update user", err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User updated"})
}

// ------------------- DELETE USER -------------------
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.Repo.DeleteUser(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			helpers.RespondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user", err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
}

// ------------------- HELPERS -------------------
func parsePagination(r *http.Request) (limit, offset int) {
	limit = config.DefaultLimit
	offset = 0

	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			if l > config.MaxLimit {
				l = config.MaxLimit
			}
			limit = l
		}
	}

	if oStr := r.URL.Query().Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
			offset = o
		}
	}
	return
}

func parseID(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "id"))
}