package controllers

import (
	"gin-rest-api/config"
	"gin-rest-api/models"
	"gin-rest-api/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserAPI struct {
	cfg *config.Config
	db  *gorm.DB
	rds *redis.Client
}

func NewUserAPI(cfg *config.Config, db *gorm.DB, rds *redis.Client) *UserAPI {
	return &UserAPI{cfg, db, rds}
}

func (a *UserAPI) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	if utils.IsExist(a.db, "users", "email", req.Email) {
		utils.StatusUnprocessable(c, map[string]any{
			"Email": "the email is already exist!",
		})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		utils.StatusServerError(c, "failed to hash password")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashPassword),
	}

	result := a.db.Create(&user)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, user)
}

func (a *UserAPI) Login(c *gin.Context) {
	var req models.LoginRequest

	if c.ShouldBindJSON(&req) != nil {
		utils.StatusBadRequest(c, "failed to read body")
		return
	}

	var user models.User

	a.db.First(&user, "email = ?", req.Email)

	if user.ID == 0 {
		utils.StatusBadRequest(c, "invalid email or password")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		utils.StatusBadRequest(c, "invalid email or password")
		return
	}

	token, err := utils.CreateToken(a.cfg, user.ID)
	if err != nil {
		utils.StatusBadRequest(c, "failed to create token")
		return
	}

	a.rds.Set(c.Request.Context(), strconv.Itoa(int(user.ID)), token.AccessToken, time.Duration(a.cfg.JWTAccessExpiry)*time.Second)
	a.rds.Set(c.Request.Context(), "refresh_"+strconv.Itoa(int(user.ID)), token.RefreshToken, time.Duration(a.cfg.JWTRefreshExpiry)*time.Second)

	utils.StatusOK(c, token)
}

func (a *UserAPI) Logout(c *gin.Context) {
	userId := utils.GetUserID(c)

	a.rds.Del(c.Request.Context(), strconv.Itoa(int(userId)))
	a.rds.Del(c.Request.Context(), "refresh_"+strconv.Itoa(int(userId)))

	utils.StatusOK(c, nil, "logout successful")
}

func (a *UserAPI) Me(c *gin.Context) {
	var user models.User

	result := a.db.First(&user, utils.GetUserID(c))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, user)
}

func (a *UserAPI) Gets(c *gin.Context) {
	var users []models.User

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := utils.Paginate(a.db, page, perPage, nil, &users)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *UserAPI) Get(c *gin.Context) {
	var user models.User

	result := a.db.First(&user, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, user)
}

func (a *UserAPI) Update(c *gin.Context) {
	var req models.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	var user models.User

	result := a.db.First(&user, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	if user.Email != req.Email && utils.IsExist(a.db, "users", "email", req.Email) {
		utils.StatusUnprocessable(c, map[string]any{
			"Email": "the email is already exist!",
		})
		return
	}

	updateUser := models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	result = a.db.Model(&user).Updates(&updateUser)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, user)
}

func (a *UserAPI) Delete(c *gin.Context) {
	var user models.User

	result := a.db.First(&user, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Delete(&user)

	utils.StatusOK(c, nil, "the user has been deleted successfully")
}

func (a *UserAPI) Trashed(c *gin.Context) {
	var users []models.User

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := utils.Paginate(a.db.Unscoped().Where("deleted_at IS NOT NULL"), page, perPage, nil, &users)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *UserAPI) EmptyTrash(c *gin.Context) {
	var user models.User

	if err := a.db.Unscoped().First(&user, c.Param("id")).Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Unscoped().Delete(&user)

	utils.StatusOK(c, nil, "the user has been deleted permanently")
}
