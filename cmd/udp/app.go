package main

import (
	"fmt"
	"golang.org/x/exp/slog"
	"net"
)

type App struct {
	logger *slog.Logger
}

func (a *App) start() {
	var err error // Explicitly defining as the net package implementation of the error interface was causing issues regarding compilation.

	sourceAddr, err := net.ResolveUDPAddr("udp", opts.Source)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not resolve source address: %s", opts.Source), slog.String("err", err.Error()))
		return
	}

	targetAddr, err := net.ResolveUDPAddr("udp", opts.Target)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not resolve target address: %s", opts.Target), slog.String("err", err.Error()))
		return
	}

	sourceConn, err := net.ListenUDP("udp", sourceAddr)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Could not listen on address: %s", opts.Source), slog.String("err", err.Error()))
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

	a.logger.Info(fmt.Sprintf("Starting udpproxy, Source at %v, Target at %v...", opts.Source, opts.Target))

	for {
		b := make([]byte, opts.Buffer)
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

func NewApp(logger *slog.Logger) *App {
	return &App{
		logger: logger,
	}
}
