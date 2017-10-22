package main

import (
	_ "github.com/PomeloCloud/demo-server/routers"
	"github.com/astaxie/beego"
	"github.com/PomeloCloud/pcfs/drone/storage"
	"os"
	"github.com/astaxie/beego/plugins/cors"
)

func cleanup() {
	os.RemoveAll("wallet/")
	os.Remove("example.out.jpg")
}

func main() {
	cleanup()
	storage.Main()
	beego.Debug("Filters init...")

	// CORS for https://foo.* origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
	beego.Run()
}

