package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	// 执行 go run code-user/main.go
	cmd := exec.Command("go", "run", "code-user/main.go")
	// io.Writer 标准输出,错误
	var out, stderr bytes.Buffer
	// 执行完cmd命令后提示用户的错误
	cmd.Stderr = &stderr
	// 执行完cmd命令的输出结果
	cmd.Stdout = &out
	// StdinPipe 返回一个管道，该管道将在命令启动时连接到命令的标准输入。
	stdinPipi, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	io.WriteString(stdinPipi, "23 11\n")
	// 根据测试的输入案列进行运行，拿到输出结果和标准的输出结果进行比对
	if err := cmd.Run(); err != nil {
		log.Fatalln(err, stderr.String())
	}
	fmt.Println(out.String())
	fmt.Println(out.String() == "34\n")
}
