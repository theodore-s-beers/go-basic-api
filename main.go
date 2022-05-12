package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/theodore-s-beers/go-basic-api/internal/database"
)

func main() {
	c := database.NewClient("db.json")

	err := c.EnsureDB()

	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		dbClient: c,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/users", apiCfg.endpointUsersHandler)
	mux.HandleFunc("/users/", apiCfg.endpointUsersHandler)

	mux.HandleFunc("/posts", apiCfg.endpointPostsHandler)
	mux.HandleFunc("/posts/", apiCfg.endpointPostsHandler)

	const addr = "localhost:8080"

	srv := http.Server{
		Handler:      mux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	srv.ListenAndServe()
}

// Types

type apiConfig struct {
	dbClient database.Client
}

type errorBody struct {
	Error string `json:"error"`
}

// Methods

func (apiCfg apiConfig) endpointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerRetrievePosts(w, r)
	case http.MethodPost:
		apiCfg.handlerCreatePost(w, r)
	case http.MethodPut:
		// None so far, I think
		respondWithError(w, http.StatusBadRequest, errors.New("method not supported"))
	case http.MethodDelete:
		apiCfg.handlerDeletePost(w, r)
	default:
		respondWithError(w, http.StatusBadRequest, errors.New("method not supported"))
	}
}

func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerGetUser(w, r)
	case http.MethodPost:
		apiCfg.handlerCreateUser(w, r)
	case http.MethodPut:
		apiCfg.handlerUpdateUser(w, r)
	case http.MethodDelete:
		apiCfg.handlerDeleteUser(w, r)
	default:
		respondWithError(w, http.StatusBadRequest, errors.New("method not supported"))
	}
}

func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserEmail string `json:"userEmail"`
		Text      string `json:"text"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decodingErr := decoder.Decode(&params)

	if decodingErr != nil {
		respondWithError(w, http.StatusBadRequest, decodingErr)
		return
	}

	_, creationErr := apiCfg.dbClient.CreatePost(params.UserEmail, params.Text)

	if creationErr != nil {
		respondWithError(w, http.StatusBadRequest, creationErr)
		return
	}

	respondWithJSON(w, http.StatusCreated, params)
}

func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decodingErr := decoder.Decode(&params)

	if decodingErr != nil {
		respondWithError(w, http.StatusBadRequest, decodingErr)
		return
	}

	_, creationErr := apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)

	if creationErr != nil {
		respondWithError(w, http.StatusBadRequest, creationErr)
		return
	}

	respondWithJSON(w, http.StatusCreated, params)
}

func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimPrefix(r.URL.Path, "/posts/")

	err := apiCfg.dbClient.DeletePost(uuid)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")

	err := apiCfg.dbClient.DeleteUser(email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")

	user, err := apiCfg.dbClient.GetUser(email)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (apiCfg apiConfig) handlerRetrievePosts(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimPrefix(r.URL.Path, "/posts/")

	posts, err := apiCfg.dbClient.GetPosts(user)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

func (apiCfg apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decodingErr := decoder.Decode(&params)

	if decodingErr != nil {
		respondWithError(w, http.StatusBadRequest, decodingErr)
		return
	}

	email := strings.TrimPrefix(r.URL.Path, "/users/")

	user, updateErr := apiCfg.dbClient.UpdateUser(email, params.Password, params.Name, params.Age)

	if updateErr != nil {
		respondWithError(w, http.StatusBadRequest, updateErr)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// Functions

func respondWithError(w http.ResponseWriter, code int, err error) {
	eb := errorBody{
		Error: err.Error(),
	}

	respondWithJSON(w, code, eb)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func userIsEligible(email, password string, age int) error {
	if email == "" {
		return errors.New("email can't be empty")
	}

	if password == "" {
		return errors.New("password can't be empty")
	}

	if age < 18 {
		return errors.New("age must be at least 18")
	}

	return nil
}
