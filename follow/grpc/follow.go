package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/follow/v1"
	"webook/follow/domain"
	"webook/follow/service"
)

type FollowRelationServiceServer struct {
	followv1.UnimplementedFollowServiceServer
	svc service.FollowService
}

func NewFollowRelationServiceServer(svc service.FollowService) *FollowRelationServiceServer {
	return &FollowRelationServiceServer{svc: svc}
}

func (f *FollowRelationServiceServer) Register(server *grpc.Server) {
	followv1.RegisterFollowServiceServer(server, f)
}

func (f *FollowRelationServiceServer) Follow(ctx context.Context, request *followv1.FollowRequest) (*followv1.FollowResponse, error) {
	err := f.svc.Follow(ctx, request.GetFollower(), request.GetFollowee())
	return &followv1.FollowResponse{}, err
}

func (f *FollowRelationServiceServer) CancelFollow(ctx context.Context, request *followv1.CancelFollowRequest) (*followv1.CancelFollowResponse, error) {
	err := f.svc.CancelFollow(ctx, request.GetFollower(), request.GetFollowee())
	return &followv1.CancelFollowResponse{}, err
}

func (f *FollowRelationServiceServer) GetFollowee(ctx context.Context, request *followv1.GetFolloweeRequest) (*followv1.GetFolloweeResponse, error) {
	dmFollowee, err := f.svc.GetFollowee(ctx, request.GetFollower(), request.GetOffset(), request.GetLimit())
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(dmFollowee))
	for _, relation := range dmFollowee {
		res = append(res, f.convertToView(relation))
	}
	return &followv1.GetFolloweeResponse{
		FollowRelations: res,
	}, nil
}

func (f *FollowRelationServiceServer) FollowInfo(ctx context.Context, request *followv1.FollowInfoRequest) (*followv1.FollowInfoResponse, error) {
	ri, err := f.svc.FollowInfo(ctx, request.GetFollower(), request.GetFollowee())
	if err != nil {
		return nil, err
	}
	return &followv1.FollowInfoResponse{
		FollowRelation: f.convertToView(ri),
	}, nil
}

func (f *FollowRelationServiceServer) GetFollower(ctx context.Context, request *followv1.GetFollowerRequest) (*followv1.GetFollowerResponse, error) {
	fe, err := f.svc.GetFollower(ctx, request.GetFollowee(), request.GetOffset(), request.GetLimit())
	if err != nil {
		return nil, err
	}
	res := make([]*followv1.FollowRelation, 0, len(fe))
	for _, relation := range fe {
		res = append(res, f.convertToView(relation))
	}
	return &followv1.GetFollowerResponse{
		FollowRelation: res,
	}, nil
}

func (f *FollowRelationServiceServer) GetFollowStatic(ctx context.Context, request *followv1.GetFollowStaticRequest) (*followv1.GetFollowStaticResponse, error) {
	fs, err := f.svc.GetFollowStatic(ctx, request.GetUid())
	if err != nil {
		return nil, err
	}
	return &followv1.GetFollowStaticResponse{
		FollowStatic: &followv1.FollowStatic{
			Followers: fs.Followers,
			Followees: fs.Followees,
		},
	}, nil
}

func (f *FollowRelationServiceServer) convertToView(relation domain.FollowRelation) *followv1.FollowRelation {
	return &followv1.FollowRelation{
		Follower: relation.Follower,
		Followee: relation.Followee,
	}
}
