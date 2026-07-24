// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"Ticket/app/user/internal/logic"
	"Ticket/app/user/internal/svc"

	"Ticket/common/response"
)

func ProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewProfileLogic(r.Context(), svcCtx)
		resp, err := l.Profile()
		response.HttpResult(w, r, resp, err)
	}
}
