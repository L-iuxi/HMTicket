// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"net/http"

	"Ticket/app/order/order-api/internal/logic/order"
	"Ticket/app/order/order-api/internal/svc"
	"Ticket/app/order/order-api/internal/types"

	"Ticket/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayorderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := order.NewPayorderLogic(r.Context(), svcCtx)
		resp, err := l.Payorder(&req)
		response.HttpResult(w, r, resp, err)
	}
}
