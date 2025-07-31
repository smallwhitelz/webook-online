package grpc

import (
	"context"
	"google.golang.org/grpc"
	"webook/api/proto/gen/tag/v1"
	"webook/tag/domain"
	"webook/tag/service"
)

type TagServiceServer struct {
	tagv1.UnimplementedTagServiceServer
	svc service.TagService
}

func NewTagServiceServer(svc service.TagService) *TagServiceServer {
	return &TagServiceServer{svc: svc}
}

func (t *TagServiceServer) Register(server *grpc.Server) {
	tagv1.RegisterTagServiceServer(server, t)
}

func (t *TagServiceServer) CreateTag(ctx context.Context, request *tagv1.CreateTagRequest) (*tagv1.CreateTagResponse, error) {
	tagId, err := t.svc.CreateTag(ctx, request.GetUid(), request.GetName())
	if err != nil {
		return nil, err
	}
	return &tagv1.CreateTagResponse{
		Tag: &tagv1.Tag{
			Id:   tagId,
			Name: request.GetName(),
			Uid:  request.GetUid(),
		},
	}, nil
}

func (t *TagServiceServer) GetTags(ctx context.Context, request *tagv1.GetTagsRequest) (*tagv1.GetTagsResponse, error) {
	tags, err := t.svc.GetTags(ctx, request.GetUid())
	if err != nil {
		return nil, err
	}
	res := make([]*tagv1.Tag, 0, len(tags))
	for _, tag := range tags {
		res = append(res, t.convertToView(tag))
	}
	return &tagv1.GetTagsResponse{
		Tags: res,
	}, nil
}

func (t *TagServiceServer) AttachTags(ctx context.Context, request *tagv1.AttachTagsRequest) (*tagv1.AttachTagsResponse, error) {
	err := t.svc.AttachTags(ctx, request.GetUid(), request.GetBiz(), request.GetBizId(), request.GetTids())
	if err != nil {
		return nil, err
	}
	return &tagv1.AttachTagsResponse{}, err
}

func (t *TagServiceServer) GetBizTags(ctx context.Context, request *tagv1.GetBizTagsRequest) (*tagv1.GetBizTagsResponse, error) {
	tags, err := t.svc.GetBizTags(ctx, request.GetUid(), request.GetBiz(), request.GetBizId())
	if err != nil {
		return nil, err
	}
	res := make([]*tagv1.Tag, 0, len(tags))
	for _, tag := range tags {
		res = append(res, t.convertToView(tag))
	}
	return &tagv1.GetBizTagsResponse{
		Tags: res,
	}, nil
}

func (t *TagServiceServer) convertToView(tag domain.Tag) *tagv1.Tag {
	return &tagv1.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}
