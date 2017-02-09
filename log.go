package main

import (
	"cloud/logger"
	"fmt"
	"runtime"
	"time"
)

type Config struct {
	Mail    //[]string
	Monitor []Server
	Setting
}

type Setting struct {
	MonitorOnline  int `yaml:"monitor_time_online"`
	MonitorOutline int `yaml:"monitor_time_outline"`
}

type Mail struct {
	Address        []string `yaml:"address"`
	Header         string   `yaml:"header"`
	MessageOnline  string   `yaml:"message_online"`
	MessageOutline string   `yaml:"message_outline"`
}

type Server struct {
	Ip     string `yaml:"ip"`
	Port   string `yaml:"port"`
	Name   string `yaml:"name"`
	Status bool   `yaml:"status"`
}

type Log struct {
	Message    string `yaml:"message"`
	Created_at int64  `yaml:"created_at"`
	Level      string `yaml:"level"`
	Source     string `yaml:"scource"`
}

func LoggerChannel() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			select {
			case v := <-ChanLogEvent:
				logger.OutputLogger(v.Level, v.Message)
			default:
			}
		}
	}()
}

func AddLogtoChan(err error) {
	var message string
	var log Log
	if err == nil {
		message := "[MONITOR]Begin\n"
		log = Log{Level: "INFO", Message: message}
	} else {
		pc, fn, line, _ := runtime.Caller(1)
		message = fmt.Sprintf("[MONITOR][%s %s:%d] %s", runtime.FuncForPC(pc).Name(), fn, line, err)
		log = Log{Level: "ERROR", Message: message}
	}

	ChanLogEvent <- log
	return
}
