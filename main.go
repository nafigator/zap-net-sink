// Package zap_net_sink provides sink functionality for zap logger.
package zap_net_sink //nolint:staticcheck // Acknowledged. ST1003: should not use underscores in package names.

import (
	"fmt"
	"net"
	"net/url"

	"go.uber.org/zap"
)

// Register udp and tcp urls with zap.
func init() { //nolint:gochecknoinits	// Acknowledged
	if err := zap.RegisterSink("udp", NewUDPSink); err != nil {
		panic(err)
	}

	if err := zap.RegisterSink("tcp", NewTCPSink); err != nil {
		panic(err)
	}
}

// NewUDPSink creates a zap sink to the given url.
func NewUDPSink(url *url.URL) (zap.Sink, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", url.Hostname(), url.Port())) //nolint:govet,noctx // Acknowledged
	if err != nil {
		return nil, fmt.Errorf("failed to setup a UDP sink - %w", err)
	}

	return &WriteSyncer{conn}, nil
}

// NewTCPSink creates a zap sink to the given url.
func NewTCPSink(url *url.URL) (zap.Sink, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", url.Hostname(), url.Port())) //nolint:govet,noctx // Acknowledged
	if err != nil {
		return nil, fmt.Errorf("failed to setup a TCP sink - %w", err)
	}

	return &WriteSyncer{conn}, nil
}

type WriteSyncer struct {
	conn net.Conn
}

func (z *WriteSyncer) Close() error {
	return z.conn.Close()
}

func (z *WriteSyncer) Write(p []byte) (int, error) {
	return z.conn.Write(p)
}

func (z *WriteSyncer) Sync() error {
	return nil
}
