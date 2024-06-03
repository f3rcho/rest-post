package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/f3rcho/rest-posts/models"
	"github.com/f3rcho/rest-posts/repository"
	"github.com/f3rcho/rest-posts/server"
	"github.com/f3rcho/rest-posts/utils"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

func InserPost(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.GetClaims(s, w, r)

		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		var postRequest = UpsertPostRequest{}
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		post := models.Post{
			Id:          id.String(),
			PostContent: postRequest.PostContent,
			UserId:      claims.UserID,
		}
		err = repository.InsertPost(r.Context(), &post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(PostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})
	}
}
