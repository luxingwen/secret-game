package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func handleOk(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func handleErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": err.Error(),
	})
}

func ParseRequest(c *gin.Context, request interface{}) error {
	err := c.ShouldBindWith(request, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "parse Request Error",
			"error":   err.Error(),
		})

		log.Println("ParseRequest Result", request)
		log.Println("ParseRequest Error", err.Error())
		return err
	}
	return nil
}

func SuccessResponse(c *gin.Context, response interface{}) {
	handleOk(c, response)
}

func CheckErr(c *gin.Context, err error) {
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
}

func Router() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	api := r.Group("/api")
	api.GET("/ws", WsHandler)
	api.POST("/wxlogin", WxLogin)
	api.POST("/wxcode", WxSetCode)

	api.Use(WxJWT())
	team := api.Group("/team")
	teamCtl := &TeamController{}
	{
		team.GET("/list", teamCtl.List)
		team.POST("/add", teamCtl.Create)
		team.POST("/join", teamCtl.JoinTeam)
		team.POST("/quit", teamCtl.QuiteTeam)
		team.GET("/myinfo", teamCtl.TeamInfo)
	}

	testCtl := &TestController{}

	test := api.Group("/test")
	{
		test.GET("/list", testCtl.List)
		test.POST("/start", testCtl.Start)
		test.POST("/answer", testCtl.Answer)
		test.POST("/hits", testCtl.GetHits)
	}

	return r

}
