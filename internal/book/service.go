package book

import (
	"context"
)

type Service interface {
	NewBook(ctx context.Context, b Book) (string, error)
	LoanBook(ctx context.Context, id string, b Book) error
	Get(ctx context.Context, id string) (*Book, error)
	ReturnBook(ctx context.Context, id string) error
}

type BookService struct {
	repo Repository
}

func NewService(r Repository) *BookService {
	return &BookService{
		repo: r,
	}
}

func (s BookService) NewBook(ctx context.Context, b Book) (string, error) {
	return s.repo.Create(ctx, b)
}
func (s BookService) Get(ctx context.Context, id string) (*Book, error) {
	return s.repo.FindId(ctx, id)
}

func (s BookService) LoanBook(ctx context.Context, id string, b Book) error {
	bencontado, err := s.repo.FindId(ctx, id)
	if err != nil {
		return err
	}

	if err := bencontado.Emprestar(b.Locatario); err != nil {
		return err
	}

	err = s.repo.Update(ctx, id, bencontado)
	if err != nil {
		return err
	}
	return nil
}

func (s BookService) ReturnBook(ctx context.Context, id string) error {
	bencontado, err := s.repo.FindId(ctx, id)
	if err != nil {
		return err
	}

	if err := bencontado.Devolver(); err != nil {
		return err
	}

	err = s.repo.Update(ctx, id, bencontado)
	if err != nil {
		return err
	}
	return nil
}
