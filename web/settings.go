package web

import (
	"SQLTrace/config"
	"encoding/json"
	"net/http"
	"strings"
)

func SaveSettings() ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := checkAndSave(r)
		byt, _ := json.Marshal(map[string]string{"result": result})
		w.Write(byt)
	}
}

func checkAndSave(request *http.Request) string {
	conf, _ := config.GetConfigCached()
	// 从请求中读取 JSON 数据
	var data struct {
		Port                string `json:"port"`
		TraceLogPath        string `json:"trace_log_path"`
		LogCleaningStrategy int    `json:"log_cleaning_strategy"`
	}
	// 解析请求中的 JSON 数据
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		return "数据解析失败, 请刷新页面重试"
	}
	conf.Port = strings.TrimSpace(data.Port)
	conf.TraceLogPath = strings.TrimSpace(data.TraceLogPath)
	conf.LogCleaningStrategy = data.LogCleaningStrategy
	// 保存到用户目录
	err = conf.SaveConfig()
	// 回写错误信息
	if err != nil {
		return err.Error()
	}
	return "ok"
}
