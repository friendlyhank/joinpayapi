package joinpayapi

//错误机制处理

type ErrCode int32

const (
	ErrCode_Success       ErrCode = 100
	ErrCode_Fail          ErrCode = 101
	ErrCode_AcceptSuccess ErrCode = 2001
	ErrCode_AcceptFail    ErrCode = 2002
	ErrCode_NoSure        ErrCode = 2003
)

var ErrCode_name = map[int32]string{
	100:  "Success",       //成功
	101:  "Fail",          //失败
	2001: "AcceptSuccess", //受理成功
	2002: "AcceptFail",    //受理失败
	2003: "NoSure ",       //不确定
}
