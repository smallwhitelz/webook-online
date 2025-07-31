package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"webook/api/proto/gen/mysqljob/v1"
	"webook/mysqljob/domain"
	"webook/mysqljob/service"
)

type MysqlCronJobServiceServer struct {
	sqljobv1.UnimplementedMysqlCronJobServiceServer
	svc service.CronJobService
}

func NewMysqlCronJobServiceServer(svc service.CronJobService) *MysqlCronJobServiceServer {
	return &MysqlCronJobServiceServer{svc: svc}
}

func (m *MysqlCronJobServiceServer) Register(server *grpc.Server) {
	sqljobv1.RegisterMysqlCronJobServiceServer(server, m)
}

func (m *MysqlCronJobServiceServer) Preempt(ctx context.Context, request *sqljobv1.PreemptRequest) (*sqljobv1.PreemptResponse, error) {
	job, err := m.svc.Preempt(ctx)
	if err != nil {
		return nil, err
	}
	return &sqljobv1.PreemptResponse{
		Job: &sqljobv1.CronJob{
			Id:         job.Id,
			Name:       job.Name,
			Expression: job.Expression,
			Executor:   job.Executor,
			Cfg:        job.Cfg,
			NextTime:   timestamppb.New(job.NextTime),
		},
	}, nil
}

func (m *MysqlCronJobServiceServer) ResetNextTime(ctx context.Context, request *sqljobv1.ResetNextTimeRequest) (*sqljobv1.ResetNextTimeResponse, error) {
	err := m.svc.ResetNextTime(ctx, domain.Job{
		Id:         request.GetJob().GetId(),
		Name:       request.GetJob().GetName(),
		Expression: request.GetJob().GetExpression(),
		Executor:   request.GetJob().GetExecutor(),
		Cfg:        request.GetJob().GetCfg(),
		NextTime:   request.GetJob().GetNextTime().AsTime(),
	})
	return &sqljobv1.ResetNextTimeResponse{}, err
}

func (m *MysqlCronJobServiceServer) AddJob(ctx context.Context, request *sqljobv1.AddJobRequest) (*sqljobv1.AddJobResponse, error) {
	err := m.svc.AddJob(ctx, domain.Job{
		Id:         request.GetJob().GetId(),
		Name:       request.GetJob().GetName(),
		Expression: request.GetJob().GetExpression(),
		Executor:   request.GetJob().GetExecutor(),
		Cfg:        request.GetJob().GetCfg(),
		NextTime:   request.GetJob().GetNextTime().AsTime(),
	})
	return &sqljobv1.AddJobResponse{}, err
}
