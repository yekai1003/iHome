package controllers

import (
	"encoding/json"
	_ "fmt"

	"iHome/utils"

	"iHome/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

//http://localhost:8080/api/v1.0/users post
//user register
func (this *UserController) UserReg() {
	beego.Info("UserReg is called....")
	var resp UtilResp
	resp = UtilResp{Errno: utils.RECODE_OK, Errmsg: utils.RecodeText(utils.RECODE_OK)}
	defer this.RetData(&resp)
	var reginfo = make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &reginfo)
	if err != nil {
		beego.Debug("marshal request body err", err)
		resp.Errno = utils.RECODE_PARAMERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	beego.Info(reginfo)
	if reginfo["mobile"] == "" || reginfo["sms_code"] == "" || reginfo["password"] == "" {
		resp.Errno = utils.RECODE_PARAMERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//同步到数据库
	o := orm.NewOrm()
	var r orm.RawSeter
	//r = o.Raw("UPDATE user SET name = ? WHERE name = ?", "testing", "slene")
	r = o.Raw("insert into user(name,password_hash,mobile) values(?,?,?)", reginfo["mobile"], reginfo["password"], reginfo["mobile"])
	_, err = r.Exec()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//代表登陆成功
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	var userName = make(map[string]interface{})
	userName["name"] = reginfo["mobile"]
	resp.Data = userName
	beego.Info(resp)
	this.SetSession("name", reginfo["mobile"])
	this.SetSession("mobile", reginfo["mobile"])

}

//8080/api/v1.0/user get 查询用户信息
func (this *UserController) GetUserInfo() {
	//从session中获取数据，首先得是登陆的用户，也就是说session中有user_id信息
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	mobile := this.GetSession("mobile")

	beego.Info("get mobile by session=====", mobile)
	//通过user_id去查询用户资料表
	o := orm.NewOrm()
	//rSql := fmt.Sprintf("select * from user where user_id=%d", userid)
	//rSql := "select * from user where user_id=" + string(userid)
	beego.Info("select * from user where mobile=?", mobile)
	r := o.Raw("select * from user where mobile=?", mobile)
	var userInfo models.User
	err := r.QueryRow(&userInfo)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	beego.Info("get user info by orm===", userInfo)

	resp.Data = userInfo
	//搞定

}

//更新用户名 api/v1.0/user/name

func (this *UserController) SetUserName() {
	beego.Info("SetUserName is called...")
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	var reqinfo = make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &reqinfo)
	if err != nil {
		resp.Errno = utils.RECODE_PARAMERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	beego.Info(reqinfo)
	//beego.Info(userName)
	//更新数据库
	o := orm.NewOrm()
	userid := this.GetSession("user_id")
	r := o.Raw("update user set name=? where id=?", reqinfo["name"], userid)
	_, err = r.Exec()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//更新session
	this.SetSession("name", reqinfo["name"])
	resp.Data = reqinfo

}

//更新实名认证 api/v1.0/user/auth

func (this *UserController) SetUserAuth() {
	beego.Info("SetUserAuth is called...")
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	var reqinfo = make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &reqinfo)
	if err != nil {
		resp.Errno = utils.RECODE_PARAMERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	beego.Info(reqinfo)
	//beego.Info(userName)
	//更新数据库
	o := orm.NewOrm()
	userid := this.GetSession("user_id")
	r := o.Raw("update user set real_name=?,id_card=? where id=?", reqinfo["real_name"], reqinfo["id_card"], userid)
	_, err = r.Exec()
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//更新session
	this.SetSession("name", reqinfo["name"])
	this.SetSession("user_id", userid)
	resp.Data = reqinfo

}

func (this *UserController) UploadAvatar() {
	beego.Info("upload avatar is called.....")
	var resp UtilResp
	resp.Errno = utils.RECODE_PARAMERR
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	file, header, err := this.GetFile("avatar")
	if err != nil {
		beego.Info("getfile  err...", err)
		return
	}
	defer file.Close()

	this.SaveToFile("avatar", header.Filename)
	beego.Info("UploadAvatar is running...", header.Filename)
	ava_url := UploadFile(header.Filename)
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	urldata := make(map[string]interface{})
	urldata["avatar_url"] = ava_url
	resp.Data = urldata
	//同步到数据库中
	o := orm.NewOrm()
	mobile := this.GetSession("mobile")
	r := o.Raw("update user set avatar_url=? where mobile=?", ava_url, mobile)
	r.Exec()
	beego.Info(resp)
	return
}

// /user/orders [GET]
func (this *UserController) GetOrders() {
	resp := UtilResp{Errno: utils.RECODE_OK, Errmsg: utils.RecodeText(utils.RECODE_OK)}
	defer this.RetData(&resp)
	//得到用户id
	user_id := this.GetSession("user_id").(int)
	//得到用户角色
	var role string
	this.Ctx.Input.Bind(&role, "role")

	if role == "" {
		resp.Errno = utils.RECODE_ROLEERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}

	o := orm.NewOrm()
	orders := []models.OrderHouse{}
	order_list := []interface{}{}

	if "landlord" == role {
		//角色为房东
		//现找到自己目前已经发布了哪些房子
		landLordHouses := []models.House{}
		o.QueryTable("house").Filter("user__id", user_id).All(&landLordHouses)
		housesIds := []int{}
		for _, house := range landLordHouses {
			housesIds = append(housesIds, house.Id)
		}
		//在从订单中找到房屋id为自己房源的id
		o.QueryTable("order_house").Filter("house__id__in", housesIds).OrderBy("-ctime").All(&orders)
	} else {
		//角色为租客
		o.QueryTable("order_house").Filter("user__id", user_id).OrderBy("-ctime").All(&orders)
	}

	for _, order := range orders {
		o.LoadRelated(&order, "User")
		o.LoadRelated(&order, "House")
		order_list = append(order_list, order.To_order_info())
	}

	data := map[string]interface{}{}
	data["orders"] = order_list

	resp.Data = data

	return
}
