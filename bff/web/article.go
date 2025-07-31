package web

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
	articlev1 "webook/api/proto/gen/article/v1"
	commentv1 "webook/api/proto/gen/comment/v1"
	intrv1 "webook/api/proto/gen/intr/v1"
	rewardv1 "webook/api/proto/gen/reward/v1"
	"webook/bff/web/jwt"
	"webook/pkg/ginx"
	"webook/pkg/logger"
)

type ArticleHandler struct {
	svc     articlev1.ArticleServiceClient
	intrSvc intrv1.InteractiveServiceClient
	reward  rewardv1.RewardServiceClient
	cmtSvc  commentv1.CommentServiceClient
	l       logger.LoggerV1
	biz     string
}

func NewArticleHandler(svc articlev1.ArticleServiceClient, l logger.LoggerV1,
	intrSvc intrv1.InteractiveServiceClient, reward rewardv1.RewardServiceClient, cmtSvc commentv1.CommentServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		intrSvc: intrSvc,
		reward:  reward,
		cmtSvc:  cmtSvc,
		biz:     "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	// 新增和修改接口
	g.POST("/edit", ginx.WrapBodyAndClaims(h.Edit))

	// 发表接口
	g.POST("/publish", ginx.WrapBodyAndClaims(h.Publish))
	g.POST("/withdraw", ginx.WrapBodyAndClaims(h.Withdraw))

	// 创作者接口
	// 返回文章详情
	g.GET("/detail/:id", h.Detail)
	// 按道理来说，这边就是get方法
	// /list?offset=?&limit=?
	// 返回文章列表
	g.POST("/list", h.List)

	// 读者接口
	pub := g.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	// 传入一个参数，true就是点赞，false就是不点赞
	pub.POST("/like", ginx.WrapBodyAndClaims(h.Like))
	pub.POST("/collect", ginx.WrapBodyAndClaims(h.Collect))
	// 打赏
	pub.POST("/reward", ginx.WrapBodyAndClaims(h.Reward))
	// 评论
	pub.POST("/addComment", ginx.WrapBodyAndClaims[AddComment, jwt.UserClaims](h.AddComment))
	pub.POST("/getCommentList", ginx.WrapBody[GetCommentListReq](h.GetCommentList))
	pub.POST("/getMoreReplies", ginx.WrapBody[GetMoreRepliesReq](h.GetMoreReplies))
}

