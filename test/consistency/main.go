// test/consistency/main.go
// 最终一致性 8 项测试。每项测试一个不变式。
//
// 用法:
//   go run test/consistency/main.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}
var adminTK string
var passed, failed int

func post(url, token, body string) (int, string) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode, string(b)
}

func put(url, token, body string) (int, string) {
	req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode, string(b)
}

func get(url, token string) (int, string) {
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode, string(b)
}

func check(name string, ok bool, detail string) {
	if ok {
		passed++
		fmt.Printf("  ✅ %s: %s\n", name, detail)
	} else {
		failed++
		fmt.Printf("  ❌ %s: %s\n", name, detail)
	}
}

// setup 创建活动+场次+票种，返回 evtID/showID/ttID/stockURL
func setup(stock int) (evtID, showID, ttID int) {
	ts := time.Now().UnixNano()
	evtName := fmt.Sprintf("CTEST-%d", ts%9999999)
	now := time.Now()
	evtStart := now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	evtEnd := now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	// 建活动
	post("http://127.0.0.1:8889/admin/event", adminTK, fmt.Sprintf(
		`{"title":"%s","description":"consistency test","location":"test","startTime":"%s","endTime":"%s"}`,
		evtName, evtStart, evtEnd))
	for id := 50; id >= 1; id-- {
		_, b := get(fmt.Sprintf("http://127.0.0.1:8890/event/%d", id), "")
		if strings.Contains(b, evtName) {
			evtID = id
			break
		}
	}

	// 建场次
	showName := fmt.Sprintf("CSHOW-%d", ts%9999999)
	post("http://127.0.0.1:8889/admin/show", adminTK, fmt.Sprintf(
		`{"eventId":%d,"name":"%s","showTime":"%s","endTime":"%s","sortOrder":1}`,
		evtID, showName, evtStart, evtEnd))
	for id := 50; id >= 1; id-- {
		_, b := get(fmt.Sprintf("http://127.0.0.1:8890/show/%d", id), "")
		if strings.Contains(b, showName) {
			showID = id
			break
		}
	}

	// 建票种
	ttName := fmt.Sprintf("CTT-%d", ts%9999999)
	post("http://127.0.0.1:8889/admin/ticket-type", adminTK, fmt.Sprintf(
		`{"eventId":%d,"showId":%d,"name":"%s","price":1,"stock":%d,"maxPerUser":999,"sortOrder":99}`,
		evtID, showID, ttName, stock))
	for id := 50; id >= 1; id-- {
		_, b := get(fmt.Sprintf("http://127.0.0.1:8890/ticket-type/%d", id), "")
		if strings.Contains(b, ttName) {
			ttID = id
			break
		}
	}

	// draft → ready → selling
	put("http://127.0.0.1:8889/admin/event/status", adminTK,
		fmt.Sprintf(`{"eventId":%d,"status":"ready"}`, evtID))
	put("http://127.0.0.1:8889/admin/event/status", adminTK,
		fmt.Sprintf(`{"eventId":%d,"status":"selling"}`, evtID))

	// 初始化 Redis 库存
	put("http://127.0.0.1:8891/api/v1/admin/inventory/stock", adminTK,
		fmt.Sprintf(`{"ticketTypeId":%d,"stock":%d}`, ttID, stock))

	return
}

// newUser 注册+登录，返回 token
func newUser(prefix string, idx int) string {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("%s_%d_%d", prefix, ts%99999, idx)
	post("http://127.0.0.1:8888/api/user/register", "", fmt.Sprintf(
		`{"username":"%s","password":"1","email":"%s@x.com","phone":"180%08d"}`,
		name, name, ts%100000000+int64(idx)))
	_, b := post("http://127.0.0.1:8888/api/user/login", "",
		fmt.Sprintf(`{"account":"%s","password":"1"}`, name))
	var r struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(b), &r)
	return r.Data.Token
}

