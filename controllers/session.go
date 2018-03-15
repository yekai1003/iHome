package controllers

import (
	"iHome/utils"

	"github.com/astaxie/beego"
)

type SessionController struct {
	beego.Controller
}

func (this *SessionController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

type SessName struct {
	Name string `json:"name"`
}

func (this *SessionController) GetName() {
	beego.Info("GetSession is called")
	var resp UtilResp
	resp = UtilResp{Errno: utils.RECODE_USERERR, Errmsg: utils.RecodeText(resp.Errno)}
	defer this.RetData(&resp)
	var sName SessName
	if sName.Name = this.GetSession("name").(string); sName.Name != "" {
		beego.Info("GetName ....", resp)
		resp.Errno = utils.RECODE_OK
		resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		resp.Data = sName
		return
	}

}

//api/v1.0/session delete
func (this *SessionController) DelSess() {
	var resp UtilResp
	resp.Errno = utils.RECODE_SESSIONERR
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	this.DelSession("name")
	this.DelSession("user_id")
	this.DelSession("mobile")
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
}
