package services

import (
	"database/sql"
	"fmt"

	"github.com/FloRichardPro/minimal-api/internal/model"
	"github.com/google/uuid"
)

type FooService struct {
	sqlConnector *sql.DB
}

func NewFooService() *FooService {
	return &FooService{
		sqlConnector: sqlConnectorSharedInstance,
	}
}

func (s *FooService) Read(fooUUID uuid.UUID) (*model.Foo, error) {
	foo := new(model.Foo)
	query := `SELECT * FROM foo where foo_uuid=UUID_TO_BIN(?)`
	if err := s.sqlConnector.QueryRow(query, fooUUID).Scan(&foo.UUID, &foo.Msg); err != nil {
		return nil, fmt.Errorf("can't read foo by uuid : %w", err)
	}

	return foo, nil
}
