package controllers

import (
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
)

func UploadFile(filename string) (linkurl string) {

	beego.Info("run ... uploadfile", filename)
	req := httplib.Post("http://up.imgapi.com/")
	req.Header("Accept-Encoding", "gzip,deflate,sdch")
	req.Header("Host", "up.imgapi.com")

	req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")

	req.Param("Token", "c8e56d278e8bf78f6e203b4619bb153a3f07a98d:lg-I4LFPp9UcAa_BaxV0m0AlqQA=:eyJkZWFkbGluZSI6MTUyMDczMTcxMywiYWN0aW9uIjoiZ2V0IiwidWlkIjoiNjM1NzM2IiwiYWlkIjoiMTQxNjQ2OSIsImZyb20iOiJmaWxlIn0=")
	req.PostFile("file", filename)
	var respmap = make(map[string]interface{})
	err := req.ToJSON(&respmap)
	if err != nil {
		beego.Info("get response err----", err)
		return
	}
	beego.Info(respmap["linkurl"])
	linkurl = respmap["linkurl"].(string)
	os.Remove(filename) //上传完成后删除本地文件
	return
}
