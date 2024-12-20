package routes

import (
	"itam_auth/internal/database"
	"itam_auth/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(storage *database.Storage) *gin.Engine {

	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	router.Use(cors.New(config))

	// @Summary Пинг-сервис
	// @Description Проверяет доступность сервера
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "pong"}
	// @Router /api/ping [get]
	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//* USER

	// @Summary Логин пользователя
	// @Description Авторизация пользователя с использованием логина и пароля
	// @Tags User
	// @Accept json
	// @Produce json
	// @Param login body models.LoginRequest true "Login credentials"
	// @Success 200 {object} map[string]string{"token": "your-jwt-token"}
	// @Failure 400 {object} map[string]string{"error": "Invalid request"}
	// @Failure 401 {object} map[string]string{"error": "Unauthorized"}
	// @Failure 500 {object} map[string]string{"error": "Internal server error"}
	// @Router /api/login [post]
	router.POST("/api/login", func(ctx *gin.Context) {
		handlers.Login(ctx, storage)
	})

	// @Summary Регистрация нового пользователя
	// @Description Регистрация нового пользователя в системе
	// @Tags User
	// @Accept json
	// @Produce json
	// @Param register body models.RegisterRequest true "User registration details"
	// @Success 201 {object} map[string]string{"message": "User registered successfully"}
	// @Failure 400 {object} map[string]string{"error": "Invalid request"}
	// @Failure 500 {object} map[string]string{"error": "Internal server error"}
	// @Router /api/register [post]
	router.POST("/api/register", func(ctx *gin.Context) {
		handlers.Register(ctx, storage)
	})

	// @Summary Получить информацию о пользователе
	// @Description Возвращает данные текущего пользователя
	// @Tags User
	// @Produce json
	// @Param user_id path string true "User ID"
	// @Success 200 {object} models.User
	// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
	// @Failure 404 {object} gin.H{"error": "User not found"}
	// @Router /api/get_user/{user_id} [get]
	router.GET("/api/get_user/:user_id", func(ctx *gin.Context) {
		handlers.GetUser(ctx, storage)
	})

	// @Summary Обновить информацию пользователя
	// @Description Обновляет профиль пользователя
	// @Tags User
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "update_user_info"}
	// @Router /api/update_user_info [patch]
	router.PATCH("/api/update_user_info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "update_user_info",
		})
	})

	// @Summary Получить роли пользователя
	// @Description Возвращает список ролей текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user_roles"}
	// @Router /api/get_user_roles [get]
	router.GET("/api/get_user_roles", func(ctx *gin.Context) {
		handlers.GetUserRoles(ctx, storage)
	})

	// @Summary Получить свойства пользователя
	// @Description Возвращает список свойств текущего пользователя
	// @Tags User
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_user_properties"}
	// @Router /api/get_user_properties [get]
	router.GET("/api/get_user_properties", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get_user_properties",
		})
	})

	//* Requests

	// @Summary Создать запрос пользователя
	// @Description Создает новый запрос от имени пользователя
	// @Tags Requests
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "create_user_request"}
	// @Router /api/create_user_request [post]
	router.POST("/api/create_user_request", func(ctx *gin.Context) {
		handlers.CreateUserRequest(ctx, storage)
	})

	// @Summary Получить запрос пользователя
	// @Description Возвращает данные о конкретном запросе пользователя
	// @Tags Requests
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_request"}
	// @Router /api/get_request [get]
	router.GET("/api/get_request", func(ctx *gin.Context) {
		handlers.GetRequest(ctx, storage)
	})

	// @Summary Получить все запросы пользователя
	// @Description Возвращает список всех запросов текущего пользователя
	// @Tags Requests
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "get_all_requests"}
	// @Router /api/get_all_requests [get]
	router.GET("/api/get_all_requests", func(ctx *gin.Context) {
		handlers.GetAllRequests(ctx, storage)
	})

	// @Summary Обновить статус запроса
	// @Description Обновляет статус указанного запроса
	// @Tags Requests
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string{"message": "update_request_status"}
	// @Router /api/update_request_status [patch]
	router.PATCH("/api/update_request_status", func(ctx *gin.Context) {
		handlers.UpdateRequestStatus(ctx, storage)
	})

	//* Achievements

	// @Summary Получить достижения пользователя
	// @Description Возвращает список достижений текущего пользователя
	// @Tags Achievements
	// @Produce json
	// @Param user_id query string true "User ID"
	// @Success 200 {array} models.Achievement
	// @Failure 400 {object} map[string]string{"error": "Invalid user ID"}
	// @Failure 500 {object} map[string]string{"error": "Error while fetching achievements"}
	// @Router /api/get_user_achievements [get]
	router.GET("/api/get_user_achievements", func(ctx *gin.Context) {
		handlers.GetAchievementsByUserID(ctx, storage)
	})

	// @Summary Создать достижение
	// @Description Создает новое достижение
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Param achievement body models.Achievement true "Achievement data"
	// @Success 201 {object} map[string]interface{"message": "Achievement created successfully", "id": string}
	// @Failure 400 {object} map[string]string{"error": "Invalid title or points"}
	// @Failure 500 {object} map[string]string{"error": "Failed to save achievement"}
	// @Router /api/create_achievement [post]
	router.POST("/api/create_achievement", func(ctx *gin.Context) {
		handlers.CreateAchievement(ctx, storage)
	})

	// @Summary Обновить достижение
	// @Description Обновляет существующее достижение
	// @Tags Achievements
	// @Accept json
	// @Produce json
	// @Param achievement body models.Achievement true "Achievement data"
	// @Success 200 {object} map[string]string{"message": "Achievement updated successfully"}
	// @Failure 400 {object} map[string]string{"error": "Invalid request"}
	// @Failure 500 {object} map[string]string{"error": "Failed to update achievement"}
	// @Router /api/update_achievement [patch]
	router.PATCH("/api/update_achievement", func(ctx *gin.Context) {
		handlers.UpdateAchievement(ctx, storage)
	})

	// @Summary Получить достижение
	// @Description Возвращает информацию о конкретном достижении
	// @Tags Achievements
	// @Produce json
	// @Param achievement_id query string true "Achievement ID"
	// @Success 200 {object} models.Achievement
	// @Failure 400 {object} map[string]string{"error": "Invalid achievement ID"}
	// @Failure 404 {object} map[string]string{"error": "Achievement not found"}
	// @Failure 500 {object} map[string]string{"error": "Error while fetching achievement"}
	// @Router /api/get_achievement [get]
	router.GET("/api/get_achievement", func(ctx *gin.Context) {
		handlers.GetAchievementByID(ctx, storage)
	})

	// @Summary Получить все достижения
	// @Description Возвращает список всех достижений
	// @Tags Achievements
	// @Produce json
	// @Success 200 {array} models.Achievement
	// @Router /api/get_all_achievements [get]
	router.GET("/api/get_all_achievements", func(ctx *gin.Context) {
		handlers.GetAllAchievements(ctx, storage)
	})

	//* Notifications (admin)

	// CreateNotification godoc
	// @Summary Create a new notification
	// @Description Create a new notification
	// @Tags Notifications
	// @Accept json
	// @Produce json
	// @Param notification body models.Notification true "Notification"
	// @Success 201 {object} gin.H{"message": "Notification created successfully"}
	// @Failure 400 {object} gin.H{"error": "Invalid request"}
	// @Failure 500 {object} gin.H{"error": "Error while saving notification"}
	// @Router /api/create_notification [post]
	router.POST("/api/create_notification", func(ctx *gin.Context) {
		handlers.CreateNotification(ctx, storage)
	})

	// UpdateNotification godoc
	// @Summary Update an existing notification
	// @Description Update an existing notification
	// @Tags Notifications
	// @Accept json
	// @Produce json
	// @Param notification body models.Notification true "Notification"
	// @Success 201 {object} gin.H{"message": "Notification updated successfully"}
	// @Failure 400 {object} gin.H{"error": "Invalid request"}
	// @Failure 500 {object} gin.H{"error": "Error while updating notification"}
	// @Router /api/update_notification [patch]
	router.PATCH("/api/update_notification", func(ctx *gin.Context) {
		handlers.UpdateNotification(ctx, storage)
	})

	// GetAllNotifications godoc
	// @Summary Get all notifications for a user
	// @Description Get all notifications for a user
	// @Tags Notifications
	// @Produce json
	// @Param user_id path string false "User ID"
	// @Success 201 {object} models.Notification
	// @Failure 400 {object} gin.H{"error": "Invalid user ID"}
	// @Failure 500 {object} gin.H{"error": "Error while fetching notifications"}
	// @Router /api/get_all_notifications [get]
	router.GET("/api/get_all_notifications", func(ctx *gin.Context) {
		handlers.GetAllNotifications(ctx, storage)
	})

	// GetNotification godoc
	// @Summary Get a notification by its ID
	// @Description Get a notification by its ID
	// @Tags Notifications
	// @Produce json
	// @Param notification_id path string true "Notification ID"
	// @Success 201 {object} models.Notification
	// @Failure 400 {object} gin.H{"error": "Invalid notification ID"}
	// @Failure 500 {object} gin.H{"error": "Error while fetching notification"}
	// @Router /api/get_notification/{notification_id} [get]
	router.GET("/api/get_notification/{notification_id}", func(ctx *gin.Context) {
		handlers.GetNotification(ctx, storage)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
