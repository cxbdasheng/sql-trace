package web

import (
	"SQLTrace/config"
	"SQLTrace/parser"
	"SQLTrace/query"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"sync"
)

var logEntries []parser.LogEntry

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

func Favicon() ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打开 favicon 文件
		file, err := os.Open("favicon.ico")
		if err != nil {
			http.Error(w, "Favicon not found", http.StatusNotFound)
			return
		}
		defer file.Close()
		// 设置响应头 Content-Type
		w.Header().Set("Content-Type", "image/x-icon")

		// 将文件内容写入响应
		http.ServeFile(w, r, "favicon.ico")

	}
}

func HandleIndex() ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			// 解析查询参数
			params, err := query.ParseQueryParams(r.URL.Query())
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// 解析模板文件
			tmpl, err := template.ParseFiles("web/index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// 初始化日志条目
			logEntries = []parser.LogEntry{}
			conf, err := config.GetConfigCached()
			if err != nil {
				return
			}

			path := conf.TraceLogPath
			// 创建上下文和取消函数
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// 创建带缓冲的 channel
			ch := make(chan parser.LogEntry, 1000)

			// 创建 WaitGroup 用于等待所有 goroutine 完成
			var wg sync.WaitGroup
			wg.Add(1)

			// 启动解析日志的 goroutine
			go func() {
				defer wg.Done()
				err := parser.ParseLogFile(path, ch, *params)
				if err != nil {
					// 如果发生错误，取消上下文
					cancel()
					fmt.Println(err)
					return
				}
			}()

			// 在主 goroutine 中处理 channel
			go func() {
				for entry := range ch {
					select {
					case <-ctx.Done():
						// 如果上下文被取消，退出循环
						return
					default:
						logEntries = append(logEntries, entry)
					}
				}
			}()

			// 等待所有 goroutine 完成
			wg.Wait()

			// 使用 sort.Slice 实现倒序
			sort.Slice(logEntries, func(i, j int) bool {
				return logEntries[i].Timestamp.After(logEntries[j].Timestamp)
			})

			// 渲染模板并写入响应
			err = tmpl.Execute(w, map[string]interface{}{
				"logs":   logEntries,
				"Title":  template.HTML(`<i> Go Template Example </i>`),
				"params": params,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}
}
