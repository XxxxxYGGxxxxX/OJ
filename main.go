package main

import "OJ/router"

// 主函数 程序入口
func main() {
	r := router.Router()
	r.Run(":8800") // 监听并在 0.0.0.0:8800 上启动服务
}
