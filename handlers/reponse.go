package handlers

type Response struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
}

var (
	Success = response(0, "ok")
	Fail    = response(-1, "fail")
)

func response(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: map[string]interface{}{},
	}
}

func (res *Response) WithMsg(message string) Response {
	return Response{
		Code: res.Code,
		Msg:  message,
		Data: res.Data,
	}
}

func (res *Response) WithData(data map[string]interface{}) Response {
	return Response{
		Code: res.Code,
		Msg:  res.Msg,
		Data: data,
	}
}
