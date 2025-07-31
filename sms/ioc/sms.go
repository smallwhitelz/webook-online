package ioc

import (
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"webook/sms/service"
	"webook/sms/service/localsms"
	"webook/sms/service/tencent"
)

// InitSmsTencentService 这里是如果你有腾讯云短信服务可以使用这个初始化
func InitSmsTencentService() service.SmsService {
	// 在这里你也可以考虑从配置文件里面读取
	//secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	//if !ok {
	//	panic("没有找到环境变量 SMS_SECRET_ID ")
	//}
	//secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	//if !ok {
	//	panic("没有找到环境变量 SMS_SECRET_KEY")
	//}
	type Config struct {
		SecretID  string `yaml:"secretId"`
		SecretKey string `yaml:"secretKey"`
	}
	var cfg Config
	err := viper.UnmarshalKey("tencentSms", &cfg)
	c, err := tencentSMS.NewClient(common.NewCredential(cfg.SecretID, cfg.SecretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "妙影科技")
}

// InitSmsService 没有腾讯云短信服务用这个模拟
func InitSmsService() service.SmsService {
	//return initSmsTencentService()
	return InitSmsMemoryService()
}

// InitSmsMemoryService 使用基于内存，输出到控制台的实现
func InitSmsMemoryService() service.SmsService {
	//rateLimitSMSService := ratelimit.NewRateLimitSMSService(localsms.NewService(),
	//	limiter.NewRedisSlidingWindowLimiter(client, time.Second, 1000))
	//prometheus.NewDecorator(rateLimitSMSService, prometheus2.SummaryOpts{})
	return localsms.NewService()
}
