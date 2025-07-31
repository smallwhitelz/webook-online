package grpc

import (
	"context"
	"encoding/json"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/grpc"
	"time"
	"webook/api/proto/gen/feed/v1"
	"webook/feed/domain"
	"webook/feed/service"
)

type FeedServiceServer struct {
	feedv1.UnimplementedFeedServiceServer
	svc service.FeedService
}

func NewFeedServiceServer(svc service.FeedService) *FeedServiceServer {
	return &FeedServiceServer{svc: svc}
}

func (f *FeedServiceServer) Register(server *grpc.Server) {
	feedv1.RegisterFeedServiceServer(server, f)
}

func (f *FeedServiceServer) CreateFeedEvent(ctx context.Context, request *feedv1.CreateFeedEventRequest) (*feedv1.CreateFeedEventResponse, error) {
	err := f.svc.CreateFeedEvent(ctx, f.convertToDomain(request.GetFeedEvent()))
	return &feedv1.CreateFeedEventResponse{}, err
}

func (f *FeedServiceServer) GetFeedEventList(ctx context.Context, request *feedv1.GetFeedEventListRequest) (*feedv1.GetFeedEventListResponse, error) {
	eventList, err := f.svc.GetFeedEventList(ctx, request.GetUid(), request.GetTimestamp(), request.GetLimit())
	if err != nil {
		return nil, err
	}
	res := slice.Map(eventList, func(idx int, src domain.FeedEvent) *feedv1.FeedEvent {
		return f.convertToView(src)
	})
	return &feedv1.GetFeedEventListResponse{
		FeedEvents: res,
	}, nil
}

func (f *FeedServiceServer) convertToDomain(event *feedv1.FeedEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.GetContent()), &ext)
	return domain.FeedEvent{
		Id:    event.GetId(),
		Uid:   event.GetUid(),
		Type:  event.GetType(),
		Ctime: time.UnixMilli(event.GetCtime()),
		Ext:   ext,
	}
}

func (f *FeedServiceServer) convertToView(event domain.FeedEvent) *feedv1.FeedEvent {
	val, _ := json.Marshal(event.Ext)
	return &feedv1.FeedEvent{
		Id:      event.Id,
		Uid:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.UnixMilli(),
	}
}
