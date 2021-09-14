package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luxingwen/secret-game/dao"
	"github.com/luxingwen/secret-game/model"
	"github.com/luxingwen/secret-game/tools"
	"github.com/medivhzhan/weapp/v2"
)

type Wxlogin struct {
	Code      string `json:"code"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
	Gender    int    `json:"gender"`
}

const (
	//AppId     = "wx8ded9e99c86ce8c2"
	//AppSecret = "8da44079253130c4ebc93b1758eacdc0"

	AppId     = "wx8e199f928cba4bf4"
	AppSecret = "9720e0397a659f75f26dc53f6bbb1205"
)

// @Summary 登录
// @Accept json
// @Produce  json
// @Param param body models.User true "{}"
// @Success 200 {string} string "{"code":0,"data":{},"msg":"ok"}"
// @Router /wx/login [post]
func WxLogin(c *gin.Context) {

	mdata := make(map[string]interface{}, 0)
	err := c.ShouldBind(&mdata)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	wxCode, err := dao.GetDao().GetWxCode(mdata["cache_key"].(string))
	if err != nil {
		handleErr(c, CodeNotFound, errors.New("没有找到sessionkey"))
		return
	}

	encryptedData := mdata["encryptedData"].(string)
	rawData := mdata["rawData"].(string)
	sign := mdata["signature"].(string)
	iv := mdata["iv"].(string)

	res, err := weapp.DecryptUserInfo(wxCode.SessionKey, rawData, encryptedData, sign, iv)
	if err != nil {
		// 处理一般错误信息
		fmt.Println("res err:", err)
		handleErr(c, CodeReqErr, err)
		return
	}

	wxUser := &model.WxUser{NickName: res.Nickname, AvatarUrl: res.Avatar, Gender: res.Gender, OpenId: wxCode.OpenID}

	oldWxUser, err := dao.GetDao().GetByOpenId(wxCode.OpenID)

	if err != nil && err.Error() == "record not found" {
		err = dao.GetDao().AddWxUser(wxUser)
		if err != nil {
			handleErr(c, CodeDBErr, err)
			return
		}
	} else if err == nil && oldWxUser != nil {
		wxUser = oldWxUser
	}

	fmt.Printf("wxUser: %d \n", wxUser.ID)

	token, err := GenerateWxToken(int(wxUser.ID))
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}
	rmdata := make(map[string]interface{}, 0)
	rmdata["token"] = token
	rmdata["uid"] = wxUser.ID

	handleOk(c, rmdata)
}

func WxSetCode(c *gin.Context) {
	req := new(Wxlogin)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	fmt.Printf("req: %#v\n", req)
	fmt.Println("code:", req.Code)

	res, err := weapp.Login(AppId, AppSecret, req.Code)
	if err != nil {
		// 处理一般错误信息
		fmt.Println("login err:", err)
		handleErr(c, CodeDBErr, err)
		return
	}

	fmt.Printf("返回结果: %#v\n", res)

	if err := res.GetResponseError(); err != nil {
		// 处理微信返回错误信息
		fmt.Println("GetResponseError:", err)
		handleErr(c, CodeDBErr, err)
		return
	}

	cacheKey := fmt.Sprintf("api_code_%s_%d", req.Code, time.Now().Unix())
	cacheKey = tools.Md5(cacheKey)
	wxCode := &model.WxCode{Code: cacheKey, SessionKey: res.SessionKey, OpenID: res.OpenID}
	err = dao.GetDao().AddWxCode(wxCode)
	if err != nil {
		handleErr(c, CodeDBErr, err)
		return
	}

	mdata := make(map[string]interface{}, 0)
	mdata["cache_key"] = cacheKey

	b, err := json.Marshal(res)
	if err != nil {
		handleErr(c, CodeReqErr, err)
		return
	}

	fmt.Println("set_code:", string(b))
	handleOk(c, mdata)

}
