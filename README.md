# http2gopher 

## 描述
一款将http请求报文转化为gopher请求报文的工具。
## 下载
go version >= 1.6
```
go install -v github.com/ShangRui-hash/http2gopher@latest 
```
go version < 1.6
```
go get github.com/ShangRui-hash/http2gopher
```
## 使用
```
http2gopher -h
```
```
NAME:
   http2gopher - 一个用来将http请求报文转换成gopher请求的工具

USAGE:
   http2gopher

VERSION:
   v0.1

AUTHOR:
   无在无不在

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value         指定http请求报文所在的文件
   --doubleURLencoded, -d         是否进行双重URL编码 (default: false)
   --doNotCheckContentLength, -n  不检查Content-Length (default: false)
   --help, -h                     show help (default: false)
   --version, -v                  print the version (default: false)
```
## 使用案例
./test/post.txt
```
POST /flag.php HTTP/1.1
Host: 127.0.0.1:80
Content-Type: application/x-www-form-urlencoded
Content-Length: 37

key=d994d9b26da17eec2060c5e40b466186
```
默认情况下工具会自动检查 body的长度，修正Content-Length 的长度
```
/bin/bash > http2gopher -f ./test/post.txt
gopher://127.0.0.1:80/_POST%20/flag.php%20HTTP/1.1%0d%0aHost:%20127.0.0.1:80%0d%0aContent-Type:%20application/x-www-form-urlencoded%0d%0aContent-Length:%2036%0d%0a%0d%0akey=d994d9b26da17eec2060c5e40b466186
```
如果要关闭自动修正，可以添加 -n 参数
```
bin/bash > http2gopher -f ./test/post.txt -n 
gopher://127.0.0.1:80/_POST%20/flag.php%20HTTP/1.1%0d%0aHost:%20127.0.0.1:80%0d%0aContent-Type:%20application/x-www-form-urlencoded%0d%0aContent-Length:%2037%0d%0a%0d%0akey=d994d9b26da17eec2060c5e40b466186
```
如果需要将%也进行URL编码，添加-d参数
```
bin/bash > http2gopher -f ./test/post.txt -d
gopher://127.0.0.1:80/_POST%2520/flag.php%2520HTTP/1.1%250d%250aHost:%2520127.0.0.1:80%250d%250aContent-Type:%2520application/x-www-form-urlencoded%250d%250aContent-Length:%252036%250d%250a%250d%250akey=d994d9b26da17eec2060c5e40b466186
```

## 注意事项:
```
gopher://<host>:<port>/<gopher-path>_后接TCP数据流
```
gopher 协议的默认端口是70 。所以最好在http请求报文中将 80 或者443 端口在Host字段显式写出来，程序会自动读取Host字段值，填充到 <host>:<port> 部分

## 参考资料
- [https://www.cxyzjd.com/article/weixin_45887311/107327706](https://www.cxyzjd.com/article/weixin_45887311/107327706)
- [https://blog.csdn.net/weixin_44037296/article/details/118387034](https://blog.csdn.net/weixin_44037296/article/details/118387034)

## QQ 交流群
<img src="https://store.heytapimage.com/cdo-portal/feedback/202301/04/ed1d5ac9f0c48af0a154037fb892024f.png" height="250px" width="250px" alt="图片.png" title="图片.png" referrerPolicy="no-referrer" />

