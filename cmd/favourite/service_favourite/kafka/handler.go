package favKafka

import (
	"aweme_kitex/cmd/favourite/service_favourite/db"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"aweme_kitex/pkg/types"
	"context"
	"strings"

	"github.com/Shopify/sarama"
)

var kafkaProducer sarama.SyncProducer
var kafkaAddConsumer sarama.Consumer
var kafkaDelConsumer sarama.Consumer

func InitKafka() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	var err error
	kafkaProducer, err = sarama.NewSyncProducer([]string{constants.KafkaAddress}, config)
	if err != nil {
		logger.Error(err)
	}

	kafkaAddConsumer, err = sarama.NewConsumer([]string{constants.KafkaAddress}, nil)
	if err != nil {
		logger.Error(err)
	}
	kafkaDelConsumer, err = sarama.NewConsumer([]string{constants.KafkaAddress}, nil)
	if err != nil {
		logger.Error(err)
	}

	go ConsumeAddRelation()
	go ConsumeDelRelation()
}

// 生产消息
func ProduceAddRelation(topic, val string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(val),
	}
	_, _, err := kafkaProducer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

// 消费取消关注消息
func ConsumeAddRelation() {
	partitionList, err := kafkaAddConsumer.Partitions(constants.KafKaFavouriteAddTopic)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, list := range partitionList { // 遍历所有分区
		// 根据每个分区创建一个消费者
		pc, err := kafkaAddConsumer.ConsumePartition("favourite_add", int32(list), sarama.OffsetNewest)
		if err != nil {
			logger.Error(err)
		}
		defer pc.AsyncClose()
		// 异步消费数据
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				params := strings.Split(string(msg.Value), "&")
				id, userId, vid := params[0], params[1], params[2]
				favour := &types.FavouriteRaw{
					Id:      id,
					UserId:  userId,
					VideoId: vid,
				}
				err := db.NewFavouriteDaoInstance().CreateFavour(context.Background(), favour, vid)
				if err != nil {
					logger.Error(err)
				}
			}
		}(pc)
	}
}

func ConsumeDelRelation() {
	partitionList, err := kafkaDelConsumer.Partitions(constants.KafKaFavouriteDelTopic)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, list := range partitionList { // 遍历所有分区
		// 根据每个分区创建一个消费者
		pc, err := kafkaDelConsumer.ConsumePartition("favourite_del", int32(list), sarama.OffsetNewest)
		if err != nil {
			logger.Error(err)
		}
		defer pc.AsyncClose()
		// 异步消费数据
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				params := strings.Split(string(msg.Value), "&")
				userId, vid := params[0], params[1]
				err := db.NewFavouriteDaoInstance().DelFavour(context.Background(), userId, vid)
				if err != nil {
					logger.Error(err)
				}
			}
		}(pc)
	}
}
