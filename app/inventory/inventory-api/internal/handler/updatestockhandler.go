// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	inventory "Ticket/app/inventory/inventory-api/internal/logic/admin"
	"Ticket/app/inventory/inventory-api/internal/svc"
	"Ticket/app/inventory/inventory-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateStockReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := inventory.NewUpdateStockLogic(r.Context(), svcCtx)
		resp, err := l.UpdateStock(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
