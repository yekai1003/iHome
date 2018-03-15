package controllers

import (
	"encoding/json"
	"fmt"
	"ihome_go1/models"
	"ihome_go1/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
)

type HouseController struct {
	beego.Controller
}
type RespHouses struct {
	Houses []models.UserHouseInfo `json:"houses"`
}

func (this *HouseController) RetData(resp interface{}) {
	beego.Info("call RetData resp=======", resp)
	this.Data["json"] = resp
	this.ServeJSON()
}

func (this *HouseController) GetHouseIndex() {
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)

	data := []interface{}{}

	beego.Debug("Index Houses....")

	//1 从缓存服务器中请求 "home_page_data" 字段,如果有值就直接返回
	//先从缓存中获取房屋数据,将缓存数据返回前端即可
	redis_config_map := map[string]string{
		"key":   "ihome_go",
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	redis_config, _ := json.Marshal(redis_config_map)

	beego.Info("call GetHouseIndex......", string(redis_config))
	cache_conn, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Debug("connect cache error")
	}
	house_page_key := "home_page_data"
	house_page_value := cache_conn.Get(house_page_key)
	if house_page_value != nil {
		beego.Debug("======= get house page info  from CACHE!!! ========")
		json.Unmarshal(house_page_value.([]byte), &data)
		resp.Data = data
		return
	}

	houses := []models.House{}

	//2 如果缓存没有,需要从数据库中查询到房屋列表
	o := orm.NewOrm()

	if _, err := o.QueryTable("house").Limit(models.HOME_PAGE_MAX_HOUSES).All(&houses); err == nil {
		for _, house := range houses {
			o.LoadRelated(&house, "Area")
			o.LoadRelated(&house, "User")
			o.LoadRelated(&house, "Images")
			o.LoadRelated(&house, "Facilities")
			data = append(data, house.To_house_info())
		}

	}

	//将data存入缓存数据
	house_page_value, _ = json.Marshal(data)
	cache_conn.Put(house_page_key, house_page_value, 3600*time.Second)

	//返回前端data
	resp.Data = data

}
func (this *HouseController) GetUserHouses() {
	var resp UtilResp
	resp.Errno = utils.RECODE_USERERR
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	//var houses_info []models.House{}
	userId := this.GetSession("user_id")
	beego.Info("call get user houses by user_id===", userId)
	//查询mysql 获得当前用户的房源信息
	o := orm.NewOrm()
	r := o.Raw("select address,area.name area_name,ctime,house.id house_id,index_image_url image_url,order_count,price,room_count,title,user.avatar_url user_avatar   from user,area,house   where user.id = house.user_id     and house.area_id=area.id and user.id=?", userId)
	var houses_info []models.UserHouseInfo
	num, err := r.QueryRows(&houses_info)
	if err != nil || num <= 0 {
		beego.Info("GetUserHouses err===", err, num)
		return
	}
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	var resphouses RespHouses
	resphouses.Houses = houses_info
	resp.Data = resphouses
}

