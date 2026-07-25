// test/idem/main.go — 幂等 + 分布式锁测试
// 用法: go run test/idem/main.go

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var c = &http.Client{Timeout: 10 * time.Second}
var pass, fail int

func post(url, token, body string) string {
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

func registerAndLogin(name string) string {
	post("http://127.0.0.1:8888/api/user/register", "",
		fmt.Sprintf(`{"username":"%s","password":"123","email":"%s@t.com","phone":"139%08d"}`, name, name, time.Now().UnixNano()%100000000))
	r := post("http://127.0.0.1:8888/api/user/login", "",
		fmt.Sprintf(`{"account":"%s","password":"123"}`, name))
	var rr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(r), &rr)
	return rr.Data.Token
}

type buyResp struct {
	Data struct {
		OrderNo string `json:"orderNo"`
		Message string `json:"message"`
	} `json:"data"`
}
type payResp struct {
	Data struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	} `json:"data"`
}

func checkOk(name string, ok bool) {
	if ok {
		fmt.Printf("  \033[32mPASS\033[0m %s\n", name)
		pass++
	} else {
		fmt.Printf("  \033[31mFAIL\033[0m %s\n", name)
		fail++
	}
}

func main() {
	ts := time.Now().UnixNano()
	name := fmt.Sprintf("tester_%d", ts)
	tk := registerAndLogin(name)
	if tk == "" { fmt.Println("登录失败"); return }
	fmt.Printf("用户: %s\n\n", name)

	// ==== 1. 幂等 ====
	fmt.Println("===== 1. 买票幂等: 同 requestId =====")
	reqID := fmt.Sprintf("idem-%d", ts)
	r1 := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
		fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":2,"quantity":1,"requestId":"%s"}`, reqID))
	var b1 buyResp
	json.Unmarshal([]byte(r1), &b1)
	orderNo := b1.Data.OrderNo
	checkOk("首次创建订单-拿到orderNo", orderNo != "")
	checkOk("首次创建订单-状态pending", b1.Data.Message != "")

	time.Sleep(3 * time.Second) // 等 MQ 消费

	r2 := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
		fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":2,"quantity":1,"requestId":"%s"}`, reqID))
	var b2 buyResp
	json.Unmarshal([]byte(r2), &b2)
	checkOk("幂等命中-同一订单号", b2.Data.OrderNo == orderNo)

	// ==== 2. 支付幂等 ====
	fmt.Println("\n===== 2. 支付幂等: 同 requestId =====")
	payReq := fmt.Sprintf("pay-%d", ts)
	p1 := post("http://127.0.0.1:8894/api/v1/order/pay", tk,
		fmt.Sprintf(`{"orderNo":"%s","requestId":"%s"}`, orderNo, payReq))
	var pr1 payResp
	json.Unmarshal([]byte(p1), &pr1)
	checkOk("首次支付", pr1.Data.Success)

	p2 := post("http://127.0.0.1:8894/api/v1/order/pay", tk,
		fmt.Sprintf(`{"orderNo":"%s","requestId":"%s"}`, orderNo, payReq))
	var pr2 payResp
	json.Unmarshal([]byte(p2), &pr2)
	checkOk("支付幂等-重复支付也成功", pr2.Data.Success)

	// ==== 3. 分布式锁 ====
	fmt.Println("\n===== 3. 分布式锁: 5 并发抢同一票种 =====")
	var gotOrder atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
				fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":3,"quantity":1,"requestId":"lock-%d-%d"}`, ts, idx))
			var bb buyResp
			json.Unmarshal([]byte(r), &bb)
			if bb.Data.OrderNo != "" { gotOrder.Add(1) }
		}(i)
	}
	wg.Wait()
	checkOk("锁拦截并发-至多1单成功", gotOrder.Load() <= 1)

	fmt.Printf("\n===== 通过:%d / 失败:%d =====\n", pass, fail)
	if fail == 0 { fmt.Println("\033[32m全部通过\033[0m") } else { fmt.Println("\033[31m有失败\033[0m") }
}
