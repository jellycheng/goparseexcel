package goparseexcel

import (
	"github.com/jellycheng/gcurl"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/sign"
)

func PostJson(reqUrl, strJson, signSecret string) (error, string) {
	signStr := sign.BodyContentMd5Sign(strJson, signSecret)
	body := ""
	if resp, err := gcurl.Post(reqUrl, gcurl.Options{
		Headers: map[string]interface{}{
			"User-Agent": "gcurl/1.0",
			"sign":       signStr,
		},
		Query: map[string]interface{}{
			"type": "import_account",
			"sign": signStr,
		},
		JSON: strJson,
	}); err != nil {
		return err, body
	} else {
		res, _ := resp.GetBody()
		return nil, res.GetContents()
	}
}

func PostJsonV2(reqUrl, strJson, signSecret string, headers map[string]interface{}, querys map[string]interface{}) (error, string) {
	signStr := sign.BodyContentMd5Sign(strJson, signSecret)
	headersMap := gosupport.NewDataManage()
	headersMap.Set("User-Agent", "gcurl/1.0")
	headersMap.Set("sign", signStr)
	for k, v := range headers {
		headersMap.Set(k, v)
	}
	queryMap := gosupport.NewDataManage()
	queryMap.Set("sign", signStr)
	for k, v := range querys {
		queryMap.Set(k, v)
	}
	if resp, err := gcurl.Post(reqUrl, gcurl.Options{
		Headers: headersMap.GetData(),
		Query:   queryMap.GetData(),
		JSON:    strJson,
	}); err != nil {
		return err, ""
	} else {
		res, _ := resp.GetBody()
		return nil, res.GetContents()
	}
}
