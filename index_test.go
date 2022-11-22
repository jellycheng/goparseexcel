package goparseexcel

import (
	"fmt"
	"testing"
)

func TestToml(t *testing.T) {

}

// go test -run=TestParseExcel
func TestParseExcel(t *testing.T) {

	if err, ret := ParseExcel("cjs.xlsx", ""); err == nil {
		fmt.Println(ret)
	} else {
		fmt.Println(err.Error())
	}

}
