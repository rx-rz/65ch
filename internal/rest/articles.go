package rest

import "net/http"

func (api *API) initializeArticleRoutes() {

}

type CreateArticleRequest struct {
	AuthorID string   `json:"author_id" validate:"required,string"`
	Title    string   `json:"title" validate:"required,string"`
	Content  string   `json:"content" validate:"required,string"`
	Tags     []string `json:"tags"`
	Category string   `json:"category"`
}

func (api *API) createArticleHandler(w http.ResponseWriter, r *http.Request) {

}
func (api *API) createArticleDraftHandler() {

}

func (api *API) publishOrArchiveArticleHandler() {

}

func (api *API) updateArticleHandler() {

}

func (api *API) deleteArticleHandler() {

}
