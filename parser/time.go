package parser

import (
	"fmt"
	"time"
)

// TimeQueryDecorator 时间查询
type TimeQueryDecorator struct {
	parser    LogParser
	startTime *time.Time
	endTime   *time.Time
}

// NewTimeQueryDecorator 创建时间查询装饰器
func NewTimeQueryDecorator(parser LogParser, startTime, endTime *time.Time) *TimeQueryDecorator {
	return &TimeQueryDecorator{
		parser:    parser,
		startTime: startTime,
		endTime:   endTime,
	}
}

// Parse 根据时间范围过滤日志
func (d *TimeQueryDecorator) Parse(line string) (LogEntry, error) {
	entry, err := d.parser.Parse(line)
	if err != nil {
		return entry, err
	}

	// 如果设置了开始时间，检查日志时间是否在开始时间之后
	if d.startTime != nil && entry.Timestamp.Before(*d.startTime) {
		return entry, fmt.Errorf("日志时间早于查询开始时间")
	}

	// 如果设置了结束时间，检查日志时间是否在结束时间之前
	if d.endTime != nil && entry.Timestamp.After(*d.endTime) {
		return entry, fmt.Errorf("日志时间晚于查询结束时间")
	}

	return entry, nil
}
