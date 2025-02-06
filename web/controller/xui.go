package controller

import (
	"github.com/gin-gonic/gin"
)

type XUIController struct {
	BaseController

	inboundController       *InboundController
	settingController       *SettingController
	xraySettingController   *XraySettingController
	fail2banController      *Fail2banController
}

func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	return a
}

func (a *XUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/panel")
	g.Use(a.checkLogin)

	g.GET("/", a.index)
	g.GET("/inbounds", a.inbounds)
	g.GET("/settings", a.settings)
	g.GET("/xray", a.xraySettings)
	g.GET("/fail2ban", a.fail2ban)

	a.inboundController = NewInboundController(g)
	a.settingController = NewSettingController(g)
	a.xraySettingController = NewXraySettingController(g)
	a.fail2banController = NewFail2banController(g)
}

func (a *XUIController) index(c *gin.Context) {
	html(c, "index.html", "pages.index.title", nil)
}

func (a *XUIController) inbounds(c *gin.Context) {
	html(c, "inbounds.html", "pages.inbounds.title", nil)
}

func (a *XUIController) settings(c *gin.Context) {
	html(c, "settings.html", "pages.settings.title", nil)
}

func (a *XUIController) xraySettings(c *gin.Context) {
	html(c, "xray.html", "pages.xray.title", nil)
}

func (a *XUIController) fail2ban(c *gin.Context) {
	html(c, "fail2ban.html", "pages.fail2ban.title", nil)
}
