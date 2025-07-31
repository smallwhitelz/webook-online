package web

import (
	"github.com/gin-gonic/gin"
	rewardv1 "webook/api/proto/gen/reward/v1"
	"webook/bff/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

type RewardHandler struct {
	client rewardv1.RewardServiceClient
	l      logger.LoggerV1
}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("reward")
	g.POST("/detail", ginx.WrapBodyAndClaims(h.GetReward))
}

type GetRewardReq struct {
	Rid int64 `json:"rid"`
}

// GetReward 个人看到自己的打赏记录
func (h *RewardHandler) GetReward(ctx *gin.Context, req GetRewardReq, uc jwt.UserClaims) (ginx.Result, error) {
	reward, err := h.client.GetReward(ctx, &rewardv1.GetRewardRequest{
		// 我这一次打赏的id
		Rid: req.Rid,
		// 要防止非法访问，我只能看到我自己的打赏记录
		// 我不能看到别人的
		Uid: uc.Uid,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		// 暂时也就是只需要状态
		Data: reward.Status.String(),
	}, nil
}
