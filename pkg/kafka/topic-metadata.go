package kafka

import (
	"github.com/Shopify/sarama"
)

type TopicDetail struct {
	sarama.TopicDetail
	PartitionAssignment map[int32]int32
	Name                string
}

func (detail *TopicDetail) Native() *sarama.TopicDetail {
	return &detail.TopicDetail
}

func FromTopicMetadata(topic *sarama.TopicMetadata) TopicDetail {
	topicDetail := TopicDetail{
		TopicDetail: sarama.TopicDetail{
			NumPartitions: int32(len(topic.Partitions)),
		},
		Name: topic.Name,
	}
	if len(topic.Partitions) > 0 {
		topicDetail.ReplicaAssignment = make(map[int32][]int32)
		topicDetail.PartitionAssignment = make(map[int32]int32)
		for _, partition := range topic.Partitions {
			topicDetail.ReplicaAssignment[partition.ID] = partition.Replicas
			topicDetail.PartitionAssignment[partition.ID] = partition.Leader
		}
		topicDetail.ReplicationFactor = int16(len(topic.Partitions[0].Replicas))
	}
	return topicDetail
}
