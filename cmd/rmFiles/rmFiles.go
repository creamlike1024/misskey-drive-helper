package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
)

// var ForceMode bool

type Misskey struct {
	Site  string
	Token string
}

// 检查文件是否已经被附加到 note 上
// func (mi *Misskey) checkIfAttachedToNotes(fileId string) (bool, error) {
// 	type jsonStruct struct {
// 		Token  string `json:"i"`
// 		FileId string `json:"fileId"`
// 	}
// 	var data = jsonStruct{
// 		Token:  mi.Token,
// 		FileId: fileId,
// 	}
// 	var dataBytes, err = json.Marshal(data)
// 	if err != nil {
// 		return false, err
// 	}
// 	resp, err := http.Post(mi.Site+"/api/drive/files/attached-notes", "application/json", bytes.NewBuffer(dataBytes))
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()
// 	// // debug 显示 resp body
// 	// body, err := io.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	return false, err
// 	// }
// 	// println(string(body))
// 	if resp.StatusCode != 200 {
// 		return false, errors.New("failed to get attached notes")
// 	}
// 	type respStruct []struct {
// 		ID string `json:"id"`
// 	}
// 	var respData respStruct
// 	err = json.NewDecoder(resp.Body).Decode(&respData)
// 	if err != nil {
// 		return false, err
// 	}
// 	return len(respData) != 0, nil
// }

// 获取文件夹下的文件
func (mi *Misskey) getFiles(folderId string) ([]string, error) {
	type jsonStruct struct {
		Token    string `json:"i"`
		Limit    int    `json:"limit"`
		FolderId string `json:"folderId"`
	}
	var data = jsonStruct{
		Token:    mi.Token,
		Limit:    100, // 一次最多获取 100 个文件
		FolderId: folderId,
	}
	var dataBytes, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}
	// 发送请求
	resp, err := http.Post(mi.Site+"/api/drive/files", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析返回的json
	type respStruct []struct {
		ID          string    `json:"id"`
		CreatedAt   time.Time `json:"createdAt"`
		Name        string    `json:"name"`
		Type        string    `json:"type"`
		IsSensitive bool      `json:"isSensitive"`
	}
	var respData respStruct
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}
	if len(respData) == 0 {
		return nil, nil
	}
	// 获取文件 id
	var files []string
	for _, v := range respData {
		files = append(files, v.ID)
	}
	return files, nil
}

// 删除单个文件
func (mi *Misskey) deleteFile(fileId string) error {
	type jsonStruct struct {
		Token  string `json:"i"`
		FileId string `json:"fileId"`
	}
	var data = jsonStruct{
		Token:  mi.Token,
		FileId: fileId,
	}
	var dataBytes, err = json.Marshal(data)
	if err != nil {
		return err
	}
	// 发送请求
	resp, err := http.Post(mi.Site+"/api/drive/files/delete", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 如果返回的状态码不是 204，返回错误
	if resp.StatusCode != 204 {
		return errors.New("failed to delete file")
	}
	return nil
}

// 显示帮助信息
func help() {
	// println("Usage: rmFiles [-f] < folderId.txt")
	// println("       -f: force mode, Delete files attached to posts")
	// println("       -h: help")
	// println("       folderId.txt: a file containing folder ids, one id per line")
	println("Usage: rmFiles < folderId.txt")
	println("       -h: help")
	println("       folderId.txt: a file containing folder ids, one id per line")
}

func main() {
	// args := os.Args[1:]
	// ForceMode = false
	// if len(args) > 1 {
	// 	help()
	// 	return
	// }
	// if len(args) == 1 {
	// 	switch args[0] {
	// 	case "-f":
	// 		ForceMode = true
	// 	case "-h":
	// 		help()
	// 		return
	// 	}
	// }
	args := os.Args[1:]
	if len(args) > 1 || (len(args) == 1 && args[0] != "-h") {
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
	var mi = Misskey{
		Site:  os.Getenv("MISSKEY_SITE"),
		Token: os.Getenv("MISSKEY_TOKEN"),
	}
	// 获取文件夹下的所有文件
	for _, folderId := range folderIds {
		println("folder " + folderId + " start emptying")
		files, err := mi.getFiles(folderId)
		if err != nil {
			panic(err)
		}
		for _, fileId := range files {
			// // 检查是否有附加到笔记上
			// attached, err := mi.checkIfAttachedToNotes(fileId)
			// if err != nil {
			// 	panic(err)
			// }
			// if attached {
			// 	if !ForceMode {
			// 		println("file " + fileId + " is attached to notes, use -f to force delete")
			// 		continue
			// 	}
			// }
			// 删除文件
			err = mi.deleteFile(fileId)
			if err != nil {
				// 只显示错误，不退出
				println("file " + fileId + " delete failed: " + err.Error())
			}
			println("file " + fileId + " deleted")
		}
	}
}
