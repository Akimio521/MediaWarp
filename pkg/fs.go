package pkg

import "os"

// 判断文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 判断是否为文件夹
func IsDir(path string) bool {
	fileInfo, _ := os.Stat(path)
	return fileInfo.IsDir()
}

// 判断是否为文件
func IsFile(path string) bool {
	fileInfo, _ := os.Stat(path)
	return !fileInfo.IsDir()
}
