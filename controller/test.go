package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/luxingwen/secret-game/dao"
	//"github.com/luxingwen/secret-game/model"
)

type TestController struct {
}

func (ctl *TestController) List(c *gin.Context) {
	var teamId int64 = 1

	res, err := dao.GetDao().TeamTestList(teamId)
	if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, res)

}
