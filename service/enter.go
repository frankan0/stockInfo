package service

import "api.frank.top/spider/service/wall"

type ServiceGroup struct {
	WallServiceGroup   wall.WallService
}

var ServiceGroupApp = new(ServiceGroup)
