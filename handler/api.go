package handler

import (
	cmn "FileStoreNetDisk_v3/common"
	cfg "FileStoreNetDisk_v3/config"
	dblayer "FileStoreNetDisk_v3/db"
	"FileStoreNetDisk_v3/meta"
	"FileStoreNetDisk_v3/mq"
	"FileStoreNetDisk_v3/store/minIO"
	"FileStoreNetDisk_v3/store/oss"
	"FileStoreNetDisk_v3/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

// LoginInAPIHandler : 登录-接口 ( API )
func LoginInAPIHandler(w http.ResponseWriter, r *http.Request) {

	// 初始化 返回结果;
	resp := util.RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "Error",
	}

	// 1. 解析账号和密码;
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	// 2. 校验用户名及密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write(resp.JSONBytes())
		return
	}

	// 3. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write(resp.JSONBytes())
		return
	}

	// 4. 登录成功后返回对应的token;
	resp = util.RespMsg{
		Code: http.StatusOK,
		Msg:  "OK",
		Data: struct {
			Username string
			Token    string
		}{
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// SelectAllAPIHandler : 查询-接口 ( API )
func SelectAllAPIHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// 1. 解析账号和密码;
		r.ParseForm()
		username := r.Form.Get("username")

		// 2. 根据用户名获取其相应的文件信息
		documents, err := dblayer.GetDirectory(username)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(documents)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	} else {

		// 解析账号和token;
		r.ParseForm()
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		url := "/static/view/api-selectAll.html" + "?username=" + username + "&token=" + token

		http.Redirect(w, r, url, http.StatusFound)
		return
	}

}

// UploadAPIHandler : 上传-接口 ( API )
func UploadAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "View internel server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == http.MethodPost {

		r.ParseForm()
		op := r.Form.Get("op")
		async := r.Form.Get("async")

		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/NetDisk/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s \n", err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err:%s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)

		// 游标重新回到文件头部
		newFile.Seek(0, 0)

		if op == "1" {
			// 文件写入OSS存储
			MinIOPath := "/minIO/" + fileMeta.FileSha1

			if async == "0" {

				if err := minIO.PutObject(fileMeta.FileSha1+path.Ext(fileMeta.FileName), newFile); err != nil {
					log.Println(err.Error())
					return
				}

				log.Println(" success write in MinIO ")
				fileMeta.Location = MinIOPath

			} else {
				// 写入异步转移任务队列
				data := mq.TransferData{
					FileHash:      fileMeta.FileSha1,
					CurLocation:   fileMeta.Location,
					DestLocation:  MinIOPath,
					DestStoreType: cmn.StoreMinIO,
				}
				pubData, _ := json.Marshal(data)
				mq.Publish(
					cfg.TransExchangeName,
					cfg.TransMinIORoutingKey,
					pubData,
				)
			}
		} else if op == "2" {
			// 文件写入OSS存储
			ossPath := "oss/" + fileMeta.FileSha1
			// 判断写入OSS为同步还是异步
			if async == "0" {
				err = oss.Bucket().PutObject(ossPath, newFile)
				if err != nil {
					fmt.Println(err.Error())
					w.Write([]byte("Upload failed!"))
					return
				}
				fileMeta.Location = ossPath
			} else {
				// 写入异步转移任务队列
				data := mq.TransferData{
					FileHash:      fileMeta.FileSha1,
					CurLocation:   fileMeta.Location,
					DestLocation:  ossPath,
					DestStoreType: cmn.StoreOSS,
				}
				pubData, _ := json.Marshal(data)
				mq.Publish(
					cfg.TransExchangeName,
					cfg.TransOSSRoutingKey,
					pubData,
				)
			}
		}

		// 更新用户文件表记录
		_ = meta.UpdateFileMetaDB(fileMeta)

		username := r.Form.Get("username")
		pid := r.Form.Get("pid")

		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize, pid)

		if suc {
			w.Write([]byte("Upload Success."))
		} else {
			w.Write([]byte("Upload Failed."))
		}
	}
}

// RenameAPIHandler : 修改-接口 ( API )
func RenameAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析 querystring - 账号 id 新名称;
	r.ParseForm()
	username := r.Form.Get("username")
	id := r.Form.Get("id")
	newname := r.Form.Get("newname")

	err := dblayer.RenameFileAPI(username, id, newname)
	resp := util.RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "Failed",
	}
	if err != nil {
		w.Write(resp.JSONBytes())
	}
	resp.Code = http.StatusOK
	resp.Msg = "OK"
	w.Write(resp.JSONBytes())

}

// DeleteAPIHandler : 删除-接口 ( API )
func DeleteAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析 querystring - 账号 id 新名称;
	r.ParseForm()
	username := r.Form.Get("username")
	id, _ := strconv.Atoi(r.Form.Get("id"))

	suc := dblayer.RemoveDocumentAPI(username, id)
	resp := util.RespMsg{
		Code: http.StatusBadRequest,
		Msg:  "Failed",
	}
	if !suc {
		w.Write(resp.JSONBytes())
	} else {
		resp.Code = http.StatusOK
		resp.Msg = "OK"
		w.Write(resp.JSONBytes())
	}
}

// CreateFolderAPIHandler : 通过ParentID,用户名,文件名 创建新的文件夹
func CreateFolderAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	id, _ := strconv.Atoi(r.Form.Get("pid"))
	documentName := r.Form.Get("documentName")

	// 2. 创建目录
	suc := dblayer.CreateDocument(username, documentName, id, 0)
	if !suc {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}

// DownloadAPIHandler : 下载-接口 ( API )
func DownloadAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	id, _ := strconv.Atoi(r.Form.Get("id"))
	token := r.Form.Get("token")

	filesha1, err := dblayer.GetFileSha1(username, id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	} else {
		url := "/file/downloadurl" + "?username=" + username + "&token=" + token + "&filehash=" + filesha1
		http.Redirect(w, r, url, http.StatusFound)
		return
	}

}
