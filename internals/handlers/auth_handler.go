package handlers

import (
	"io"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/tacheraSasi/go-api-starter/internals/config"
	"github.com/tacheraSasi/go-api-starter/internals/dtos"
	"github.com/tacheraSasi/go-api-starter/internals/models"
	"github.com/tacheraSasi/go-api-starter/internals/services"
	"github.com/tacheraSasi/go-api-starter/pkg/jwt"
	"github.com/tacheraSasi/go-api-starter/pkg/styles"
	"github.com/tacheraSasi/go-api-starter/ui/pages"
)

type AuthHandler struct {
	service services.AuthService
	cfg     *config.Config
}

func NewAuthHandler(service services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: service,
		cfg:     cfg,
	}
}

// ValidateRequest validates the request body against the provided struct
func (h *AuthHandler) ValidateRequest(c *gin.Context, obj any) {
	dtos.Validate(c, obj)
}

func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *AuthHandler) RegisterWebRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/register", h.RegisterPage)
		authGroup.GET("/login", h.LoginPage)
		authGroup.GET("/forgot-password", h.ForgotPasswordPage)
		authGroup.GET("/reset-password", h.ResetPasswordPage)
		authGroup.POST("/logout", h.Logout)
	}
}

func (h *AuthHandler) RegisterPage(c *gin.Context) {
	templ.Handler(pages.Register(pages.RegisterProps{AppName: "Go API Starter"})).ServeHTTP(c.Writer, c.Request)
}
func (h *AuthHandler) LoginPage(c *gin.Context) {
	templ.Handler(pages.Login(pages.LoginProps{AppName: "Go API Starter"})).ServeHTTP(c.Writer, c.Request)
}
func (h *AuthHandler) ForgotPasswordPage(c *gin.Context) {
	templ.Handler(pages.ForgotPassword(pages.ForgotPasswordProps{AppName: "Go API Starter"})).ServeHTTP(c.Writer, c.Request)
}
func (h *AuthHandler) ResetPasswordPage(c *gin.Context) {
	templ.Handler(pages.ResetPassword(pages.ResetPasswordProps{AppName: "Go API Starter", Token: c.Query("token")})).ServeHTTP(c.Writer, c.Request)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var reqDto dtos.ForgotPasswordRequest
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted() {
		return
	}

	token, err := h.service.RequestPasswordReset(reqDto.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create password reset request"})
		return
	}

	response := gin.H{
		"message": "If your email exists, a reset link has been generated",
	}

	if token != "" {
		response["reset_token"] = token
		response["reset_url"] = "/auth/reset-password?token=" + token
	}

	c.JSON(200, response)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var reqDto dtos.ResetPasswordRequest
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted() {
		return
	}

	if reqDto.Password != reqDto.PasswordConfirmation {
		c.JSON(400, gin.H{"error": "Password confirmation does not match"})
		return
	}

	if err := h.service.ResetPassword(reqDto.Token, reqDto.Password); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Password reset successful"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var reqDto dtos.RegisterRequest
	var requestBody = c.Request.Body
	defer requestBody.Close()
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted() {
		return
	}
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println(styles.Request.Render(string(bodyBytes)))
	err = h.service.Register(&models.User{
		Email:    reqDto.Email,
		Password: reqDto.Password,
		Name:     reqDto.Name,
	})
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := h.service.GetUserByEmail(reqDto.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve user after registration",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Registration successful",
		"user":    user,
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var reqDto dtos.LoginRequest
	var requestBody = c.Request.Body
	defer requestBody.Close()
	h.ValidateRequest(c, &reqDto)
	if c.IsAborted() {
		return
	}
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return
	}
	log.Println(styles.Request.Render(string(bodyBytes)))

	user, err := h.service.Login(reqDto.Email, reqDto.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	token, err := jwt.GenerateToken(user, []byte(h.cfg.JWTSecret), h.cfg.JWTExpiresIn)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{"error": "Authorization header is missing"})
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]

	claims, err := jwt.ValidateToken(tokenString, []byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	expiresAt := claims.ExpiresAt.Time

	if err := h.service.Logout(tokenString, expiresAt); err != nil {
		c.JSON(500, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(200, gin.H{"message": "Logout successful"})
}
