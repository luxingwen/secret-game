package controller

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/luxingwen/secret-game/dao"
	"github.com/luxingwen/secret-game/model"
	"github.com/luxingwen/secret-game/tools"
)

type TeamController struct {
}

func (ctl *TeamController) Create(c *gin.Context) {
	team := new(model.Team)
	err := c.ShouldBind(&team)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	team.Created = time.Now()
	team.LeaderId = int64(c.GetInt("wxUserId"))
	err = dao.GetDao().AddTeam(team)

	if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		handleErr(c, CodeExist, errors.New("队伍名称已经存在"))
		return
	} else if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	handleOk(c, "ok")

}

func (ctl *TeamController) List(c *gin.Context) {
	search := new(model.TeamListSearch)
	err := c.ShouldBind(&search)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	search.UserId = int64(c.GetInt("wxUserId"))
	res, err := dao.GetDao().List(search)
	if err != nil {
		handleErr(c, CodeDBErr, err)
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
		handleErr(c, CodeReqErr, err)
		return
	}

	wxUser, err := dao.GetDao().GetWxUser(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	mdata := make(map[string]interface{}, 0)
	mdata["nickname"] = wxUser.NickName
	mdata["uid"] = wxUser.ID
	mdata["avatar_url"] = wxUser.AvatarUrl
	NotifyTeams(uid, "quit_team", mdata)
	err = dao.GetDao().BeforeJoinTeamQuitTeam(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	err = dao.GetDao().JoinTeam(uid, req.TeamId)
	if err != nil {
		handleErr(c, CodeDBErr, err)
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
		handleErr(c, CodeReqErr, err)
		return
	}

	wxUser, err := dao.GetDao().GetWxUser(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	mdata := make(map[string]interface{}, 0)
	mdata["nickname"] = wxUser.NickName
	mdata["uid"] = wxUser.ID
	mdata["avatar_url"] = wxUser.AvatarUrl
	NotifyTeams(uid, "quit_team", mdata)

	err = dao.GetDao().QuitTeam(uid, req.TeamId)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	handleOk(c, "ok")
}

func (ctl *TeamController) TeamInfo(c *gin.Context) {
	uid := c.GetInt("wxUserId")

	res, err := dao.GetDao().GetTeamInfo(uid)
	if err != nil && err.Error() != "record not found" {
		handleErr(c, CodeDBErr, err)
		return
	}
	handleOk(c, res)
}

// 查询队伍聊天
func (ctl *TeamController) TeamChatList(c *gin.Context) {
	uid := c.GetInt("wxUserId")

	res, err := dao.GetDao().GetTeamInfo(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	chatList, err := tools.TeamChatList(int(res.Id))
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	handleOk(c, chatList)
}

type Chat struct {
	Content string `json:"content"`
}

// 聊天
func (ctl *TeamController) TeamChat(c *gin.Context) {
	uid := c.GetInt("wxUserId")
	var chat Chat
	err := c.ShouldBind(&chat)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	res, err := dao.GetDao().GetTeamInfo(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	tools.TeamChat(uid, int(res.Id), chat.Content)
	wxUser, err := dao.GetDao().GetWxUser(uid)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	mdata := make(map[string]interface{}, 0)
	mdata["nickname"] = wxUser.NickName
	mdata["uid"] = wxUser.ID
	mdata["avatar_url"] = wxUser.AvatarUrl

	mdata["content"] = chat.Content

	NotifyTeams(uid, "team_chat", mdata)

	handleOk(c, "ok")
}

// 上传头像
func (ctl *TeamController) HeaderImg(c *gin.Context) {
	var saveUrl string
	// 头像上传
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
	}
	if file != nil {
		saveUrl = tools.GetHeadImgUrl(file.Filename)
		c.SaveUploadedFile(file, saveUrl)
	}
	handleOk(c, saveUrl)
}
