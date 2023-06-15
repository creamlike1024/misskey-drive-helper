package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

type Misskey struct {
	Site  string
	Token string
}

// 删除文件夹
func (mi *Misskey) deleteFolder(folderId string) error {
	type jsonStruct struct {
		Token    string `json:"i"`
		FolderId string `json:"folderId"`
	}
	var data = jsonStruct{
		Token:    mi.Token,
		FolderId: folderId,
	}
	var dataBytes, err = json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(mi.Site+"/api/drive/folders/delete", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		type errorStruct struct {
			Error struct {
				Message string `json:"message"`
				Code    string `json:"code"`
				ID      string `json:"id"`
			} `json:"error"`
		}
		// 解析返回的json
		var respData errorStruct
		err = json.NewDecoder(resp.Body).Decode(&respData)
		if err != nil {
			return err
		}
		return errors.New(respData.Error.Message)
	}
	return nil
}

// 显示帮助信息
func help() {
	println("Usage: rmFolder < folderId.txt")
	println("       -h: help")
	println("       folderId.txt: a file containing folder ids, one id per line")
	println("Folder must be empty!")
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		help()
		return
	}
	// 从标准输入读取文件夹 id
	reader := bufio.NewReader(os.Stdin)
	var folderIds []string
	// 读取所有行，直到 EOF
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// 使用 strings.TrimSpace() 去除行尾的换行符
		id := line[:len(line)-1]
		id = strings.TrimSpace(id)
		folderIds = append(folderIds, id)
	}
	var mi Misskey
	mi.Site = os.Getenv("MISSKEY_SITE")
	mi.Token = os.Getenv("MISSKEY_TOKEN")
	for _, folderId := range folderIds {
		err := mi.deleteFolder(folderId)
		if err != nil {
			println(err.Error())
		} else {
			println("Folder " + folderId + " deleted")
		}
	}
}
