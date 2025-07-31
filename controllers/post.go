package controllers

import (
	"gin-rest-api/models"
	"gin-rest-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type PostAPI struct {
	db *gorm.DB
}

func NewPostAPI(db *gorm.DB) *PostAPI {
	return &PostAPI{db}
}

func (a *PostAPI) Create(c *gin.Context) {
	var req models.PostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	if !utils.IsExist(a.db, "categories", "id", req.CategoryId) {
		utils.StatusUnprocessable(c, map[string]any{
			"CategoryId": "the category does not exist!",
		})
		return
	}

	post := models.Post{
		Title:      req.Title,
		Body:       req.Body,
		CategoryID: req.CategoryId,
		UserID:     utils.GetUserID(c),
	}

	result := a.db.Create(&post)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, post)
}

func (a *PostAPI) Gets(c *gin.Context) {
	var posts []models.Post

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	preloadFunc := func(query *gorm.DB) *gorm.DB {
		return query.Preload("Category", func(db *gorm.DB) *gorm.DB {
			return a.db.Select("id, name, slug")
		}).Preload("User", func(db *gorm.DB) *gorm.DB {
			return a.db.Select("id, name")
		})
	}

	result, err := utils.Paginate(a.db, page, perPage, preloadFunc, &posts)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *PostAPI) Get(c *gin.Context) {
	var post models.Post

	result := a.db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return a.db.Select("id, name, slug")
	}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return a.db.Select("id, name")
	}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return a.db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return a.db.Select("id, name")
		}).Select("id, post_id, user_id, body, created_at")
	}).First(&post, c.Param("id"))

	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, post)
}

func (a *PostAPI) Update(c *gin.Context) {
	var req models.PostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	var post models.Post

	result := a.db.First(&post, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	updatePost := models.Post{
		Title:      req.Title,
		Body:       req.Body,
		CategoryID: req.CategoryId,
		UserID:     utils.GetUserID(c),
	}

	result = a.db.Model(&post).Updates(&updatePost)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, updatePost)
}

func (a *PostAPI) Delete(c *gin.Context) {
	var post models.Post

	result := a.db.First(&post, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Delete(&post)

	utils.StatusOK(c, nil, "the post has been deleted successfully")
}

func (a *PostAPI) Trashed(c *gin.Context) {
	var posts []models.Post

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := utils.Paginate(a.db.Unscoped().Where("deleted_at IS NOT NULL"), page, perPage, nil, &posts)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *PostAPI) EmptyTrash(c *gin.Context) {
	var post models.Post

	if err := a.db.Unscoped().First(&post, c.Param("id")).Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Unscoped().Delete(&post)

	utils.StatusOK(c, nil, "the post has been deleted permanently")
}
