package upload

import (
	"net/http"

	"corenet/internal/logic/upload"
	"corenet/internal/svc"
	"corenet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func InitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InitReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := upload.NewInitLogic(r.Context(), svcCtx)
		resp, err := l.Init(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
