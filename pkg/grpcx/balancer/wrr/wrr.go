package wrr

import (
	"context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"sync"
)

const Name = "custom_weighted_round_robin"

// 在grpc中注册我们的新的实现
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &PickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

// PickerBuilder 平滑的加权轮训算法
type PickerBuilder struct {
}

func (p *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for sc, sci := range info.ReadySCs {
		md, _ := sci.Address.Metadata.(map[string]any)
		weightVal, _ := md["weight"]
		weight, _ := weightVal.(float64)
		conns = append(conns, &weightConn{
			SubConn:       sc,
			weight:        int(weight),
			currentWeight: int(weight),
		})
	}
	return &Picker{
		conns: conns,
	}
}

type Picker struct {
	conns []*weightConn
	lock  sync.Mutex
}

// Pick 真正执行负载均衡的
// 要考虑并发安全
func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 考虑到万一一个节点都没有呢
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var maxCC *weightConn
	for _, c := range p.conns {
		// 计算总权重
		total += c.weight
		// 每一个节点的当前权重加上自己的初始权重
		c.currentWeight = c.currentWeight + c.weight
		// 选出当前权重最大的节点
		if maxCC == nil || maxCC.currentWeight < c.currentWeight {
			maxCC = c
		}
	}
	// 请求后 当前节点的当前权重减去总权重
	maxCC.currentWeight = maxCC.currentWeight - total
	return balancer.PickResult{
		SubConn: maxCC.SubConn,
		// 调用后的回调
		Done: func(info balancer.DoneInfo) {
			// 要在这里进一步调整weight/currentWeight
			// failover 要在这里做文章
			// 根据调用结果的具体错误信息进行容错
			// 1. 如果要是触发了限流了，
			// 1.1 你可以考虑直接挪走这个节点，后面再挪回来
			// 1.2 你可以考虑直接将 weight/currentWeight 调整到极低
			// 2. 触发了熔断呢？
			// 3. 降级呢？
		},
	}, nil
}

// PickV1 动态调整权重
func (p *Picker) PickV1(info balancer.PickInfo) (balancer.PickResult, error) {
	// 考虑到万一一个节点都没有呢
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var res *weightConn
	// 全局锁，性能较低
	//p.lock.Lock()
	//defer p.lock.Unlock()
	for _, c := range p.conns {
		// 会有并发安全问题，分段锁可以提升效率，在操作节点的时候加锁，并发性能会比较好
		c.mutex.Lock()
		// 计算总权重
		total = total + c.efficientWeight
		// 每一个节点的当前权重加上自己的动态权重
		c.currentWeight = c.currentWeight + c.efficientWeight
		// 选出当前权重最大的节点
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}
	// 请求后 当前节点的当前权重减去总权重
	res.mutex.Lock()
	res.currentWeight = res.currentWeight - total
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		// 调用后的回调
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			defer res.mutex.Unlock()
			// 要在这里进一步调整weight/currentWeight
			// 要做一个最低值，不能让动态权重一直减，兜底
			if info.Err != nil && res.efficientWeight == 1 {
				return
			}
			// MaxUint32 可以替换为你认为的最大值。
			// 例如说你预期节点的权重是在 100 - 200 之间
			// 那么你可以设置经过动态调整之后的权重不会超过 500。
			if info.Err == nil && res.efficientWeight >= 500 {
				return
			}
			if info.Err != nil {
				res.efficientWeight--
			} else {
				res.efficientWeight++
			}
		},
	}, nil
}

// PickV2 考虑熔断降级限流的方式
func (p *Picker) PickV2(info balancer.PickInfo) (balancer.PickResult, error) {
	// 考虑到万一一个节点都没有呢
	if len(p.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}

	var total int
	var res *weightConn
	// 全局锁，性能较低
	//p.lock.Lock()
	//defer p.lock.Unlock()
	for _, c := range p.conns {
		// 会有并发安全问题，分段锁可以提升效率，在操作节点的时候加锁，并发性能会比较好
		c.mutex.Lock()
		// 计算总权重
		total = total + c.efficientWeight
		// 每一个节点的当前权重加上自己的动态权重3
		c.currentWeight = c.currentWeight + c.efficientWeight
		// 选出当前权重最大的节点
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}
	// 请求后 当前节点的当前权重减去总权重
	res.mutex.Lock()
	res.currentWeight = res.currentWeight - total
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		// 调用后的回调
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			defer res.mutex.Unlock()
			// 要在这里进一步调整weight/currentWeight
			// 要做一个最低值，不能让动态权重一直减，兜底
			if info.Err != nil && res.efficientWeight == 0 {
				return
			}
			// MaxUint32 可以替换为你认为的最大值。
			// 例如说你预期节点的权重是在 100 - 200 之间
			// 那么你可以设置经过动态调整之后的权重不会超过 500。
			switch info.Err {
			case nil:
				if res.efficientWeight == math.MaxInt32 {
					return
				}
				// 增加权重
				res.efficientWeight++
			case context.DeadlineExceeded:
				// 超时可以考虑动态调整。
				// 比如说第一次超时是降低 1，第二次连续超时是降低 2
				// 因为一次超时可能是偶发，连续超时一定是出了问题
				res.efficientWeight = res.efficientWeight - 10
			default:
				// 检测服务端错误
				code := status.Code(info.Err)
				switch code {
				// 假定我们服务端返回Unavailable代表熔断
				case codes.Unavailable:
					// 直接降低到 1，我们可以预期接下来几乎不会选中它。
					// 但是本身没有降低到 0，所以它又存在被选中的机会，
					// 那么后续会慢慢恢复过来
					res.efficientWeight = 1
					// ResourceExhausted 限流
				case codes.ResourceExhausted:
					// 直接减半，可以快速降低选中该节点的概率
					res.efficientWeight = res.efficientWeight / 2
					// 假定我们服务端返回这个代表降级
				case codes.Aborted:
					// 降级可以考虑和限流采用类似的策略，你也可以调整减少的幅度
					res.efficientWeight = res.efficientWeight / 2
				default:
					if res.efficientWeight == 1 {
						// 降无可降了
						return
					}
					res.efficientWeight--
				}
			}
		},
	}, nil
}

type weightConn struct {
	balancer.SubConn
	mutex sync.Mutex
	// 初始权重
	weight int
	// 当前权重
	currentWeight int

	// 动态调整过后的权重
	efficientWeight int

	// 可以用来标记不可用
	available bool
}
