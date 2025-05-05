package parser

import (
	"fmt"
	"strings"
)

// ArgumentQueryDecorator 线程查询装饰器
type ArgumentQueryDecorator struct {
	parser   LogParser
	argument string
}

func NewArgumentQueryDecorator(parser LogParser, argument string) *ArgumentQueryDecorator {
	return &ArgumentQueryDecorator{
		parser:   parser,
		argument: argument,
	}
}

func (d *ArgumentQueryDecorator) Parse(line string) (LogEntry, error) {
	entry, err := d.parser.Parse(line)
	if err != nil {
		return entry, err
	}
	if d.argument == "" {
		return entry, nil
	}
	if !strings.Contains(strings.ToLower(entry.Argument), strings.ToLower(d.argument)) {
		return entry, fmt.Errorf("查询不包含: 模糊匹配值 %s, 实际值 %s", d.argument, entry.Argument)
	}
	return entry, nil
}
