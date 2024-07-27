package notes_usecase

import (
	"context"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"
)

type Repository interface {
	ListNotes(ctx context.Context) ([]*model.Note, error)
	SaveNote(ctx context.Context, note *model.Note) (int, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
