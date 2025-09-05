# Webook Online - 微服务架构的内容平台
## 项目概述
Webook Online 是一个基于微服务架构的现代内容平台，采用 Go 语言开发，融合了 DDD（领域驱动设计）思想。该平台提供用户内容创作、社交互动、内容发现等核心功能，通过微服务拆分实现高可用、高扩展性的系统架构。

## 核心功能
* 用户系统：账户管理、用户信息、OAuth2 认证
* 内容系统：文章创作、评论互动、标签管理
* 社交系统：关注关系、动态流、互动数据
* 搜索系统：内容搜索与发现
* 互动系统：点赞、收藏、分享
* 奖励系统：用户激励与货币化
* 分布式任务调度系统：基于MySQL的抢占式任务调度
* 数据迁移：不停机数据迁移
## 技术栈
* 后端：Go + Gin + Gorm + gRPC
* 存储：MySQL + Redis + MongoDB + OSS
* 消息队列：Kafka
* 服务发现：etcd
* 依赖注入：Wire
* 容器化：Docker + Kubernetes
* 可观测性：Prometheus + Grafana + Zipkin + ELK 
* 其他中间件或开源工具：Canal、OpenIM
## 项目结构
```
webook-online/
├── account/         # 账户微服务
├── article/         # 文章微服务
├── bff/             # Backend For Frontend 层
├── code/            # 验证码微服务
├── comment/         # 评论微服务
├── config/          # 配置文件
├── feed/            # 动态流微服务
├── follow/          # 关注关系微服务
├── interactive/     # 互动微服务
├── oauth2/          # OAuth2 认证微服务
├── payment/         # 支付微服务
├── pkg/             # 公共工具包
├── ranking/         # 排行榜微服务
├── reward/          # 奖励微服务
├── search/          # 搜索微服务
├── sms/             # 短信微服务
├── tag/             # 标签微服务
└── user/            # 用户微服务
```
* 分布式任务调度系统在mysqljob下
* 数据迁移在pkg/migrator下
## 快速开始
### 环境要求
* Go 1.20+ 
* Docker 与 Docker Compose 
* Kubernetes (可选，用于生产环境)
### 本地开发环境搭建
1. 克隆项目
```bash
git clone https://github.com/smallwhitelz/webook-online.git
cd webook-online
```
2. 启动依赖服务
```bash
docker-compose up -d mysql8 redis mongo etcd kafka elasticsearch
```
**PS:** 这里需要哪些依赖服务都在docker-compose.yaml文件中，这里的启动只是举例

3. 运行服务
```bash
# 启动单个微服务(以用户服务为例)
cd user
# 每个服务都会有一些不一样的配置，至少端口一定是不同的，所以要加上--config
go run . --config=config/dev.yaml 
```
### 使用Docker部署
```bash
# 构建镜像
docker build -t webook:latest .
 
# 运行容器
docker run -p 8080:8080 webook:latest
```
### Kubernetes部署
```bash
kubectl apply -f webook-deployment.yaml
kubectl apply -f webook-service.yaml
```
## 监控与可观测性
项目集成了完整的监控体系：
* Prometheus + Grafana：性能指标监控 
* Zipkin：分布式追踪 
* ELK：日志收集与分析

**PS:** 这里可以详细查看user服务，里面集成了详细的grpc观测手段