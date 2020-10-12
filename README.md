# SKYNET
###### Go语言开发的简易聊天室

#### 简介：

**SKYNET** 是完全利用 **Go** 开发的聊天室程序，开发的初衷是为了熟悉 **Go** 语言对多线程以及网络服务器开发的支持。**SKYNET** 的客户端UI使用了 **tui-go**, 详情请参考 `https://github.com/marcusolsson/tui-go`。**SKYNET** 服务端主要参照 The Go Programming Language 8.10章节。

#### 客户端使用方法：

* 首先使用 `git clone` 将此文件夹下载到本地；
* `cd ./tui/cmd/ `来到这个目录下；
* 编译主程序 `go build main.go`；如果本地没有 **Go** 环境请在 `https://golang.org/dl/`下载；
* 运行编译好的 `main` 程序，`./main -server server-address:8000`

#### 服务端使用方法：

* 选择你想使用的web-service，创建linux操作系统的实例；
* 在服务器上安装 **Go** 语言环境；
* 将本地代码上传到服务器，或者直接使用 **Git**下载代码到服务器并解压；
* `cd ./server/main` 来到这个目录下；
* 编译主程序 `go build main.go` 并运行 `./main`
* 服务器默认监听`8000`端口；

