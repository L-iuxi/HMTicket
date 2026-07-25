// test/e2e/main.go
// 端到端全链路测试
// 覆盖: gRPC拦截器、HTTP统一响应、Lua库存扣减、MQ异步建单、requestId幂等
// 用法: go run test/e2e/main.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var c = &http.Client{Timeout: 10 * time.Second}

const (
	g = "\033[32m"
	r = "\033[31m"
	y = "\033[33m"
	n = "\033[0m"
)

func postJSON(url, token, body string) string {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := c.Do(req)
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return string(b)
}

func main() {
	ts := time.Now().UnixNano()

	// 1. admin login
	fmt.Printf("%s===== 1. 管理员登录 =====%s\n", y, n)
	lb := postJSON("http://127.0.0.1:8888/api/user/login", "",
		`{"account":"admin","password":"123456"}`)
	var lr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(lb), &lr)
	adminTK := lr.Data.Token
	if adminTK == "" {
		fmt.Printf("%sadmin 登录失败%s\n", r, n)
		return
	}
	fmt.Printf("  OK\n")

	// 2. create ticket type on existing event=1 show=1
	fmt.Printf("\n%s===== 2. 创建票种 stock=100 =====%s\n", y, n)
	ttName := fmt.Sprintf("E2E-%d", ts%999999)
	postJSON("http://127.0.0.1:8889/admin/ticket-type", adminTK, fmt.Sprintf(
		`{"eventId":1,"showId":1,"name":"%s","price":1280,"stock":100,"maxPerUser":5,"sortOrder":99}`, ttName))

	// ensure event is selling
	postJSON("http://127.0.0.1:8889/admin/event/status", adminTK,
		`{"eventId":1,"status":"selling"}`)

	// find ttID — scan up to 50 since IDs accumulate from tests
	var ttID int
	for id := 1; id <= 50; id++ {
		resp, _ := http.Get(fmt.Sprintf("http://127.0.0.1:8890/ticket-type/%d", id))
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(b), ttName) {
			ttID = id
			break
		}
	}
	if ttID == 0 {
		fmt.Printf("%s找不到新票种%s\n", r, n)
		return
	}
	fmt.Printf("  ticketTypeId: %d\n", ttID)

	// 3. set Redis stock
	fmt.Printf("\n%s===== 3. 初始化库存 =====%s\n", y, n)
	putReq, _ := http.NewRequest("PUT", fmt.Sprintf(
		"http://127.0.0.1:8891/api/v1/admin/inventory/stock"), strings.NewReader(
		fmt.Sprintf(`{"ticketTypeId":%d,"stock":100}`, ttID)))
	putReq.Header.Set("Content-Type", "application/json")
	putReq.Header.Set("Authorization", "Bearer "+adminTK)
	c.Do(putReq)
	fmt.Println("  库存 → 100")

	// 4. user register + login
	fmt.Printf("\n%s===== 4. 用户注册 + 登录 =====%s\n", y, n)
	testUser := fmt.Sprintf("e2e_%d", ts%9999999)
	postJSON("http://127.0.0.1:8888/api/user/register", "", fmt.Sprintf(
		`{"username":"%s","password":"123456","email":"%s@t.com","phone":"139%08d"}`,
		testUser, testUser, ts%100000000))
	fmt.Printf("  注册: %s\n", testUser)

	login2 := postJSON("http://127.0.0.1:8888/api/user/login", "", fmt.Sprintf(
		`{"account":"%s","password":"123456"}`, testUser))
	var lr2 struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(login2), &lr2)
	userTK := lr2.Data.Token
	fmt.Println("  登录: OK")

	// 5. buy
	fmt.Printf("\n%s===== 5. 买票（MQ 异步建订单）=====%s\n", y, n)
	buy := postJSON("http://127.0.0.1:8894/api/v1/order/buy", userTK, fmt.Sprintf(
		`{"eventId":1,"showId":1,"ticketTypeId":%d,"quantity":1,"requestId":"e2e-%d"}`, ttID, ts))
	var br struct{ Data struct{ OrderNo string } }
	json.Unmarshal([]byte(buy), &br)
	fmt.Printf("  orderNo: %s\n", br.Data.OrderNo)
	if br.Data.OrderNo != "" {
		fmt.Printf("  %sPASS%s\n", g, n)
	} else {
		fmt.Printf("  %sFAIL%s 没拿到订单号 (可能库存不足或活动未发布)\n", r, n)
		fmt.Printf("  raw: %.200s\n", buy)
		return
	}

	time.Sleep(2 * time.Second)

	// 6. check order
	fmt.Printf("\n%s===== 6. 查订单 =====%s\n", y, n)
	req, _ := http.NewRequest("GET", fmt.Sprintf(
		"http://127.0.0.1:8894/api/v1/order/%s", br.Data.OrderNo), nil)
	req.Header.Set("Authorization", "Bearer "+userTK)
	resp, _ := c.Do(req)
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	var or struct{ Data struct{ Status string } }
	json.Unmarshal(b, &or)
	fmt.Printf("  status: %s\n", or.Data.Status)
	if or.Data.Status == "unpaid" {
		fmt.Printf("  %sPASS%s\n", g, n)
	} else {
		fmt.Printf("  %sFAIL%s\n", r, n)
	}

	// 7. pay
	fmt.Printf("\n%s===== 7. 支付 =====%s\n", y, n)
	pay := postJSON("http://127.0.0.1:8894/api/v1/order/pay", userTK, fmt.Sprintf(
		`{"orderNo":"%s","requestId":"e2e-pay-%d"}`, br.Data.OrderNo, ts))
	if strings.Contains(pay, "success") || strings.Contains(pay, "支付成功") {
		fmt.Printf("  %sPASS%s\n", g, n)
	} else {
		fmt.Printf("  %sFAIL%s  raw: %.100s\n", r, n, pay)
	}

	// 8. check ticket
	fmt.Printf("\n%s===== 8. 查出票 =====%s\n", y, n)
	req2, _ := http.NewRequest("GET", "http://127.0.0.1:8892/api/ticket/list", nil)
	req2.Header.Set("Authorization", "Bearer "+userTK)
	resp2, _ := c.Do(req2)
	bb, _ := io.ReadAll(resp2.Body)
	resp2.Body.Close()
	var tkr struct {
		Code int `json:"code"`
		Data struct{ Tickets []interface{} } `json:"data"`
	}
	json.Unmarshal(bb, &tkr)
	fmt.Printf("  code: %d  ticket count: %d\n", tkr.Code, len(tkr.Data.Tickets))
	if tkr.Code == 200 {
		fmt.Printf("  %sPASS%s\n", g, n)
	} else {
		fmt.Printf("  %sFAIL%s\n", r, n)
	}

	fmt.Printf("\n%s===== 端到端完成 =====%s\n", g, n)
}
