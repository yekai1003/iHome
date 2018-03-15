package controllers

import (
	"encoding/json"
	"fmt"

	"iHome/models"
	"iHome/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) RetData(resp interface{}) {
	beego.Info("resp ....", resp)
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *LoginController) Login() {
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	var loginfo = make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &loginfo)
	if err != nil {
		resp.Errno = utils.RECODE_LOGINERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//判断用户和密码是否有效
	var user models.User
	//验证用户和密码是否正确
	o := orm.NewOrm()
	err = o.QueryTable("user").Filter("mobile", loginfo["mobile"]).Filter("password_hash", loginfo["password"]).One(&user)
	if err != nil {
		resp.Errno = utils.RECODE_USERERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	beego.Info("query res is ....", user)
	//	var UserName = make(map[string]interface{})
	//	UserName["name"] = user.Name
	//	resp.Data = UserName
	beego.Info(resp)
	this.SetSession("name", user.Name)
	this.SetSession("mobile", user.Mobile)
	this.SetSession("user_id", user.Id)
	fmt.Printf("%+v\n", resp)
}
