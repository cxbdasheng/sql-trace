package parser

import (
	"SQLTrace/query"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LogParser 定义日志解析器接口
type LogParser interface {
	Parse(line string) (LogEntry, error)
}

// BaseLogParser 基础日志解析器
type BaseLogParser struct{}

// Parse 实现基础的日志解析功能
func (p *BaseLogParser) Parse(line string) (LogEntry, error) {
	return parseLogLine(line)
}

type LogEntry struct {
	Timestamp time.Time
	ThreadID  int
	Command   string
	Argument  string
}

// ParseLogFile 解析日志文件
func ParseLogFile(path string, ch chan<- LogEntry, query query.Query) error {
	// 打开日志文件
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// 创建基础解析器
	baseParser := &BaseLogParser{}

	// 创建装饰器链
	parser := NewArgumentQueryDecorator(
		NewCommandFuzzyQueryDecorator(
			NewThreadQueryDecorator(
				NewTimeQueryDecorator(
					baseParser,
					&query.StartTime,
					&query.EndTime,
				),
				query.ThreadID,
			),
			query.Types,
		),
		query.Argument,
	)

	// 逐行读取文件
	scanner := bufio.NewScanner(file)
	isHeader := true
	headerLines := 0
	var currentEntry *LogEntry
	var multiLineBuffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// 跳过日志头信息
		if isHeader {
			if strings.Contains(line, "Time") && strings.Contains(line, "Command") {
				headerLines++
				if headerLines >= 2 {
					isHeader = false
				}
				continue
			}
			if strings.Contains(line, "started with:") || strings.Contains(line, "Tcp port:") {
				continue
			}
		}

		// 检查是否是新的日志条目
		if strings.HasPrefix(line, "202") || strings.HasPrefix(line, "20") {
			// 如果有多行缓冲，先处理之前的条目
			if currentEntry != nil {
				if multiLineBuffer.Len() > 0 {
					currentEntry.Argument = multiLineBuffer.String()
				}
				ch <- *currentEntry
				multiLineBuffer.Reset()
			}

			// 解析新条目
			entry, err := parser.Parse(line)
			if err != nil {
				// 如果是过滤条件不匹配，则跳过该日志
				if strings.Contains(err.Error(), "不匹配") ||
					strings.Contains(err.Error(), "不包含") ||
					strings.Contains(err.Error(), "早于") ||
					strings.Contains(err.Error(), "晚于") {
					currentEntry = nil
					continue
				}
				// 如果是解析错误，记录错误并继续
				fmt.Printf("解析日志行出错: %v, 行内容: %s\n", err, line)
				currentEntry = nil
				continue
			}
			currentEntry = &entry
			if entry.Argument != "" {
				multiLineBuffer.WriteString(entry.Argument)
			}
		} else if currentEntry != nil {
			// 继续收集多行内容
			if multiLineBuffer.Len() > 0 {
				multiLineBuffer.WriteString("\n")
			}
			multiLineBuffer.WriteString(line)
		}
	}

	// 处理最后一个条目
	if currentEntry != nil {
		if multiLineBuffer.Len() > 0 {
			currentEntry.Argument = multiLineBuffer.String()
		}
		ch <- *currentEntry
	}

	// 检查扫描过程中是否有错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	close(ch)
	return nil
}

func parseLogLine(line string) (LogEntry, error) {
	var entry LogEntry
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return entry, fmt.Errorf("无效的日志行: %s", line)
	}

	// 尝试解析时间戳（支持两种格式）
	var timestamp time.Time
	var err error
	var commandStartIndex int

	// 首先尝试 RFC3339Nano 格式（MySQL 8.0+）
	timestamp, err = time.Parse(time.RFC3339Nano, fields[0])
	if err != nil {
		// 如果失败，尝试 MySQL 5.7 格式
		if len(fields) >= 2 {
			timestamp, err = time.Parse("2006-01-02 15:04:05", fields[0]+" "+fields[1])
			if err != nil {
				return entry, fmt.Errorf("解析时间戳失败: %v", err)
			}
			commandStartIndex = 3 // 时间戳占两个字段，线程ID在第三个字段
		}
	} else {
		commandStartIndex = 2 // 时间戳占一个字段，线程ID在第二个字段·
	}

	// 处理线程ID
	threadIDStr := strings.TrimRight(fields[commandStartIndex-1], ":")
	entry.ThreadID, err = strconv.Atoi(threadIDStr)
	if err != nil {
		return entry, fmt.Errorf("解析线程 ID 失败: %v", err)
	}

	// 处理命令和参数
	if len(fields) > commandStartIndex {
		entry.Command = fields[commandStartIndex]
		if len(fields) > commandStartIndex+1 {
			entry.Argument = strings.Join(fields[commandStartIndex+1:], " ")
		} else {
			entry.Argument = "" // 确保对于没有参数的命令（如 Quit）也设置空字符串
		}
	} else {
		return entry, fmt.Errorf("无效的命令格式: %s", line)
	}

	entry.Timestamp = timestamp
	return entry, nil
}
