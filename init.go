package main

import ()

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
