package main

import (
	_ "ihome_go1/models"
	_ "ihome_go1/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
