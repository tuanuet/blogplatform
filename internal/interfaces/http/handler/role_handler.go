package handler

import (
	"errors"
	"net/http"

	"github.com/aiagent/boilerplate/internal/application/dto"
	"github.com/aiagent/boilerplate/internal/application/usecase"
	"github.com/aiagent/boilerplate/pkg/response"
	"github.com/aiagent/boilerplate/pkg/validator"
	"github.com/gin-gonic/gin"
	validatorLib "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// RoleHandler handles role-related HTTP requests
type RoleHandler struct {
	roleUseCase usecase.RoleUseCase
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleUseCase usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{
		roleUseCase: roleUseCase,
	}
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.CreateRoleRequest true "Create role request"
// @Success 201 {object} response.Response{data=dto.RoleResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	validate := validatorLib.New()
	if err := validate.Struct(req); err != nil {
		errs := validator.FormatValidationErrors(err)
		if len(errs) > 0 {
			response.ValidationError(c, errs[0].Message)
			return
		}
	}

	role, err := h.roleUseCase.CreateRole(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleAlreadyExists) {
			response.Conflict(c, "Role already exists")
			return
		}
		response.InternalServerError(c, "Failed to create role")
		return
	}

	response.Success(c, http.StatusCreated, role)
}

// GetRole godoc
// @Summary Get a role by ID
// @Description Get a role by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Role ID"
// @Success 200 {object} response.Response{data=dto.RoleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	role, err := h.roleUseCase.GetRole(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalServerError(c, "Failed to get role")
		return
	}

	response.Success(c, http.StatusOK, role)
}

// ListRoles godoc
// @Summary List all roles
// @Description List all available roles
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response{data=[]dto.RoleResponse}
// @Router /api/v1/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.roleUseCase.ListRoles(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, "Failed to list roles")
		return
	}

	response.Success(c, http.StatusOK, roles)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update a role (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Update role request"
// @Success 200 {object} response.Response{data=dto.RoleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	role, err := h.roleUseCase.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalServerError(c, "Failed to update role")
		return
	}

	response.Success(c, http.StatusOK, role)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Role ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	if err := h.roleUseCase.DeleteRole(c.Request.Context(), id); err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalServerError(c, "Failed to delete role")
		return
	}

	c.Status(http.StatusNoContent)
}

// SetPermission godoc
// @Summary Set permission for a role
// @Description Set permission for a role on a specific resource (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Role ID"
// @Param request body dto.SetPermissionRequest true "Set permission request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/roles/{id}/permissions [post]
func (h *RoleHandler) SetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	var req dto.SetPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	validate := validatorLib.New()
	if err := validate.Struct(req); err != nil {
		errs := validator.FormatValidationErrors(err)
		if len(errs) > 0 {
			response.ValidationError(c, errs[0].Message)
			return
		}
	}

	if err := h.roleUseCase.SetPermission(c.Request.Context(), id, req); err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalServerError(c, "Failed to set permission")
		return
	}

	response.Success(c, http.StatusOK, map[string]string{"message": "Permission set successfully"})
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Get all roles assigned to a user
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId path string true "User ID"
// @Success 200 {object} response.Response{data=dto.UserRolesResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/users/{userId}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	roles, err := h.roleUseCase.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		response.InternalServerError(c, "Failed to get user roles")
		return
	}

	response.Success(c, http.StatusOK, roles)
}

// AssignRole godoc
// @Summary Assign role to user
// @Description Assign a role to a user (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId path string true "User ID"
// @Param request body dto.AssignRoleRequest true "Assign role request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{userId}/roles [post]
func (h *RoleHandler) AssignRole(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.roleUseCase.AssignRole(c.Request.Context(), userID, req.RoleID); err != nil {
		if errors.Is(err, usecase.ErrRoleNotFound) {
			response.NotFound(c, "Role not found")
			return
		}
		response.InternalServerError(c, "Failed to assign role")
		return
	}

	response.Success(c, http.StatusOK, map[string]string{"message": "Role assigned successfully"})
}

// RemoveRole godoc
// @Summary Remove role from user
// @Description Remove a role from a user (admin only)
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId path string true "User ID"
// @Param roleId path string true "Role ID"
// @Success 204 "No Content"
// @Failure 400 {object} response.Response
// @Router /api/v1/users/{userId}/roles/{roleId} [delete]
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	if err := h.roleUseCase.RemoveRole(c.Request.Context(), userID, roleID); err != nil {
		response.InternalServerError(c, "Failed to remove role")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetMyPermission godoc
// @Summary Get my permission on a resource
// @Description Get the current user's permission on a specific resource
// @Tags Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param resource query string true "Resource name (e.g., blogs, categories)"
// @Success 200 {object} response.Response{data=dto.UserPermissionResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/permissions [get]
func (h *RoleHandler) GetMyPermission(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID")
		return
	}

	resource := c.Query("resource")
	if resource == "" {
		response.BadRequest(c, "Resource query parameter is required")
		return
	}

	perm, err := h.roleUseCase.GetUserPermission(c.Request.Context(), uid, resource)
	if err != nil {
		response.InternalServerError(c, "Failed to get permission")
		return
	}

	response.Success(c, http.StatusOK, perm)
}
