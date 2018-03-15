package main

import (
	_ "iHome/models"
	_ "iHome/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
