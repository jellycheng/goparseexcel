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

## 分析excel
```
package main

import (
	"fmt"
	"github.com/jellycheng/goparseexcel"
	"strings"
)

func main() {
	excelFile := "/Users/jelly/Desktop/合同订单-20230608-导出数据.xlsx"
	sheetName := ""
	count := 0 // 控制解析行数
	if err, rows := goparseexcel.ParseExcel(excelFile, sheetName); err == nil {
		for rowKey, row := range rows { // 循环行
			fmt.Println("rowKey=", rowKey)
			count++
			if count > 4 {
				break
			}
			for colkey, colCell := range row { // 循环列
				colName := goparseexcel.GetExcelNo(colkey) // excel列号 A、B、C...
				colCell = strings.TrimSpace(colCell)       // 内容并去掉前后空格
				fmt.Println("列号=", colName, " | 列内容=", colCell)
			}
		}
	} else {
		fmt.Println(err.Error())
	}

}

```
