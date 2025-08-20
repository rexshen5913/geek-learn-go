package web

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type Context struct {

	// 如果用戶直接使用這個，那麼他們就繞開了 respData 和 respStatusCode
	// 這樣的話，部分 middleware 就無法正常運作
	Resp http.ResponseWriter
	Req  *http.Request

	// 這個主要是讓 middleware 讀寫用的
	RespData       []byte
	RespStatusCode int

	PathParams map[string]string

	queryValues  url.Values
	MatchedRoute string

	tplEngine TemplateEngine
}

func (c *Context) Render(tpl string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tpl, data)
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
		return err
	}
	c.RespStatusCode = http.StatusOK
	return nil
}

// 泛型 QueryValueV2[int]("key1") 希望返回 int
// 這是做不到的，因為 Method could not have type parameters
// 因為方法本身不能再額外宣告自己的型別參數
// func (c *Context) QueryValueV2[T any](key string) (T, error) {
// 	var val T
// 	val, ok := c.queryValues[key]
// 	if !ok {
// 		return val, errors.New("web: key not found")
// 	}
// 	return val, nil
// }

func QueryValueV3[T any](c *Context, key string, parse func(string) (T, error)) (T, error) {
	var zero T
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	s, ok := c.queryValues[key] // 假設是 map[string]string
	if !ok {
		return zero, errors.New("web: key not found")
	}
	return parse(s[0])
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) RespStatusOK(val any) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(status int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.Header().Set("Content-Type", "application/json")
	c.Resp.Header().Set("Content-Length", strconv.Itoa(len(data)))

	c.RespStatusCode = status
	c.RespData = data

	// 不要立即写入响应，让框架的 flushResp 统一处理
	// 这样中间件才能正常工作
	return nil
}

func (c *Context) BindJSON(val any) error {
	// if c.Req.Body == nil {
	// 	return errors.New("web: body 為 nil")
	// }
	if c.Req.Body == nil {
		return errors.New("web Body is nil")
	}
	decoder := json.NewDecoder(c.Req.Body)
	// 使用 Number 类型, 避免精度丢失, 例如 1.0 會被轉換成 1
	// 默認是用 float64 來接收
	// decoder.UseNumber()

	// 如果 JSON 中存在未定義的欄位, 則返回錯誤
	// 這個選項可以防止 JSON 中存在未定義的欄位, 例如:
	// 假設 User 結構體中只有 name 和 age 兩個欄位,
	// 但是 JSON 中多了 email 欄位, 這樣會返回錯誤
	// decoder.DisallowUnknownFields()

	// 解碼 JSON 到 val 中
	return decoder.Decode(val)
}

func (c *Context) BindXML(val any) error {
	if c.Req.Body == nil {
		return errors.New("web Body is nil")
	}
	decoder := xml.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.Req.FormValue(key), nil
}

// Query 比起 Form 是沒有緩存, 每次都會重新解析
func (c *Context) QueryValueV1(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key not found"),
		}
	}
	return StringValue{
		val: vals[0],
	}

	// return c.queryValues.Get(key), nil
}

// Query 比起 Form 是沒有緩存, 每次都會重新解析
func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return "", errors.New("web: key not found")
	}
	return vals[0], nil

	// return c.queryValues.Get(key), nil
}

func (c *Context) PathValueV1(key string) StringValue {
	val, ok := c.PathParams[key]
	if !ok {
		return StringValue{
			err: errors.New("web: key not found"),
		}
	}
	return StringValue{
		val: val,
	}
}

func (c *Context) PathValue(key string) (string, error) {
	vals, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key not found")
	}
	return vals, nil
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}

// 一般來說，Context 並非線程安全的，若是真的有這個需求，可以自己實現一個線程安全的 Context
// 例如:
type SafeContext struct {
	mu  sync.RWMutex
	ctx *Context
}

func (sc *SafeContext) SetCookie(cookie *http.Cookie) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.ctx.SetCookie(cookie)
}

func (sc *SafeContext) RespStatusOK(val any) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.ctx.RespStatusOK(val)
}
