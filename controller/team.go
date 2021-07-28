package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"

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

	fmt.Println("--->", c.GetInt("wxUserId"))
	team.LeaderId = int64(c.GetInt("wxUserId"))
	err = dao.GetDao().AddTeam(team)

	if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		handleErr(c, errors.New("队伍名称已经存在"))
		return
	} else if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, "ok")

}

func (ctl *TeamController) List(c *gin.Context) {
	search := new(model.TeamListSearch)
	err := c.ShouldBind(&search)
	fmt.Println(search)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	search.UserId = int64(c.GetInt("wxUserId"))
	res, err := dao.GetDao().List(search)
	if err != nil {
		handleErr(c, err)
		return
	}

	handleOk(c, res)
}

type ReqJoinTeam struct {
	TeamId int `json:"team_id"`
}

func (ctl *TeamController) JoinTeam(c *gin.Context) {
	uid := c.GetInt("wxUserId")

	var req ReqJoinTeam

	err := c.ShouldBind(&req)
	if err != nil {
		handleErr(c, err)
		return
	}

	wxUser, err := dao.GetDao().GetWxUser(uid)
	if err != nil {
		handleErr(c, err)
		return
	}
	mdata := make(map[string]interface{}, 0)
	mdata["nickname"] = wxUser.NickName
	mdata["uid"] = wxUser.ID
	mdata["avatar_url"] = wxUser.AvatarUrl
	NotifyTeams(uid, "quit_team", mdata)
	err = dao.GetDao().BeforeJoinTeamQuitTeam(uid)
	if err != nil {
		handleErr(c, err)
		return
	}

	err = dao.GetDao().JoinTeam(uid, req.TeamId)
	if err != nil {
		handleErr(c, err)
		return
	}

	NotifyTeams(uid, "join_team", mdata)
	handleOk(c, "ok")
}

type ReqQuitTeam struct {
	TeamId int `json:"team_id"`
}

func (ctl *TeamController) QuiteTeam(c *gin.Context) {
	uid := c.GetInt("wxUserId")
	var req ReqQuitTeam

	err := c.ShouldBind(&req)
	if err != nil {
		handleErr(c, err)
		return
	}

	wxUser, err := dao.GetDao().GetWxUser(uid)
	if err != nil {
		handleErr(c, err)
		return
	}
	mdata := make(map[string]interface{}, 0)
	mdata["nickname"] = wxUser.NickName
	mdata["uid"] = wxUser.ID
	mdata["avatar_url"] = wxUser.AvatarUrl
	NotifyTeams(uid, "quit_team", mdata)

	err = dao.GetDao().QuitTeam(uid, req.TeamId)
	if err != nil {
		handleErr(c, err)
		return
	}

	handleOk(c, "ok")
}

func (ctl *TeamController) TeamInfo(c *gin.Context) {
	uid := c.GetInt("wxUserId")

	res, err := dao.GetDao().GetTeamInfo(uid)
	if err != nil {
		handleErr(c, err)
		return
	}
	handleOk(c, res)
}
