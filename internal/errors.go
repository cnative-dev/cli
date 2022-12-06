package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-resty/resty/v2"
)

type ErrResp struct {
	Code int `json:"code"`
}

// NotFoundError: { code: 404000 },
// ServerError: { code: 500000 },
// ClientError: { code: 400000 },
// ArgumentError: { code: 400001 },
// DuplicateError: { code: 409000 }

func translate(code int) string {
	switch code {
	case 404000:
		return "请求资源不存在"
	case 500000:
		return "服务端出错了"
	case 400000:
		return "客户端出错了"
	case 400001:
		return "参数错误"
	case 409000:
		return "资源冲突"
	default:
		return "出错了。"
	}
}

func HandleError(response *resty.Response, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "出错了, %s\n", err)
	} else {
		definedGlobalError := &ErrResp{}
		content, _ := ioutil.ReadAll(response.RawResponse.Body)
		if marshalError := json.Unmarshal(content, definedGlobalError); marshalError == nil {
			fmt.Fprintln(os.Stderr, translate(definedGlobalError.Code))
		} else {
			fmt.Fprintf(os.Stderr, "服务器出错了, [%s]%s\n", response.Status(), response.Body())
		}
	}
}