//用户发布房源信息
func (this *HouseController) UserPublishInfo() {
	beego.Info("run UserPublishInfo-----")
	var resp UtilResp
	resp.Errno = utils.RECODE_PARAMERR
	resp.Errmsg = utils.RecodeText(resp.Errno)
	beego.Info("call UserPublishInfo.....", resp)
	defer this.RetData(&resp)
	var house_info models.House
	var mapInfo = make(map[string]interface{})
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &mapInfo)
	if err != nil {
		beego.Info("UserPublishInfo json unmarshal err", err)
		return
	}
	beego.Info(mapInfo)
	var user models.User
	user.Id = this.GetSession("user_id").(int)

	house_info.Title = mapInfo["title"].(string)
	house_info.Address = mapInfo["address"].(string)
	house_info.Price, _ = strconv.Atoi(mapInfo["price"].(string))
	house_info.Price = house_info.Price * 100
	house_info.Area = &models.Area{Id: 0}
	house_info.Area.Id, _ = strconv.Atoi(mapInfo["area_id"].(string))
	house_info.Room_count, _ = strconv.Atoi(mapInfo["room_count"].(string))
	house_info.Acreage, _ = strconv.Atoi(mapInfo["acreage"].(string))
	house_info.Unit = mapInfo["unit"].(string)
	house_info.Capacity, _ = strconv.Atoi(mapInfo["capacity"].(string))
	house_info.Beds = mapInfo["beds"].(string)
	house_info.Deposit, _ = strconv.Atoi(mapInfo["deposit"].(string))
	house_info.Deposit = house_info.Deposit * 100
	house_info.Min_days, _ = strconv.Atoi(mapInfo["min_days"].(string))
	house_info.Max_days, _ = strconv.Atoi(mapInfo["max_days"].(string))
	house_info.User = &user
	//house_info.Facilities = mapInfo["facility"].([]int)
	beego.Info("run UserPublishInfo-----", house_info)
	houseid, err := orm.NewOrm().Insert(&house_info)
	if err != nil {
		beego.Info("insert into house err ", err)
		return
	}

	beego.Info("facility.....", mapInfo["facility"])

	//还要搞定房子和房子内设施的多对多关系
	facilities := []*models.Facility{}
	for _, fid := range mapInfo["facility"].([]interface{}) {
		id, _ := strconv.Atoi(fid.(string))
		facility := &models.Facility{Id: id}
		facilities = append(facilities, facility)
	}
	// 第一个参数的对象，主键必须有值
	// 第二个参数为对象需要操作的 M2M 字段
	// QueryM2Mer 的 api 将作用于 Id 为 1 的 House
	m2mhouse_facility := orm.NewOrm().QueryM2M(&house_info, "Facilities")

	num, err := m2mhouse_facility.Add(facilities)
	if err != nil || num <= 0 {
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//设置成功返回信息
	dataid := make(map[string]interface{})
	dataid["house_id"] = houseid
	resp.Data = dataid
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
}

