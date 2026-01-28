package book

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Repository interface {
	Create(ctx context.Context, b Book) (string, error)
	FindId(ctx context.Context, id string) (*Book, error)
	Update(ctx context.Context, id string, b *Book) error
}

type BookRepository struct {
	client *firestore.Client
}

func NewRepository(client *firestore.Client) *BookRepository {
	return &BookRepository{client: client}
}

func (r *BookRepository) Create(ctx context.Context, b Book) (string, error) {
	collection := r.client.Collection("biblioteca")
	docref := collection.NewDoc()

	data := Book{
		ID:        docref.ID,
		Titulo:    b.Titulo,
		Autor:     b.Autor,
		Status:    "disponivel",
		Locatario: b.Locatario,
	}

	if _, err := docref.Set(ctx, data); err != nil {
		return "", err
	}

	return docref.ID, nil
}

func (r BookRepository) FindId(ctx context.Context, id string) (*Book, error) {
	doc, err := r.client.Collection("biblioteca").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var b Book
	if err := doc.DataTo(&b); err != nil {
		return nil, err
	}

	return &Book{
		ID:        b.ID,
		Titulo:    b.Titulo,
		Autor:     b.Autor,
		Status:    b.Status,
		Locatario: b.Locatario,
	}, nil
}

func (r *BookRepository) Update(ctx context.Context, id string, b *Book) error {
	docref := r.client.Collection("biblioteca").Doc(id)
	docSnap, err := docref.Get(ctx)
	if err != nil {
		return err
	}

	var before Book
	if err := docSnap.DataTo(&before); err != nil {
		return err
	}

	var updates []firestore.Update
	if b.Titulo != "" {
		updates = append(updates, firestore.Update{Path: "titulo", Value: b.Titulo})
	}
	if b.Autor != "" {
		updates = append(updates, firestore.Update{Path: "autor", Value: b.Autor})
	}
	if b.Locatario != "" {
		updates = append(updates, firestore.Update{Path: "locatario", Value: b.Locatario})
	}
	if b.Status != "" {
		updates = append(updates, firestore.Update{Path: "status", Value: b.Status})
	}

	if len(updates) > 0 {
		if _, err := docref.Update(ctx, updates); err != nil {
			return err
		}
	}

	docSnapAfter, err := docref.Get(ctx)
	if err != nil {
		return err
	}

	var after Book
	if err := docSnapAfter.DataTo(&after); err != nil {
		return err
	}

	return nil
}
