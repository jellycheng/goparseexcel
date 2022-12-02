# goparseexcel
```

```
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/jellycheng/goparseexcel.svg)](https://pkg.go.dev/github.com/jellycheng/goparseexcel)
[![Goproxy.cn](https://goproxy.cn/stats/github.com/jellycheng/goparseexcel/badges/download-count.svg)](https://goproxy.cn/stats/github.com/jellycheng/goparseexcel/badges/download-count.svg)

## Requirements
```
gosupport library requires Go version >=1.16

```

## get依赖
```
方式1：
    go get github.com/jellycheng/goparseexcel

方式2：
GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get github.com/jellycheng/goparseexcel

```

## Documentation
[https://pkg.go.dev/github.com/jellycheng/goparseexcel](https://pkg.go.dev/github.com/jellycheng/goparseexcel)

## 示例
```
package main

import (
	"errors"
	"fmt"
	"github.com/jellycheng/goparseexcel"
	"github.com/jellycheng/gosupport"
)

func main() {

	if err := goparseexcel.DataProcessMode("example.toml", func(tomlCfg gosupport.H, dto goparseexcel.ApiBodyDto) error {		
		if apiConfig, isOk := tomlCfg["api_config"].(map[string]interface{}); isOk {
			signSecret := gosupport.ToStr(apiConfig["sign_secret"])
			reqUrl := gosupport.ToStr(apiConfig["api_url"])
            var errOne error
            for _, tmpVal := range dto.Data {//遍历excel记录
                reqOne := goparseexcel.ApiBodyOneDto{
                    Header: dto.Header,
                    Data:   tmpVal,
                }
                // 发起请求,单条
                if err, ret := goparseexcel.PostJson(reqUrl, gosupport.ToJson(reqOne), signSecret); err != nil {
                    errOne = err
                } else {
                    fmt.Println(ret)
                }
            }
            return errOne
		} else {
			return errors.New("缺少api_config配置")
		}

	}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}

```
