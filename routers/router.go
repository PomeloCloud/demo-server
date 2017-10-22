package routers

import (
	"github.com/PomeloCloud/demo-server/controllers"
	"github.com/astaxie/beego"
)

func init() {

    beego.Router("/", &controllers.MainController{})
	beego.Router("/storage", &controllers.StorageController{})
	beego.Router("/storage/show", &controllers.StorageShowController{})
	beego.Router("/api/file/list", &controllers.GetFileController{})
	beego.Router("/api/block/list", &controllers.GetFileBlocksListController{})
	beego.Router("/api/file/upload", &controllers.GetFileBlocksListController{})
}
