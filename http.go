package goparseexcel

import (
	"github.com/jellycheng/gcurl"
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
