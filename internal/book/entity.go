package book

import (
	"errors"
)

type Book struct {
	ID        string `firestore:"id" json:"id"`
	Titulo    string `firestore:"titulo" json:"titulo"`
	Autor     string `firestore:"autor" json:"autor"`
	Status    Status `firestore:"status" json:"status"`
	Locatario string `firestore:"locatario" json:"locatario"`
}

type Status string

const (
	Disponivel Status = "disponivel"
	Emprestado Status = "emprestado"
)

func (s Status) validate() error {
	switch s {
	case Disponivel, Emprestado:
		return nil
	default:
		return errors.New("invalid status")
	}
}
func NewBook(b Book) (*Book, error) {
	if err := b.Status.validate(); err != nil {
		return nil, err
	}
	return &Book{
		Titulo: b.Titulo,
		Autor:  b.Autor,
		Status: "disponivel",
	}, nil
}

func (b *Book) Emprestar(nome string) error {
	switch b.Status {
	case "emprestado":
		return errors.New("Livro indisponivel")
	case "disponivel":
		b.Status = "emprestado"
		b.Locatario = nome
		return nil
	default:
		return errors.New("erro ao alugar")
	}
}
