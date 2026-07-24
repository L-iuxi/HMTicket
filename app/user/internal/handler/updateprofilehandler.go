// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"Ticket/app/user/internal/logic"
	"Ticket/app/user/internal/svc"
	"Ticket/app/user/internal/types"

	"Ticket/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProfileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.HttpResult(w, r, nil, err)
			return
		}

		l := logic.NewUpdateProfileLogic(r.Context(), svcCtx)
		resp, err := l.UpdateProfile(&req)
		response.HttpResult(w, r, resp, err)
	}
}
