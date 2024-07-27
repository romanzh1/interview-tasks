package notes_usecase

import (
	"context"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"
)

func (s *Service) ListNotes(ctx context.Context) ([]*model.Note, error) {
	return s.repository.ListNotes(ctx)
}
