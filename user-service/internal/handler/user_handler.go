package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"user-service/internal/model"
	"user-service/internal/service"
	"user-service/pkg/middleware"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService, middleware middleware.AuthMiddleware, r *gin.Engine) {
	h := &userHandler{userService: userService}
	authRoute := r.Group("/auth")
	authRoute.POST("/register", h.register)
	authRoute.POST("/login", h.login)

	userRoute := r.Group("/users")
	userRoute.GET("/profile", middleware.ValidateAndExtractJwt(), h.getProfile)
	userRoute.PATCH("/profile", middleware.ValidateAndExtractJwt(), h.updateProfile)
	userRoute.PUT("/password", middleware.ValidateAndExtractJwt(), h.updatePassword)
	userRoute.GET("/:id", middleware.ValidateAndExtractJwt(), h.getUserById)
}

func (h *userHandler) getUserById(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	role := userClaims["role"].(string)
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id must be integer"})
		return
	}
	res, err := h.userService.GetUserById(c, i)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	res.Password = ""
	c.JSON(http.StatusOK, res)
}

func (h *userHandler) register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.userService.Register(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res.Password = ""
	c.JSON(http.StatusOK, res)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *userHandler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, token, role, err := h.userService.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId":      userId,
		"accessToken": token,
		"role":        role,
	})
}

func (h *userHandler) getProfile(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	id := int(userClaims["userId"].(float64))
	user, err := h.userService.GetUserById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) updateProfile(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	id := int(userClaims["userId"].(float64))
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = int(id)
	updatedUser, err := h.userService.UpdateUserById(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updatedUser.Password = ""
	c.JSON(http.StatusOK, updatedUser)
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *userHandler) updatePassword(c *gin.Context) {
	claims, _ := c.Get(middleware.JWTClaimsContextKey)
	userClaims := claims.(jwt.MapClaims)
	id := int(userClaims["userId"].(float64))
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.UpdatePasswordById(c, id, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "password changed successfully",
		})
	}
}
