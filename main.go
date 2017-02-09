package main

import (
	"cloud"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

const (
	Path = "monitor.yaml"
)

var ChanLogEvent chan Log

func main() {
	ChanLogEvent = make(chan Log, 1)
	Online()
	Outline()
	LoggerChannel()
	AddLogtoChan(nil)
	for {
		time.Sleep(10 * time.Second)
	}

}

func Outline() {
	go func() {
		for {

			var conf Config

			result := readConf()
			if err := yaml.Unmarshal([]byte(result), &conf); err != nil {
				AddLogtoChan(err)
			}
			for num, val := range conf.Monitor {
				if val.Ip != "" && !val.Status {
					if err := Ping(val.Ip, val.Port); err == nil {
						conf.Monitor[num].Status = true
						WriteConf(conf)
						Send(val.Name, conf.Mail.MessageOnline)
					}
				}
			}
			time.Sleep(time.Duration(conf.Setting.MonitorOutline) * time.Second)
		}

	}()
}

func Online() {
	go func() {
		for {

			var conf Config

			result := readConf()
			if err := yaml.Unmarshal([]byte(result), &conf); err != nil {
				AddLogtoChan(err)
			}

			for num, val := range conf.Monitor {
				if val.Ip != "" && val.Status {
					if err := Ping(val.Ip, val.Port); err != nil {
						Response(conf, val, num)

					}
				}
			}
			time.Sleep(time.Duration(conf.Setting.MonitorOnline) * time.Second)
		}

	}()
}

func Ping(ip, port string) error {
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Second*1)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func readConf() string {
	fi, err := os.Open(Path)
	if err != nil {
		AddLogtoChan(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		AddLogtoChan(err)
	}
	return string(fd)
}

func WriteConf(conf Config) {
	d, err := yaml.Marshal(&conf)
	if err != nil {
		AddLogtoChan(err)
	}

	str := "---\n" + string(d) + "\n"
	yaml := []byte(str)

	fi, err := os.Open(Path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	err = ioutil.WriteFile(Path, yaml, 0666)
	if err != nil {
		panic(err)
	}
}

func Send(name, message string) {
	var conf Config
	mails := make([]string, 0)

	result := readConf()
	if err := yaml.Unmarshal([]byte(result), &conf); err != nil {
		AddLogtoChan(err)
	}

	for _, val := range conf.Mail.Address {
		mails = append(mails, val)
	}
	cloud.Sendto(mails, name+" "+message, conf.Mail.Header)

}

func Response(conf Config, val Server, num int) error {
	count := 0
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		if err := Ping(val.Ip, val.Port); err != nil {
			count += 1
		}
	}
	if count > 3 {

		conf.Monitor[num].Status = false
		WriteConf(conf)
		Send(val.Name, conf.Mail.MessageOutline)
		AddLogtoChan(errors.New(val.Name))
	}
	return nil
}
