// test/rate/main.go — 三层限流测试
// 用法: go run test/rate/main.go

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

func checkOk(name string, ok bool, detail string) {
	if ok {
		fmt.Printf("  \033[32mPASS\033[0m %s\n", name)
		pass++
	} else {
		fmt.Printf("  \033[31mFAIL\033[0m %s %s\n", name, detail)
		fail++
	}
}

func main() {
	ts := time.Now().UnixNano()

	// register
	name := fmt.Sprintf("tester_%d", ts)
	post("http://127.0.0.1:8888/api/user/register", "",
		fmt.Sprintf(`{"username":"%s","password":"123","email":"%s@t.com","phone":"139%08d"}`, name, name, ts%100000000))
	r := post("http://127.0.0.1:8888/api/user/login", "",
		fmt.Sprintf(`{"account":"%s","password":"123"}`, name))
	var rr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(r), &rr)
	tk := rr.Data.Token
	if tk == "" { fmt.Println("登录失败"); return }
	fmt.Printf("用户: %s\n\n", name)

	// ==== 1. 用户固定窗口 5req/s ====
	fmt.Println("===== 1. 用户限流: 1s 内 10 次, 期望 ≥4 次被拒 =====")
	rejected := 0
	for i := 0; i < 10; i++ {
		body := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
			fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":2,"quantity":1,"requestId":"rl-%d-%d"}`, ts, i))
		if strings.Contains(body, "过于频繁") {
			rejected++
		}
	}
	fmt.Printf("  拒绝: %d\n", rejected)
	checkOk("拒绝 ≥4 次", rejected >= 4, fmt.Sprintf("(got %d)", rejected))

	time.Sleep(1 * time.Second)

	// ==== 2. 令牌桶 ====
	fmt.Println("\n===== 2. 令牌桶: cap=50 rate=10/s =====")
	// drain first
	fmt.Println("  清空令牌桶...")
	for i := 0; i < 55; i++ {
		post("http://127.0.0.1:8894/api/v1/order/buy", tk,
			fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":3,"quantity":1,"requestId":"tb-%d-%d"}`, ts, i))
	}
	r2 := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
		fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":3,"quantity":1,"requestId":"tb-empty-%d"}`, ts))
	checkOk("令牌清空-被拒", strings.Contains(r2, "过于频繁"), fmt.Sprintf("(raw: %.80s)", r2))

	time.Sleep(1 * time.Second)

	// ==== 3. 支付限流 ====
	fmt.Println("\n===== 3. 支付限流: 10s 内付两次, 第二次被拒 =====")
	buyR := post("http://127.0.0.1:8894/api/v1/order/buy", tk,
		fmt.Sprintf(`{"eventId":1,"showId":1,"ticketTypeId":4,"quantity":1,"requestId":"pl-buy-%d"}`, ts))
	var bb struct{ Data struct{ OrderNo string } }
	json.Unmarshal([]byte(buyR), &bb)
	if bb.Data.OrderNo != "" {
		time.Sleep(2 * time.Second)

		p1 := post("http://127.0.0.1:8894/api/v1/order/pay", tk,
			fmt.Sprintf(`{"orderNo":"%s","requestId":"pl-1-%d"}`, bb.Data.OrderNo, ts))
		checkOk("第一次支付通过", !strings.Contains(p1, "过于频繁"), fmt.Sprintf("(raw: %.80s)", p1))

		p2 := post("http://127.0.0.1:8894/api/v1/order/pay", tk,
			fmt.Sprintf(`{"orderNo":"%s","requestId":"pl-2-%d"}`, bb.Data.OrderNo, ts))
		checkOk("第二次支付被限流", strings.Contains(p2, "过于频繁"), fmt.Sprintf("(raw: %.80s)", p2))
	} else {
		fmt.Println("  跳过: 无订单（缺库存）")
	}
	fmt.Printf("\n===== 通过:%d / 失败:%d =====\n", pass, fail)
	if fail == 0 { fmt.Println("\033[32m全部通过\033[0m") }
}
