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

func (s *FooService) ReadAll() ([]model.Foo, error) {
	query := `SELECT * FROM foo`
	rows, err := s.sqlConnector.Query(query)
	if err != nil {
		return nil, fmt.Errorf("can't read all foos : %w", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			fmt.Println("Can't close sql rows : ", err.Error())
		}
	}()

	var foos []model.Foo
	for rows.Next() {
		foo := new(model.Foo)
		if err = rows.Scan(&foo.UUID, &foo.Msg); err != nil {
			return nil, fmt.Errorf("can't read all foos : %w", err)
		}
		foos = append(foos, *foo)
	}

	return foos, nil
}

func (s *FooService) Read(fooUUID uuid.UUID) (*model.Foo, error) {
	foo := new(model.Foo)
	query := `SELECT * FROM foo where foo_uuid=UUID_TO_BIN(?)`
	if err := s.sqlConnector.QueryRow(query, fooUUID).Scan(&foo.UUID, &foo.Msg); err != nil {
		return nil, fmt.Errorf("can't read foo by uuid : %w", err)
	}

	return foo, nil
}

func (s *FooService) Write(foo *model.PostFoo) error {
	tx, err := s.sqlConnector.Begin()
	if err != nil {
		return fmt.Errorf("can't begin write sql transacation : %w", err)
	}

	query := "INSERT INTO foo(msg) VALUES(?)"
	res, err := tx.Exec(query, foo.Msg)
	if err != nil {
		return fmt.Errorf("can't insert new foo : %w", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("can't insert new foo : : %w", err)
	}

	if rowCount != 1 {
		_ = tx.Rollback()
		return fmt.Errorf("invalid insert : %w", ErrNoRowsAffected)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("can't commit tx : %w", err)
	}

	return nil
}

func (s *FooService) Update(foo *model.Foo) (*model.Foo, error) {
	tx, err := s.sqlConnector.Begin()
	if err != nil {
		return nil, fmt.Errorf("can't begin sql transacation : %w", err)
	}

	query := "UPDATE foo SET msg=? WHERE foo_uuid=UUID_TO_BIN(?)"
	res, err := tx.Exec(query, foo.Msg)
	if err != nil {
		return nil, fmt.Errorf("can't insert new foo : %w", err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("can't insert new foo : : %w", err)
	}

	if rowCount != 1 {
		_ = tx.Rollback()
		return nil, fmt.Errorf("invalid insert : %w", ErrNoRowsAffected)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit tx : %w", err)
	}

	return s.Read(foo.UUID)
}
