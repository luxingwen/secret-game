package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/luxingwen/secret-game/dao"
	"github.com/luxingwen/secret-game/model"
)

type TeamController struct {
}

func (ctl *TeamController) Create(c *gin.Context) {
	team := new(model.Team)
	err := c.ShouldBind(&team)
	if err != nil {
		handleErr(c, err)
		return
	}
	team.LeaderId = 1
	err = dao.GetDao().AddTeam(team)
	if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, "ok")

}

func (ctl *TeamController) List(c *gin.Context) {
	res, err := dao.GetDao().List()
	if err != nil {
		handleErr(c, err)
		return
	}

	handleOk(c, res)
}

func (ctl *TeamController) JoinTeam(c *gin.Context) {
	var teamId int64 = 1
	var uid int64 = 2
	err := dao.GetDao().JoinTeam(uid, teamId)
	if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, "ok")
}

func (ctl *TeamController) QuiteTeam(c *gin.Context) {
	var teamId int64 = 1
	var uid int64 = 2
	err := dao.GetDao().JoinTeam(uid, teamId)
	if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, "ok")
}
