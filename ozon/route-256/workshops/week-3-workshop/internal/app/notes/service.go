package notes

import (
	"context"

	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"
	servicepb "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/pkg/api/notes/v1"
)

/*
*
Небольшая хитрость для проверки себя
Поскольку мы приводим nil к указателю на Service
То на этапе компиляции мы получим ошибку, если Service не удовлетворяет интерфейсу генерированного NotesServer
*/
var _ servicepb.NotesServer = (*Service)(nil)

type NoteService interface {
	ListNotes(ctx context.Context) ([]*model.Note, error)
	SaveNote(ctx context.Context, note *model.Note) (int, error)
}

type Service struct {
	servicepb.UnimplementedNotesServer
	impl NoteService
}

func NewService(impl NoteService) *Service {
	return &Service{impl: impl}
}
