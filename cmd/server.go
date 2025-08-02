package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

)



func StartServer(port int16) {
	// 设置服务器监听地址和端口
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("HTTP 报文打印服务器正在监听 %s...\n", addr)

	// 创建 HTTP 服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. 创建缓冲区来存储完整的 HTTP 请求
		var rawRequest bytes.Buffer

		// 2. 打印请求行
		rawRequest.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.URL.RequestURI(), r.Proto))

		// 3. 打印请求头
		for name, values := range r.Header {
			for _, value := range values {
				rawRequest.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
			}
		}

		// 4. 添加空白行分隔头部和正文
		rawRequest.WriteString("\r\n")

		// 5. 读取并打印请求正文
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("读取正文错误: %v", err)
		}
		if len(body) > 0 {
			rawRequest.Write(body)
		}

		// 6. 在控制台打印原始请求
		fmt.Println("===== 接收到的 HTTP 请求 =====")
		fmt.Println(rawRequest.String())
		fmt.Println("============================")

		// 7. 返回确认响应
		w.WriteHeader(http.StatusOK)
		w.Write(rawRequest.Bytes())
	})

	// 启动 HTTP 服务器
	log.Fatal(http.ListenAndServe(addr, nil))
}
