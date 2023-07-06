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
	query := `SELECT foo_uuid, msg, phone, email FROM foo`
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
		if err = rows.Scan(&foo.UUID, &foo.Msg, &foo.Phone, &foo.Email); err != nil {
			return nil, fmt.Errorf("can't read all foos : %w", err)
		}
		foos = append(foos, *foo)
	}

	return foos, nil
}

func (s *FooService) Read(fooUUID uuid.UUID) (*model.Foo, error) {
	foo := new(model.Foo)
	query := `SELECT foo_uuid, msg, phone, email FROM foo where foo_uuid=UUID_TO_BIN(?)`
	if err := s.sqlConnector.QueryRow(query, fooUUID).Scan(&foo.UUID, &foo.Msg, &foo.Phone, &foo.Email); err != nil {
		return nil, fmt.Errorf("can't read foo by uuid : %w", err)
	}

	return foo, nil
}

func (s *FooService) Write(foo *model.PostFoo) error {
	tx, err := s.sqlConnector.Begin()
	if err != nil {
		return fmt.Errorf("can't begin write sql transacation : %w", err)
	}

	query := "INSERT INTO foo(msg, phone, email) VALUES(?,?,?)"
	res, err := tx.Exec(query, foo.Msg, foo.Phone, foo.Email)
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

	query := "UPDATE foo SET msg=?, phone=?,email=? WHERE foo_uuid=UUID_TO_BIN(?)"
	res, err := tx.Exec(query, foo.Msg, foo.Phone, foo.Email, foo.UUID)
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
		_ = tx.Rollback()
		return nil, fmt.Errorf("can't commit tx : %w", err)
	}

	return s.Read(foo.UUID)
}
