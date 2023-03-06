package service_feed

import (
	videoKafka "aweme_kitex/cmd/feed/service_feed/kafka"
	constants "aweme_kitex/pkg/constant"
	"context"
)

type ChangeCommentCountService struct {
	ctx context.Context
}

// NewChangeCommentCountService new ChangeCommentCountService
func NewChangeCommentCountService(ctx context.Context) *ChangeCommentCountService {
	return &ChangeCommentCountService{
		ctx: ctx,
	}
}

func (c *ChangeCommentCountService) ChangeCommentCount(vid string, action int64) error {
	// 用kafka去写 防止热点视频/时间打崩数据库
	if action == 1 {
		err := videoKafka.ProduceFollowMsg(constants.KafKaVideoCommentAddTopic, vid)
		if err != nil {
			return err
		}
	} else if action == 2 {
		err := videoKafka.ProduceFollowMsg(constants.KafKaFavouriteDelTopic, vid)
		if err != nil {
			return err
		}
	}
	return nil
}
