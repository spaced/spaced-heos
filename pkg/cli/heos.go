package cli

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
)

type HeosCli struct {
	conn net.Conn
	buf  *bufio.Reader
}

type HeosCommand struct {
	Group   string
	Command string
	Values  url.Values
}

var HEOS_TELNET_PORT = 1255

func CreateConnection(ip net.IP) (*HeosCli, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", ip, HEOS_TELNET_PORT))
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(conn)
	return &HeosCli{conn: conn, buf: buf}, nil
}

//heos://command_group/command?attribute1=value1&attribute2=value2&...&attributeN=val
func (cli *HeosCli) SendCommand(cmd *HeosCommand) ([]byte, error) {
	u := url.URL{
		Scheme:      "heos",
		Opaque:      "",
		User:        nil,
		Host:        cmd.Group,
		Path:        cmd.Command,
		RawPath:     "",
		ForceQuery:  false,
		RawQuery:    cmd.Values.Encode(),
		Fragment:    "",
		RawFragment: "",
	}
	fmt.Fprintln(cli.conn, u.String())
	return cli.buf.ReadBytes('\n')
}
