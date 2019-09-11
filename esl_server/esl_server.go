package main

import (
	esl "canglong/esl_server/goesl"
	"strings"
	"time"
)

func main() {
	for {
		client, err := esl.NewClient("192.168.85.132", 8021, "123qwe", 300)
		if err != nil {
			esl.Error("Error while creating new client: %s", err)
			time.Sleep(10*time.Second)
			continue
		}

		go client.Handle()

		client.Send("events json ALL")

		//client.BgApi(fmt.Sprintf("originate %s %s", "sofia/gateway/kam_registrar/900000100001", "&echo"))


		for {
			msg, err := client.ReadMessage()

			if err != nil {
				client.Close()
				// If it contains EOF, we really dont care...
				if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
					esl.Error("Error while reading Freeswitch message: %s", err)
				}
				esl.Error("Error : %s", err)
				break
			}

			esl.Debug("Got new message: ", msg)
		}

		time.Sleep(10*time.Second)
	}
}
