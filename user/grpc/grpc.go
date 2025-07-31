package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"webook/api/proto/gen/user/v1"
	domain2 "webook/oauth2/domain"
	"webook/user/domain"
	"webook/user/service"
)

type UserServiceServer struct {
	userv1.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserServiceServer(svc service.UserService) *UserServiceServer {
	return &UserServiceServer{svc: svc}
}

func (u *UserServiceServer) Register(server *grpc.Server) {
	userv1.RegisterUserServiceServer(server, u)
}

func (u *UserServiceServer) Signup(ctx context.Context, request *userv1.SignupRequest) (*userv1.SignupResponse, error) {
	err := u.svc.Signup(ctx, u.convertToDomain(request.GetUser()))
	return &userv1.SignupResponse{}, err
}

func (u *UserServiceServer) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := u.svc.Login(ctx, request.GetEmail(), request.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userv1.LoginResponse{
		User: u.convertToView(user),
	}, nil
}

func (u *UserServiceServer) FindById(ctx context.Context, request *userv1.FindByIdRequest) (*userv1.FindByIdResponse, error) {
	user, err := u.svc.FindById(ctx, request.GetUid())
	if err != nil {
		return nil, err
	}
	return &userv1.FindByIdResponse{
		User: u.convertToView(user),
	}, nil
}

func (u *UserServiceServer) UpdateNonSensitiveInfo(ctx context.Context, request *userv1.UpdateNonSensitiveInfoRequest) (*userv1.UpdateNonSensitiveInfoResponse, error) {
	err := u.svc.UpdateNonSensitiveInfo(ctx, u.convertToDomain(request.GetUser()))
	return &userv1.UpdateNonSensitiveInfoResponse{}, err
}

func (u *UserServiceServer) FindOrCreate(ctx context.Context, request *userv1.FindOrCreateRequest) (*userv1.FindOrCreateResponse, error) {
	user, err := u.svc.FindOrCreate(ctx, request.GetPhone())
	if err != nil {
		return nil, err
	}
	return &userv1.FindOrCreateResponse{
		User: u.convertToView(user),
	}, nil
}

func (u *UserServiceServer) FindOrCreateByWechat(ctx context.Context, request *userv1.FindOrCreateByWechatRequest) (*userv1.FindOrCreateByWechatResponse, error) {
	user, err := u.svc.FindOrCreateByWechat(ctx, domain2.WechatInfo{
		UnionId: request.GetWechatInfo().GetUnionId(),
		OpenId:  request.GetWechatInfo().GetOpenId(),
	})
	if err != nil {
		return nil, err
	}
	return &userv1.FindOrCreateByWechatResponse{
		User: u.convertToView(user),
	}, nil
}

func (u *UserServiceServer) convertToDomain(user *userv1.User) domain.User {
	return domain.User{
		Id:          user.GetId(),
		Email:       user.GetEmail(),
		Password:    user.GetPassword(),
		Nickname:    user.GetNickName(),
		Birthday:    user.GetBirthday().AsTime(),
		Description: user.GetDescription(),
		Phone:       user.GetPhone(),
		Ctime:       user.GetCtime().AsTime(),
		WechatInfo: domain2.WechatInfo{
			UnionId: user.GetWechatInfo().GetUnionId(),
			OpenId:  user.GetWechatInfo().GetOpenId(),
		},
	}
}

func (u *UserServiceServer) convertToView(user domain.User) *userv1.User {
	return &userv1.User{
		Id:          user.Id,
		Email:       user.Email,
		Password:    user.Password,
		NickName:    user.Nickname,
		Description: user.Description,
		Phone:       user.Phone,
		Birthday:    timestamppb.New(user.Birthday),
		Ctime:       timestamppb.New(user.Ctime),
		WechatInfo: &userv1.WechatInfo{
			OpenId:  user.WechatInfo.OpenId,
			UnionId: user.WechatInfo.UnionId,
		},
	}
}
