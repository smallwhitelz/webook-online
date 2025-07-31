package grpc

import (
	"context"
	"google.golang.org/grpc"
	rewardv1 "webook/api/proto/gen/reward/v1"
	"webook/reward/domain"
	"webook/reward/service"
)

type RewardServiceServer struct {
	rewardv1.UnimplementedRewardServiceServer
	svc service.RewardService
}

func NewRewardServiceServer(svc service.RewardService) *RewardServiceServer {
	return &RewardServiceServer{svc: svc}
}

func (r *RewardServiceServer) Register(s *grpc.Server) {
	rewardv1.RegisterRewardServiceServer(s, r)
}

func (r *RewardServiceServer) GetReward(ctx context.Context, request *rewardv1.GetRewardRequest) (*rewardv1.GetRewardResponse, error) {
	rw, err := r.svc.GetReward(ctx, request.GetRid(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &rewardv1.GetRewardResponse{
		Status: rewardv1.RewardStatus(rw.Status),
	}, nil
}

func (r *RewardServiceServer) UpdateReward(ctx context.Context, request *rewardv1.UpdateRewardRequest) (*rewardv1.UpdateRewardResponse, error) {
	err := r.svc.UpdateReward(ctx, request.GetBizTradeNo(), domain.RewardStatus(request.GetStatus()))
	return &rewardv1.UpdateRewardResponse{}, err
}

func (r *RewardServiceServer) PreReward(ctx context.Context, request *rewardv1.PreRewardRequest) (*rewardv1.PreRewardResponse, error) {
	res, err := r.svc.PreReward(ctx, domain.Reward{
		Biz:       request.GetBiz(),
		BizId:     request.GetBizId(),
		BizName:   request.GetBizName(),
		TargetUid: request.GetTargetUid(),
		Uid:       request.GetUid(),
		Amt:       request.GetAmt(),
	})
	if err != nil {
		return nil, err
	}
	return &rewardv1.PreRewardResponse{
		CodeUrl: res.Url,
		Rid:     res.Rid,
	}, nil
}
