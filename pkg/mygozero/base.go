package mygozero

import (
	"bytes"
	"encoding/json"
	"github.com/tiptok/gocomm/common"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type RestBase struct {
}

func (controller *RestBase) JsonUnmarshal(r *http.Request, v interface{}) error {
	var reader io.Reader
	var buf *bytes.Buffer = bytes.NewBuffer(nil)
	if r.ContentLength > 0 && strings.Contains(r.Header.Get(httpx.ContentType), httpx.JsonContentType) {
		read := io.TeeReader(r.Body, buf)
		reader = io.LimitReader(read, 8<<20) //8M
		r.Body = ioutil.NopCloser(buf)
	} else {
		reader = strings.NewReader("{}")
	}
	return jsonx.UnmarshalFromReader(reader, v)
}

func (controller *RestBase) BodyKeys(r *http.Request, firstCaseToUpper bool) []string {
	var bodyKV map[string]json.RawMessage
	controller.JsonUnmarshal(r, &bodyKV)
	if len(bodyKV) == 0 {
		return []string{}
	}
	var list []string
	for k, _ := range bodyKV {
		list = append(list, common.CamelCase(k, true))
	}
	return list
}

func (controller *RestBase) Resp(w http.ResponseWriter, msg interface{}) {
	httpx.OkJson(w, msg)
}

func (controller *RestBase) Response(w http.ResponseWriter, msg interface{}, err error) {
	if err != nil {
		httpx.Error(w, err)
		return
	}
	httpx.OkJson(w, msg)
}

func (controller *RestBase) GetLimitInfo(r *http.Request) (offset int, limit int) {
	var pageInfo struct {
		OffSet int `form:"offset"`
		Limit  int `form:"limit"`
	}
	httpx.ParseForm(r, &pageInfo)
	offset = pageInfo.OffSet
	limit = pageInfo.Limit
	if offset > 0 {
		offset = (offset - 1) * limit
	}
	return
}
