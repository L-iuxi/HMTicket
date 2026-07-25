package response

import (
	"net/http"
	"strings"

	"Ticket/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// Response 统一 HTTP 响应体
type Response struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// HttpResult 统一 HTTP 响应输出
// 自动区分成功 / CodeError / gRPC status / 普通 error
func HttpResult(w http.ResponseWriter, r *http.Request, resp interface{}, err error) {
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, Response{
			Code: xerr.OK,
			Msg:  "success",
			Data: resp,
		})
		return
	}

	// 1. 直接返回的 *CodeError
	var ce *xerr.CodeError
	if xerr.AsCodeError(err, &ce) {
		httpx.WriteJson(w, http.StatusOK, Response{
			Code: ce.GetErrCode(),
			Msg:  ce.GetErrMsg(),
		})
		return
	}

	// 2. gRPC 返回的 status.Error（RPC 调用失败）
	if st, ok := status.FromError(err); ok {
		code := grpcMsgToCode(st.Message())
		httpx.WriteJson(w, http.StatusOK, Response{
			Code: code,
			Msg:  grpcMsgToText(st.Message()),
		})
		return
	}

	// 3. 普通 error
	httpx.WriteJson(w, http.StatusInternalServerError, Response{
		Code: xerr.SERVER_COMMON_ERROR,
		Msg:  err.Error(),
	})
}

// grpcMsgToCode 从 "ErrCode:300001, ErrMsg:活动不存在" 提取 code
func grpcMsgToCode(msg string) uint32 {
	// 格式: "ErrCode:300001, ErrMsg:活动不存在"
	if idx := strings.Index(msg, "ErrCode:"); idx >= 0 {
		start := idx + len("ErrCode:")
		end := strings.Index(msg[start:], ",")
		if end < 0 {
			end = len(msg[start:])
		}
		codeStr := msg[start : start+end]
		var code uint32
		for _, c := range codeStr {
			if c >= '0' && c <= '9' {
				code = code*10 + uint32(c-'0')
			}
		}
		if code > 0 {
			return code
		}
	}
	return xerr.SERVER_COMMON_ERROR
}

// grpcMsgToText 从 "ErrCode:300001, ErrMsg:活动不存在" 提取中文消息
func grpcMsgToText(msg string) string {
	if idx := strings.Index(msg, "ErrMsg:"); idx >= 0 {
		return msg[idx+len("ErrMsg:"):]
	}
	return msg
}
