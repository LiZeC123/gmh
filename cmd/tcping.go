package cmd

import (
	"context"
	"fmt"
	"github.com/LiZeC123/gmh/util"
	"net"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

const defaultTimeout = 1

func TcpingCommand() *cli.Command {
	return &cli.Command{
		Name:  "tcping",
		Usage: "Probe TCP port connectivity",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "host",
				Min:  1,
				Max:  -1,
			},
		},
		Flags: []cli.Flag{
			&cli.Uint16Flag{
				Name:     "port",
				Aliases:  []string{"p"},
				Required: true,
			},
			&cli.Uint8Flag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Value:    defaultTimeout,
				Required: false,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			hosts := c.StringArgs("host")

			port := c.Uint16("port")
			timeout := c.Uint8("timeout")

			for _, host := range hosts {
				err := Tcping(host, port, timeout)
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}

func Tcping(host string, port uint16, timeout uint8) error {
	target := net.JoinHostPort(host, strconv.Itoa(int(port)))

	total, fail := 0, 0
	for range 4 {
		rst, err := doOneConnect(target, time.Duration(timeout)*time.Second)
		if err != nil {
			fmt.Printf("Probing %v - No response - time=%v (err=%v)\n", target, rst, err)
			fail++
		} else {
			fmt.Printf("Probing %v - Port is open - time=%v\n", target, rst)
		}
		total++
	}
	fmt.Printf("Ping statistics for %v\n", target)
	fmt.Printf("\t %v probes sent.\n", total)
	fmt.Printf("\t %v successful, %v failed.  (%v%% fail)\n\n", total-fail, fail, float32(fail*100)/float32(total))

	return nil
}

func doOneConnect(target string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	duration := time.Since(start)
	if err == nil {
		util.CloseWithLog(conn)
	}

	return duration, err
}
