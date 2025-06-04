package router

import (
	"errors"
	"net/http"

	"github.com/khlyazzat/user-crud-k8s-helm/internal/dto"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/user"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/values"

	"github.com/gin-gonic/gin"
)

type userClient struct {
	service user.User
}

func NewUserClient(a user.User) Router {
	return &userClient{
		service: a,
	}
}

func (c *userClient) RegisterRouter(g *gin.RouterGroup) {
	userGroup := g.Group("/user")
	userGroup.POST("/create", c.CreateUser)
	userGroup.GET("/get/:id", c.GetUser)
	userGroup.DELETE("/delete/:id", c.DeleteUser)
	userGroup.PUT("/update/:id", c.UpdateUser)
}

func (c *userClient) RegisterAdminRouter(_ *gin.RouterGroup) {}

func (c *userClient) CreateUser(ctx *gin.Context) {
	var body dto.AddUserRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}
	res, err := c.service.AddUser(ctx, &body)
	if errors.Is(err, values.ErrEmailExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "email already exists"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": res.ID})
}

func (c *userClient) GetUser(ctx *gin.Context) {
	var path dto.UserPathRequest
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}
	res, err := c.service.GetUser(ctx, &dto.GetUserRequest{ID: path.ID})
	if errors.Is(err, values.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *userClient) DeleteUser(ctx *gin.Context) {
	var path dto.UserPathRequest
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}
	err := c.service.DeleteUser(ctx, &dto.DeleteUserRequest{ID: path.ID})
	if errors.Is(err, values.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *userClient) UpdateUser(ctx *gin.Context) {
	var path dto.UserPathRequest
	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}
	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	res, err := c.service.UpdateUser(ctx, path.ID, &req)
	if errors.Is(err, values.ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
