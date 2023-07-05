package controllers

import (
	"fmt"

	"github.com/FloRichardPro/minimal-api/internal/model"
	"github.com/FloRichardPro/minimal-api/internal/services"
	"github.com/google/uuid"
)

type IFooController interface {
	Read(fooUUID uuid.UUID) (*model.Foo, error)
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

func (ctl *FooController) Read(fooUUID uuid.UUID) (*model.Foo, error) {
	foo, err := ctl.fooService.Read(fooUUID)
	if err != nil {
		return nil, fmt.Errorf("can't read foo from data source : %w", err)
	}

	return foo, nil
}
