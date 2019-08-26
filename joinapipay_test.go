package joinpayapi

import (
	"fmt"
	"testing"

	"git.biezao.com/ant/xmiss/foundation/cache"

	"git.biezao.com/ant/xmiss/foundation"
	"git.biezao.com/ant/xmiss/foundation/db"
)

//订单支付测试
func TestJoinPayMakeOrder(t *testing.T) {
	//初始化配置
	foundation.Init()

	user := &db.User{Uid: 15614, Openid: "o_9Tw0OD6SYcMlb059RFW0YitKRU"}
	key, _ := cache.GetKey(3)
	orderno := "103409782611613087"
	body := "商城商品"
	producetname := "雪梨子"
	joinPayMakeOrderRes, _ := JoinPayMakeOrder(user, body, orderno, producetname, 1, key)
	fmt.Println(joinPayMakeOrderRes)
}

//订单支付测试
func TestJoinPayRefundToUser(t *testing.T) {
	//初始化配置
	foundation.Init()

	JoinPayRefundToUser()
}

//订单支付测试
func TestJoinPaySingleRePay(t *testing.T) {
	joinPaySingleRePayRes, _ := JoinPaySingleRePay()
	fmt.Println(joinPaySingleRePayRes)
}
