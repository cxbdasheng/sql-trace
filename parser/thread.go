package parser

import "fmt"

// ThreadQueryDecorator 线程查询装饰器
type ThreadQueryDecorator struct {
	parser   LogParser
	threadID int
}

// NewThreadQueryDecorator 创建线程查询装饰器
func NewThreadQueryDecorator(parser LogParser, threadID int) *ThreadQueryDecorator {
	return &ThreadQueryDecorator{
		parser:   parser,
		threadID: threadID,
	}
}

// Parse 根据线程ID过滤日志
func (d *ThreadQueryDecorator) Parse(line string) (LogEntry, error) {
	entry, err := d.parser.Parse(line)
	if err != nil {
		return entry, err
	}

	// 如果线程ID为0，则不过滤
	if d.threadID == 0 {
		return entry, nil
	}

	// 检查日志的线程ID是否匹配
	if entry.ThreadID != d.threadID {
		return entry, fmt.Errorf("线程ID不匹配: 期望 %d, 实际 %d", d.threadID, entry.ThreadID)
	}

	return entry, nil
}
