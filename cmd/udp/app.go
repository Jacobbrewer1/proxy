package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jacobbrewer1/reverse-proxy/cmd/udp/config"
	"github.com/jacobbrewer1/reverse-proxy/pkg/dataacess"
	"golang.org/x/exp/slog"
	"net"
	"time"
)

const redisTargetKey string = "udp_proxy_target"
const defaultBufferSize int = 10240

type App struct {
	logger *slog.Logger
	cfg    *config.Config
}

func (a *App) getTargets() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	db := dataacess.NewRedisDalWithContext(ctx, 0)
	if db == nil {
		a.logger.Error("Redis client came back nil")
		return "", nil
	}
	const errMsg string = "Error fetching url from redis"
	target, err := db.GetValue(redisTargetKey)
	if errors.Is(err, redis.Nil) {
		a.logger.Error(errMsg, slog.String("err", fmt.Sprintf("redis key %s was not found", redisTargetKey)))
	} else if err != nil {
		a.logger.Error(errMsg, slog.String("err", err.Error()))
		return "", err
	}
	return target, err
}

func (a *App) start() {
	var err error // Explicitly defining as the net package implementation of the error interface was causing issues regarding compilation.

	listeningPort := fmt.Sprintf(":%s", a.cfg.ListeningPort)
	target, err := a.getTargets()
	if err != nil {
		a.logger.Error("error retrieving targets", slog.String("err", err.Error()))
		return
	}

	sourceAddr, err := net.ResolveUDPAddr("udp", listeningPort)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not resolve source address: %s", listeningPort), slog.String("err", err.Error()))
		return
	}

	targetAddr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not resolve target address: %s", target), slog.String("err", err.Error()))
		return
	}

	sourceConn, err := net.ListenUDP("udp", sourceAddr)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not listen on address: %s", a.cfg.ListeningPort), slog.String("err", err.Error()))
		return
	}
	defer func() {
		if err := sourceConn.Close(); err != nil {
			a.logger.Error("Error closing source connection", slog.String("err", err.Error()))
		}
	}()

	targetConn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not connect to target address: %s", targetAddr), slog.String("err", err.Error()))
		return
	}
	defer func(c *net.UDPConn) {
		if err := c.Close(); err != nil {
			a.logger.Error("Error closing target connection", slog.String("err", err.Error()))
		}
	}(targetConn)

	a.logger.Info(fmt.Sprintf("Starting %s, source at %v, target at %v...", config.AppName, listeningPort, target))

	for {
		b := make([]byte, defaultBufferSize)
		n, addr, err := sourceConn.ReadFromUDP(b)
		if err != nil {
			a.logger.Error("Could not receive a packet", slog.String("err", err.Error()))
			continue
		}
		a.logger.Debug("Packet received",
			slog.String("address", addr.String()),
			slog.Int("num_of_bytes", n),
			slog.String("content", string(b)),
		)
		if _, err := targetConn.Write(b[0:n]); err != nil {
			a.logger.Error("Could not forward packet", slog.String("err", err.Error()))
		}
	}
}

func NewApp(logger *slog.Logger, cfg *config.Config) *App {
	return &App{
		logger: logger,
		cfg:    cfg,
	}
}
