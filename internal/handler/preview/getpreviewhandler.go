package preview

import (
	"net/http"

	"corenet/internal/logic/preview"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := preview.NewGetPreviewLogic(r.Context(), svcCtx)
		resp, err := l.GetPreview()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
