package book

import (
	"encoding/json"
	"net/http"

	"github.com/ArtoIi/BIBLIOTECA/internal/web"
)

type BookHandler struct {
	service Service
}

func NewHandler(s Service) *BookHandler {
	return &BookHandler{service: s}

}

func (h BookHandler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /book", h.Create)
	mux.HandleFunc("GET /book/{id}", h.Get)
	mux.HandleFunc("PUT /book/{id}", h.Loan)
	mux.HandleFunc("PUT /book/{id}/return", h.Return)

}

func (h BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	endPointId := r.PathValue("id")

	b, err := h.service.Get(r.Context(), endPointId)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, b)

}
func (h BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var b Book

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.NewBook(r.Context(), b)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
	}

	web.Respond(w, http.StatusOK, id)
}

func (h BookHandler) Loan(w http.ResponseWriter, r *http.Request) {
	endPointId := r.PathValue("id")
	var recivedB Book
	if err := json.NewDecoder(r.Body).Decode(&recivedB); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.service.LoanBook(r.Context(), endPointId, recivedB); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}
	web.Respond(w, http.StatusOK, nil)
}

func (h BookHandler) Return(w http.ResponseWriter, r *http.Request) {
	endPointId := r.PathValue("id")

	if err := h.service.ReturnBook(r.Context(), endPointId); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}
	web.Respond(w, http.StatusOK, nil)
}
