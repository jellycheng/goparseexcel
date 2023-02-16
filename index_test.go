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

// go test -run=TestPingWhereIdin
func TestPingWhereIdin(t *testing.T) {
	if err := DataProcessMode("example.toml", func(tomlCfg gosupport.H, dto ApiBodyDto) error {
		tmpSlice := []string{}
		fieldName01 := "deptName" // 要取的列，对应的字段值
		for _, v := range dto.Data {
			if tmp, ok := v[fieldName01]; ok && tmp != "" {
				tmpSlice = append(tmpSlice, tmp)
			}
		}
		// 去重
		tmpSlice = gosupport.RemoveRepeatByString(tmpSlice)
		// 拼接，19591,19595,19599
		s, _ := gosupport.SliceJointoString(tmpSlice, ",", false)
		fmt.Println(s)

		return nil
	}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}

// go test -run=TestPingWhereStrin
func TestPingWhereStrin(t *testing.T) {
	if err := DataProcessMode("example.toml", func(tomlCfg gosupport.H, dto ApiBodyDto) error {
		tmpSlice := []string{}
		fieldName01 := "deptName" // 要取的列，对应的字段值
		for _, v := range dto.Data {
			if tmp, ok := v[fieldName01]; ok && tmp != "" {
				tmpSlice = append(tmpSlice, tmp)
			}
		}
		// 去重
		tmpSlice = gosupport.RemoveRepeatByString(tmpSlice)
		// 拼接，'19591','19595','19599'
		s := ""
		if tmpstr, _ := gosupport.SliceJointoString(tmpSlice, "','", false); tmpstr != "" {
			s = "'" + tmpstr + "'"
		}
		fmt.Println(s)

		return nil
	}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}

// go test -run=TestPinSql
func TestPinSql(t *testing.T) {
	// sql模板
	sqlTpl := "INSERT INTO `t_shortlink_token` (`id`, `appid`, `token`, `mark`, `is_delete`, `create_time`, `update_time`, `delete_time`) " +
		"VALUES (null, '{appid}', '{token}', '{mark}', 0, {create_time}, {update_time}, 0);"

	// 提取code
	codes := gosupport.ExtractCode(sqlTpl)
	// 默认值
	defaultData := map[string]string{
		"create_time":gosupport.ToStr(gosupport.TimeNow()),
		"update_time":gosupport.ToStr(gosupport.TimeNow()),
		"mark":"xxx原因新增记录",
	}

	// example
	if err := DataProcessMode("cjs.toml", func(tomlCfg gosupport.H, dto ApiBodyDto) error {
		for _, vData := range dto.Data {
			newSql := sqlTpl
			if len(codes) == 0 {
				break
			}
			for _, c := range codes {
				if val, ok := vData[c]; ok {
					newSql, _ = gosupport.Replace4code(newSql, c, val)
				} else if val2,ok2 := defaultData[c]; ok2 {
					newSql, _ = gosupport.Replace4code(newSql, c, val2)
				}else {
					newSql, _ = gosupport.Replace4code(newSql, c, "")
				}
			}
			fmt.Println(newSql)
			// 写文件
			gosupport.FilePutContents("./cjs.sql", newSql + gosupport.GO_EOL)

		}

		return nil
	}); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}


