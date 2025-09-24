package taskcenter

import (
	"net/http"

	"corenet/internal/logic/taskcenter"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetJobHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := taskcenter.NewGetJobLogic(r.Context(), svcCtx)
		resp, err := l.GetJob()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
