# go-meta-refresh
a golang meta refresh tool
模拟浏览器当页面输出带有refresh的meta时
`<meta http-equiv=refresh content="50; url=?">`
自动进行刷新的操作

## 下载安装

源码编译或下载release包

```
$ git clone https://github.com/xiaogouxo/go-meta-refresh.git
```

### 使用说明
编译
```
go build main.go
```
配置
```
1.在同级目录下新建conf.json
2.配置网址
{
  "urls": [
    "http://www.xxx.com",
    "http://myserver.com/sync-data.php"
  ]
}
```
执行`mac/linux下支持两种运行方式，windows下直接运行exe`
```bash
$ ./main -h
NAME:
   go-meta-refresh - Easy meta refresh tool

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   master

COMMANDS:
     execute  命令行执行
     daemon   进驻后台
     help, h  Shows a list of commands or help for one command
```

