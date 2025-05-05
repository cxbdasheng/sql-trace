package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	StartTime time.Time
	EndTime   time.Time
	Types     []string
	ThreadID  int
	Argument  string
}

// ParseQueryParams 解析查询参数
func ParseQueryParams(query url.Values) (*Query, error) {
	params := &Query{
		StartTime: time.Now().AddDate(-10, 0, 0),
		EndTime:   time.Now(),
		Types:     []string{},
		ThreadID:  0,
		Argument:  "",
	}

	// 解析开始时间
	if startTimeStr := query.Get("start_time"); startTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
			params.StartTime = t
		}
	}

	// 解析结束时间
	if endTimeStr := query.Get("end_time"); endTimeStr != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
			params.EndTime = t
		}
	}

	// 检查时间范围是否有效
	if params.StartTime.After(params.EndTime) {
		return nil, fmt.Errorf("开始时间不能大于结束时间")
	}

	// 解析命令类型数组
	if types, ok := query["type"]; ok {
		if len(types) == 1 {
			// 如果只有一个值，尝试按逗号分割
			params.Types = strings.Split(types[0], ",")
		} else {
			// 如果有多个值，直接使用
			params.Types = types
		}
	}

	// 解析线程ID
	if threadIDStr := query.Get("id"); threadIDStr != "" {
		if id, err := strconv.Atoi(threadIDStr); err == nil {
			params.ThreadID = id
		}
	}

	if argument := query.Get("argument"); argument != "" {
		params.Argument = argument
	}
	return params, nil
}
