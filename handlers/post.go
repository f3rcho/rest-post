package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/f3rcho/rest-posts/models"
	"github.com/f3rcho/rest-posts/repository"
	"github.com/f3rcho/rest-posts/server"
	"github.com/f3rcho/rest-posts/utils"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

type UpsertPostRequest struct {
	PostContent string `json:"post_content"`
}

type PostResponse struct {
	Id          string `json:"id"`
	PostContent string `json:"post_content"`
}

type MessageResponse struct {
	Message string `json:"message"`
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

		var postMessage = models.WebSocketMessage{
			Type:    "Post_Created",
			Payload: post,
		}
		s.Hub().BroadCast(postMessage, nil)
		json.NewEncoder(w).Encode(PostResponse{
			Id:          post.Id,
			PostContent: post.PostContent,
		})
	}
}

func GetPostById(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		post, err := repository.GetPostById(r.Context(), params["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}
func DeletePostById(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		claims, err := utils.GetClaims(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		err = repository.DeletePostById(r.Context(), params["id"], claims.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(MessageResponse{
			Message: "Resource deleted",
		})
	}
}

func ListPosts(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		pagination := utils.Pagination(pageStr, limitStr)
		posts, err := repository.ListPosts(r.Context(), &pagination)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(posts)
	}
}

func UpdatePost(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		claims, err := utils.GetClaims(s, w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var postRequest = UpsertPostRequest{}
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		post := models.Post{
			PostContent: postRequest.PostContent,
			Id:          params["id"],
		}

		err = repository.UpdatePost(r.Context(), &post, claims.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}
