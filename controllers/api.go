package controllers

import (
	"github.com/astaxie/beego"
	"github.com/PomeloCloud/pcfs/drone/storage"
	"github.com/PomeloCloud/pcfs/proto"
	"encoding/base64"
	"github.com/PomeloCloud/pcfs/server"
	pb "github.com/PomeloCloud/pcfs/proto"
	"context"
	"log"
	"fmt"
	"github.com/PomeloCloud/BFTRaft4go/utils"
)


type File struct {
	Filename string `json:"filename"`
	Contract string `json:"contract"`
}
type FileBlockNode struct {
	Role    string `json:"role"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type GetFileController struct {
	beego.Controller
}

func (this *GetFileController) Get() {
	home := storage.FS.Home()
	dir := storage.FS.Ls(home)
	fileList := []File{}
	for _, f := range dir.Items {
		if f.Type == client.DirectoryItem_FILE {
			fileList = append(fileList, File{
				Filename: f.File.Name,
				Contract: base64.StdEncoding.EncodeToString(f.File.Key),
			})
		}
	}
	this.Data["json"] = fileList
	this.ServeJSON()
}

type GetFileBlocksListController struct {
	beego.Controller
}

func (this *GetFileBlocksListController) Get() {
	file := storage.FS.Home() + "/" + this.Ctx.Input.Query("filename")
	stream, err := storage.FS.NewStream(file)
	if err != nil {
		panic(err)
	}
	blockNodeList := [][]FileBlockNode{}
	for _, block := range stream.Meta.Blocks {
		nodes := []FileBlockNode{}
		for i, hostId := range block.Hosts {
			node := FileBlockNode{}
			blockData, blockError := server.GetPeerRPC(node.Address).GetBlock(context.Background(), &pb.GetBlockRequest{
				Group: server.STASH_GROUP, Index: block.Index, File: stream.Meta.Key,
			})
			blockOnline := true
			if blockData == nil || blockError != nil {
				blockOnline = true
			}
			status := "Online"
			if !blockOnline {
				status = "Offline"

			}
			strKey := fmt.Sprint(hostId, "-", block.Index, "-", i, "-", stream.Meta.Key)
			hash, _ := utils.SHA1Hash([]byte(strKey))
			node.Address = base64.StdEncoding.EncodeToString(hash)
			node.Status = status
			if i == 0 {
				node.Role = "Master"
			} else {
				node.Role = "Slave"
			}
			nodes = append(nodes, node)
		}
		blockNodeList = append(blockNodeList, nodes)
	}
	this.Data["json"] = blockNodeList
	this.ServeJSON()
}

type UploadFileController struct {
	beego.Controller
}

func (this *UploadFileController) Post() {
	file, header, err := this.GetFile("file")
	if err != nil {
		panic(err)
	}
	filePath := storage.FS.Home() + "/" + header.Filename
	stream, err := storage.FS.NewStream(filePath)
	if err != nil {
		panic(err)
	}
	bufferSize := 1024
	for true {
		b := make([]byte, bufferSize)
		n, err := file.Read(b)
		if err != nil {
			panic(err)
		}
		wb := b[0:n]
		stream.Write(&wb)
		if n < bufferSize {
			break
		}
	}
	stream.LandWrite()
	log.Println("upload file succeed")
	this.Data["json"] = map[string]string {"success": "true"}
	this.ServeJSON()
}