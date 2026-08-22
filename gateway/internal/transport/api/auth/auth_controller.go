package auth

import (
	"gateway/internal/domain/contract"
	"net/http"
	authv1 "proto/gen/auth/v1"

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

	token, err := c.authService.Login(ctx.Request.Context(), input.Email, input.Password)
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

	resp, err := c.authService.CreateUser(ctx.Request.Context(), input.Name, input.Email, input.Password, input.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, UserOutput{
		ID:    resp.Id,
		Name:  resp.Name,
		Email: resp.Email,
	})
}

func (c *AuthController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := c.authService.DeleteUser(ctx.Request.Context(), id)
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

	var permissions []*authv1.Permission
	for _, p := range input.Permissions {
		permissions = append(permissions, &authv1.Permission{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	resp, err := c.authService.CreateRole(ctx.Request.Context(), input.Name, input.Description, permissions)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var createdPermissions []Permission
	for _, p := range resp.Permissions {
		createdPermissions = append(createdPermissions, Permission{
			Action: p.Action,
			Path:   p.Path,
		})
	}

	ctx.JSON(http.StatusCreated, RoleOutput{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Permissions: createdPermissions,
	})
}

func (c *AuthController) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := c.authService.DeleteRole(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, DeleteRoleOutput{Msg: resp.Msg})
}

func (c *AuthController) FindUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := c.authService.FindUserByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, UserOutput{
		ID:     resp.Id,
		Email:  resp.Email,
		Name:   resp.Name,
		RoleID: resp.RoleId,
	})

}

func (c *AuthController) FindRoleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, err := c.authService.FindRoleByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	perm := make([]Permission, len(resp.Permissions))

	for i, v := range resp.Permissions {
		perm[i] = Permission{
			Action: v.Action,
			Path:   v.Path,
		}
	}

	ctx.JSON(http.StatusOK, RoleOutput{
		ID:          resp.Id,
		Name:        resp.Name,
		Description: resp.Description,
		Permissions: perm,
	})

}

func (c *AuthController) FindAllRoles(ctx *gin.Context) {
	resp, err := c.authService.FindAllRoles(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	roles := make([]RoleOutput, len(resp.Roles))

	for i, v := range resp.Roles {
		roles[i] = RoleOutput{
			ID:          v.Id,
			Description: v.Description,
			Name:        v.Name,
		}
	}

	ctx.JSON(http.StatusOK, FindAllRolesOutput{
		Roles: roles,
	})
}
