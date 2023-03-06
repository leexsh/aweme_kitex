package userKafka

import (
	"aweme_kitex/cmd/user/service_user/db"
	constants "aweme_kitex/pkg/constant"
	"aweme_kitex/pkg/logger"
	"context"
	"strings"

	"github.com/Shopify/sarama"
)

var kafkaProducer sarama.SyncProducer
var kafkaFollowAddConsumer sarama.Consumer
var kafkaFollowDelConsumer sarama.Consumer

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

	kafkaFollowAddConsumer, err = sarama.NewConsumer([]string{constants.KafkaAddress}, nil)
	if err != nil {
		logger.Error(err)
	}

	kafkaFollowDelConsumer, err = sarama.NewConsumer([]string{constants.KafkaAddress}, nil)
	if err != nil {
		logger.Error(err)
	}

	go ConsumeFollowAddMsg()
	go ConsumeFollowDelMsg()
}

// 生产消息
func ProduceFollowMsg(topic, val string) error {
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

// 消费增加关注数目
func ConsumeFollowAddMsg() {
	partitionList, err := kafkaFollowAddConsumer.Partitions(constants.KafKaUserAddRelationTopic)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, list := range partitionList { // 遍历所有分区
		// 根据每个分区创建一个消费者
		pc, err := kafkaFollowAddConsumer.ConsumePartition("follow_add", int32(list), sarama.OffsetNewest)
		if err != nil {
			logger.Error(err)
		}
		defer pc.AsyncClose()
		// 异步消费数据
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				params := strings.Split(string(msg.Value), "&")
				userId, toUserId := params[0], params[1]
				err := db.NewUserDaoInstance().IncreaseFollowCount(context.Background(), userId, toUserId)
				if err != nil {
					logger.Error("create relation err: " + err.Error())
				}
			}
		}(pc)
	}
}

// 消费减少关注数目
func ConsumeFollowDelMsg() {
	partitionList, err := kafkaFollowDelConsumer.Partitions(constants.KafKaUserDelRelationTopic)
	if err != nil {
		logger.Error(err)
		return
	}
	for _, list := range partitionList { // 遍历所有分区
		// 根据每个分区创建一个消费者
		pc, err := kafkaFollowDelConsumer.ConsumePartition("follow_del", int32(list), sarama.OffsetNewest)
		if err != nil {
			logger.Error(err)
		}
		defer pc.AsyncClose()
		// 异步消费数据
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				params := strings.Split(string(msg.Value), "&")
				userId, toUserId := params[0], params[1]
				err := db.NewUserDaoInstance().DecreaseFollowCount(context.Background(), userId, toUserId)
				if err != nil {
					logger.Error("create relation err: " + err.Error())
				}
			}
		}(pc)
	}
}
