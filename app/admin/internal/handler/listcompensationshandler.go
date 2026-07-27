package handler

import (
	"net/http"

	admin "Ticket/app/admin/internal/logic"
	"Ticket/app/admin/internal/svc"
	"Ticket/app/admin/internal/types"
	"Ticket/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListCompensationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCompensationsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := admin.NewListCompensationsLogic(r.Context(), svcCtx)
		resp, err := l.ListCompensations(&req)
		response.HttpResult(w, r, resp, err)
	}
}
