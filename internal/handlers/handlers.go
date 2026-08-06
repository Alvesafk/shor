package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Alvesafk/shor/internal/app"
	"github.com/Alvesafk/shor/internal/models"
)

type Connection struct {
	db  *app.FireDB
	ctx context.Context
}

func NewConnection(db *app.FireDB, ctx context.Context) *Connection {
	return &Connection{db, ctx}
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Content any    `json:"content"`
}

func (r Response) WriteJSON(w http.ResponseWriter, h int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(h)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(r)
}

func (c Connection) GetHelloWorld(w http.ResponseWriter, r *http.Request) {
	Response{
		Message: "This is a Hello, World!",
		Status:  "ok",
		Content: "Hello, World!",
	}.WriteJSON(w, http.StatusOK)
}

func (c Connection) PostURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Response{
			Message: "method not allowed",
			Status:  "failed",
		}.WriteJSON(w, http.StatusMethodNotAllowed)

		return
	}

	var u struct {
		Url string `json:"url"`
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil || u.Url == "" {
		Response{
			Message: "invalid json request",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	url, alreadyRegistered, err := c.db.URLAlreadyRegistered(u.Url, c.ctx)
	if err != nil {
		Response{
			Message: "error on checking if url is already registered",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	if alreadyRegistered {
		Response{
			Message: "url already registered",
			Status:  "failed",
			Content: url,
		}.WriteJSON(w, http.StatusConflict)

		return
	}

	shortCode, err := generateShortURL()
	if err != nil {
		Response{
			Message: "error on generating short code",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	for {
		exist, err := c.db.ShortURLExists(shortCode, c.ctx)
		if err != nil {
			Response{
				Message: "error on checking short code",
				Status:  "failed",
			}.WriteJSON(w, http.StatusInternalServerError)
		}

		if !exist {
			break
		}

		shortCode, err = generateShortURL()
		if err != nil {
			Response{
				Message: "error on generating short code",
				Status:  "failed",
			}.WriteJSON(w, http.StatusInternalServerError)

			return
		}
	}

	url = models.GenURLStruct(u.Url, shortCode)

	err = c.db.CreateURL(*url, c.ctx)
	if err != nil {
		Response{
			Message: "error on creating db entry",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  "ok",
		Content: url,
	}.WriteJSON(w, http.StatusOK)
}
