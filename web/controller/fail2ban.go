package controller

import (
    "fmt"
	"x-ui/web/service"
	"github.com/gin-gonic/gin"
)

type Fail2banController struct {
    Fail2banService     service.Fail2banService
}

func NewFail2banController(g *gin.RouterGroup) *Fail2banController {
	a := &Fail2banController{}
	a.initRouter(g)
	return a
}

func (a *Fail2banController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/fail2ban")

	g.POST("/", a.getFail2banStatus)
	g.POST("/install", a.installFail2banService)
	g.GET("/config", a.getFail2banConfig)
}

func (a *Fail2banController) getFail2banStatus(c *gin.Context) {
	fail2banStatus, err := a.Fail2banService.GetStatus()
    if err != nil {
        jsonMsg(c, "Error getting status", err)
        return
    }
	jsonObj(c, fail2banStatus, err)
}

func (a *Fail2banController) installFail2banService(c *gin.Context) {
	err := a.Fail2banService.InstallService()
    if err != nil {
        jsonMsg(c, "Error install", err)
        return
    }
	jsonObj(c, nil, err)
}

func (a *Fail2banController) getFail2banConfig(c *gin.Context) {
	fail2banConfig, err := a.Fail2banService.GetConfig()
    if err != nil {
        jsonMsg(c, "Error getting status", err)
        return
    }
    fmt.Println(fail2banConfig)
	jsonObj(c, fail2banConfig, err)
}