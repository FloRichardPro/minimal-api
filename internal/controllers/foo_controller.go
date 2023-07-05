package controllers

import (
	"fmt"

	"github.com/FloRichardPro/minimal-api/internal/model"
	"github.com/FloRichardPro/minimal-api/internal/services"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type IFooController interface {
	Read(fooUUID uuid.UUID) (*model.Foo, error)
	ReadAll() ([]model.Foo, error)
	Write(foo *model.PostFoo) error
	Update(foo *model.Foo) (*model.Foo, error)
}

type FooControllerConf struct {
}

type FooController struct {
	fooService services.FooService
}

func NewFooController() IFooController {
	return &FooController{
		fooService: *services.NewFooService(),
	}
}

func (ctl *FooController) ReadAll() ([]model.Foo, error) {
	foos, err := ctl.fooService.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("can't read all foos from data source : %w", err)
	}

	return foos, nil
}

func (ctl *FooController) Read(fooUUID uuid.UUID) (*model.Foo, error) {
	foo, err := ctl.fooService.Read(fooUUID)
	if err != nil {
		return nil, fmt.Errorf("can't read foo from data source : %w", err)
	}

	return foo, nil
}

func (ctl *FooController) Write(foo *model.PostFoo) error {
	if err := ctl.fooService.Write(foo); err != nil {
		return fmt.Errorf("can't write foo to data source : %w", err)
	}

	return nil
}

func (ctl *FooController) Update(foo *model.Foo) (*model.Foo, error) {
	oldFoo, err := ctl.fooService.Read(foo.UUID)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve not updated foo from data source : %w", err)
	}

	oldFooAsMap := make(map[string]any)

	if err := mapstructure.Decode(oldFoo, oldFooAsMap); err != nil {
		return nil, fmt.Errorf("can't decode foo to map : %w", err)
	}

	updatedFoo, err := ctl.fooService.Update(foo)
	if err != nil {
		return nil, fmt.Errorf("can't write foo to data source : %w", err)
	}

	return updatedFoo, nil
}
