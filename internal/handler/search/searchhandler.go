package search

import (
	"net/http"

	"corenet/internal/logic/search"
	"corenet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := search.NewSearchLogic(r.Context(), svcCtx)
		err := l.Search()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
