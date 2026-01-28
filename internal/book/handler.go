package book

import (
	"encoding/json"
	"net/http"
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

}

func (h BookHandler) Get(w http.ResponseWriter, r *http.Request) {
	endPointId := r.PathValue("id")

	b, err := h.service.Get(r.Context(), endPointId)
	if err != nil {
		http.Error(w, "livro nao encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)

}
func (h BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var b Book

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "pedido errado", http.StatusBadRequest)
		return
	}

	id, err := h.service.NewBook(r.Context(), b)
	if err != nil {
		http.Error(w, "error ao salvar", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)

}

func (h BookHandler) Loan(w http.ResponseWriter, r *http.Request) {
	endPointId := r.PathValue("id")
	var recivedB Book
	if err := json.NewDecoder(r.Body).Decode(&recivedB); err != nil {
		http.Error(w, "pedido errado", http.StatusBadRequest)
		return
	}

	if err := h.service.LoanBook(r.Context(), endPointId, recivedB); err != nil {
		http.Error(w, "erro ao salvar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
