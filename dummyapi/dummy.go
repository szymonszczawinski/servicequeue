package dummyapi

import (
	"servicequeue/service"
	"servicequeue/userapi"
)

type IDummyService interface {
	service.IProxy
	VoidMethod(msg string)
	NonVoidMethod(msg string) string
	RegisterUserListener(listener userapi.IUserListener)
}
