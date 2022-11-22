package goparseexcel

import (
	"errors"
	"github.com/jellycheng/gosupport"
	"github.com/xuri/excelize/v2"
)

func ParseExcel(excelFile, sheetName string) (error, [][]string) {
	var err error
	var ret [][]string
	if !gosupport.FileExists(excelFile) {
		err = errors.New("excel 文件不存在：" + excelFile)
		return err, ret
	}

	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		return err, ret
	}
	defer func() {
		_ = f.Close()
	}()
	if sheetName == "" {
		sheetName = f.GetSheetName(f.GetActiveSheetIndex())
	}

	ret, err = f.GetRows(sheetName)
	if err != nil {
		return err, ret
	}

	return err, ret
}