// Edit 接收Article 输入，返回一个ID，文章ID
func (h *ArticleHandler) Edit(ctx *gin.Context, req ArticleEditReq, uc jwt.UserClaims) (ginx.Result, error) {
	resp, err := h.svc.Save(ctx, &articlev1.SaveRequest{
		Article: &articlev1.Article{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Author: &articlev1.Author{
				Id: uc.Uid,
			},
		},
	})
	if err != nil {
		h.l.Error("保存文章数据失败",
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Data: resp.GetId(),
	}, nil
}

func (h *ArticleHandler) Publish(ctx *gin.Context, req PublishReq, uc jwt.UserClaims) (ginx.Result, error) {
	resp, err := h.svc.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Title,
			Author: &articlev1.Author{
				Id: uc.Uid,
			},
		},
	})
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, fmt.Errorf("发表文章失败 aid %d, uid %d %w", req.Id, uc.Uid, err)
	}
	return ginx.Result{
		Data: resp.GetId(),
	}, nil
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context, req ArticleWithdrawReq, uc jwt.UserClaims) (ginx.Result, error) {
	_, err := h.svc.Withdraw(ctx, &articlev1.WithdrawRequest{
		Uid: uc.Uid,
		Aid: req.Id,
	})
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.l.Warn("查询文章失败，id格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}
	resp, err := h.svc.GetById(ctx, &articlev1.GetByIdRequest{
		Aid: id,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("查询文章失败",
			logger.Int64("id", id),
			logger.Error(err))
		return
	}
	article := resp.GetArticle()
	uc := ctx.MustGet("user").(jwt.UserClaims)
	if article.GetAuthor().GetId() != uc.Uid {
		// 有人在搞鬼
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("非法查询",
			logger.Int64("id", id),
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}
	vo := ArticleVo{
		Id:    article.GetId(),
		Title: article.GetTitle(),
		//Abstract: art.Abstract(),
		Content:  article.GetContent(),
		AuthorId: article.GetAuthor().GetId(),

		Status: uint8(article.GetStatus()),
		Ctime:  article.GetCtime().AsTime().Local().Format(time.DateTime),
		Utime:  article.GetUtime().AsTime().Local().Format(time.DateTime),
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: vo,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page
	if err := ctx.Bind(&page); err != nil {
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	resp, err := h.svc.GetByAuthor(ctx, &articlev1.GetByAuthorRequest{
		Uid:    uc.Uid,
		Offset: int32(page.Offset),
		Limit:  int32(page.Limit),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("查找文章列表失败",
			logger.Error(err),
			logger.Int("offset", page.Offset),
			logger.Int("limit", page.Limit),
			logger.Int64("uid", uc.Uid))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slice.Map[*articlev1.Article, ArticleVo](resp.GetArticles(), func(idx int, src *articlev1.Article) ArticleVo {
			return ArticleVo{
				Id:    src.GetId(),
				Title: src.GetTitle(),
				// 摘要
				Abstract: h.Abstract(src.GetContent()),
				//Content: src.Content,
				// 创作者列表 不需要显示创作者
				Status: uint8(src.GetStatus()),
				Ctime:  src.GetCtime().AsTime().Local().Format(time.DateTime),
				Utime:  src.GetUtime().AsTime().Local().Format(time.DateTime),
			}
		}),
	})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.l.Warn("查询文章失败，id格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)
	var (
		eg       errgroup.Group
		artResp  *articlev1.GetPubByIdResponse
		intrResp *intrv1.GetResponse
	)

	// 在引入goroutine的时候，一定要重新传一个超时控制
	// 个人遇见的问题：只开了bff、user、art的server，没有开intr
	// 结果导致一直请求，没有返回，随后才发现没起intr的server
	eg.Go(func() error {
		var er error
		artCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 这里引入kafka，在获取单个文章详情的时候
		// 就发送一条消息到kafka，从而增加阅读数
		artResp, er = h.svc.GetPubById(artCtx, &articlev1.GetPubByIdRequest{
			Aid: id,
			Uid: uc.Uid,
		})
		return er
	})

	eg.Go(func() error {

		var er error
		intrCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		intrResp, er = h.intrSvc.Get(intrCtx, &intrv1.GetRequest{
			Biz:   h.biz,
			BizId: id,
			Uid:   uc.Uid,
		})
		return er
	})

	// 等待结果
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.l.Error("查询文章失败，系统错误",
			logger.Int64("aid", id),
			logger.Int64("uid", uc.Uid),
			logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: ArticleVo{
			Id:         artResp.GetArticle().GetId(),
			Title:      artResp.GetArticle().GetTitle(),
			Content:    artResp.GetArticle().GetContent(),
			AuthorId:   artResp.GetArticle().GetAuthor().GetId(),
			AuthorName: artResp.GetArticle().GetAuthor().GetName(),
			ReadCnt:    intrResp.GetIntr().GetReadCnt(),
			CollectCnt: intrResp.GetIntr().GetCollectCnt(),
			LikeCnt:    intrResp.GetIntr().GetLikeCnt(),
			Liked:      intrResp.GetIntr().GetLiked(),
			Collected:  intrResp.GetIntr().GetCollected(),

			Status: uint8(artResp.GetArticle().GetStatus()),
			Ctime:  artResp.GetArticle().GetCtime().AsTime().Local().Format(time.DateTime),
			Utime:  artResp.GetArticle().GetUtime().AsTime().Local().Format(time.DateTime),
		},
	})
}

func (h *ArticleHandler) Like(ctx *gin.Context, req ArticleLikeReq, uc jwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		// 点赞
		_, err = h.intrSvc.Like(ctx, &intrv1.LikeRequest{
			Biz:   h.biz,
			BizId: req.Id,
			Uid:   uc.Uid,
		})
	} else {
		// 取消点赞
		_, err = h.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz: h.biz,
			Id:  req.Id,
			Uid: uc.Uid,
		})
	}
	if err != nil {
		h.l.Error("点赞/取消点赞失败",
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id),
			logger.Error(err))
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *ArticleHandler) Collect(ctx *gin.Context, req ArticleCollectReq, uc jwt.UserClaims) (ginx.Result, error) {
	_, err := h.intrSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz:   h.biz,
		BizId: req.Id,
		Cid:   req.Cid,
		Uid:   uc.Uid,
	})
	if err != nil {
		h.l.Error("收藏失败",
			logger.Int64("uid", uc.Uid),
			logger.Int64("aid", req.Id),
			logger.Error(err))
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	return ginx.Result{
		Msg: "OK",
	}, nil
}

func (h *ArticleHandler) Reward(ctx *gin.Context, req ArticleRewardReq, uc jwt.UserClaims) (ginx.Result, error) {
	// 如果实现了多种支付方式，例如wx和支付宝
	// 那么可以在这里分发
	//h.reward.Wechat
	//h.reward.Alipay
	artResp, err := h.svc.GetPubById(ctx, &articlev1.GetPubByIdRequest{
		Aid: req.Id,
		Uid: uc.Uid,
	})
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	article := artResp.GetArticle()
	resp, err := h.reward.PreReward(ctx, &rewardv1.PreRewardRequest{
		Biz:       "article",
		BizId:     article.GetId(),
		BizName:   article.GetTitle(),
		TargetUid: article.GetAuthor().GetId(),
		Uid:       uc.Uid,
		Amt:       req.Amt,
	})
	if err != nil {
		return ginx.Result{
			Msg:  "系统错误",
			Code: 5,
		}, err
	}
	return ginx.Result{
		Data: map[string]any{
			"codeUrl": resp.CodeUrl,
			"rid":     resp.Rid,
		},
	}, nil
}

func (h *ArticleHandler) AddComment(ctx *gin.Context, req AddComment, uc jwt.UserClaims) (ginx.Result, error) {
	// 构建gRPC请求
	grpcReq := &commentv1.CreateCommentRequest{
		Comment: &commentv1.Comment{
			Uid:     uc.Uid,
			Biz:     req.Biz,
			BizId:   req.BizId,
			Content: req.Content,
		},
	}

	// 设置根评论/父评论ID
	if req.RootID > 0 {
		grpcReq.Comment.RootComment = &commentv1.Comment{Id: req.RootID}
	}
	if req.ParentID > 0 {
		grpcReq.Comment.ParentComment = &commentv1.Comment{Id: req.ParentID}
	}

	_, err := h.cmtSvc.CreateComment(ctx, grpcReq)
	if err != nil {
		return ginx.Result{
			Msg:  "评论失败",
			Code: 5,
		}, err
	}
	return ginx.Result{Msg: "评论成功"}, nil
}

func (h *ArticleHandler) GetCommentList(ctx *gin.Context, req GetCommentListReq) (ginx.Result, error) {
	resp, err := h.cmtSvc.GetCommentList(ctx, &commentv1.GetCommentListRequest{
		Biz:   req.Biz,
		BizId: req.BizId,
		MinId: req.MinId,
		Limit: req.Limit,
	})
	if err != nil {
		return ginx.Result{
			Msg:  "获取评论列表失败",
			Code: 5,
		}, err
	}
	vo := h.commentsToVo(resp.Comments)
	fmt.Println(vo)
	return ginx.Result{
		Data: vo,
	}, nil
}

func (h *ArticleHandler) GetMoreReplies(ctx *gin.Context, req GetMoreRepliesReq) (ginx.Result, error) {
	resp, err := h.cmtSvc.GetMoreReplies(ctx, &commentv1.GetMoreRepliesRequest{
		Rid:   req.Rid,
		MaxId: req.MaxId,
		Limit: req.Limit,
	})
	if err != nil {
		return ginx.Result{
			Msg:  "获取回复失败",
			Code: 5,
		}, err
	}
	return ginx.Result{Data: h.commentsToVo(resp.Replies)}, nil
}

func (h *ArticleHandler) commentsToVo(comments []*commentv1.Comment) []CommentVo {
	return slice.Map(comments, func(idx int, src *commentv1.Comment) CommentVo {
		return CommentVo{
			Id: src.GetId(),
			Commentator: User{
				Id:   src.GetUid(),
				Name: "", // 需要后续补充用户名
			},
			Biz:     src.GetBiz(),
			BizId:   src.GetBizId(),
			Content: src.GetContent(),
			// 注意，这里这样写会导致时间和正确的时间相差8小时
			Ctime:    src.GetCtime().AsTime().Local().Format(time.DateTime),
			Utime:    src.GetUtime().AsTime().Local().Format(time.DateTime),
			ParentID: src.GetParentComment().GetId(),
		}
	})
}

func (h *ArticleHandler) Abstract(content string) string {
	str := []rune(content)
	// 只取部分作为摘要
	if len(str) > 128 {
		str = str[:128]
	}
	return string(str)
}
