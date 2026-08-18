// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"Ticket/app/event/event-api/internal/logic"
	"Ticket/app/event/event-api/internal/svc"
	"Ticket/app/event/event-api/internal/types"
	"Ticket/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetTicketTypeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTicketTypeListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := logic.NewGetTicketTypeListLogic(r.Context(), svcCtx)
		resp, err := l.GetTicketTypeList(&req)
		response.HttpResult(w, r, resp, err)
	}
}
