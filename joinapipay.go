package joinpayapi

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"git.biezao.com/ant/xmiss/foundation/cache"

	"git.biezao.com/ant/xmiss/foundation/uniqueid"
	"git.biezao.com/ant/xmiss/foundation/util/str"
	"git.biezao.com/ant/xmiss/foundation/vars"

	"git.biezao.com/ant/xmiss/externalapi/wxapi"
	"git.biezao.com/ant/xmiss/foundation/db"
	xhttp "git.biezao.com/ant/xmiss/foundation/http"
	"git.biezao.com/ant/xmiss/foundation/profile"
	"github.com/astaxie/beego/logs"
)

const (
	//聚合支付接口
	UNI_PAY_URL = "https://www.joinpay.com/trade/uniPayApi.action"
	//订单退款
	REFUND_URL = "https://www.joinpay.com/trade/refund.action"

	//单笔代付接口
	SINGLEPAY_URL = "https://www.joinpay.com/payment/pay/singlePay"
)

var (
	ErrOrderPaid = errors.New("订单已支付")
)

type JoinPayMakeOrderReq struct {
	Version          string `json:"p0_Version"`
	MerchantNo       string `json:"p1_MerchantNo"`       //商户编号
	OrderNo          string `json:"p2_OrderNo"`          //订单编号
	Amount           string `json:"p3_Amount"`           //支付金额
	Cur              string `json:"p4_Cur"`              //交易币种
	ProductName      string `json:"p5_ProductName"`      //商品名称
	ProductDesc      string `json:"p6_ProductDesc"`      //商品描述
	Mp               string `json:"p7_Mp"`               //公用传回参数
	ReturnUrl        string `json:"p8_ReturnUrl"`        //商户页面通知地址
	NotifyUrl        string `json:"p9_NotifyUrl"`        //服务器异步通知地址
	FrpCode          string `json:"q1_FrpCode"`          //交易类型
	MerchantBankCode string `json:"q2_MerchantBankCode"` //银行商户编码
	SubMerchantNo    string `json:"q3_SubMerchantNo"`    //子商户号
	IsShowPic        string `json:"q4_IsShowPic"`        //是否展示图片
	OpenId           string `json:"q5_OpenId"`           //微信Openid
	AuthCode         string `json:"q6_AuthCode"`         //付款码数字
	AppId            string `json:"q7_AppId"`            //APPID
	TerminalNo       string `json:"q8_TerminalNo"`       //终端号
	TransactionModel string `json:"q9_TransactionModel"` //微信H5模式
}

