package errors

type errCode int

type ErrMsg struct {
	Code errCode     `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"date"`
}

const (
	CodeSuccess = 1000 + iota
	CodeInvalidPath
	CodeKeySetFaild
	CodeKeyNotFound
	CodeKeyDelFaild
)

var codeMap = map[errCode]string{
	CodeSuccess: "成功",

	CodeInvalidPath: "无效的请求路径",

	CodeKeySetFaild: "未成功设置key",
	CodeKeyNotFound: "未找到指定的key",
	CodeKeyDelFaild: "未成功删除key",
}

func NewerrMsg(code errCode, data interface{}) *ErrMsg {
	return &ErrMsg{
		Code: code,
		Msg:  codeMap[code],
		Data: data,
	}
}
