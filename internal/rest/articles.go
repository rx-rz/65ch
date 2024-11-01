package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
)

func (api *API) initializeArticleRoutes() {

}

type CreateArticleRequest struct {
	AuthorID string   `json:"author_id" validate:"required,string"`
	Title    string   `json:"title" validate:"required,string"`
	Content  string   `json:"content" validate:"required,string"`
	Tags     []string `json:"tags"`
	Category string   `json:"category" validate:"required"`
}

func (api *API) createArticleHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequest
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
	_, err = api.models.Users.FindByID(req.AuthorID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	err = api.models.Articles.Create(&data.Article{
		AuthorID: req.AuthorID,
		Content:  req.Content,
		Tags:     req.Tags,
		Category: req.Category,
	})
}
func (api *API) createArticleDraftHandler() {

}

func (api *API) publishOrArchiveArticleHandler() {

}

func (api *API) updateArticleHandler() {

}

func (api *API) deleteArticleHandler() {

}