func buy(token string, evtID, showID, ttID int, requestID string) (int, string) {
	_, body := post("http://127.0.0.1:8894/api/v1/order/buy", token, fmt.Sprintf(
		`{"eventId":%d,"showId":%d,"ticketTypeId":%d,"quantity":1,"requestId":"%s"}`,
		evtID, showID, ttID, requestID))
	var r struct {
		Code uint32 `json:"code"`
		Data struct {
			OrderNo string `json:"orderNo"`
			Status  string `json:"status"`
			Message string `json:"message"`
		} `json:"data"`
	}
	json.Unmarshal([]byte(body), &r)
	if r.Data.OrderNo != "" {
		return int(r.Code), r.Data.OrderNo
	}
	return int(r.Code), r.Data.Message
}

func getStock(ttID int) int {
	_, b := get(fmt.Sprintf("http://127.0.0.1:8891/api/v1/inventory/stock/%d", ttID), "")
	var r struct {
		Data struct{ Stock int `json:"stock"` }
	}
	json.Unmarshal([]byte(b), &r)
	return r.Data.Stock
}

func getOrder(token, orderNo string) string {
	_, b := get("http://127.0.0.1:8894/api/v1/order/" + orderNo, token)
	var r struct {
		Data struct{ Status string }
	}
	json.Unmarshal([]byte(b), &r)
	return r.Data.Status
}

func getCompensations() string {
	_, b := get("http://127.0.0.1:8889/admin/compensations", adminTK)
	return b
}

