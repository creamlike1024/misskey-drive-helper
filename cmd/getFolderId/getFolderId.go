package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Misskey struct {
	Site  string
	Token string
}

// 获取文件夹 id
func (mi *Misskey) getFolderId(folderName string) ([]string, error) {
	type jsonStruct struct {
		Token    string      `json:"i"`
		Name     string      `json:"name"`
		ParentId interface{} `json:"parentId"`
	}
	var data = jsonStruct{
		Token:    mi.Token,
		Name:     folderName,
		ParentId: nil,
	}
	var dataBytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// 发送请求
	resp, err := http.Post(mi.Site+"/api/drive/folders/find", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析返回的json
	type respStruct []struct {
		ID           string    `json:"id"`
		CreatedAt    time.Time `json:"createdAt"`
		Name         string    `json:"name"`
		FoldersCount int       `json:"foldersCount"`
		FilesCount   int       `json:"filesCount"`
		ParentID     string    `json:"parentId"`
		Parent       struct {
		} `json:"parent"`
	}
	var respData respStruct
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	if len(respData) == 0 {
		return nil, fmt.Errorf("folder not found")
	}
	// 将所有的folder id放入数组
	var folderId []string
	for _, v := range respData {
		folderId = append(folderId, v.ID)
	}
	return folderId, nil

}

// 打印帮助信息
func help() {
	fmt.Println("Usage: getFolderId [folderName]")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		return
	}
	// 将所有 folderName 放入数组
	var folderName []string
	for _, v := range args {
		folderName = append(folderName, v)
	}
	var mi = Misskey{
		Site:  os.Getenv("MISSKEY_SITE"),
		Token: os.Getenv("MISSKEY_TOKEN"),
	}
	var folderId []string
	for _, v := range folderName {
		id, err := mi.getFolderId(v)
		if err != nil {
			panic(err)
		}
		folderId = append(folderId, id...)
	}
	// 打印folder id
	for _, v := range folderId {
		fmt.Println(v)
	}
}
