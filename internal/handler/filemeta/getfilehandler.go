package filemeta

import (
	"net/http"

	"corenet/internal/logic/filemeta"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := filemeta.NewGetFileLogic(r.Context(), svcCtx)
		resp, err := l.GetFile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
