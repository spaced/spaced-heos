package player

import (
	"encoding/json"
	"errors"
	"github.com/spaced/heos-cli/pkg/cli"
	"github.com/spaced/heos-cli/pkg/model"
	"net/url"
	"strconv"
)

type PlayersResponse struct {
	Heos    model.Heos
	Payload []HeosPlayer `json:"payload,omitempty"`
}

type HeosPlayer struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	Network string `json:"network"`
	Ip      string `json:"ip"`
	Pid     int    `json:"pid"`
	Serial  string `json:"serial"`
}

func Players(c *cli.HeosCli) ([]HeosPlayer, error) {
	cmd := cli.HeosCommand{
		Group:   "player",
		Command: "get_players",
		Values:  nil,
	}
	raw, err := c.SendCommand(&cmd)
	if err != nil {
		return nil, err
	}
	response := new(PlayersResponse)
	_ = json.Unmarshal(raw, response)

	if response.Heos.Result != "success" {
		return nil, errors.New(response.Heos.Message)
	}

	return response.Payload, nil
}

func (pl *HeosPlayer) GetState(c *cli.HeosCli) (string, error) {
	q, err := pl.get(c, "get_play_state")
	if err != nil {
		return "unknown", err
	}
	return q.Get("state"), nil
}

func (pl *HeosPlayer) GetVolume(c *cli.HeosCli) (int, error) {
	q, err := pl.get(c, "get_volume")
	if err != nil {
		return 0, err
	}
	level, err := strconv.Atoi(q.Get("level"))
	if err != nil {
		return 0, err
	}
	return level, nil
}

func (pl *HeosPlayer) get(c *cli.HeosCli, command string) (url.Values, error) {
	cmd := cli.HeosCommand{
		Group:   "player",
		Command: command,
		Values:  url.Values{"pid": []string{strconv.Itoa(pl.Pid)}},
	}
	raw, err := c.SendCommand(&cmd)
	if err != nil {
		return nil, err
	}
	response := new(PlayersResponse)
	_ = json.Unmarshal(raw, response)

	if response.Heos.Result != "success" {
		return nil, errors.New(response.Heos.Message)
	}
	q, err := url.ParseQuery(response.Heos.Message)
	if err != nil {
		panic(err)
	}
	return q, nil
}

func (pl *HeosPlayer) Play(c *cli.HeosCli) bool {
	return pl.setPlayState(c, "play")
}

func (pl *HeosPlayer) Stop(c *cli.HeosCli) bool {
	return pl.setPlayState(c, "stop")
}

func (pl *HeosPlayer) Pause(c *cli.HeosCli) bool {
	return pl.setPlayState(c, "pause")
}

func (pl *HeosPlayer) setPlayState(c *cli.HeosCli, state string) bool {
	values := url.Values{}
	values.Add("state", state)
	return pl.set(c, "set_play_state", values)
}

func (pl *HeosPlayer) SetVolume(c *cli.HeosCli, level int) bool {
	if level < 0 || level > 100 {
		return false
	}
	values := url.Values{}
	values.Add("level", strconv.Itoa(level))
	return pl.set(c, "set_volume", values)
}

func (pl *HeosPlayer) volumeUp(c *cli.HeosCli) bool {
	return pl.set(c, "volume_up", url.Values{})
}

func (pl *HeosPlayer) volumeDown(c *cli.HeosCli) bool {
	return pl.set(c, "volume_up", url.Values{})
}

func (pl *HeosPlayer) set(c *cli.HeosCli, command string, values url.Values) bool {
	values.Add("pid", strconv.Itoa(pl.Pid))
	cmd := cli.HeosCommand{
		Group:   "player",
		Command: command,
		Values:  values,
	}
	raw, err := c.SendCommand(&cmd)
	if err != nil {
		return false
	}
	response := new(PlayersResponse)
	_ = json.Unmarshal(raw, response)

	return response.Heos.Result == "success"
}
