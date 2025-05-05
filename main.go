package main

import (
	"SQLTrace/config"
	"SQLTrace/web"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 配置文件路径
var configFilePath = flag.String("c", config.GetConfigFilePathDefault(), "Custom configuration file path")

// 监听地址
var listen = flag.String("l", "", "Listen address")

// 日志文件路径
var traceLogPath = flag.String("t", "", "Trace log file path")

func main() {
	flag.Parse()

	// 设置配置文件路径
	if *configFilePath != "" {
		absPath, err := filepath.Abs(*configFilePath)
		if err != nil {
			log.Fatalf("Failed to get absolute path: %v", err)
		}
		os.Setenv(config.PathENV, absPath)
	}
	// 获取配置
	conf, _ := config.GetConfigCached()
	if *listen != "" {
		conf.Port = *listen
	} else {
		if conf.Port == "" {
			conf.Port = config.GetDefaultPort()
		}
		*listen = conf.Port
	}
	if *traceLogPath != "" {
		conf.TraceLogPath = *traceLogPath
	}
	conf.SaveConfig()
	run()
}

func run() {
	go func() {
		// 启动web服务
		err := runWebServer()
		if err != nil {
			log.Fatalf("Web server error: %v", err)
		}
	}()

	// 运行定时器
	RunTimer(time.Duration(300) * time.Second)
}

func RunTimer(delay time.Duration) {
	for {
		time.Sleep(delay)
	}
}

func runWebServer() error {
	// 设置路由
	http.HandleFunc("/favicon.ico", web.Favicon())
	http.HandleFunc("/", web.HandleIndex())
	http.HandleFunc("/settings", web.SaveSettings())
	// 启动服务器
	l, err := net.Listen("tcp", *listen)
	if err != nil {
		return fmt.Errorf("could not listen: %v", err)
	}

	log.Printf("Server started on %s", *listen)
	return http.Serve(l, nil)
}
