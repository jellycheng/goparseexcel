package goparseexcel

type ExcelDto struct {
	Cols     []string            // excel列，["A", "B", "C"]
	Header   map[string]string   // 表头 A=》表头内容即表头名，B=》表头内容
	RowsData []map[string]string // 多条行记录，A=》内容，B=》内容
}

type ApiBodyDto struct {
	Header map[string]string   // 表头 参数名1=》表头名,company_name=>公司名称
	Data   []map[string]string // 多条数据，参数名1=》内容，参数名2=》内容，company_name=>xxx公司
}

type ApiBodyOneDto struct {
	Header map[string]string // 表头 参数名1=》表头名,company_name=>公司名称
	Data   map[string]string // 单条数据，参数名1=》内容，参数名2=》内容，company_name=>xxx公司
}
