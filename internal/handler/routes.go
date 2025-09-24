package handler

import (
	"net/http"

	"corenet/internal/handler/filemeta"
	"corenet/internal/handler/preview"
	"corenet/internal/handler/search"
	taskcenter "corenet/internal/handler/taskcenter"
	"corenet/internal/handler/upload"
	"corenet/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/upload/init",
			Handler: upload.InitHandler(serverCtx),
		},
		{
			Method:  http.MethodPut,
			Path:    "/upload/chunk",
			Handler: upload.PutChunkHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/upload/status",
			Handler: upload.StatusHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/upload/complete",
			Handler: upload.CompleteHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/files/:fileId",
			Handler: filemeta.GetFileHandler(serverCtx),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/files/:fileId",
			Handler: filemeta.DeleteFileHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/files/:fileId/preview",
			Handler: preview.GetPreviewHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/search",
			Handler: search.SearchHandler(serverCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/tasks/:jobId",
			Handler: taskcenter.GetJobHandler(serverCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/tasks/ocr",
			Handler: taskcenter.CreateOCRHandler(serverCtx),
		},
	})
}
