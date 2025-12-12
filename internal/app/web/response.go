package gono_web

const (
	ResponseCode_RequestError  = -2
	ResponseCode_InternalError = -1
	ResponseCode_OK            = 0
	ResponseCode_JwtExpired    = 1
)

const (
	ResponseMsg_OK = "OK"
)

type Response struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}
