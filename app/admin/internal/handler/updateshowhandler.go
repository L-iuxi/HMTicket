// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	admin "Ticket/app/admin/internal/logic"
	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateShowReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewUpdateShowLogic(r.Context(), svcCtx)
		resp, err := l.UpdateShow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
