package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"encoding/hex"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"

	"gopkg.in/resty.v1"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 10*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")
)

var ids = make(chan string, 1)

func loop() {
	for {
		id := <-ids
		resp, err := resty.R().Get("http://omaraa.ddns.net:62027/db/beacons/"+id)
		if err!=nil {
			panic(err)
		}
		if resp.StatusCode() == 404 {
			var ans string
			fmt.Printf("device not found. add %s(y/n)?", id)
			fmt.Scanf("%s",&ans)

			if ans == "y" {
				putDevice(id)
			}
		}
	}
}

func putDevice(id string) {
	_, err := resty.R().
		SetBody(Beacon{
		ID: id,
	}).
		Put("http://omaraa.ddns.net:62027/db/beacons/"+id)

	if err!=nil {
		panic(err)
	}

	fmt.Println("Added!")
}

func main() {
	flag.Parse()

	d, err := linux.NewDevice()
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	go loop()

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, advHandler, nil))
}

var devices = map[string]bool{}

type Beacon struct {
	ID string `json:"_id"`
}

func advHandler(a ble.Advertisement) {
	if len(a.ServiceData())> 0 {
		data := a.ServiceData()[0].Data


		if len(data) > 0 {
			if data[0] & 0x0F == 0x02 {
				id := hex.EncodeToString(data[1:9])
				if !devices[id] {
					devices[id] = true
					ids <- id
				}

			}
		}
	}
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Println("devices:",devices)
	default:
		log.Fatalf(err.Error())
	}
}
