// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"Ticket/app/ticket/ticket-api/internal/logic"
	"Ticket/app/ticket/ticket-api/internal/svc"
	"Ticket/app/ticket/ticket-api/internal/types"
	"Ticket/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTicketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTicketReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := logic.NewGetTicketLogic(r.Context(), svcCtx)
		resp, err := l.GetTicket(&req)
		response.HttpResult(w, r, resp, err)
	}
}
