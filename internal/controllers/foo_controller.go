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
	Update(fooUUID uuid.UUID, fooFieldsToUpdate map[string]any) (*model.Foo, error)
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

func (ctl *FooController) Update(fooUUID uuid.UUID, fooFieldsToUpdate map[string]any) (*model.Foo, error) {
	// Reads the original foo
	foo, err := ctl.fooService.Read(fooUUID)
	if err != nil {
		return nil, fmt.Errorf("can't retrieve not updated foo from data source : %w", err)
	}

	// Transform original fool to a map
	fooAsMap := make(map[string]any)
	if err := mapstructure.Decode(foo, &fooAsMap); err != nil {
		return nil, fmt.Errorf("can't decode foo to map : %w", err)
	}

	// Apply changes from fooFieldsToUpdate map
	for fieldName, value := range fooFieldsToUpdate {
		fooAsMap[fieldName] = value
	}

	// Transform the original update foo map to a foo object
	if err := mapstructure.Decode(fooAsMap, &foo); err != nil {
		return nil, fmt.Errorf("can't decode foo to map : %w", err)
	}

	foo.UUID = fooUUID

	fmt.Println("Update foo = ", fooAsMap)

	// Update foo
	updatedFoo, err := ctl.fooService.Update(foo)
	if err != nil {
		return nil, fmt.Errorf("can't write foo to data source : %w", err)
	}

	return updatedFoo, nil
}
