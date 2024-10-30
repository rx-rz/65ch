package rest

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"strconv"
)

func (api *API) initializeCategoryRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/categories", api.createCategoryHandler)
	api.router.HandlerFunc(http.MethodGet, "/v1/categories", api.getCategoriesHandler)
	api.router.HandlerFunc(http.MethodPatch, "/v1/categories", api.updateCategoryNameHandler)
	api.router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", api.deleteCategoryHandler)
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

func (api *API) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	existingCategory, _ := api.models.Categories.GetByName(req.Name)
	if existingCategory != nil {
		api.conflictResponse(w, fmt.Sprintf("Category with name %s already exists", req.Name))
		return
	}
	category, err := api.models.Categories.Create(req.Name)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusCreated, envelope{"category": category}, "Category created successfully")

}

func (api *API) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := api.models.Categories.GetAll()
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"categories": categories}, "")
}

type UpdateCategoryRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (api *API) updateCategoryNameHandler(w http.ResponseWriter, r *http.Request) {
	var req UpdateCategoryRequest
	err := api.readJSON(w, r, req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); err != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	updateInfo, err := api.models.Categories.UpdateName(data.Category{
		Name: req.Name,
		ID:   req.ID,
	})
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "Category name updated successfully")
}

type DeleteCategoryRequest struct {
	ID int `json:"id" validate:"required"`
}

func (api *API) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	param, err := api.readParam(r, "id")
	if err != nil {
		api.badRequestResponse(w, err, "ID param not provided")
		return
	}
	id, err := strconv.Atoi(param)
	if err != nil {
		api.badRequestResponse(w, err, "ID param must be an integer")
		return
	}
	deleteInfo, err := api.models.Categories.DeleteByID(id)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": deleteInfo}, "Category deleted successfully")
}
