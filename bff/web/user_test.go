package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	codev1 "webook/api/proto/gen/code/v1"
	grpcmocks2 "webook/api/proto/gen/code/v1/mocks"
	userv1 "webook/api/proto/gen/user/v1"
	grpcmocks "webook/api/proto/gen/user/v1/mocks"
	"webook/internal/errs"
	"webook/pkg/ginx"
	service2 "webook/user/service"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		// mock
		mock func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient)

		// 构造请求，预期中的输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期中的输出
		wantCode int
		wantBody ginx.Result
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupRequest{
					User: &userv1.User{
						Email:    "1234@qq.com",
						Password: "hello#world123",
					},
				}).Return(&userv1.SignupResponse{}, nil)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "1234@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Msg: "注册成功",
			},
		},
		{
			name: "Bind出错",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "非法邮箱格式",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: errs.UserInvalidInput,
				Msg:  "非法邮箱格式",
			},
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world3"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: errs.UserInvalidInput,
				Msg:  "两次输入的密码不相等",
			},
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello23",
"confirmPassword": "hello23"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: errs.UserInvalidInput,
				Msg:  "密码必须包含字母、数字、特殊字符",
			},
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupRequest{
					User: &userv1.User{
						Email:    "123@qq.com",
						Password: "hello#world123",
					},
				}).Return(&userv1.SignupResponse{}, errors.New("db 错误"))
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: errs.UserInternalServerError,
				Msg:  "系统错误",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (userv1.UserServiceClient, codev1.CodeServiceClient) {
				userSvc := grpcmocks.NewMockUserServiceClient(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), &userv1.SignupRequest{
					User: &userv1.User{
						Email:    "123@qq.com",
						Password: "hello#world123",
					},
				}).Return(nil, service2.ErrDuplicateEmail)
				codeSvc := grpcmocks2.NewMockCodeServiceClient(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: ginx.Result{
				Code: errs.UserDuplicateEmail,
				Msg:  "邮箱冲突",
			},
		},
	}
	ginx.InitCount(prometheus.CounterOpts{
		Namespace: "geektime_zl",
		Subsystem: "webook",
		Name:      "biz_code",
		Help:      "统计业务错误码	",
	})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造handler
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, nil, codeSvc)
			// 准备服务器和构造路由
			server := gin.Default()
			hdl.RegisterRoutes(server)

			// 准备Req和记录的 recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res ginx.Result
			err := json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantBody, res)
		})
	}
}

func TestUserEmailPattern(t *testing.T) {
	testCases := []struct {
		name  string
		email string
		match bool
	}{
		{
			name:  "不带@",
			email: "123456",
			match: false,
		},
		{
			name:  "带@ 但是没后缀",
			email: "123456@",
			match: false,
		},
		{
			name:  "合法邮箱",
			email: "123456@qq.com",
			match: true,
		},
	}

	h := NewUserHandler(nil, nil, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := h.emailRexExp.MatchString(tc.email)
			require.NoError(t, err)
			assert.Equal(t, tc.match, match)
		})
	}
}

//func TestHTTP(t *testing.T) {
//	// 没有body就传nil
//	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte("我的请求体")))
//	assert.NoError(t, err)
//	recorder := httptest.NewRecorder()
//	assert.Equal(t, http.StatusOK, recorder.Code)
//}
