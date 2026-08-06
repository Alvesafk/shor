package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Alvesafk/shor/internal/app"
	"github.com/Alvesafk/shor/internal/models"
)

// TODO: use a constant for Status field on the Response struct.
// TODO: make a validate function to validate if an shortURL exists in the db.

type Connection struct {
	db  *app.FireDB
	ctx context.Context // TODO: drop ctx field on connection struct for context of the request
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

			return
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

func (c Connection) GetURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	if shortUrlString == "" {
		Response{
			Message: "invalid json request",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	url, err := c.db.GetURLByShortCode(shortUrlString, c.ctx)
	if err != nil {
		if errors.Is(err, app.UrlNotFound) {
			Response{
				Message: "url was not found",
				Status:  "failed",
			}.WriteJSON(w, http.StatusNotFound)

			return
		}

		Response{
			Message: "could not find url",
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

func (c Connection) UpdateURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	if shortUrlString == "" {
		Response{
			Message: "invalid json request",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	exist, err := c.db.ShortURLExists(shortUrlString, c.ctx)
	if err != nil {
		if errors.Is(err, app.UrlNotFound) {
			Response{
				Message: "url not found",
				Status:  "failed",
			}.WriteJSON(w, http.StatusNotFound)

			return
		}

		Response{
			Message: "could not find url",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	if !exist {
		Response{
			Message: "short url does not exist",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

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

	url, err := c.db.UpdateURL(u.Url, shortUrlString, c.ctx)
	if err != nil {
		Response{
			Message: "could not update url",
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

func (c Connection) DeleteURL(w http.ResponseWriter, r *http.Request) {
	shortUrlString := r.PathValue("shortUrl")
	if shortUrlString == "" {
		Response{
			Message: "invalid json request",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	exist, err := c.db.ShortURLExists(shortUrlString, c.ctx)
	if err != nil {
		if errors.Is(err, app.UrlNotFound) {
			Response{
				Message: "url not found",
				Status:  "failed",
			}.WriteJSON(w, http.StatusNotFound)

			return
		}

		Response{
			Message: "could not find url",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	if !exist {
		Response{
			Message: "short url does not exist",
			Status:  "failed",
		}.WriteJSON(w, http.StatusBadRequest)

		return
	}

	if err := c.db.DeleteURL(shortUrlString, c.ctx); err != nil {
		Response{
			Message: "could not delete the url",
			Status:  "failed",
		}.WriteJSON(w, http.StatusInternalServerError)

		return
	}

	Response{
		Message: "success",
		Status:  "ok",
	}.WriteJSON(w, http.StatusNoContent)
}
