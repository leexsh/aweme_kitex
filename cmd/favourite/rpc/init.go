package favRPC

import favKafka "aweme_kitex/cmd/favourite/service_favourite/kafka"

func InitRPC() {
	initUserRpc()
	initFeedRpc()
	initRelationRpc()
	favKafka.InitKafka()
}
