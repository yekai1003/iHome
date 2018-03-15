package routers

import (
	"ihome_go1/controllers"
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func ignoreStaticPath() {

	//透明static

	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}
func init() {
	ignoreStaticPath()
	beego.Router("/", &controllers.MainController{})
	//添加各种路由的方法--根据业务来决定
	//:8080/api/v1.0/areas - GET 获得区域信息
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	// :8080/api/v1.0/session GET/delete 用户登出
	beego.Router("/api/v1.0/session", &controllers.SessionController{}, "get:GetName;delete:DelSess")
	//:8080/api/v1.0/users post 用户注册
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:UserReg")
	//api/v1.0/sessions post - login 登陆
	beego.Router("/api/v1.0/sessions", &controllers.LoginController{}, "post:Login")
	// api/v1.0/houses/index get 获得房源
	beego.Router("/api/v1.0/houses/index", &controllers.HouseController{}, "get:GetHouseIndex")
	// 8080/api/v1.0/user get 查询用户信息
	beego.Router("/api/v1.0/user", &controllers.UserController{}, "get:GetUserInfo")
	//更新用户名 api/v1.0/user/name put 更新用户名
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:SetUserName")
	//8080/api/v1.0/user/auth get 获得用户信息
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{}, "get:GetUserInfo")
	//8080/api/v1.0/user/auth post 更新用户验证信息
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{}, "post:SetUserAuth")
	//用户头像上传
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")
	//查看当前用户房源信息 //api/v1.0/user/houses get
	beego.Router("/api/v1.0/user/houses", &controllers.HouseController{}, "get:GetUserHouses;post:UserPublishInfo")
	//发布房源
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "post:UserPublishInfo;get:GetHouseInfo")
	//http://localhost:8080/api/v1.0/houses/12/images 上传房屋图片
	beego.Router("/api/v1.0/houses/:id/images", &controllers.HouseController{}, "post:UploadHouseImage")
	//请求某房源全部信息
	beego.Router("/api/v1.0/houses/:id", &controllers.HouseController{}, "get:GetOneHouseInfo")
	//用户请求房源首页列表信息
	beego.Router("/api/v1.0/houses/index", &controllers.HouseController{}, "get:GetHouseIndex")
	//用户订单查询
	beego.Router("/api/v1.0/user/orders", &controllers.UserController{}, "get:GetOrders")
	//发布订单
	beego.Router("/api/v1.0/orders", &controllers.OrderController{}, "post:PublishOrders")
	//房东用户接受/拒绝 订单请求
	beego.Router("/api/v1.0/orders/:id/status", &controllers.OrderController{}, "put:OrderStatus")
	//用户发送订单评价信息
	beego.Router("/api/v1.0/orders/:id/comment", &controllers.OrderController{}, "put:OrderComment")
}
