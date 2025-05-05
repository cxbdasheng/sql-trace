package parser

import (
	"fmt"
	"strings"
)

// CommandFuzzyQueryDecorator 命令模糊查询装饰器
type CommandFuzzyQueryDecorator struct {
	parser          LogParser
	commandKeywords []string
}

// NewCommandFuzzyQueryDecorator 创建命令模糊查询装饰器
func NewCommandFuzzyQueryDecorator(parser LogParser, commandKeywords []string) *CommandFuzzyQueryDecorator {
	return &CommandFuzzyQueryDecorator{
		parser:          parser,
		commandKeywords: commandKeywords,
	}
}

// Parse 实现命令模糊查询装饰器的解析方法
func (d *CommandFuzzyQueryDecorator) Parse(line string) (LogEntry, error) {
	entry, err := d.parser.Parse(line)
	if err != nil {
		return entry, err
	}

	// 如果没有指定关键词，返回所有日志
	if len(d.commandKeywords) == 0 {
		return entry, nil
	}
	// 检查日志的命令是否匹配任何一个关键词
	for _, keyword := range d.commandKeywords {
		if strings.Contains(strings.ToLower(entry.Command), strings.ToLower(keyword)) {
			return entry, nil
		}
	}
	return entry, fmt.Errorf("类型不匹配: 期望 %v, 实际 %s", d.commandKeywords, entry.Command)
}
