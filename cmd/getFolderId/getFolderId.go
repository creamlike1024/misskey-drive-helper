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

type Folder []struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	Name         string    `json:"name"`
	FoldersCount int       `json:"foldersCount"`
	FilesCount   int       `json:"filesCount"`
	ParentID     string    `json:"parentId"`
	Parent       struct {
	} `json:"parent"`
}

// 打印帮助信息
func help() {
	fmt.Println("Usage: getFolderId [-r] folderName [folderName2] ...")
	fmt.Println("Options:")
	fmt.Println("  -r\t\trecursively search for folderName")
}

func main() {
	RECURSIVE := false
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		return
	}
	// 如果第一个参数为 -r ，则递归列出所有文件夹
	if args[0] == "-r" {
		args = args[1:]
		RECURSIVE = true
	}
	// 将所有 folderName 放入数组
	var folderName []string
	folderName = append(folderName, args...)
	var mi = Misskey{
		Site:  os.Getenv("MISSKEY_SITE"),
		Token: os.Getenv("MISSKEY_TOKEN"),
	}
	// var folderId []string
	for _, v := range folderName {
		// 首先搜索根目录下的文件夹
		folders, err := mi.listFolder(nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !RECURSIVE {
			// 非递归模式，只列出根目录下的文件夹
			fmt.Println(folders[v])
		} else {
		// 递归模式，所有名为 v 的文件夹都会被列出，无论它们在哪个目录下
		// 此时已搜索过根目录，根据之前返回的 folder id 继续搜索，遍历整个文件夹树
		// 假如要搜索代目录是 test1 
		
	}
	// // 打印folder id
	// for _, v := range folderId {
	// 	fmt.Println(v)
	// }
}

// 列出文件夹
func (mi *Misskey) listFolder(folderId interface{}) (map[string]string, error) {
	type jsonStruct struct {
		Token    string      `json:"i"`
		Limit    int         `json:"limit"`
		FolderId interface{} `json:"FolderId"`
	}
	var data = jsonStruct{
		Token:    mi.Token,
		Limit:    100, // api 限制一次最多列出100个文件
		FolderId: folderId,
	}
	var dataBytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// 发送请求
	resp, err := http.Post(mi.Site+"/api/drive/folders", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 如果返回的状态码不是200，返回错误
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	// 解析返回的json
	var respData Folder
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	if len(respData) == 0 {
		// 如果 folderId 不正确或者文件夹为空，respData 为空数组，返回空 map
		return make(map[string]string), nil
	}
	folders := make(map[string]string)
	for _, v := range respData {
		folders[v.Name] = v.ID
	}
	return folders, nil
}

// // 获取文件夹 id
// func (mi *Misskey) getFolderId(folderName string, parentId interface{}) ([]string, error) {
// 	type jsonStruct struct {
// 		Token    string      `json:"i"`
// 		Name     string      `json:"name"`
// 		ParentId interface{} `json:"parentId"`
// 	}
// 	var data = jsonStruct{
// 		Token:    mi.Token,
// 		Name:     folderName,
// 		ParentId: parentId,
// 	}
// 	var dataBytes, err = json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 发送请求
// 	resp, err := http.Post(mi.Site+"/api/drive/folders/find", "application/json", bytes.NewBuffer(dataBytes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	// 解析返回的json
// 	var respData Folder
// 	err = json.NewDecoder(resp.Body).Decode(&respData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(respData) == 0 {
// 		return nil, fmt.Errorf("folder not found")
// 	}
// 	// 将所有的folder id放入数组
// 	var folderId []string
// 	for _, v := range respData {
// 		folderId = append(folderId, v.ID)
// 	}
// 	return folderId, nil
// }
