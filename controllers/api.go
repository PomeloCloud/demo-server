package controllers

import (
	"github.com/astaxie/beego"
	"github.com/PomeloCloud/pcfs/drone/storage"
	"fmt"
	"strconv"
	"github.com/PomeloCloud/pcfs/proto"
	"encoding/base64"
	"github.com/PomeloCloud/pcfs/server"
	pb "github.com/PomeloCloud/pcfs/proto"
	"context"
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

func (c *GetFileController) Get() {
	c.TplName = "index.html"
}

func (this *GetFileController) GetFileList() {
	home := fmt.Sprint("/", strconv.Itoa(int(storage.FS.Network.BFTRaft.Id)))
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
	this.Data["list"] = fileList
	this.ServeJSON()
}

func (this *GetFileController) GetFileBlocksList() {
	file := fmt.Sprint("/", strconv.Itoa(int(storage.FS.Network.BFTRaft.Id)), "/", this.Ctx.Input.Param("filename"))
	stream, err := storage.FS.NewStream(file)
	if err != nil {
		panic(err)
	}
	blockNodeList := [][]FileBlockNode{}
	for _, block := range stream.Meta.Blocks {
		nodes := []FileBlockNode{}
		for _, hostId := range block.Hosts {
			host := storage.FS.Network.BFTRaft.GetHostNTXN(hostId)
			node := FileBlockNode{}
			block, blockError := server.GetPeerRPC(node.Address).GetBlock(context.Background(), &pb.GetBlockRequest{
				Group: server.STASH_GROUP, Index: block.Index, File: stream.Meta.Key,
			})
			blockOnline := true
			if block == nil || blockError != nil {
				blockOnline = true
			}
			status := "Online"
			if !blockOnline {
				status = "Offline"
			}
			node.Address = host.ServerAddr
			node.Status = status
			node.Role = "Replica"
			nodes = append(nodes, node)
		}
		blockNodeList = append(blockNodeList, nodes)
	}
	this.Data["list"] = blockNodeList
	this.ServeJSON()
}