package api

import (
	"MediaWarp/schemas/schemas_alist"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type AlistServer struct {
	URL      string
	Username string
	Password string
	Token    string
}

// 登录Alist（获取Token）
func (alistServer *AlistServer) AuthLogin() (err error) {
	funcInfo := "Alist登录"
	url := alistServer.URL + "/api/auth/login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "username": "%s",
    "password": "%s"
}`, alistServer.Username, alistServer.Password))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s失败: %w", funcInfo, err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return
	}

	var (
		authLogin     schemas_alist.AuthLogin
		alistResponse schemas_alist.AlistResponse
	)

	err = json.Unmarshal(body, &alistResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return
	}
	if alistResponse.Code != 200 {
		err = errors.New(alistResponse.Message)
		return
	}
	dataBytes, err := json.Marshal(alistResponse.Data)
	if err != nil {
		err = fmt.Errorf("序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	// 将 JSON 字符串解析为 authLogin 结构体
	err = json.Unmarshal(dataBytes, &authLogin)
	if err != nil {
		err = fmt.Errorf("反序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	alistServer.Token = authLogin.Token
	return
}

// 获取某个文件/目录信息
func (alistServer *AlistServer) FsGet(path string) (fsGet schemas_alist.FsGet, err error) {
	funcInfo := "Alist获取某个文件/目录信息"
	url := alistServer.URL + "/api/fs/get"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "path": "%s",
    "password": "",
    "page": 1,
    "per_page": 0,
    "refresh": false
}`, path))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return
	}
	req.Header.Add("Authorization", alistServer.Token)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s信息失败: %w", funcInfo, err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return
	}

	var alistResponse schemas_alist.AlistResponse

	err = json.Unmarshal(body, &alistResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return
	}
	if alistResponse.Code != 200 {
		err = errors.New(alistResponse.Message)
		return
	}
	dataBytes, err := json.Marshal(alistResponse.Data)
	if err != nil {
		err = fmt.Errorf("序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	err = json.Unmarshal(dataBytes, &fsGet)
	if err != nil {
		err = fmt.Errorf("反序列化%s响应体data失败: %w", funcInfo, err)
		return
	}
	return
}
