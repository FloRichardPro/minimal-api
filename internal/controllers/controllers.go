package controllers

var (
	Config Conf

	FooControllerInstance IFooController
)

type Conf struct {
}

func Init() error {
	FooControllerInstance = NewFooController()

	return nil
}
