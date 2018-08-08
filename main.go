package main

import (
	"os"
	"encoding/json"
	"fmt"
	"sync"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"time"
	"log"
	"flag"
	"github.com/sevlyar/go-daemon"
	"gopkg.in/urfave/cli.v1"
	"runtime"
)

const (
	usage = "Easy meta refresh tool"
)

var wg sync.WaitGroup
var (
	version   = "master"
	wait	  = 90
)
var (
	logF = flag.String("log", "run.log", "Log file name")
)

type configuration struct {
	Urls    []string
}

func main() {
	initConf()
	initLog()
	var ostype = runtime.GOOS
	if ostype == "windows"{
		performOne()
		return
	}
	app := cli.NewApp()
	app.Version = version
	app.Name = "go-meta-refresh"
	app.Usage = usage

	app.Commands = []cli.Command{
		cli.Command{
			Name: "execute",
			Usage:"命令行执行",
			Action: func(c *cli.Context) error {
				performOne()
				return nil
			},
		},
		cli.Command{
			Name: "daemon",
			Usage:"进驻后台",
			Action: func(c *cli.Context) error {
				daemonize(os.Args)
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func daemonize(arguments []string){
	cntxt := &daemon.Context{
		PidFileName: "gometa.pid",
		PidFilePerm: 0644,
		LogFileName: "run.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        arguments,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	performOne()
}

var conf = configuration{}

func initConf()  {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("conf.json文件不存在或配置不正确")
	}
	if len(conf.Urls) == 0 {
		fmt.Println("Error:网址未配置")
		return
	}
}

func performOne() {

	for _, v := range conf.Urls {
		wg.Add(1)
		go postScrape(v,v)
	}
	wg.Wait() //阻塞等待所有组内成员都执行完毕退栈
	log.Println("所有任务执行完毕")
}

func initLog(){
	// 定义一个文件
	outfile, err := os.OpenFile(*logF, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(*outfile, "open failed")
		os.Exit(1)
	}
	log.SetOutput(outfile)  //设置log的输出文件，不设置log输出默认为stdout
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //设置答应日志每一行前的标志信息，这里设置了日期，打印时间，当前go文件的文件名
}

func postScrape(v string,t string) {
	_, er := url.Parse(t)
	if er != nil {
		log.Println("非正常URL:"+t)
		//记录日志
		fmt.Println("Error:", er)
	}else{
		doc, errs := goquery.NewDocument(t)
		if errs != nil {
			log.Println("获取URL资源错误:"+t)
			// 出现错误写入日志
			fmt.Println("Error:", errs)
		}else{
			log.Println("正在执行:"+t)
			doc.Find("meta[http-equiv=\"refresh\"]").Each(func(index int, item *goquery.Selection) {
				content, exists := item.Attr("content")
				if exists {
					params :=strings.Split(content, ";")
					seconds,error := strconv.Atoi(params[0])
					if error != nil{
						log.Println("解析时间出错:"+params[0])
						// 出现错误写入日志
						fmt.Println("Error:", error)
						seconds = wait
					}
					if seconds<=0{
						seconds = wait
					}
					url := strings.Replace(strings.Replace(params[1], " ", "", -1),"url=",v,-1)
					time.Sleep(time.Duration(seconds) * time.Second)
					wg.Add(1)//为同步等待组增加一个成员
					go postScrape(v,url)
				}
			})
		}
	}

	defer wg.Done()
}




