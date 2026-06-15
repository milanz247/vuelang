package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"vuelang/app/requests"
	"vuelang/app/services"
	"vuelang/internal/framework/response"
)

// UserController handles CRUD for /api/v1/users (admin only).
type UserController struct {
	userSvc *services.UserService
}

func NewUserController(userSvc *services.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

// Index  GET /api/v1/users
func (ctrl *UserController) Index(c *gin.Context) {
	users, err := ctrl.userSvc.All(c.Request.Context())
	if err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, users, "OK")
}

// Show  GET /api/v1/users/:id
func (ctrl *UserController) Show(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "Invalid ID", nil)
		return
	}
	user, err := ctrl.userSvc.FindByID(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c)
		return
	}
	if user == nil {
		response.NotFound(c, "User")
		return
	}
	response.Success(c, user, "OK")
}

// Store  POST /api/v1/users
func (ctrl *UserController) Store(c *gin.Context) {
	var req requests.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	user, err := ctrl.userSvc.Create(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			response.Conflict(c, "Email is already in use")
			return
		}
		response.ServerError(c)
		return
	}
	response.Created(c, user, "User created")
}

// Update  PUT /api/v1/users/:id
func (ctrl *UserController) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "Invalid ID", nil)
		return
	}
	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	user, err := ctrl.userSvc.Update(c.Request.Context(), id, req.Name, req.Email, req.IsActive)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			response.Conflict(c, "Email is already in use")
			return
		}
		response.ServerError(c)
		return
	}
	if user == nil {
		response.NotFound(c, "User")
		return
	}
	response.Success(c, user, "User updated")
}

// Destroy  DELETE /api/v1/users/:id
func (ctrl *UserController) Destroy(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		response.Error(c, 400, "Invalid ID", nil)
		return
	}
	if err := ctrl.userSvc.Delete(c.Request.Context(), id); err != nil {
		response.ServerError(c)
		return
	}
	response.Success(c, nil, "User deleted")
}

func parseID(c *gin.Context) (uint, error) {
	v, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(v), err
}
