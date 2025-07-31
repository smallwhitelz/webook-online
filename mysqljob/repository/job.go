package repository

import (
	"context"
	"time"
	"webook/mysqljob/domain"
	"webook/mysqljob/repository/dao"
)

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, time time.Time) error
	AddJob(ctx context.Context, j domain.Job) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

func NewPreemptJobRepository(dao dao.JobDAO) CronJobRepository {
	return &PreemptJobRepository{dao: dao}
}

func (p *PreemptJobRepository) AddJob(ctx context.Context, j domain.Job) error {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.dao.Preempt(ctx)
	return domain.Job{
		Id:         j.Id,
		Expression: j.Expression,
		Executor:   j.Executor,
		Name:       j.Name,
		Cfg:        j.Cfg,
		NextTime:   time.UnixMilli(j.NextTime),
	}, err
}

func (p *PreemptJobRepository) Release(ctx context.Context, jid int64) error {
	return p.dao.Release(ctx, jid)
}

func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptJobRepository) UpdateNextTime(ctx context.Context, id int64, time time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, time)
}

func (p *PreemptJobRepository) toEntity(j domain.Job) dao.Job {
	return dao.Job{
		Id:         j.Id,
		Name:       j.Name,
		Executor:   j.Executor,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		NextTime:   j.NextTime.UnixMilli(),
	}
}
