//
// File: event.go
// Created by Dizzrt on 2023/01/31.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.
//

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"oset/model"
	"strconv"
	"time"

	"github.com/Dizzrt/etlog"
	"github.com/Dizzrt/etstream/kafka"
	"github.com/Dizzrt/go-sse"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	kafkaWrite *kafka.KafkaWriter
	sseServer  *sse.Server
)

func InitEvent() {
	sseServer = sse.NewServer(nil)

	kconfig := kafka.KafkaConfig{
		SaramaConfig: kafka.DefaultProducerConfig(),
		Host:         viper.GetString("kafka.host"),
		Topic:        "events",
	}

	kw, err := kafka.NewKafkaWriter(kconfig, nil, nil)
	if err != nil {
		etlog.L().Panic("failed to create kafka writer", zap.Error(err))
	}

	kafkaWrite = kw
}

func ReportEvent(ctx *gin.Context) {
	event := model.Event{}
	err := ctx.BindJSON(&event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "bind json error",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to receive event, because bind json failed", zap.String("event", event.Event), zap.String("raw_data", event.Data), zap.Error(err))
		return
	}
	event.Time = time.Now()

	said := ctx.Param("aid")
	aid, err := strconv.Atoi(said)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "invalid aid",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to receive event, because invalid aid", zap.String("event", event.Event), zap.String("raw_data", event.Data), zap.String("target_aid", said), zap.Error(err))
		return
	}
	event.Aid = aid

	data := make(map[string]interface{})
	err = json.Unmarshal([]byte(event.Data), &data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "parse data failed",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to receive event, because parse reported data failed", zap.String("event", event.Event), zap.String("raw_data", event.Data), zap.Error(err))
		return
	}

	if event.Did == 0 {
		etlog.L().Warn("report without did (did is 0)", zap.Any("raw_event", event))
	}

	jevent, err := json.Marshal(event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "convert data failed",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to receive event, beacause convert event to json failed", zap.Any("raw_event", event), zap.Error(err))
		return
	}

	sseServer.SendMessage(fmt.Sprintf("/event/tool/realtime/%d/%d", aid, event.Did), sse.SimpleMessage(string(jevent)))

	kafkaWrite.Write(jevent)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func RegisterRealtimeEvent(ctx *gin.Context) {
	said := ctx.Param("aid")
	sdid := ctx.Param("did")

	aid, err := strconv.Atoi(said)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "invalid aid",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to register realtime event service, because invalid aid", zap.String("target_aid", said), zap.String("target_did", sdid), zap.Error(err))
		return
	}

	did, err := strconv.Atoi(sdid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "invalid did",
			"error": err.Error(),
		})
		ctx.Abort()

		etlog.L().Warn("unable to register realtime event service, because invalid did", zap.String("target_aid", said), zap.String("target_did", sdid), zap.Error(err))
		return
	}

	etlog.L().Info("registerd realtime event", zap.Int("aid", aid), zap.Int("did", did))
	sseServer.ServeHTTP(ctx.Writer, ctx.Request)
}
