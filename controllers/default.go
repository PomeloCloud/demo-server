package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "index.html"
}


type StorageController struct {
	beego.Controller
}

func (c *StorageController) Get() {
	c.TplName = "storage.html"
}

type StorageShowController struct {
	beego.Controller
}

func (c *StorageShowController) Get() {
	c.TplName = "show.html"
}