package repository

import (
	"context"
	"webook/article/domain"
)

// HistoryRecordRepository 阅读历史
type HistoryRecordRepository interface {
	AddRecord(ctx context.Context, record domain.HistoryRecord) error
}
