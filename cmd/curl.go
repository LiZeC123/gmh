package cmd

import (
	"io"
	"net/http"
	"sync"
	"time"
)

func DoCurl(url string, timeout uint8, retry uint8) (body string, err error) {
	req, _ := http.NewRequest("GET", url, nil)

	// 覆盖默认标识，添加浏览器特征头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	return string(bytes), nil
}

type Task struct {
	Urls        []string
	Concurrency uint16
	Timeout     uint8
	Retry       uint8
	UrlOnly     bool
}

type TaskRst struct {
	Data string
	Err  error
}

// TODO: 支持log能力, 详细信息可输出到log文件
func DoCurlTask(task Task) (out chan TaskRst) {
	out = make(chan TaskRst)
	sem := make(chan struct{}, task.Concurrency)
	var wg sync.WaitGroup

	go func() {
		defer close(out)

		for _, url := range task.Urls {
			sem <- struct{}{}
			wg.Add(1)

			go func(u string) {
				defer wg.Done()
				data, err := DoCurl(url, task.Timeout, task.Retry)

				if task.UrlOnly {
					data = url
				}

				out <- TaskRst{
					Data: data,
					Err:  err,
				}
			}(url)
		}
		wg.Wait()
	}()

	return out

}
