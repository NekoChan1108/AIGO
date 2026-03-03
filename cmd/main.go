package main

import (
	"AIGO/config"
	"AIGO/internal/logic"
	"AIGO/internal/router"
	"AIGO/pkg/log"
	"AIGO/pkg/mq"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	errChan := make(chan error, 2)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		if err := logic.Consume(ctx, mq.KafkaConsumer); err != nil {
			errChan <- fmt.Errorf("consume message failed: %w", err)
		}
	}()
	go func() {
		if err := router.Router.Run(fmt.Sprintf("0.0.0.0%s", config.Cfg.AppCfg.Port)); err != nil {
			errChan <- fmt.Errorf("run router failed: %w", err)
		}
	}()
	select {
	case err := <-errChan:
		log.Error(err)
		cancelFunc()
		os.Exit(0)
	case sig := <-sigChan:
		cancelFunc()
		log.Infof("receive signal %s", sig.String())
		os.Exit(0)
	case <-ctx.Done():
		log.Info("context done")
		os.Exit(0)
	}
}

func init() {
	isEnabled := true
	isColorEnabled := true
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Runtime Caller failed")
	}
	bannerPath := filepath.Join(filepath.Dir(path), "banner.txt")
	bannerData, err := os.ReadFile(bannerPath)
	if err != nil {
		log.Panic(fmt.Errorf("read banner file failed: %w", err))
	}
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(string(bannerData)))
}
