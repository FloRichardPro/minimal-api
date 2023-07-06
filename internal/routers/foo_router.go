package routers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/FloRichardPro/minimal-api/internal/controllers"
	"github.com/FloRichardPro/minimal-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

var (
	onceInitFooRouter sync.Once

	// singletonFooRouter is a singleton instance of FooRouter.
	singletonFooRouter *FooRouter
)

type FooRouter struct {
	ctlFoo   controllers.IFooController
	validate *validator.Validate
}

func (r *FooRouter) GetAll(c *gin.Context) {
	foos, err := r.ctlFoo.ReadAll()
	if err != nil {
		Log.Error("FooRouter.GetAll fail : can't get all foos : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, foos)
}

func (r *FooRouter) Get(c *gin.Context) {
	fooUUIDAsString := c.Param("foo_uuid")
	fooUUID, err := uuid.Parse(fooUUIDAsString)
	if err != nil {
		Log.Error("FooRouter.Get fail : can't parse foo uuid path parameter", zap.String("foo_uuid", fooUUIDAsString), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid foo_uuid path parameter")
		return
	}

	foo, err := r.ctlFoo.Read(fooUUID)
	if err != nil {
		Log.Error("FooRouter.Get fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid resource id")
		return
	}

	c.JSON(http.StatusOK, foo)
}

func (r *FooRouter) Post(c *gin.Context) {
	foo := new(model.PostFoo)
	if err := c.ShouldBindJSON(foo); err != nil {
		Log.Error("FooRouter.Post fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := r.validate.Struct(foo); err != nil {
		Log.Error("FooRouter.Post fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := r.ctlFoo.Write(foo); err != nil {
		Log.Error("FooRouter.Post fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "bad request")
		return
	}

	c.JSON(http.StatusCreated, "resource created")

}

func (r *FooRouter) Put(c *gin.Context) {
	fooUUIDAsString := c.Param("foo_uuid")
	fooUUID, err := uuid.Parse(fooUUIDAsString)
	if err != nil {
		Log.Error("FooRouter.Get Put : can't parse foo uuid path parameter", zap.String("foo_uuid", fooUUIDAsString), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid foo_uuid path parameter")
		return
	}

	fooFieldsToUpdate := map[string]any{}
	if err := c.ShouldBindJSON(&fooFieldsToUpdate); err != nil {
		Log.Error("FooRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "bad request")
		return
	}

	if len(fooFieldsToUpdate) <= 0 {
		c.JSON(http.StatusNoContent, "no fields to update")
		return
	}

	// Decode incoming request body to a Foo object.
	fooUpdater := new(model.UpdateFoo)
	config := &mapstructure.DecoderConfig{
		ErrorUnused: true, // Extra fields throw an error.
		Result:      &fooUpdater,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		Log.Error("FooRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
		return
	}

	if err := decoder.Decode(fooFieldsToUpdate); err != nil {
		Log.Error("FooRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Can't decode request body")
		return
	}

	fmt.Println("Foo update : ", *fooUpdater)

	// Iterate over request body keys, at this point there is only Foo fields.
	fieldsToUpdate := make([]string, 0)
	for fieldName := range fooFieldsToUpdate {
		fieldsToUpdate = append(fieldsToUpdate, fieldName)
	}

	// Only validate the fields that appears in request body
	if err := r.validate.StructPartial(fooUpdater, fieldsToUpdate...); err != nil {
		Log.Error("FooRouter.Put fail : invalid request body fields: %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "Can't decode request body")
		return
	}

	updatedFoo, err := r.ctlFoo.Update(fooUUID, fooFieldsToUpdate)
	if err != nil {
		Log.Error("FooRouter.Put fail : %w", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, "can't update the resource")
		return
	}

	c.JSON(http.StatusOK, updatedFoo)
}

func (r *FooRouter) Patch(c *gin.Context) {

}

// GetInstanceFooRouter get singleton instance of FooRouter.
func GetInstanceFooRouter() *FooRouter {
	if singletonFooRouter == nil {
		onceInitFooRouter.Do(
			func() {
				singletonFooRouter = &FooRouter{
					ctlFoo:   controllers.FooControllerInstance,
					validate: ValidateInstance,
				}
			},
		)
	}

	return singletonFooRouter
}
