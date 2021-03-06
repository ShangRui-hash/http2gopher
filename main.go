package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var (
	httpRequestFile         string
	isDoubleURLencoded      bool
	doNotCheckContentLength bool
)

func main() {
	author := cli.Author{
		Name: "无在无不在",
	}
	app := &cli.App{
		Name:      "http2gopher",
		Usage:     "一个用来将http请求报文转换成gopher请求报文的工具",
		UsageText: "http2gopher",
		Version:   "v0.1",
		Authors:   []*cli.Author{&author},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "指定http请求报文所在的文件",
				Destination: &httpRequestFile,
				Required:    true,
			},
			&cli.BoolFlag{
				Name:        "doubleURLencoded",
				Aliases:     []string{"d"},
				Usage:       "是否进行双重URL编码",
				Value:       false,
				Destination: &isDoubleURLencoded,
			},
			&cli.BoolFlag{
				Name:        "doNotCheckContentLength",
				Aliases:     []string{"n"},
				Usage:       "不检查Content-Length",
				Value:       false,
				Destination: &doNotCheckContentLength,
			},
		},
		Action: run,
	}
	//启动app
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
	}
}

func run(c *cli.Context) (err error) {
	//读取http请求文件
	buffer, err := ioutil.ReadFile(httpRequestFile)
	if err != nil {
		logrus.Error("ioutil.ReadFile failed,err:", err)
		return err
	}
	//转换成gopher请求
	gopherRequest, err := http2gopher(buffer)
	if err != nil {
		logrus.Error("http2gopher failed,err:", err)
		return err
	}
	//输出gopher请求
	fmt.Println(gopherRequest)
	return nil
}

func getBody(request *http.Request, buffer []byte) (body string) {
	//获取body
	if buf, err := ioutil.ReadAll(request.Body); err == nil {
		return string(buf)
	}
	flag := false
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	for scanner.Scan() {
		line := scanner.Text()
		if flag {
			body = body + line + "\n"
		}
		if line == "" {
			flag = true
		}
	}
	return body
}

func http2gopher(buffer []byte) (gopherRequest string, err error) {
	//解析http请求报文
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buffer)))
	if err != nil {
		logrus.Error("http.ReadRequest failed,err:", err)
		return "", err
	}
	//初始化gopher
	gopherRequest = fmt.Sprintf("gopher://%s/_", request.Host)
	gopherRequest += fmt.Sprintf("%s %s %s", request.Method, request.RequestURI, request.Proto)
	gopherRequest += "%0d%0a"
	gopherRequest += fmt.Sprintf("Host: %s", request.Host)
	gopherRequest += "%0d%0a"
	body := getBody(request, buffer)
	//转化请求头
	if !doNotCheckContentLength {
		request.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))
	}
	for k, v := range request.Header {
		value := strings.Join(v, ",")
		gopherRequest += fmt.Sprintf("%s: %s", k, value)
		gopherRequest += "%0d%0a"
	}
	gopherRequest += "%0d%0a"
	//转换body
	lines := strings.Split(body, "\n")
	linesLen := len(lines)
	for i := range lines {
		gopherRequest = gopherRequest + strings.TrimSpace(lines[i])
		//追加%0d%0a
		if i != linesLen-1 {
			gopherRequest = gopherRequest + "%0d%0a"
		}
	}
	//将空格替换为%20
	gopherRequest = strings.Replace(gopherRequest, " ", "%20", -1)
	//将&替换为%26
	gopherRequest = strings.Replace(gopherRequest, "&", "%26", -1)
	//将#替换为%23
	gopherRequest = strings.Replace(gopherRequest, "#", "%23", -1)
	//将=替换为%3d
	gopherRequest = strings.Replace(gopherRequest, "=", "%3d", -1)
	//如果开启了双重URL编码
	if isDoubleURLencoded {
		gopherRequest = strings.Replace(gopherRequest, "%", "%25", -1)
	}
	return gopherRequest, nil
}
