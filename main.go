// File: main.go
// Created by Dizzrt on 2023/01/17.
//
// Copyright (C) 2023 The oset Authors.
// This source code is licensed under the MIT license found in
// the LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"oset/common/oset"
	"oset/router"
	"syscall"
	"time"

	"github.com/Dizzrt/etlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func gracefullyShutdown(server http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println()
		etlog.L().Warn("server shutdown error", zap.String("err", err.Error()))
	}
}

func main() {
	oset.Init()
	defer oset.Defer()

	// init router
	r := gin.Default()
	router.CollectRoutes(r)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	// start listening and serving
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("server listen err: " + err.Error())
		}
	}()

	// waiting for shutdown signal
	gracefullyShutdown(server)
}
