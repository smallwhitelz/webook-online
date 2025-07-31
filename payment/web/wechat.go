package web

import (
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
	"webook/payment/service/wechat"
	"webook/pkg/logger"
)

type WechatHandler struct {
	// wx提供的处理支付回调通知的方式
	handler   *notify.Handler
	l         logger.LoggerV1
	nativeSvc *wechat.NativePaymentService
}

func (h *WechatHandler) RegisterRoutes(server *gin.Engine) {
	server.Any("/pay/callback", h.HandleNative)
}

func (h *WechatHandler) HandleNative(ctx *gin.Context) {
	// 用来接收解密后的数据
	transaction := new(payments.Transaction)
	_, err := h.handler.ParseNotifyRequest(ctx, ctx.Request, transaction)
	if err != nil {
		ctx.String(http.StatusBadRequest, "参数解析失败")
		h.l.Error("解析微信支付回调失败", logger.Error(err))
		// 在这里可以进一步考虑增加告警和监控
		// 很大概率是黑客在攻击
		return
	}
	// 处理回调
	err = h.nativeSvc.HandleCallback(ctx, transaction)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		// 说明你处理回调失败了
		h.l.Error("处理微信支付回调失败", logger.Error(err),
			logger.String("biz_trade_no", *transaction.OutTradeNo))
		return
	}
	ctx.String(http.StatusOK, "OK")
}
