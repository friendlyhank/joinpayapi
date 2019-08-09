package joinpayapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

//汇聚支付回调接口
type JoinPayCallBack struct {
	MerchantNo     string `json:"r1_MerchantNo"`
	Openid         string `json:"r2_OrderNo"`
	OrderNo        int64  `json:"r3_Amount"`
	Amount         string `json:"r4_Cur"`
	Mp             string `json:"r5_Mp"`
	Status         string `json:"r6_Status"`
	TrxNo          string `json:"r7_TrxNo"`
	BankOrderNo    string `json:"r8_BankOrderNo"`
	PayTime        string `json:"ra_PayTime"`
	DealTime       string `json:"rb_DealTime"`
	BankCode       string `json:"rc_BankCode"`
	OpenId         string `json:"rd_OpenId"`
	DiscountAmount string `json:"re_DiscountAmount"`
	Hmac           string `json:"hmac"`
}

func PostCallBack(ctx *context.Context) {

	body, _ := ioutil.ReadAll(ctx.Request.Body)

	logs.Info("joinpayapi|CallBack|Body|%v|", string(body))

	joinPayCallBack := &JoinPayCallBack{}
	_ = json.Unmarshal(body, joinPayCallBack)

	fmt.Println(joinPayCallBack)
}

func GetCallBack(ctx *context.Context) {
}
