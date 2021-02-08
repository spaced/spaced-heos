package main

import (
	"fmt"
	"github.com/spaced/heos-cli/pkg/cli"
	"github.com/spaced/heos-cli/pkg/discovery"
	"github.com/spaced/heos-cli/pkg/player"
)

func main() {
	ips, err := discovery.DiscoverDevices()
	if err != nil {
		return
	}
	//grap first discovered device
	if len(ips) > 0 {

		c, err := cli.CreateConnection(ips[0])
		if err != nil {
			fmt.Println(err)
		}

		players, err := player.Players(c)
		for _, p := range players {
			if p.Name == "Küche" {
				s, _ := p.GetState(c)
				if s == "stop" {
					fmt.Println("Speaker was stopped, play now..")
					p.Play(c)

				}
				level, _ := p.GetVolume(c)
				fmt.Printf("Current volume: %v", level)
			}
		}

	}
}
