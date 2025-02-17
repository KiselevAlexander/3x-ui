package controller

import (
	"errors"
	"time"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	"x-ui/web/entity"
	"x-ui/web/service"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type updateUserForm struct {
	OldUsername string `json:"oldUsername" form:"oldUsername"`
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewUsername string `json:"newUsername" form:"newUsername"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

type updateSecretForm struct {
	LoginSecret string `json:"loginSecret" form:"loginSecret"`
}

type SettingController struct {
	settingService service.SettingService
	userService    service.UserService
	panelService   service.PanelService
}

type ApiTokenResponse struct {
	Token string `json:"token"`
}

func NewSettingController(g *gin.RouterGroup) *SettingController {
	a := &SettingController{}
	a.initRouter(g)
	return a
}

func (a *SettingController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/setting")

	g.POST("/all", a.getAllSetting)
	g.POST("/defaultSettings", a.getDefaultSettings)
	g.POST("/update", a.updateSetting)
	g.POST("/updateUser", a.updateUser)
	g.POST("/restartPanel", a.restartPanel)
	g.GET("/getDefaultJsonConfig", a.getDefaultXrayConfig)
	g.POST("/updateUserSecret", a.updateSecret)
	g.POST("/getUserSecret", a.getUserSecret)

	g.GET("/apiToken", a.getApiToken)
	g.POST("/apiToken", a.generateApiToken)
	g.DELETE("/apiToken", a.removeApiToken)
}

func (a *SettingController) getAllSetting(c *gin.Context) {
	allSetting, err := a.settingService.GetAllSetting()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	jsonObj(c, allSetting, nil)
}

func (a *SettingController) getDefaultSettings(c *gin.Context) {
	result, err := a.settingService.GetDefaultSettings(c.Request.Host)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	jsonObj(c, result, nil)
}

func (a *SettingController) updateSetting(c *gin.Context) {
	allSetting := &entity.AllSetting{}
	err := c.ShouldBind(allSetting)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), err)
		return
	}
	err = a.settingService.UpdateAllSetting(allSetting)
	jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), err)
}

func (a *SettingController) updateUser(c *gin.Context) {
	form := &updateUserForm{}
	err := c.ShouldBind(form)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), err)
		return
	}
	user := session.GetSessionUser(c)
	if user.Username != form.OldUsername || user.Password != form.OldPassword {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifyUser"), errors.New(I18nWeb(c, "pages.settings.toasts.originalUserPassIncorrect")))
		return
	}
	if form.NewUsername == "" || form.NewPassword == "" {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifyUser"), errors.New(I18nWeb(c, "pages.settings.toasts.userPassMustBeNotEmpty")))
		return
	}
	err = a.userService.UpdateUser(user.Id, form.NewUsername, form.NewPassword)
	if err == nil {
		user.Username = form.NewUsername
		user.Password = form.NewPassword
		session.SetSessionUser(c, user)
	}
	jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifyUser"), err)
}

func (a *SettingController) restartPanel(c *gin.Context) {
	err := a.panelService.RestartPanel(time.Second * 3)
	jsonMsg(c, I18nWeb(c, "pages.settings.restartPanel"), err)
}

func (a *SettingController) updateSecret(c *gin.Context) {
	form := &updateSecretForm{}
	err := c.ShouldBind(form)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifySettings"), err)
	}
	user := session.GetSessionUser(c)
	err = a.userService.UpdateUserSecret(user.Id, form.LoginSecret)
	if err == nil {
		user.LoginSecret = form.LoginSecret
		session.SetSessionUser(c, user)
	}
	jsonMsg(c, I18nWeb(c, "pages.settings.toasts.modifyUser"), err)
}

func (a *SettingController) getUserSecret(c *gin.Context) {
	loginUser := session.GetSessionUser(c)
	user := a.userService.GetUserSecret(loginUser.Id)
	if user != nil {
		jsonObj(c, user, nil)
	}
}

func (a *SettingController) getDefaultXrayConfig(c *gin.Context) {
	defaultJsonConfig, err := a.settingService.GetDefaultXrayConfig()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.settings.toasts.getSettings"), err)
		return
	}
	jsonObj(c, defaultJsonConfig, nil)
}

func (a *SettingController) getApiToken(c *gin.Context) {
	response := &ApiTokenResponse{}
	token, err := a.settingService.GetApiToken()
	if err != nil {
		jsonObj(c, response , err)
		return
	}

	response.Token = token

	jsonObj(c, response , nil)
}

func (a *SettingController) generateApiToken(c *gin.Context) {
	response := &ApiTokenResponse{}
	randomBytes := make([]byte, 32)

	_, err := rand.Read(randomBytes)
	if err != nil {
		jsonObj(c, nil, err)
		return
	}

	hash := sha512.Sum512(randomBytes)
	response.Token = hex.EncodeToString(hash[:])

	saveErr := a.settingService.SaveApiToken(response.Token)

	if saveErr != nil {
		jsonObj(c, nil, saveErr)
		return
	}

	jsonMsgObj(c, I18nWeb(c, "pages.settings.security.apiTokenGeneratedSuccess"), response, nil)
}

func (a *SettingController) removeApiToken(c *gin.Context) {
	err := a.settingService.RemoveApiToken()

	if err != nil {
		jsonObj(c, nil, err)
		return
	}

	jsonMsg(c, "Removed", nil)
}
