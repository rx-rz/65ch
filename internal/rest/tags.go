package rest

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"strconv"
)

func (api *API) initializeTagRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/tags", api.createTagHandler)
	api.router.HandlerFunc(http.MethodPatch, "/v1/tags", api.updateTagHandler)
	api.router.HandlerFunc(http.MethodGet, "/v1/tags", api.getTagsHandler)
	api.router.HandlerFunc(http.MethodDelete, "/v1/tags/:id", api.deleteTagHandler)

}

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}

func (api *API) createTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req CreateTagRequest
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
	existingTag, _ := api.models.Tags.GetByName(ctx, req.Name)
	if existingTag != nil {
		api.conflictResponse(w, fmt.Sprintf("Tag with name %s already exists", req.Name))
		return
	}
	tag, err := api.models.Tags.Create(ctx, req.Name)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusCreated, envelope{"tag": tag}, "Tag created successfully")
}

type UpdateTagRequest struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (api *API) updateTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req UpdateTagRequest
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
	updateInfo, err := api.models.Tags.UpdateName(ctx, &data.Tag{
		Name: req.Name,
		ID:   req.ID,
	})
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "Tag updated successfully")
}

func (api *API) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	tags, err := api.models.Tags.GetAll(ctx)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"tags": tags}, "")
}

func (api *API) deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

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

	deleteInfo, err := api.models.Tags.Delete(ctx, id)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": deleteInfo}, "Tag deleted successfully")
}
