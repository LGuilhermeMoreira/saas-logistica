package auth

import (
	"gateway/internal/domain/contract"
	"net/http"
	authv1 "proto/gen/auth/v1"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService contract.AuthService
}

func NewAuthController(authService contract.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.Login(input.Email, input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, LoginOutput{Token: token})
}

func (c *AuthController) CreateUser(ctx *gin.Context) {
	var input CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authField := getAuthField(ctx)

	resp, err := c.authService.CreateUser(input.Name, input.Email, input.Password, input.RoleID, authField)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, CreateUserOutput{
		ID:    resp.Id,
		Name:  resp.Name,
		Email: resp.Email,
	})
}

func (c *AuthController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	authField := getAuthField(ctx)

	resp, err := c.authService.DeleteUser(id, authField)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, DeleteUserOutput{Msg: resp.Msg})
}

func (c *AuthController) CreateRole(ctx *gin.Context) {
	var input CreateRoleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authField := getAuthField(ctx)
	var permissions []*authv1.Permission
	for _, p := range input.Permissions {
		permissions = append(permissions, &authv1.Permission{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	resp, err := c.authService.CreateRole(input.Name, input.Description, permissions, authField)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var createdPermissions []PermissionInput
	for _, p := range resp.Permissions {
		createdPermissions = append(createdPermissions, PermissionInput{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	ctx.JSON(http.StatusCreated, CreateRoleOutput{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Permissions: createdPermissions,
	})
}

func (c *AuthController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")
	authField := getAuthField(ctx)

	resp, err := c.authService.DeleteRole(id, authField)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, DeleteRoleOutput{Msg: resp.Msg})
}

func getAuthField(c *gin.Context) *authv1.AuthField {
	return &authv1.AuthField{
		Token:  extractToken(c.GetHeader("Authorization")),
		Path:   c.Request.URL.Path,
		Action: c.Request.Method,
	}
}

func extractToken(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}
