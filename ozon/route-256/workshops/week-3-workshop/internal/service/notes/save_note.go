package notes_usecase

import (
	"context"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"
)

func (s *Service) SaveNote(ctx context.Context, note *model.Note) (int, error) {
	return s.repository.SaveNote(ctx, note)
}
