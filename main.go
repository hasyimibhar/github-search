package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hasyimibhar/github-search/api"
	"github.com/hasyimibhar/github-search/common"
	"github.com/hasyimibhar/github-search/github"
	"upper.io/db.v3/mysql"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: github-search /path/to/config.yml")
		os.Exit(1)
	}

	config, err := LoadConfig(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to load config file: %s", err))
		os.Exit(1)
	}

	logger := common.NewStandardLogger(os.Stderr, config.LogLevel)

	githubClient := github.NewClient(&github.ClientConfig{
		ClientID:     config.Github.ClientID,
		ClientSecret: config.Github.ClientSecret,
	}, logger)

	dbSession, err := mysql.Open(mysql.ConnectionURL{
		Host:     config.Database.Host,
		Database: config.Database.Database,
		User:     config.Database.User,
		Password: config.Database.Password,
	})
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to initialize database: %s", err))
		os.Exit(1)
	}

	defer dbSession.Close()

	// This is required to fix intermittent invalid connection errors
	dbSession.SetConnMaxLifetime(time.Second * 10)

	apiServer, err := api.NewServer(&api.Config{
		HTTPPort:    config.HTTPPort,
		CORSEnabled: config.CORSEnabled,
	}, githubClient, dbSession, logger)
	if err != nil {
		fmt.Println("failed to start api server:", err)
		os.Exit(1)
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		apiServer.Shutdown(shutdownCtx)
	}()

	sig := waitForSignal()
	logger.Infof("received signal %s", sig)
	logger.Info("shutting down")

	os.Exit(0)
}

func waitForSignal() os.Signal {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	return <-signalCh
}
