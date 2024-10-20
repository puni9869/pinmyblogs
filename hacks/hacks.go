package main

import (
	"bufio"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net"
	"strings"
	"time"
)

const PORT = "6652"

func main() {
	log := logger.NewLogger()
	var err error
	l, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.WithError(err).Error("failed to listen.")
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		log.WithError(err).Error("failed to accept the connection request.")
		return
	}
	for {
		netData, err := bufio.NewReader(c).ReadString(byte('\n'))
		if err != nil {
			log.WithError(err).Error("failed to get data from stream.")
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			c.Write([]byte("BYE..."))
			return
		}

		log.Infof("-> %s", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC1123Z) + "\n"

		c.Write([]byte(myTime))
	}
}
