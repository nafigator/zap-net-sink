package zap_net_sink_test

import (
	"encoding/json"
	"fmt"
	"net"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Example() {
	var err error
	var laddr *net.UDPAddr
	var udpServer *net.UDPConn

	// Setup UDP server
	laddr, err = net.ResolveUDPAddr("udp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}
	udpServer, err = net.ListenUDP("udp", laddr)
	if err != nil {
		panic(err)
	}
	defer func() {
		if defErr := udpServer.Close(); defErr != nil {
			panic(defErr)
		}
	}()

	// Pass all JSON messages that the UDP server receives to a map[string]interface{} channel.
	udpDec := json.NewDecoder(udpServer)
	udpC := make(chan map[string]any, 1)
	go func() {
		for {
			var v map[string]any
			if defErr := udpDec.Decode(&v); defErr != nil {
				return
			}
			udpC <- v
		}
	}()

	writer, cleanup, err := zap.Open("udp://127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer cleanup()

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		writer,
		zap.DebugLevel,
	)).Sugar()
	logger.Infow("hello", "from", "Example()")

	msg := <-udpC
	msg["ts"] = "overridden to make tests pass"
	fmt.Printf("%#v\n", msg)
	// Output:
	// map[string]interface {}{"from":"Example()", "level":"info", "msg":"hello", "ts":"overridden to make tests pass"}
}
