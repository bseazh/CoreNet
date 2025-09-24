package taskcenter

import (
	"net/http"

	"corenet/internal/logic/taskcenter"
	"corenet/internal/svc"
	"corenet/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateOCRHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOCRReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := taskcenter.NewCreateOCRLogic(r.Context(), svcCtx)
		resp, err := l.CreateOCR(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
