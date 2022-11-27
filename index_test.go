package goparseexcel

import (
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
	if err := DataProcessMode("example.toml"); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}

}
