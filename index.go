package goparseexcel

import (
	"errors"
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/pelletier/go-toml/v2"
	"github.com/xuri/excelize/v2"
	"strings"
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

// GetExcelNo 下标从0开始，A-Z，AA-AZ，AAA-AAZ
func GetExcelNo(column int) string {
	ret := ""
	start := 65
	i := column % 26
	column = column / 26
	if column >= 1 {
		ret += GetExcelNo(column - 1)
	}
	return ret + string(rune(start+i))
}

// ParseExcelContent 解析excel内容,headerRow表头行号，rowsIgnore忽略的行记录、多个逗号分割，colsIgnore忽略列号、多个逗号分割
func ParseExcelContent(excelFile, sheetName string, headerRow int, rowsIgnore string, colsIgnore string) (error, ExcelDto) {
	excelDto := ExcelDto{
		Cols:     make([]string, 0),
		Header:   map[string]string{},
		RowsData: make([]map[string]string, 0),
	}
	if headerRow <= 0 {
		headerRow = 1
	}
	rowsIgnoreSlice := strings.Split(rowsIgnore, ",")
	colsIgnoreSlice := strings.Split(colsIgnore, ",")
	if err, rows := ParseExcel(excelFile, sheetName); err == nil {
		for rowKey, row := range rows { // 行
			tmpCon := map[string]string{}
			isIgnoreRow := false               // 是否忽略行
			for colkey, colCell := range row { //列
				colName := GetExcelNo(colkey)
				if rowKey == (headerRow - 1) { // 表头
					isIgnoreRow = true
					if gosupport.StrInSlice(colName, colsIgnoreSlice) { //忽略列
						continue
					}
					excelDto.Header[colName] = colCell
					if !gosupport.StrInSlice(colName, excelDto.Cols) {
						excelDto.Cols = append(excelDto.Cols, colName)
					}
				} else { // 内容
					if gosupport.StrInSlice(gosupport.ToStr(rowKey+1), rowsIgnoreSlice) { //忽略行
						isIgnoreRow = true
						continue
					}
					if gosupport.StrInSlice(colName, colsIgnoreSlice) { //忽略列
						continue
					}
					tmpCon[colName] = colCell
					if !gosupport.StrInSlice(colName, excelDto.Cols) {
						excelDto.Cols = append(excelDto.Cols, colName)
					}
				}

			}
			if !isIgnoreRow { // 不是忽略行
				excelDto.RowsData = append(excelDto.RowsData, tmpCon)
			}
		}
		return nil, excelDto
	} else {
		return err, excelDto
	}
}

func ParseToml(tomlFile string) (map[string]interface{}, error) {
	var err error
	var cfg map[string]interface{}
	if !gosupport.FileExists(tomlFile) {
		err = errors.New("toml 文件不存在：" + tomlFile)
		return cfg, err
	}
	con, _ := gosupport.FileGetContents(tomlFile)
	if err = toml.Unmarshal([]byte(con), &cfg); err == nil {
		return cfg, nil
	} else {
		return cfg, err
	}
}

func DataProcessMode(tomlFile string, callbackFunc func(gosupport.H, ApiBodyDto) error) error {
	if cfg, err := ParseToml(tomlFile); err == nil {
		tmp := gosupport.H(cfg)
		if defaultCfg, ok := tmp["default"].(map[string]interface{}); ok {
			if dataProcessMode, ok2 := defaultCfg["data_process_mode"]; ok2 {
				switch gosupport.StrTo(gosupport.ToStr(dataProcessMode)).MustInt() {
				case 1:
					if err, ret := DataProcessMode1(tmp); err == nil {
						//fmt.Println(gosupport.ToJson(ret))
						return callbackFunc(tmp, ret)
					} else {
						return err
					}
				case 3:
					if err, ret := DataProcessMode3(tmp); err == nil {
						return callbackFunc(tmp, ret)
					} else {
						return err
					}
				default:
					return errors.New(fmt.Sprintf("不支持default.data_process_mode=%d", gosupport.StrTo(gosupport.ToStr(dataProcessMode)).MustInt()))
				}
			} else {
				return errors.New("缺少default.data_process_mode配置")
			}
		} else {
			return errors.New("缺少default配置")
		}
	} else {
		return err
	}
}

// DataProcessMode1 组织field_mapping配置的表头
func DataProcessMode1(tomlCfg gosupport.H) (error, ApiBodyDto) {
	var ret = ApiBodyDto{
		Header: map[string]string{},
		Data:   make([]map[string]string, 0),
	}
	// 解析excel内容
	if defaultCfg, ok := tomlCfg["default"].(map[string]interface{}); ok {
		var headerRow int = 1 // 第几行用做表头
		if tmpVal, isOk := defaultCfg["header_row"]; isOk {
			headerRow = gosupport.StrTo(gosupport.ToStr(tmpVal)).MustInt()
		}
		var rowsIgnore string = "" // 忽略行,多个用逗号分割
		if tmpVal, isOk := defaultCfg["rows_ignore"]; isOk {
			rowsIgnore = gosupport.ToStr(tmpVal)
		}
		var colsIgnore string = "" // 忽略列,多个用逗号分割
		if tmpVal, isOk := defaultCfg["cols_ignore"]; isOk {
			colsIgnore = gosupport.ToStr(tmpVal)
		}
		if excelFile, ok1 := defaultCfg["excel_file"]; ok1 {
			if sheetName, ok2 := defaultCfg["excel_sheetname"]; ok2 {
				if err, excelDto := ParseExcelContent(excelFile.(string), sheetName.(string), headerRow, rowsIgnore, colsIgnore); err == nil {
					//fmt.Println(gosupport.ToJson(excelDto))
					// field_mapping
					if fieldMappingCfg, ok3 := tomlCfg["field_mapping"].(map[string]interface{}); ok3 {
						// 组装数据
						colsMap := map[string]string{}
						for _, v := range excelDto.Cols {
							if paramName, isOk := fieldMappingCfg[v]; isOk {
								ret.Header[paramName.(string)] = v
								colsMap[v] = paramName.(string)
							} else if hName, isOk := excelDto.Header[v]; isOk {
								if paramName, isOk := fieldMappingCfg[hName]; isOk {
									ret.Header[paramName.(string)] = hName
									colsMap[v] = paramName.(string)
								}
							}
						}
						for _, v := range excelDto.RowsData {
							dataMap := map[string]string{}
							for colName, fieldName := range colsMap {
								if v2, isOk := v[colName]; isOk {
									dataMap[fieldName] = v2
								} else {
									dataMap[fieldName] = ""
								}
							}
							ret.Data = append(ret.Data, dataMap)
						}
						return nil, ret
					} else {
						return errors.New("field_mapping未配置"), ret
					}

				} else {
					return err, ret
				}
			} else {
				return errors.New("default.excel_sheetname未配置"), ret
			}
		} else {
			return errors.New("default.excel_file未配置"), ret
		}
	}

	return nil, ret
}

// DataProcessMode3 使用excel列作为参数
func DataProcessMode3(tomlCfg gosupport.H) (error, ApiBodyDto) {
	var ret = ApiBodyDto{
		Header: map[string]string{},
		Data:   make([]map[string]string, 0),
	}
	// 解析excel内容
	if defaultCfg, ok := tomlCfg["default"].(map[string]interface{}); ok {
		var headerRow int = 1 // 第几行用做表头
		if tmpVal, isOk := defaultCfg["header_row"]; isOk {
			headerRow = gosupport.StrTo(gosupport.ToStr(tmpVal)).MustInt()
		}
		var rowsIgnore string = "" // 忽略行,多个用逗号分割
		if tmpVal, isOk := defaultCfg["rows_ignore"]; isOk {
			rowsIgnore = gosupport.ToStr(tmpVal)
		}
		var colsIgnore string = "" // 忽略列,多个用逗号分割
		if tmpVal, isOk := defaultCfg["cols_ignore"]; isOk {
			colsIgnore = gosupport.ToStr(tmpVal)
		}
		if excelFile, ok1 := defaultCfg["excel_file"]; ok1 {
			if sheetName, ok2 := defaultCfg["excel_sheetname"]; ok2 {
				if err, excelDto := ParseExcelContent(excelFile.(string), sheetName.(string), headerRow, rowsIgnore, colsIgnore); err == nil {
					//fmt.Println(gosupport.ToJson(excelDto))
					// 组装数据
					for _, v := range excelDto.Cols {
						ret.Header[v] = excelDto.Header[v]
					}
					ret.Data = excelDto.RowsData

					return nil, ret
				} else {
					return err, ret
				}
			} else {
				return errors.New("default.excel_sheetname未配置"), ret
			}
		} else {
			return errors.New("default.excel_file未配置"), ret
		}
	}

	return nil, ret
}