//UploadHouseImage 上传房屋照片
func (this *HouseController) UploadHouseImage() {
	var resp UtilResp
	resp.Errno = utils.RECODE_OK
	resp.Errmsg = utils.RecodeText(resp.Errno)
	defer this.RetData(&resp)
	houseid := this.Ctx.Input.Param(":id")
	beego.Info("get param houseid ====", houseid)
	//上传图片
	file, header, err := this.GetFile("house_image")
	if err != nil {
		beego.Info("getfile  err...", err)
		resp.Errno = utils.RECODE_THIRDERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//defer file.Close()

	this.SaveToFile("house_image", header.Filename)
	file.Close()
	beego.Info("UploadHouseImage is running...", header.Filename)
	img_url := UploadFile(header.Filename)
	urldata := make(map[string]interface{})
	urldata["url"] = img_url
	resp.Data = urldata
	beego.Info(resp)
	//保存到数据库
	//update house set index_image_url=url where id=houseid and length(index_image_url)=0
	r := orm.NewOrm().Raw("update house set index_image_url=? where id=? and length(index_image_url)=0", img_url, houseid)
	_, err = r.Exec()
	if err != nil {
		beego.Info("update house info err===", err)
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
	//insert into house_image(url,house_id) values('url',15)
	r = orm.NewOrm().Raw("insert into house_image(url,house_id) values(?,?)", img_url, houseid)
	_, err = r.Exec()
	if err != nil {
		beego.Info("update house info err===", err)
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}
}

// /houses/:id [GET]
func (this *HouseController) GetOneHouseInfo() {
	resp := UtilResp{Errno: utils.RECODE_OK, Errmsg: utils.RecodeText(utils.RECODE_OK)}
	defer this.RetData(&resp)
	data := make(map[string]interface{})

	user_id := this.GetSession("user_id")
	beego.Debug("user_id = ", user_id)
	if user_id == nil {
		user_id = -1
	}
	house_id := this.Ctx.Input.Param(":id")

	//先从缓存中获取房屋数据,将缓存数据返回前端即可
	redis_config_map := map[string]string{
		"key":   "ihome_go",
		"conn":  utils.G_redis_addr + ":" + utils.G_redis_port,
		"dbNum": utils.G_redis_dbnum,
	}
	redis_config, _ := json.Marshal(redis_config_map)

	cache_conn, err := cache.NewCache("redis", string(redis_config))
	if err != nil {
		beego.Debug("connect cache error")
	}
	house_info_key := fmt.Sprintf("house_info_%s", house_id)
	house_info_value := cache_conn.Get(house_info_key)
	if house_info_value != nil {
		beego.Debug("======= get house info desc  from CACHE!!! ========")
		data["user_id"] = user_id
		house_info := map[string]interface{}{}
		json.Unmarshal(house_info_value.([]byte), &house_info)
		data["house"] = house_info
		resp.Data = data
		return
	}

	beego.Debug("======= no house info desc CACHE!!!  SAVE house desc to CACHE !========")
	//如果缓存没有房屋数据,那么从数据库中获取数据,再存入缓存中,然后返回给前端
	//1 根据house_id 关联查询数据库
	o := orm.NewOrm()

	//根据house_id查询house
	/*
		// ---- 方法1  关联关系查询 ----
			if err := o.QueryTable("house").Filter("id", house_id).RelatedSel().One(&house); err != nil {
				rep.Errno = utils.RECODE_NODATA
				rep.Errmsg = utils.RecodeText(rep.Errno)
				return
			}
			//查询house_id为 id的house_image
			if _, err := o.QueryTable("house_image").Filter("House", house_id).RelatedSel().All(&house.Images); err != nil {

			}

			//查询house_id为 id的Facility
			if _, err := o.QueryTable("facility").Filter("Houses__House__Id", house_id).All(&house.Facilities); err != nil {

			}
	*/
	// --- 方法2  载入关系查询 -----
	house := models.House{}
	house.Id, _ = strconv.Atoi(house_id)
	o.Read(&house)
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")

	//2 将该房屋的json格式数据保存在redis缓存数据库
	house_info_value, _ = json.Marshal(house.To_one_house_desc())
	cache_conn.Put(house_info_key, house_info_value, 3600*time.Second)

	//3 返回数据
	data["user_id"] = user_id
	data["house"] = house.To_one_house_desc()

	resp.Data = data
	return
}

// /houese?aid=1&sd=2017-11-09&ed=2017-11-11&sk=new&p=1 [GET]
func (this *HouseController) GetHouseInfo() {
	resp := UtilResp{Errno: utils.RECODE_OK, Errmsg: utils.RecodeText(utils.RECODE_OK)}
	defer this.RetData(&resp)

	beego.Debug()
	var aid int
	this.Ctx.Input.Bind(&aid, "aid")
	var sd string
	this.Ctx.Input.Bind(&sd, "sd")
	var ed string
	this.Ctx.Input.Bind(&ed, "ed")
	var sk string
	this.Ctx.Input.Bind(&sk, "sk")
	var page int
	this.Ctx.Input.Bind(&page, "p")

	beego.Debug(aid, sd, ed, sk, page)

	//把时间从str转换成字符串格式

	//校验开始时间一定要早于结束时间

	//判断page的合法性 一定是大于0的整数

	//尝试从redis中获取数据, 有则返回

	//如果没有 从mysql中查询
	houses := []models.House{}

	o := orm.NewOrm()

	qs := o.QueryTable("house")

	num, err := qs.Filter("area_id", aid).All(&houses)
	if err != nil {
		resp.Errno = utils.RECODE_PARAMERR
		resp.Errmsg = utils.RecodeText(resp.Errno)
		return
	}

	total_page := int(num)/models.HOUSE_LIST_PAGE_CAPACITY + 1
	house_page := 1

	house_list := []interface{}{}
	for _, house := range houses {
		o.LoadRelated(&house, "Area")
		o.LoadRelated(&house, "User")
		o.LoadRelated(&house, "Images")
		o.LoadRelated(&house, "Facilities")
		house_list = append(house_list, house.To_house_info())
	}

	data := map[string]interface{}{}
	data["houses"] = house_list
	data["total_page"] = total_page
	data["current_page"] = house_page

	resp.Data = data

	return
}
