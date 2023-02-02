//
// File: kafka.go
// Created by Dizzrt on 2023/02/02.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package stream

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
)

var (
	logProducer sarama.SyncProducer
	// etProducer  sarama.SyncProducer
)

func InitKafka() {
	host := viper.GetString("kafka.host")
	port := viper.GetString("kafka.port")

	logProducerConfig := sarama.NewConfig()
	logProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	logProducerConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	logProducerConfig.Producer.Return.Successes = true
	logProducerConfig.Version = sarama.V3_2_2_0

	logProducer_, err := sarama.NewSyncProducer([]string{fmt.Sprintf("%s:%s", host, port)}, logProducerConfig)
	if err != nil {
		panic("producer closed, err: " + err.Error())
	}

	logProducer = logProducer_
}

func SendLog(log string) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = "log_proxy"
	msg.Value = sarama.StringEncoder(log)

	_, _, err := logProducer.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err: " + err.Error())
	}
}
