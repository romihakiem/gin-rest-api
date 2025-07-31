package controllers

import (
	"gin-rest-api/models"
	"gin-rest-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type CategoryAPI struct {
	db *gorm.DB
}

func NewCategoryAPI(db *gorm.DB) *CategoryAPI {
	return &CategoryAPI{db}
}

func (a *CategoryAPI) Create(c *gin.Context) {
	var req models.CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	if utils.IsExist(a.db, "categories", "name", req.Name) ||
		utils.IsExist(a.db, "categories", "slug", slug.Make(req.Name)) {
		utils.StatusUnprocessable(c, map[string]any{
			"Name": "the name is already exist!",
		})
		return
	}

	category := models.Category{
		Name: req.Name,
	}

	result := a.db.Create(&category)
	if result.Error != nil {
		utils.StatusServerError(c, "cannot create category")
		return
	}

	utils.StatusOK(c, category)
}

func (a *CategoryAPI) Gets(c *gin.Context) {
	var categories []models.Category

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := utils.Paginate(a.db, page, perPage, nil, &categories)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *CategoryAPI) Get(c *gin.Context) {
	var category models.Category

	result := a.db.First(&category, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, category)
}

func (a *CategoryAPI) Update(c *gin.Context) {
	var req models.CategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	var category models.Category

	result := a.db.First(&category, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	if (category.Name != req.Name &&
		utils.IsExist(a.db, "categories", "name", req.Name)) ||
		(category.Name != req.Name &&
			utils.IsExist(a.db, "categories", "slug", slug.Make(req.Name))) {
		utils.StatusUnprocessable(c, map[string]any{
			"Name": "the name is already exist!",
		})
		return
	}

	updateCategory := models.Category{
		Name: req.Name,
		Slug: slug.Make(req.Name),
	}

	result = a.db.Model(&category).Updates(updateCategory)
	if err := result.Error; err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, updateCategory)
}

func (a *CategoryAPI) Delete(c *gin.Context) {
	var category models.Category

	result := a.db.First(&category, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Delete(&category)

	utils.StatusOK(c, nil, "the category has been deleted successfully")
}

func (a *CategoryAPI) Trashed(c *gin.Context) {
	var categories []models.Category

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	perPageStr := c.DefaultQuery("perPage", "10")
	perPage, _ := strconv.Atoi(perPageStr)

	result, err := utils.Paginate(a.db.Unscoped().Where("deleted_at IS NOT NULL"), page, perPage, nil, &categories)
	if err != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, result)
}

func (a *CategoryAPI) EmptyTrash(c *gin.Context) {
	result := a.db.Unscoped().Delete(&models.Category{}, c.Param("id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, nil, "the category has been deleted permanently")
}
