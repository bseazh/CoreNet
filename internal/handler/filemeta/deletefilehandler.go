package filemeta

import (
	"net/http"

	"corenet/internal/logic/filemeta"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := filemeta.NewDeleteFileLogic(r.Context(), svcCtx)
		err := l.DeleteFile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
