// test/bench/main.go
// 压测: N 个用户抢 M 张票，验证 0 超卖
//
// 用法:
//   go run test/bench/main.go -c 200 -n 50        # 200 人抢 50 张
//   go run test/bench/main.go -c 100 -n 100        # 100 人抢 100 张
//
// 流程:
//   管理员登录 → 创建新票种 stock=N → N 个用户注册登录 → 并发买 → 验证卖出≤N

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	concurrency = flag.Int("c", 200, "并发用户数")
	tickets     = flag.Int("n", 50, "库存票数")
	timeout     = flag.Duration("t", 30*time.Second, "重试超时")
)

var client = &http.Client{Timeout: 10 * time.Second}

func post(url, token, body string) (int, string) {
	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, ""
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp.StatusCode, string(b)
}

func put(url, token, body string) {
	req, _ := http.NewRequest("PUT", url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func main() {
	flag.Parse()
	ts := time.Now().UnixNano()

	fmt.Printf("===== %d 人抢 %d 张票 =====\n\n", *concurrency, *tickets)

	// 1. admin login
	fmt.Print("[1/4] admin 登录... ")
	_, loginBody := post("http://127.0.0.1:8888/api/user/login", "",
		`{"account":"admin","password":"123456"}`)
	var lr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(loginBody), &lr)
	adminTK := lr.Data.Token
	if adminTK == "" {
		fmt.Println("FAIL")
		return
	}
	fmt.Println("OK")

	// 2. create new ticket type with exact stock
	fmt.Printf("[2/4] 创建票种 stock=%d maxPerUser=999... ", *tickets)
	ttName := fmt.Sprintf("BENCH-%d", ts%9999999)
	_, _ = post("http://127.0.0.1:8889/admin/ticket-type", adminTK,
		fmt.Sprintf(`{"eventId":1,"showId":1,"name":"%s","price":1,"stock":%d,"maxPerUser":999,"sortOrder":99}`, ttName, *tickets))
	put("http://127.0.0.1:8889/admin/event/status", adminTK, `{"eventId":1,"status":"selling"}`)

	// find ttID
	var ttID int
	for id := 1; id <= 30; id++ {
		resp, _ := http.Get(fmt.Sprintf("http://127.0.0.1:8890/ticket-type/%d", id))
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(b), ttName) {
			ttID = id
			break
		}
	}
	put("http://127.0.0.1:8891/api/v1/admin/inventory/stock", adminTK,
		fmt.Sprintf(`{"ticketTypeId":%d,"stock":%d}`, ttID, *tickets))
	fmt.Printf("ID=%d  stock=%d\n", ttID, *tickets)

	// 3. register users
	fmt.Printf("[3/4] 注册 %d 个用户... ", *concurrency)
	tokens := make([]string, 0, *concurrency)
	for i := 0; i < *concurrency; i++ {
		name := fmt.Sprintf("bench_%d_%d", ts%99999, i)
		post("http://127.0.0.1:8888/api/user/register", "", fmt.Sprintf(
			`{"username":"%s","password":"1","email":"%s@x.com","phone":"180%08d"}`,
			name, name, ts%100000000+int64(i)))
		_, lb := post("http://127.0.0.1:8888/api/user/login", "",
			fmt.Sprintf(`{"account":"%s","password":"1"}`, name))
		var r struct{ Data struct{ Token string } }
		json.Unmarshal([]byte(lb), &r)
		if r.Data.Token != "" {
			tokens = append(tokens, r.Data.Token)
		}
		if (i+1)%50 == 0 {
			fmt.Printf("%d ", i+1)
		}
	}
	fmt.Printf("= %d 个 token\n", len(tokens))

	// 4. concurrent buy with retry
	fmt.Printf("[4/4] %d 人开抢！ 超时=%v\n\n", len(tokens), *timeout)
	var sold, totalReq atomic.Int64
	var wg sync.WaitGroup
	start := time.Now()

	for i, tk := range tokens {
		wg.Add(1)
		go func(token string, idx int) {
			defer wg.Done()
			seq := 0
			for time.Since(start) < *timeout {
				seq++
				totalReq.Add(1)
				_, body := post("http://127.0.0.1:8894/api/v1/order/buy", token, fmt.Sprintf(
					`{"eventId":1,"showId":1,"ticketTypeId":%d,"quantity":1,"requestId":"bench-%d-%d-%d"}`,
					ttID, ts, idx, seq))
				var r struct{ Data struct{ OrderNo string } }
				json.Unmarshal([]byte(body), &r)
				if r.Data.OrderNo != "" {
					n := sold.Add(1)
					if n%10 == 0 || n <= 3 || n >= int64(*tickets-2) {
						fmt.Printf("  #%d 卖出 (用户%d 尝试%d次 %.2fs)\n", n, idx, seq, time.Since(start).Seconds())
					}
					return
				}
				// 被拒，等 100ms 再试
				time.Sleep(100 * time.Millisecond)
			}
		}(tk, i)
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()

	time.Sleep(3 * time.Second) // wait for MQ

	// 5. result
	fmt.Printf("\n========== 结果 ==========\n")
	fmt.Printf("库存: %d  用户: %d  耗时: %.1fs  总请求: %d\n", *tickets, len(tokens), elapsed, totalReq.Load())
	fmt.Printf("卖出: %d 张\n", sold.Load())

	switch {
	case int(sold.Load()) == *tickets:
		fmt.Println("\n✅ 精准卖完，0 超卖！")
	case int(sold.Load()) < *tickets:
		fmt.Printf("\n⚠️ 卖出 %d/%d（被防护拦截 %d 人）\n", sold.Load(), *tickets, *tickets-int(sold.Load()))
	default:
		fmt.Printf("\n❌ 超卖 %d 张！\n", int(sold.Load())-*tickets)
	}
}
