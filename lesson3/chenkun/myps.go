package main

/*
 实现类似 ps 的功能（读取/proc/下的pid目录名，及pid目录下cmdline文件内容）
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	procdir     = "/proc/"
	proccommand = "/cmdline"
)

func main() {
	f, err := os.Open(procdir) // 打开目录
	if err != nil {
		log.Fatal(err)
	}
	infos, _ := f.Readdir(-1)    // 读取目录下的文件或目录（结果类型是数组）
	for _, info := range infos { // 遍历结果，对目录为数字的（确定它是PID），读取目录下cmdline文件内容（该PID对应的命令）
		if info.IsDir() {
			infoname, err := strconv.Atoi(info.Name()) // 字符串转换为数字，如果不出err 说明其为数字
			if err == nil {                            // 没有异常
				filebuf, _ := ioutil.ReadFile(procdir + strconv.Itoa(infoname) + proccommand) // 读取cmdline文件内容
				fmt.Println(infoname, string(filebuf))                                        // 打印PID，PID的命令名
			}
		}
	}
	f.Close()
}
