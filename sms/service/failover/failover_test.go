package failover

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"webook/sms/service"
)

func TestFailOverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name  string
		mocks func(ctrl *gomock.Controller) []service.SmsService

		wantErr error
	}{
		{
			name: "一次发送成功",
			mocks: func(ctrl *gomock.Controller) []service.SmsService {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []service.SmsService{
					svc0,
				}
			},
		},
		{
			name: "第二次发送成功",
			mocks: func(ctrl *gomock.Controller) []service.SmsService {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []service.SmsService{
					svc0,
					svc1,
				}
			},
		},
		{
			name: "全部失败",
			mocks: func(ctrl *gomock.Controller) []service.SmsService {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("发送失败"))
				return []service.SmsService{
					svc0,
					svc1,
				}
			},
			wantErr: errors.New("轮询了所有服务商，但是发送都失败了"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailOverSMSService(tc.mocks(ctrl))
			err := svc.Send(context.Background(), "abc", []string{"123"}, "123456")
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