type JoinPayMakeOrderData struct {
	Appid     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type JoinPayMakeOrderRes struct {
	Version          string `json:"p0_Version"`
	MerchantNo       string `json:"r1_MerchantNo"`
	OrderNo          string `json:"r2_OrderNo"`
	Amount           string `json:"r3_Amount"`
	Cur              string `json:"r4_Cur"`
	Mp               string `json:"r5_Mp"`
	FrpCode          string `json:"r6_FrpCode"`
	TrxNo            string `json:"r7_TrxNo"`
	MerchantBankCode string `json:"r8_MerchantBankCode"`
	SubMerchantNo    string `json:"r9_SubMerchantNo"`
	Code             int32  `json:"ra_Code"`
	CodeMsg          string `json:"rb_CodeMsg"`
	Result           string `json:"rc_Result"`
	Pic              string `json:"rd_Pic"`
	Hmac             string `json:"hmac"`
}

type JoinPayRefundToUserReq struct {
}

type JoinPayRefundToUserRes struct {
	MerchantNo    string `json:"r1_MerchantNo"`
	OrderNo       string `json:"r2_OrderNo"`
	RefundOrderNo string `json:"r3_RefundOrderNo"`
	RefundAmount  string `json:"r4_RefundAmount"`
	RefundTrxNo   string `json:"r5_RefundTrxNo"`
	Status        string `json:"ra_Status"`
	Code          int64  `json:"ra_Code"`
	CodeMsg       string `json:"rb_CodeMsg"`
	Hmac          string `json:"hmac"`
}

type JoinPaySingleRePayReq struct {
	UserNo                string `json:"userNo"`
	ProductCode           string `json:"productCode"`
	RequestTime           string `json:"requestTime"`
	MerchantOrderNo       string `json:"merchantOrderNo"`
	ReceiverAccountNoEnc  string `json:"receiverAccountNoEnc"`
	ReceiverNameEnc       string `json:"receiverNameEnc"`
	ReceiverAccountType   string `json:"receiverAccountType"`
	ReceiverBankChannelNo string `json:"receiverBankChannelNo"`
	PaidAmount            string `json:"paidAmount"`
	Currency              string `json:"currency"`
	IsChecked             string `json:"isChecked"`
	PaidDesc              string `json:"paidDesc"`
	PaidUse               string `json:"paidUse"`
	CallbackUrl           string `json:"callbackUrl"`
	FirstProductCode      string `json:"firstProductCode"`
	Hmac                  string `json:"hmac"`
}

type JoinPaySingleRePayData struct {
	ErrorCode       string `json:"errorCode"`
	ErrorDesc       string `json:"errorDesc"`
	UserNo          string `json:"userNo"`
	MerchantOrderNo string `json:"merchantOrderNo"`
	Hmac            string `json:"hmac"`
}

type JoinPaySingleRePayRes struct {
	StatusCode int32                  `json:"statusCode"`
	Message    string                 `json:"message"`
	Data       JoinPaySingleRePayData `json:"data"`
}

// JoinPayMakeOrder -汇聚订单支付订单
func JoinPayMakeOrder(user *db.User, body, orderno, productname string, amount int64, key *db.Key) (joinPayMakeOrderRes *JoinPayMakeOrderRes, err error) {
	defer profile.TimeTrack(time.Now(), "[JoinPay-API] JoinPayMakeOrder")

	if amount <= 0 {
		err = errors.New("支付的金额有误")
		return
	}

	// 获取支付环境
	env := wxapi.GetPayEnv(key)
	if env == nil {
		err = fmt.Errorf("kid:%v,keyname:%v,获取支付环境失败", key.Kid, key.Name)
		return
	}

	var slicekey = []string{"p0_Version", "p1_MerchantNo", "p2_OrderNo", "p3_Amount", "p4_Cur", "p5_ProductName", "p6_ProductDesc", "p7_Mp", "p8_ReturnUrl", "p9_NotifyUrl",
		"q1_FrpCode", "q2_MerchantBankCode", "q3_SubMerchantNo", "q4_IsShowPic", "q5_OpenId", "q6_AuthCode", "q7_AppId", "q8_TerminalNo", "q9_TransactionModel", "hmac"}

	values := &url.Values{}
	values.Add("p0_Version", "1.0")
	values.Add("p1_MerchantNo", env.Mchid)
	values.Add("p2_OrderNo", orderno)
	values.Add("p3_Amount", str.GetMoneyYuan(amount))
	values.Add("p4_Cur", "1")
	values.Add("p5_ProductName", productname)
	values.Add("p6_ProductDesc", body)
	values.Add("p7_Mp", "")
	values.Add("p8_ReturnUrl", env.NotifyURL)
	values.Add("p9_NotifyUrl", env.NotifyURL)
	values.Add("q1_FrpCode", "WEIXIN_XCX") //WEIXIN_GZH公众号 WEIXIN_XCX 小程序支付
	values.Add("q2_MerchantBankCode", "")
	values.Add("q3_SubMerchantNo", "")
	values.Add("q4_IsShowPic", "1")
	values.Add("q5_OpenId", user.Logintoken)
	values.Add("q6_AuthCode", "")
	values.Add("q7_AppId", env.Appid) //
	values.Add("q8_TerminalNo", "")
	values.Add("q9_TransactionModel", "")

	if !vars.IsProd() { // 开发环境 , 测试环境
		values.Set("p3_Amount", "0.01")
	}

	sign := &Sign{Values: values}
	values.Add("hmac", sign.Sign(env.APIKey, slicekey, "Md5"))

	joinPayMakeOrderRes = &JoinPayMakeOrderRes{}
	err = xhttp.PostJSON(UNI_PAY_URL, values, nil, joinPayMakeOrderRes)
	if err != nil {
		logs.Error("|JoinApi|joinapipay|JoinPayMakeOrder|%v", err)
		joinPayMakeOrderRes = nil
	}

	return
}

//JoinPayRefundToUser -汇聚支付订单退款
func JoinPayRefundToUser(order *db.Order, refundno string, amount int64, remark string) (joinPayRefundToUserRes *JoinPayRefundToUserRes, err error) {
	defer profile.TimeTrack(time.Now(), "[JoinPay-API] JoinPayRefundToUser")

	var (
		key *db.Key
	)

	if amount <= 0 {
		err = errors.New("支付的金额有误")
		return
	}

	// 查询个人退款接口请求参数
	if key, err = cache.GetKey(order.Kid); err != nil || key == nil {
		return nil, fmt.Errorf("支付失败：%v", "key不存在")
	}

	// 获取支付环境
	env := wxapi.GetPayEnv(key)
	if env == nil {
		err = fmt.Errorf("kid:%v,keyname:%v,获取支付环境失败", key.Kid, key.Name)
		return
	}

	var slicekey = []string{"p1_MerchantNo", "p2_OrderNo", "p3_RefundOrderNo", "p4_RefundAmount", "p5_RefundReason", "p6_NotifyUrl", "hmac"}

	values := &url.Values{}
	values.Add("p1_MerchantNo", env.Mchid)
	values.Add("p2_OrderNo", order.Orderno)
	values.Add("p3_RefundOrderNo", refundno)
	values.Add("p4_RefundAmount", str.GetMoneyYuan(amount))
	values.Add("p5_RefundReason", remark)
	values.Add("p6_NotifyUrl", env.NotifyURL)
	sign := &Sign{Values: values}
	values.Add("hmac", sign.Sign(env.APIKey, slicekey, "Md5"))

	joinPayRefundToUserRes = &JoinPayRefundToUserRes{}
	err = xhttp.PostJSON(REFUND_URL, values, nil, joinPayRefundToUserRes)

	if err != nil {
		logs.Error("|JoinApi|joinapipay|JoinPayRefundToUser|%v", err)
		joinPayRefundToUserRes = nil
	}
	return
}

//JoinPaySingleRePay -单笔代付金额
func JoinPaySingleRePay(uid int64, receiverAccountNoEnc string, receiverNameEnc string, amount int64, remark string, key *db.Key) (joinPaySingleRePayRes *JoinPaySingleRePayRes, err error) {
	defer profile.TimeTrack(time.Now(), "[JoinPay-API] JoinPaySingleRePay")

	if amount <= 0 {
		err = errors.New("支付的金额有误")
		return
	}

	// 获取支付环境
	env := wxapi.GetPayEnv(key)
	if env == nil {
		err = fmt.Errorf("kid:%v,keyname:%v,获取支付环境失败", key.Kid, key.Name)
		return
	}

	var slicekey = []string{"UserNo", "ProductCode", "RequestTime", "MerchantOrderNo", "ReceiverAccountNoEnc", "ReceiverNameEnc", "ReceiverAccountType", "ReceiverBankChannelNo",
		"PaidAmount", "Currency", "IsChecked", "PaidDesc", "PaidUse", "CallbackUrl", "FirstProductCode", "hmac"}

	values := &url.Values{}
	values.Add("UserNo", env.Mchid)
	values.Add("ProductCode", "BANK_PAY_DAILY_ORDER") //朝夕付 BANK_PAY_DAILY_ORDER  BANK_PAY_MAT_ENDOWMENT_ORDER 任意付
	values.Add("RequestTime", time.Now().Format("2006-01-02 15:04:05"))
	values.Add("MerchantOrderNo", uniqueid.GenerateOrderRefundNo(uid))
	values.Add("ReceiverAccountNoEnc", receiverAccountNoEnc) //银行卡号
	values.Add("ReceiverNameEnc", receiverNameEnc)
	values.Add("ReceiverAccountType", "201")
	//values.Add("ReceiverBankChannelNo", "")
	values.Add("PaidAmount", str.GetMoneyYuan(amount))
	values.Add("Currency", "201")
	values.Add("IsChecked", "202") //是否复核 201复核 202不复核
	values.Add("PaidDesc", remark)
	values.Add("PaidUse", "201")
	values.Add("CallbackUrl", env.NotifyURL)
	values.Add("FirstProductCode", "BANK_PAY_DAILY_ORDER")

	if !vars.IsProd() { // 开发环境 , 测试环境
		values.Set("PaidAmount", "0.01")
	}

	sign := &Sign{Values: values}
	values.Add("hmac", sign.Sign(env.APIKey, slicekey, "Md5"))

	joinPaySingleRePayReq := &JoinPaySingleRePayReq{
		UserNo:                values.Get("UserNo"),
		ProductCode:           values.Get("ProductCode"),
		RequestTime:           values.Get("RequestTime"),
		MerchantOrderNo:       values.Get("MerchantOrderNo"),
		ReceiverAccountNoEnc:  values.Get("ReceiverAccountNoEnc"),
		ReceiverNameEnc:       values.Get("ReceiverNameEnc"),
		ReceiverAccountType:   values.Get("ReceiverAccountType"),
		ReceiverBankChannelNo: values.Get("ReceiverBankChannelNo"),
		PaidAmount:            values.Get("PaidAmount"),
		Currency:              values.Get("Currency"),
		IsChecked:             values.Get("IsChecked"), //是否复核 201复核 202不复核
		PaidDesc:              values.Get("PaidDesc"),
		PaidUse:               values.Get("PaidUse"),
		CallbackUrl:           values.Get("CallbackUrl"),
		FirstProductCode:      values.Get("FirstProductCode"),
		Hmac:                  values.Get("hmac"),
	}

	joinPaySingleRePayRes = &JoinPaySingleRePayRes{}
	err = xhttp.PostJSON(SINGLEPAY_URL, &url.Values{}, joinPaySingleRePayReq, joinPaySingleRePayRes)

	if err != nil {
		logs.Error("|JoinApi|joinapipay|JoinPaySingleRePay|%v", err)
		joinPaySingleRePayRes = nil
	}

	return
}
