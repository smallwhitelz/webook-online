package service

//改起来非常简单，这里就先不动了
//就是调用方法，放入适合的参数以及返回适合的grpc值即可
//玩到这里这种操作基本没有任何难度了
//func TestBatchRankingService_TopN(t *testing.T) {
//	const batchSize = 2
//	now := time.Now()
//	testCases := []struct {
//		name string
//		mock func(ctrl *gomock.Controller) (intrv1.InteractiveServiceClient, articlev1.ArticleServiceClient)
//
//		wantArts []domain.Article
//		wantErr  error
//	}{
//		{
//			name: "成功获取",
//			mock: func(ctrl *gomock.Controller) (intrv1.InteractiveServiceClient, articlev1.ArticleServiceClient) {
//				intrSvc := grpcmocks.NewMockInteractiveServiceClient(ctrl)
//				artSvc := grpcartmocks.NewMockArticleServiceClient(ctrl)
//				// 先模拟批量获取数据
//				// 先模拟第一批
//				artSvc.EXPECT().ListPub(gomock.Any(), &articlev1.ListPubRequest{
//					StartTime: timestamppb.New(time.Now()),
//					Offset:    0,
//					Limit:     2,
//				}).
//					Return([]domain.Article{
//						{Id: 1, Utime: now},
//						{Id: 2, Utime: now},
//					}, nil)
//				// 模拟第二批
//				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 2, 2).
//					Return([]domain.Article{
//						{Id: 3, Utime: now},
//						{Id: 4, Utime: now},
//					}, nil)
//				// 模拟第三批
//				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 4, 2).
//					// 没数据了
//					Return([]domain.Article{}, nil)
//
//				// 第一批的点赞数据
//				intrSvc.EXPECT().GetByIds(gomock.Any(), &intrv1.GetByIdsRequest{
//					Biz: "article",
//					Ids: []int64{1, 2},
//				}, gomock.Any()).Return(&intrv1.GetByIdsResponse{
//					Intrs: map[int64]*intrv1.Interactive{
//						1: {LikeCnt: 1},
//						2: {LikeCnt: 2},
//					},
//				}, nil)
//				// 第二批的点赞数据
//				intrSvc.EXPECT().GetByIds(gomock.Any(), &intrv1.GetByIdsRequest{
//					Biz: "article",
//					Ids: []int64{3, 4},
//				}, gomock.Any()).
//					Return(&intrv1.GetByIdsResponse{
//						Intrs: map[int64]*intrv1.Interactive{
//							3: {LikeCnt: 3},
//							4: {LikeCnt: 4},
//						},
//					}, nil)
//				return intrSvc, artSvc
//			},
//			wantErr: nil,
//			wantArts: []domain.Article{
//				{Id: 4, Utime: now},
//				{Id: 3, Utime: now},
//				{Id: 2, Utime: now},
//			},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			intrSvc, artSvc := tc.mock(ctrl)
//			svc := &BatchRankingService{
//				intrSvc:   intrSvc,
//				artSvc:    artSvc,
//				n:         3,
//				batchSize: batchSize,
//				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
//					return float64(likeCnt)
//				},
//			}
//			arts, err := svc.topN(context.Background())
//			assert.Equal(t, tc.wantErr, err)
//			assert.Equal(t, tc.wantArts, arts)
//		})
//	}
//}
