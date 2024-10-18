package util

import (
	"log"
	"os"
	"path/filepath"
)

// ExecPath
// @author: [kulisi](https://github.com/kulisi)
// @function: ExecPath
// @description: 获取可执行文件的绝对路径
// @return: string
func ExecPath() string {
	file, err := os.Executable()
	if err != nil {
		log.Fatalf("Executable file path error : %s\n", err.Error())
	}
	path := filepath.Dir(file)
	return path
}

// ExecPathJoin
// @author: [kulisi](https://github.com/kulisi)
// @function: ExecPathJoin
// @description: 将文件转换成可执行文件运行目录下文件
// @param: target string
// @return: string
func ExecPathJoin(target string) string {
	file, err := os.Executable()
	if err != nil {
		// log.Printf("executeble file path error: %s\r\n", err.Error())
		log.Fatalf("executeble file path error : %s\n", err.Error())
	}
	_, target = filepath.Split(target)
	return filepath.Join(filepath.Dir(file), target)
}

// PathExists
// 判断目录是否存在
func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if !fi.IsDir() {
		return false, nil
	}
	return true, nil
}
