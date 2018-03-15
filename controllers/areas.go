package controllers

import (
	"encoding/json"
	"ihome_go1/models"
	"time"

	"ihome_go1/utils"

	//	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	//	redigo "github.com/garyburd/redigo/redis"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *AreaController) GetAreas() {
	resp := UtilResp{Errno: "0", Errmsg: "OK"}
	var areas []models.Area
	defer this.RetData(&resp)
	//获得redis缓存数据
	//如果没有 则继续查询mysql
	//如果有，则将redis结果返回给前端 over

	ncache, err := cache.NewCache("redis", `{"key":"ihome_go1","conn":"101.200.170.171:6382","dbNum":"0"}`)
	if err != nil {
		beego.Info("redis conn err", err)
	} else {
		//查询数据，并且返回
		area_info_value := ncache.Get("area_info_str")
		if area_info_value == nil {
			beego.Info("area infos not in redis")

		} else {
			var area_infos interface{}
			json.Unmarshal(area_info_value.([]byte), &area_infos)
			beego.Info(area_infos)
			resp.Data = area_infos
			return //缓存中有不再执行后面的逻辑
		}

	}
	//查询数据库 mysql
	o := orm.NewOrm()
	num, err := o.QueryTable("area").All(&areas)
	if err != nil {
		//代表错误
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	if num <= 0 {
		//代表没有数据
		resp.Errno = utils.RECODE_NODATA
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//将数据返回给 前端

	resp.Data = areas

	//将数据放到redis数据库
	area_str, _ := json.Marshal(&areas)
	ncache.Put("area_info_str", area_str, time.Second*300)

}
