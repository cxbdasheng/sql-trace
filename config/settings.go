package config

type Settings struct {
	// 启动端口
	Port string
	// 日志路径
	TraceLogPath string
	// 日志清理策略
	LogCleaningStrategy int
}
