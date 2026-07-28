// test/bench/main.go
// 压测: N 个用户抢 M 张票，验证 0 超卖
//
// 用法:
//   go run test/bench/main.go -c 200 -n 50                              # 直连各服务端口
//   go run test/bench/main.go -c 200 -n 50 -addr 127.0.0.1:8090         # 走 Nginx

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
	raw         = flag.Bool("raw", false, "裸测：去掉 sleep，打满 QPS")
	addr        = flag.String("addr", "", "Nginx 地址，如 127.0.0.1:8090（走 Nginx 统一入口）")
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 500,
		MaxConnsPerHost:     500,
		IdleConnTimeout:     90 * time.Second,
	},
}

func post(u, token, body string) (int, string) {
	req, _ := http.NewRequest("POST", u, strings.NewReader(body))
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

func put(u, token, body string) (int, string) {
	req, _ := http.NewRequest("PUT", u, strings.NewReader(body))
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

// baseURL 根据 -addr 返回服务地址
// goViaNginx: 全部走 http://addr:80 + path
// 直连: 各服务用自己的端口
func baseURL(directPort int, path string) string {
	if *addr != "" {
		return fmt.Sprintf("http://%s%s", *addr, path)
	}
	return fmt.Sprintf("http://127.0.0.1:%d%s", directPort, path)
}

func main() {
	flag.Parse()
	ts := time.Now().UnixNano()

	mode := "直连各服务"
	if *addr != "" {
		mode = fmt.Sprintf("Nginx → %s:80", *addr)
	}
	fmt.Printf("===== %d 人抢 %d 张票 | %s =====\n\n", *concurrency, *tickets, mode)

	// 1. admin login
	fmt.Print("[1/4] admin 登录... ")
	_, loginBody := post(baseURL(8888, "/api/user/login"), "",
		`{"account":"admin","password":"123456"}`)
	var lr struct{ Data struct{ Token string } }
	json.Unmarshal([]byte(loginBody), &lr)
	adminTK := lr.Data.Token
	if adminTK == "" {
		fmt.Println("FAIL")
		return
	}
	fmt.Println("OK")

	// 2. 创建新活动 + 新场次 + 新票种
	fmt.Printf("[2/4] 创建活动/场次/票种 stock=%d... ", *tickets)
	evtName := fmt.Sprintf("BENCH-%d", ts%9999999)
	now := time.Now()
	evtStart := now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05")
	evtEnd := now.Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	post(baseURL(8889, "/admin/event"), adminTK, fmt.Sprintf(
		`{"title":"%s","description":"bench","location":"test","startTime":"%s","endTime":"%s"}`,
		evtName, evtStart, evtEnd))

	var evtID int
	for id := 100; id >= 1; id-- {
		resp, _ := http.Get(baseURL(8890, fmt.Sprintf("/event/%d", id)))
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(b), evtName) {
			evtID = id
			break
		}
	}

	showName := fmt.Sprintf("SHOW-%d", ts%9999999)
	post(baseURL(8889, "/admin/show"), adminTK, fmt.Sprintf(
		`{"eventId":%d,"name":"%s","showTime":"%s","endTime":"%s","sortOrder":1}`,
		evtID, showName, evtStart, evtEnd))

	var showID int
	for id := 100; id >= 1; id-- {
		resp, _ := http.Get(baseURL(8890, fmt.Sprintf("/show/%d", id)))
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(b), showName) {
			showID = id
			break
		}
	}

	ttName := fmt.Sprintf("BENCH-TT-%d", ts%9999999)
	post(baseURL(8889, "/admin/ticket-type"), adminTK, fmt.Sprintf(
		`{"eventId":%d,"showId":%d,"name":"%s","price":1,"stock":%d,"maxPerUser":999,"sortOrder":99}`,
		evtID, showID, ttName, *tickets))

	code, putBody := put(baseURL(8889, "/admin/event/status"), adminTK,
		fmt.Sprintf(`{"eventId":%d,"status":"ready"}`, evtID))
	fmt.Printf("draft→ready: code=%d body=%s\n", code, strings.TrimSpace(putBody))
	code, putBody = put(baseURL(8889, "/admin/event/status"), adminTK,
		fmt.Sprintf(`{"eventId":%d,"status":"selling"}`, evtID))
	fmt.Printf("ready→selling: code=%d body=%s\n", code, strings.TrimSpace(putBody))

	evtResp, _ := http.Get(baseURL(8890, fmt.Sprintf("/event/%d", evtID)))
	evtBody, _ := io.ReadAll(evtResp.Body)
	evtResp.Body.Close()
	fmt.Printf("活动状态确认: %s\n", strings.TrimSpace(string(evtBody)))

	var ttID int
	for id := 100; id >= 1; id-- {
		resp, _ := http.Get(baseURL(8890, fmt.Sprintf("/ticket-type/%d", id)))
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if strings.Contains(string(b), ttName) {
			ttID = id
			break
		}
	}

	_, _ = put(baseURL(8891, "/api/v1/admin/inventory/stock"), adminTK,
		fmt.Sprintf(`{"ticketTypeId":%d,"stock":%d}`, ttID, *tickets))

	resp, _ := http.Get(baseURL(8891, fmt.Sprintf("/api/v1/inventory/stock/%d", ttID)))
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("evt=%d show=%d tt=%d stock=%s\n", evtID, showID, ttID, strings.TrimSpace(string(b)))

	// 3. register users
	fmt.Printf("[3/4] 注册 %d 个用户... ", *concurrency)
	tokens := make([]string, 0, *concurrency)
	for i := 0; i < *concurrency; i++ {
		name := fmt.Sprintf("bench_%d_%d", ts%99999, i)
		post(baseURL(8888, "/api/user/register"), "", fmt.Sprintf(
			`{"username":"%s","password":"1","email":"%s@x.com","phone":"180%08d"}`,
			name, name, ts%100000000+int64(i)))
		_, lb := post(baseURL(8888, "/api/user/login"), "",
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

	// 4. concurrent buy
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
				_, body := post(baseURL(8894, "/api/v1/order/buy"), token, fmt.Sprintf(
					`{"eventId":%d,"showId":%d,"ticketTypeId":%d,"quantity":1,"requestId":"bench-%d-%d-%d"}`,
					evtID, showID, ttID, ts, idx, seq))
				var r struct{ Data struct{ OrderNo string } }
				json.Unmarshal([]byte(body), &r)
				if r.Data.OrderNo != "" {
					n := sold.Add(1)
					if n%10 == 0 || n <= 3 || n >= int64(*tickets-2) {
						fmt.Printf("  #%d 卖出 (用户%d 尝试%d次 %.2fs)\n", n, idx, seq, time.Since(start).Seconds())
					}
					return
				}
				if !*raw {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(tk, i)
	}
	wg.Wait()
	elapsed := time.Since(start).Seconds()

	time.Sleep(3 * time.Second)

	// 5. result
	qps := float64(totalReq.Load()) / elapsed
	successQPS := float64(sold.Load()) / elapsed

	fmt.Printf("\n========== 结果 ==========\n")
	fmt.Printf("库存: %d  用户: %d  耗时: %.1fs  总请求: %d  QPS: %.0f  成功QPS: %.1f\n",
		*tickets, len(tokens), elapsed, totalReq.Load(), qps, successQPS)
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
