package main

import (
	_ "github.com/PomeloCloud/demo-server/routers"
	"github.com/astaxie/beego"
	"github.com/PomeloCloud/pcfs/drone/storage"
	"os"
)

func cleanup() {
	os.RemoveAll("wallet/")
	os.Remove("example.out.jpg")
}

func main() {
	cleanup()
	storage.Main()
	beego.Run()
}

