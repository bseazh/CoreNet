package handler

import (
	dblayer "FileStoreNetDisk_v3/db"
	"FileStoreNetDisk_v3/util"
	"net/http"
	"strconv"
)

// GetCurPathID : 通过ParentID,用户名,文件名 获取当前目录的ID
func GetCurPathID(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	str_id := r.Form.Get("parentID")
	id, _ := strconv.Atoi(str_id)
	documentName := r.Form.Get("documentName")
	//	token := r.Form.Get("token")

	// // 2. 验证token是否有效
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	// 3. 查询用户信息
	id, err := dblayer.QueryUserID(username, id, documentName)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			ParentID     int
			DocumentName string
		}{
			ParentID:     id,
			DocumentName: documentName,
		},
	}
	w.Write(resp.JSONBytes())
}

// CreateFolder : 通过ParentID,用户名,文件名 创建新的文件夹
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	id, _ := strconv.Atoi(r.Form.Get("parentID"))
	documentName := r.Form.Get("documentName")
	//	token := r.Form.Get("token")

	// 2. 创建目录
	suc := dblayer.CreateDocument(username, documentName, id, 0)
	if !suc {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// OpenFolder : 通过ParentID,用户名,文件名 打开文件夹
func OpenFolder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	parentID, _ := strconv.Atoi(r.Form.Get("parentID"))
	documentName := r.Form.Get("documentName")
	//	token := r.Form.Get("token")

	// 2. 打开文件夹
	data, err := dblayer.OpenFolder(username, documentName, parentID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Write(util.NewRespMsg(0, "OK", data).JSONBytes())
}

// GoUpFolder : 通过ParentID,用户名,文件名 打开文件夹
func GoUpFolder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	parentID, _ := strconv.Atoi(r.Form.Get("parentID"))
	//	token := r.Form.Get("token")

	// 2. 打开文件夹

	data, err := dblayer.GoUpFolder(username, parentID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Write(util.NewRespMsg(0, "OK", data).JSONBytes())
}

// DeleteFolder : 通过ParentID,用户名,文件名 打开文件夹
func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	parentID, _ := strconv.Atoi(r.Form.Get("parentID"))
	documentName := r.Form.Get("documentName")
	documentSize, _ := strconv.Atoi(r.Form.Get("documentSize"))
	//	token := r.Form.Get("token")

	//log.Printf("Handler/document/DeleteFolder: username : %s , parentID : %s , documentName : %s\n",
	//	username,r.Form.Get("parentID"),documentName)

	// 2. 删除文件夹 以及里面的内容;
	suc := dblayer.RemoveDocument(username, documentName, parentID, 0, int64(documentSize))
	if !suc {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// 测试接口
// GetDocumentList : 通过 ParentID,用户名,文件名  获取其文件夹下所有文件( 以及文件夹 )
func GetDocumentList(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	parentID, _ := strconv.Atoi(r.Form.Get("parentID"))

	documents, _ := dblayer.GetDocumentList(username, parentID)

	resp := util.RespMsg{
		Code: 200,
		Msg:  "document List : OK",
		Data: documents,
	}

	w.Write(resp.JSONBytes())
}

//GetDocumentID : 通过 ParentID,用户名,文件名 获取文件夹的ID;
func GetDocumentID(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	parentID, _ := strconv.Atoi(r.Form.Get("parentID"))
	documentName := r.Form.Get("documentName")

	ID, _ := dblayer.GetDocumentID(username, parentID, documentName)

	resp := util.RespMsg{
		Code: 200,
		Msg:  "document List : OK",
		Data: ID,
	}

	w.Write(resp.JSONBytes())
}
