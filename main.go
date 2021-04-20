package main

import (
	// "FileStoreNetDisk_v3/assets"
	cfg "FileStoreNetDisk_v3/config"
	"FileStoreNetDisk_v3/handler"
	"fmt"
	"net/http"
)

func main() {
	// 静态资源处理
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(assets.AssetFS())))
	http.Handle("/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 文件存取接口
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.UploadSucHandler))
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))

	// 秒传接口
	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(
		handler.TryFastUploadHandler))

	http.HandleFunc("/file/downloadurl", handler.HTTPInterceptor(
		handler.DownloadURLHandler))

	// 分块上传接口
	http.HandleFunc("/file/mpupload/init",
		handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart",
		handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete",
		handler.HTTPInterceptor(handler.CompleteUploadHandler))

	// 用户相关接口
	http.HandleFunc("/", handler.SignInHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	// 文档相关接口
	http.HandleFunc("/document/info", handler.HTTPInterceptor(handler.GetCurPathID))
	http.HandleFunc("/document/createFolder", handler.HTTPInterceptor(handler.CreateFolder))
	http.HandleFunc("/document/openFolder", handler.HTTPInterceptor(handler.OpenFolder))
	http.HandleFunc("/document/goUpFolder", handler.HTTPInterceptor(handler.GoUpFolder))
	http.HandleFunc("/document/getDocumentList", handler.HTTPInterceptor(handler.GetDocumentList))
	http.HandleFunc("/document/deleteFolder", handler.HTTPInterceptor(handler.DeleteFolder))
	http.HandleFunc("/document/getDocumentID", handler.HTTPInterceptor(handler.GetDocumentID))

	// API 服务接口
	http.HandleFunc("/api/loginIn", handler.LoginInAPIHandler)
	http.HandleFunc("/api/selectAll", handler.HTTPInterceptor(handler.SelectAllAPIHandler))
	http.HandleFunc("/api/upload", handler.HTTPInterceptor(handler.UploadAPIHandler))
	http.HandleFunc("/api/rename", handler.HTTPInterceptor(handler.RenameAPIHandler))
	http.HandleFunc("/api/delete", handler.HTTPInterceptor(handler.DeleteAPIHandler))
	http.HandleFunc("/api/download", handler.HTTPInterceptor(handler.DownloadAPIHandler))
	http.HandleFunc("/api/createFolder", handler.HTTPInterceptor(handler.CreateFolderAPIHandler))

	fmt.Printf("上传服务启动中，开始监听监听[%s]...\n", cfg.UploadServiceHost)
	// 启动服务并监听端口
	err := http.ListenAndServe(cfg.UploadServiceHost, nil)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
