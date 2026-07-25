// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1
package handler

import (
	"net/http"

	"Ticket/app/inventory/inventory-api/internal/logic/inventory"
	"Ticket/app/inventory/inventory-api/internal/svc"
	"Ticket/app/inventory/inventory-api/internal/types"

	"Ticket/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetStockReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := inventory.NewGetStockLogic(r.Context(), svcCtx)
		resp, err := l.GetStock(&req)
		response.HttpResult(w, r, resp, err)
	}
}
