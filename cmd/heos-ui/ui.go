package main

import (
	"fmt"
	"github.com/spaced/heos-cli/pkg/cli"
	"github.com/spaced/heos-cli/pkg/discovery"
	"github.com/spaced/heos-cli/pkg/player"
	"github.com/therecipe/qt/core"
	"log"
	"os"

	"github.com/therecipe/qt/widgets"
)

func createQPlayerGroupBox(c *cli.HeosCli, player player.HeosPlayer) *widgets.QGroupBox {
	currentLevel, _ := player.GetVolume(c)
	currentState, _ := player.GetState(c)

	box := widgets.NewQGroupBox(nil)
	box.SetTitle(player.Name)
	box.SetLayout(widgets.NewQHBoxLayout())

	slider := widgets.NewQSlider2(core.Qt__Horizontal, nil)
	slider.SetTickInterval(10)
	slider.SetSingleStep(1)
	slider.SetValue(currentLevel)
	slider.ConnectSliderMoved(func(level int) {
		log.Printf("Volume of %s to %v ..", player.Name, level)
		player.SetVolume(c, level)
	})
	box.Layout().AddWidget(slider)

	tb := widgets.NewQToolButton(nil)
	if currentState == "stop" {
		tb.SetChecked(false)
		tb.SetIcon(tb.Style().StandardIcon(widgets.QStyle__SP_MediaPlay, nil, nil))
	} else {
		tb.SetChecked(true)
		tb.SetIcon(tb.Style().StandardIcon(widgets.QStyle__SP_MediaStop, nil, nil))
	}
	tb.SetCheckable(true)
	tb.ConnectClicked(func(checked bool) {
		if checked {
			log.Printf("Stop %s ..", player.Name)
			player.Stop(c)
			tb.SetIcon(tb.Style().StandardIcon(widgets.QStyle__SP_MediaPlay, nil, nil))
		} else {
			log.Printf("Play %s ..", player.Name)
			player.Play(c)
			tb.SetIcon(tb.Style().StandardIcon(widgets.QStyle__SP_MediaStop, nil, nil))
		}
	})
	box.Layout().AddWidget(tb)
	return box
}

func createHeosConnection() *cli.HeosCli {
	ips, err := discovery.DiscoverDevices()
	if err != nil {
		log.Fatal("failed to discover", err)
	}
	//grap first discovered device
	if len(ips) == 0 {
		log.Fatal("no devices found")
	}

	c, err := cli.CreateConnection(ips[0])
	if err != nil {
		fmt.Println(err)
	}
	return c
}

func getPlayers(c *cli.HeosCli) []player.HeosPlayer {
	players, err := player.Players(c)
	if err != nil {
		log.Fatal("failed to get players", err)
	}
	return players
}

func main() {

	c := createHeosConnection()

	// needs to be called once before you can start using the QWidgets
	app := widgets.NewQApplication(len(os.Args), os.Args)

	// create a window
	// with a minimum size of 250*200
	// and sets the title to "Hello Widgets Example"
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("Heos Client UI")

	// create a regular widget
	// give it a QVBoxLayout
	// and make it the central widget of the window
	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	for _, p := range getPlayers(c) {
		box := createQPlayerGroupBox(c, p)
		widget.Layout().AddWidget(box)
	}

	// make the window visible
	window.Show()

	// start the main Qt event loop
	// and block until app.Exit() is called
	// or the window is closed by the user
	app.Exec()
}
