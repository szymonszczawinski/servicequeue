package userapi

import "servicequeue/service"

type IUserService interface {
	service.IProxy
	VoidMethod()
}

type IUserListener interface {
	service.IProxy
	OnUser(user string)
}
