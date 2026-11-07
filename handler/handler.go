package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"week_06_web/app_err"
	"week_06_web/auth"
	"week_06_web/middleware"
	"week_06_web/model"
	"week_06_web/response"
	"week_06_web/service"
)

type UserHandler struct {
	service *service.UserService
	tokens  *auth.TokenManager
}

func NewUserHandler(s *service.UserService, tokens *auth.TokenManager) *UserHandler {
	return &UserHandler{service: s, tokens: tokens}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, app_err.BadRequest("BAD_REQUEST", "invalid request body", err))
		return
	}

	user, err := h.service.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, app_err.BadRequest("BAD_REQUEST", "invalid request body", err))
		return
	}

	user, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	token, expiresAt, err := h.tokens.Generate(user.Id)
	if err != nil {
		response.Error(c, app_err.Internal(err))
		return
	}

	auth.SetCookie(c, token, expiresAt)

	response.OK(c, model.LoginResponse{User: user, Token: token})
}

func (h *UserHandler) Logout(c *gin.Context) {
	auth.ClearCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) Me(c *gin.Context) {
	userId := c.GetString(middleware.UserIDKey)
	if userId == "" {
		response.Error(c, app_err.Unauthorized("UNAUTHORIZED", "unauthorized", nil))
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userId := c.GetString(middleware.UserIDKey)
	if userId == "" {
		response.Error(c, app_err.Unauthorized("UNAUTHORIZED", "unauthorized", nil))
		return
	}

	if c.Param("id") != userId {
		response.Error(c, app_err.Forbidden("FORBIDDEN", "cannot change another user's password", nil))
		return
	}

	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, app_err.BadRequest("BAD_REQUEST", "invalid request body", err))
		return
	}

	if err := h.service.ChangePassword(c.Request.Context(), userId, req.OldPassword, req.NewPassword); err != nil {
		response.Error(c, err)
		return
	}

	auth.ClearCookie(c)
	c.Status(http.StatusNoContent)
}
