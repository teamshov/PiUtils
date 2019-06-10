package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"

	"gopkg.in/resty.v1"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 100*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")
)

var ids = make(chan string, 1)

func loop() {
	for {
		id := <-ids
		fmt.Printf("\nID: " + id)
		resp, err := resty.R().Get("http://omaraa.ddns.net:62027/db/beacons/" + id)
		fmt.Println(resp.StatusCode())
		if err != nil {
			panic(err)
		}
		if resp.StatusCode() != 200 {
			var ans string
			fmt.Printf("device not found. add %s(y/n)?", id)
			fmt.Scanf("%s", &ans)

			if ans == "y" {
				var x float32
				var y float32
				fmt.Printf("Enter x coordinate: ")
				fmt.Scanf("%f", &x)
				fmt.Printf("Enter y coordinate: ")
				fmt.Scanf("%f", &y)
				putDevice(id, x, y)
			}
		}
	}
}

func putDevice(id string, x float32, y float32) {
	_, err := resty.R().
		SetBody(Beacon{
			ID: id,
			building: "eb2",
			floor: "L1",
			temp: "0.0",
			xpos: x,
			ypos: y,
		}).
		Put("http://omaraa.ddns.net:62027/db/beacons/" + id)

	if err != nil {
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

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	go loop()
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.Scan(ctx, *dup, advHandler, nil))
}

var devices = map[string]bool{}

type Beacon struct {
	ID string `json:"_id"`
	building string `json:"building"`
	floor string `json:"floor"`
	temp string `json:"temp"`
	xpos float32 `json:"xpos"`
	ypos float32 `json:"ypos"`
}

func advHandler(a ble.Advertisement) {
	if len(a.ServiceData()) > 0 {
		data := a.ServiceData()[0].Data

		if len(data) > 0 {
			if data[0]&0x0F == 0x02 {
				//fmt.Println(data)
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
		fmt.Println("devices:", devices)
	default:
		log.Fatalf(err.Error())
	}
}
