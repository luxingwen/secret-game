package dao

import (
	"github.com/luxingwen/secret-game/model"
)

func (d *Dao) AddWxUser(user *model.WxUser) (err error) {
	err = d.DB.Table(TableWxUser).Create(user).Error
	return
}

func (d *Dao) GetByOpenId(openId string) (wxUser *model.WxUser, err error) {
	wxUser = new(model.WxUser)
	err = d.DB.Table(TableWxUser).Where("open_id = ?", openId).First(&wxUser).Error
	return
}

func (d *Dao) GetWxUser(id int) (r *model.WxUser, err error) {
	r = new(model.WxUser)
	err = d.DB.Table(TableWxUser).Where("id = ?", id).First(&r).Error
	return
}

func (d *Dao) AddWxCode(wxCode *model.WxCode) (err error) {
	err = d.DB.Table(TableWxCode).Create(wxCode).Error
	return
}

func (d *Dao) GetWxCode(code string) (res *model.WxCode, err error) {
	res = new(model.WxCode)
	err = d.DB.Table(TableWxCode).Where("code = ?", code).First(&res).Error
	return
}
