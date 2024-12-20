package handlers

import (
	"context"
	"itam_auth/internal/database"
	"itam_auth/internal/models"
	"itam_auth/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TelegramAuth struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context, storage *database.Storage) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	user, err := auth.RegisterUser(ctx, storage, req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user", "err": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user": user})
}

func Login(c *gin.Context, storage *database.Storage) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tokenString, err := auth.AuthenticateUser(ctx, storage, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// TELEGRAM WEB APP
// func Login(c *gin.Context, storage *database.Storage) {
// 	var auth TelegramAuth
// 	if err := c.ShouldBindJSON(&auth); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	data := url.Values{
// 		"id":         {fmt.Sprintf("%d", auth.ID)},
// 		"first_name": {auth.FirstName},
// 		"last_name":  {auth.LastName},
// 		"username":   {auth.Username},
// 		"photo_url":  {auth.PhotoURL},
// 		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
// 	}

// 	if !utils.ValidateTelegramAuth(data, "") {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
// 		return
// 	}

// 	// Генерация UUID
// 	userID := uuid.New()

// 	ctx := context.Background()
// 	user, err := storage.GetUserByID(ctx, userID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	tokenString, err := utils.GenerateJWT(user.Email)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }

// TELEGRAM WEB APP
// func Register(c *gin.Context, storage *database.Storage) {
// 	var auth TelegramAuth

// 	if err := c.ShouldBindJSON(&auth); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	data := url.Values{
// 		"id":         {fmt.Sprintf("%d", auth.ID)},
// 		"first_name": {auth.FirstName},
// 		"last_name":  {auth.LastName},
// 		"username":   {auth.Username},
// 		"photo_url":  {auth.PhotoURL},
// 		"auth_date":  {fmt.Sprintf("%d", auth.AuthDate)},
// 	}

// 	if !utils.ValidateTelegramAuth(data, "your_bot_token") {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Telegram authorization"})
// 		return
// 	}

// 	userID := uuid.New()

// 	user := models.User{
// 		ID:           userID,
// 		Name:         fmt.Sprintf("%s %s", auth.FirstName, auth.LastName),
// 		Email:        fmt.Sprintf("%d@telegram.com", auth.ID),
// 		Telegram:     auth.Username,
// 		PasswordHash: "",
// 		PhotoURL:     auth.PhotoURL,
// 		CreatedAt:    time.Now(),
// 		UpdatedAt:    time.Now(),
// 	}

// 	ctx := context.Background()

// 	_, err := storage.SaveUser(ctx, user)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while saving user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
// }

func GetUser(c *gin.Context, storage *database.Storage) {
	userID := c.Query("user_id")
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	user, err := storage.GetUserByID(ctx, uuidUserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUserRoles(c *gin.Context, storage *database.Storage) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	uuidUserID, errUUID := uuid.Parse(userID)
	if errUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()
	roles, err := storage.GetUserRoles(ctx, uuidUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching user roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func GetUserPermissions(c *gin.Context, storage *database.Storage) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	uuidUserID, errUUID := uuid.Parse(userID)
	if errUUID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := context.Background()

	userRoles, err := storage.GetUserRoles(ctx, uuidUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching user roles"})
		return
	}

	var permissions []models.RolePermission
	for _, role := range userRoles {
		rolePermissions, err := storage.GetRolePermissions(ctx, role.RoleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching role permissions"})
			return
		}
		permissions = append(permissions, rolePermissions...)
	}

	c.JSON(http.StatusOK, permissions)
}
