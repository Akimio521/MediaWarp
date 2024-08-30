package pkg

import (
	"errors"
	"io"
	"os"
)

// 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) { // isnotexist来判断，是不是不存在的错误
		// 如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err // 如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// 判断路径是否为文件夹
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// 判断路径是否为文件
func IsFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return !fileInfo.IsDir(), nil
}

// 读取文件内容
func GetFileContent(filepath string) ([]byte, error) {
	isFile, err := IsFile(filepath)
	if err != nil {
		return nil, err
	}

	if !isFile {
		return nil, errors.New(filepath + "不是文件")
	}

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContent, nil

}
