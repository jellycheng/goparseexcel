package goparseexcel

import (
	"errors"
	"fmt"
	"github.com/jellycheng/gosupport"
	"testing"
)

// go test -run=TestParseExcel
func TestParseExcel(t *testing.T) {
	contents := []map[string]string{}
	if err, rows := ParseExcel("cjs.xlsx", ""); err == nil {
		for _, row := range rows { // 行
			tmpCon := map[string]string{}
			for colkey, colCell := range row { //列
				colName := GetExcelNo(colkey)
				tmpCon[colName] = colCell
			}
			contents = append(contents, tmpCon)
		}
		fmt.Println(gosupport.ToJson(contents))

	} else {
		fmt.Println(err.Error())
	}

}

// go test -run=TestParseToml
func TestParseToml(t *testing.T) {
	if cfg, err := ParseToml("example.toml"); err == nil {
		fmt.Println(cfg)
	} else {
		fmt.Println(err)
	}

}

// go test -run=DataProcessMode
func TestDataProcessMode(t *testing.T) {
	if err := DataProcessMode("example.toml", func(tomlCfg gosupport.H, dto ApiBodyDto) error {
		defaultCfg := tomlCfg["default"].(map[string]interface{})
		excelHandleWay := gosupport.StrTo(gosupport.ToStr(defaultCfg["excel_handle_way"])).MustInt()
		// 发起请求
		if apiConfig, isOk := tomlCfg["api_config"].(map[string]interface{}); isOk {
			signSecret := gosupport.ToStr(apiConfig["sign_secret"])
			reqUrl := gosupport.ToStr(apiConfig["api_url"])
			if excelHandleWay == ExcelHandleWayOne { // 单条
				var errOne error
				for _, tmpVal := range dto.Data {
					reqOne := ApiBodyOneDto{
						Header: dto.Header,
						Data:   tmpVal,
					}
					if err, ret := PostJson(reqUrl, gosupport.ToJson(reqOne), signSecret); err != nil {
						errOne = err
					} else {
						fmt.Println(ret)
					}
				}
				return errOne
			} else if excelHandleWay == ExcelHandleWayAll { //所有
				err, ret := PostJson(reqUrl, gosupport.ToJson(dto), signSecret)
				fmt.Println(ret)
				return err
			}
			return nil
		} else {
			return errors.New("缺少api_config配置")
		}

	}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}
