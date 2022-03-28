package wall

import (
	"api.frank.top/spider/global"
	"api.frank.top/spider/model"
)

type WallService struct {

}


func (wallService *WallService) AddWall(w model.Wall) (err error) {
	err = global.GVA_DB.Create(&w).Error
	return err
}

func (wallService *WallService) BatchAdd(walls []model.Wall)(err error)  {
	err = global.GVA_DB.Create(&walls).Error
	return err
}
