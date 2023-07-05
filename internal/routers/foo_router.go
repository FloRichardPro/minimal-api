package routers

import (
	"net/http"
	"sync"

	"github.com/FloRichardPro/minimal-api/internal/controllers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	onceInitFooRouter sync.Once

	// singletonFooRouter is a singleton instance of FooRouter.
	singletonFooRouter *FooRouter
)

type FooRouter struct {
	ctlFoo controllers.IFooController
}

func (r *FooRouter) GetAll(c *gin.Context) {

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
	}

	c.JSON(http.StatusOK, foo)
}

func (r *FooRouter) Post(c *gin.Context) {

}

func (r *FooRouter) Put(c *gin.Context) {

}

func (r *FooRouter) Patch(c *gin.Context) {

}

// GetInstanceFooRouter get singleton instance of FooRouter.
func GetInstanceFooRouter() *FooRouter {
	if singletonFooRouter == nil {
		onceInitFooRouter.Do(
			func() {
				singletonFooRouter = &FooRouter{
					ctlFoo: controllers.FooControllerInstance,
				}
			},
		)
	}

	return singletonFooRouter
}