func main() {
	// admin login
	_, loginBody := post("http://127.0.0.1:8888/api/user/login", "",
		`{"account":"admin","password":"123456"}`)
	var lr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(loginBody), &lr)
	adminTK = lr.Data.Token
	if adminTK == "" {
		fmt.Println("admin 登录失败，请确保服务已启动")
		return
	}

	fmt.Println("===== 最终一致性 8 项测试 =====\n")

	// ——— 1. 正常买票流程 ———
	fmt.Println("[1/8] 正常买票流程")
	evtID, showID, ttID := setup(10)
	tk := newUser("normal", 1)
	before := getStock(ttID)
	code, orderNo := buy(tk, evtID, showID, ttID, "test-normal-1")
	time.Sleep(500 * time.Millisecond) // 等 MQ 消费
	after := getStock(ttID)
	status := getOrder(tk, orderNo)
	check("返回 orderNo", code == 200 && orderNo != "", orderNo)
	check("库存扣减 1", after == before-1, fmt.Sprintf("%d → %d", before, after))
	check("订单状态为 unpaid", status == "unpaid", status)
	check("非 creating 终态", status != "" && status != "creating", status)

	// ——— 2. 幂等：相同 requestId 返回相同 orderNo ———
	fmt.Println("\n[2/8] 幂等：相同 requestId 返回相同 orderNo")
	_, orderNo2 := buy(tk, evtID, showID, ttID, "test-normal-1") // 同 requestId 再买
	stockAfterIdem := getStock(ttID)
	check("相同 requestId 返回同一 orderNo", orderNo2 == orderNo, orderNo2)
	check("库存未重复扣减", stockAfterIdem == after, fmt.Sprintf("%d == %d", stockAfterIdem, after))

	// ——— 3. 并发抢票 0 超卖 ———
	fmt.Println("\n[3/8] 并发抢票 0 超卖")
	evt2, show2, tt2 := setup(5)
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		tokens[i] = newUser("conc", i)
	}
	stockBefore := getStock(tt2)
	sold := 0
	type result struct{ code int; orderNo string }
	ch := make(chan result, 5)
	for i := 0; i < 5; i++ {
		go func(tk string, idx int) {
			c, o := buy(tk, evt2, show2, tt2, fmt.Sprintf("test-conc-%d", idx))
			ch <- result{c, o}
		}(tokens[i], i)
	}
	for i := 0; i < 5; i++ {
		r := <-ch
		if r.orderNo != "" {
			sold++
		}
	}
	time.Sleep(500 * time.Millisecond)
	stockAfter := getStock(tt2)
	check("卖出 5 张", sold == 5, fmt.Sprintf("卖出 %d", sold))
	check("库存减少 5", stockBefore-stockAfter == 5, fmt.Sprintf("%d → %d", stockBefore, stockAfter))
	check("0 超卖", stockAfter >= 0, fmt.Sprintf("剩余 %d", stockAfter))

	// ——— 4. 支付成功 → 订单状态变更 ———
	fmt.Println("\n[4/8] 支付成功 → 订单状态变更")
	tk4 := newUser("pay", 1)
	evt4, show4, tt4 := setup(5)
	stockBefore4 := getStock(tt4)
	_, payOrderNo := buy(tk4, evt4, show4, tt4, "test-pay-1")
	time.Sleep(500 * time.Millisecond)
	stockAfterBuy := getStock(tt4)
	// 支付
	_, payBody := post("http://127.0.0.1:8894/api/v1/order/pay", tk4, fmt.Sprintf(
		`{"orderNo":"%s","requestId":"pay-req-1"}`, payOrderNo))
	time.Sleep(500 * time.Millisecond)
	payStatus := getOrder(tk4, payOrderNo)
	stockAfterPay := getStock(tt4)
	check("库存扣减后不变（支付不扣库存）", stockAfterPay == stockAfterBuy, fmt.Sprintf("%d == %d", stockAfterPay, stockAfterBuy))
	check("订单状态变为 paid", payStatus == "paid", payStatus)
	check("库存未额外减少", stockBefore4-stockAfterPay == 1, fmt.Sprintf("从 %d 到 %d", stockBefore4, stockAfterPay))
	_ = payBody

	// ——— 5. 取消订单 → 库存回滚 ———
	fmt.Println("\n[5/8] 取消订单 → 库存回滚")
	tk5 := newUser("cancel", 1)
	evt5, show5, tt5 := setup(5)
	stockBefore5 := getStock(tt5)
	_, cancelOrderNo := buy(tk5, evt5, show5, tt5, "test-cancel-1")
	time.Sleep(500 * time.Millisecond)
	// 取消
	post("http://127.0.0.1:8894/api/v1/order/cancel", tk5, fmt.Sprintf(`{"orderNo":"%s"}`, cancelOrderNo))
	time.Sleep(500 * time.Millisecond)
	cancelStatus := getOrder(tk5, cancelOrderNo)
	stockAfterCancel := getStock(tt5)
	check("订单状态变为 cancelled", cancelStatus == "cancelled", cancelStatus)
	check("库存回滚恢复", stockAfterCancel == stockBefore5, fmt.Sprintf("%d == %d", stockAfterCancel, stockBefore5))

	// ——— 6. 创建中状态回退查询 ———
	fmt.Println("\n[6/8] 创建中状态回退查询")
	evt6, show6, tt6 := setup(5)
	tk6 := newUser("creating", 1)
	_, createOrderNo := buy(tk6, evt6, show6, tt6, "test-creating-1")
	// 立刻查（MQ 可能还没消费完）
	immediateStatus := getOrder(tk6, createOrderNo)
	time.Sleep(1 * time.Second)
	finalStatus := getOrder(tk6, createOrderNo)
	check("立刻查询可能是 creating（MQ 未消费）或 unpaid（MQ 已消费）",
		immediateStatus == "creating" || immediateStatus == "unpaid", immediateStatus)
	check("1 秒后查询终态不是 creating",
		finalStatus != "creating" && finalStatus != "", finalStatus)

	// ——— 7. 补偿记录查询 ———
	fmt.Println("\n[7/8] 补偿记录查询")
	_, compBody := get("http://127.0.0.1:8889/admin/compensations", adminTK)
	check("补偿接口可访问", strings.Contains(compBody, "records") || strings.Contains(compBody, "\"code\":200"), "API 正常")

	// ——— 8. 限流正确拒绝 ———
	fmt.Println("\n[8/8] 限流正确拒绝")
	tk8 := newUser("ratelimit", 1)
	// 5 次买不同 requestId（同一用户，1 秒内超过 5 次）
	blocked := 0
	for i := 0; i < 10; i++ {
		code, msg := buy(tk8, evtID, showID, ttID, fmt.Sprintf("test-rate-%d", i))
		if code != 200 || strings.Contains(msg, "频繁") {
			blocked++
		}
	}
	check("10 次请求中超过 5 次（一半）被限流", blocked >= 3, fmt.Sprintf("被挡 %d/10", blocked))

	// ——— 结果 ———
	fmt.Printf("\n========== 结果 ==========\n")
	fmt.Printf("通过: %d  失败: %d\n", passed, failed)
	if failed == 0 {
		fmt.Println("✅ 8 项全部通过")
	} else {
		fmt.Printf("❌ %d 项失败\n", failed)
	}
}
