package notes

import (
	"context"
	"gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/internal/model"
	servicepb "gitlab.ozon.dev/go/classroom-12/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) SaveNote(ctx context.Context, in *servicepb.SaveNoteRequest) (*servicepb.SaveNoteResponse, error) {

	// time.Sleep(time.Second * 2)
	// валидацию можно вызвать явно через err := in.Validate()
	// Validate - дает информацию об одной ошибке, ValidateAll - о всех, перечисленных в proto
	id, err := s.impl.SaveNote(ctx, repackNote(in))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &servicepb.SaveNoteResponse{NoteId: uint64(id)}, nil
}

func repackNote(in *servicepb.SaveNoteRequest) *model.Note {
	return &model.Note{
		Title:   in.Info.Title,
		Content: in.Info.Content,
	}
}
