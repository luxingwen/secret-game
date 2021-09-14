package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/luxingwen/secret-game/dao"
	//"github.com/luxingwen/secret-game/model"
)

type TestController struct {
}

func (ctl *TestController) List(c *gin.Context) {

	uid := c.GetInt("wxUserId")

	teamId, err := dao.GetDao().GetTeamIdByUserUid(uid)

	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	res, err := dao.GetDao().TeamTestList(teamId)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	handleOk(c, res)

}

func (ctl *TestController) Start(c *gin.Context) {
	uid := c.GetInt("wxUserId")
	team, err := dao.GetDao().GetTeamByLeaderId(int64(uid))
	if err == gorm.ErrRecordNotFound {
		handleErr(c, CodePermissions, errors.New("你不是队长"))
		return
	}

	err = dao.GetDao().TeamStartGame(team.Id)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	err = dao.GetDao().GenTeamTest(team.Id)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	mdata := make(map[string]interface{}, 0)
	mdata["status"] = 1
	mdata["id"] = team.Id

	NotifyTeams(uid, "start_test", mdata)

	handleOk(c, mdata)

}

type ReqAnswer struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
}

func (ctl *TestController) Answer(c *gin.Context) {
	//uid := c.GetInt("wxUserId")

	var req ReqAnswer
	err := c.ShouldBind(&req)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	subject, err := dao.GetDao().GetSubjectByTestId(req.Id)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	if subject.Answer == req.Content {
		err = dao.GetDao().TeatTestUpdateAnswerStatus(req.Id)
		if err != nil {
			return
		}
		handleOk(c, "回答正确")
		return
	}

	handleErr(c, 1, errors.New("回答错误"))
}

type RequestHit struct {
	Id int `json:"id"`
}

func (ctl *TestController) GetHits(c *gin.Context) {
	uid := c.GetInt("wxUserId")
	var req RequestHit
	err := c.ShouldBind(&req)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	hit, err := dao.GetDao().GetTeamTestHitsById(req.Id)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	mdata := make(map[string]interface{}, 0)

	mdata["id"] = req.Id
	mdata["hit"] = hit

	NotifyTeams(uid, "test_hit", mdata)
	handleOk(c, mdata)
}
