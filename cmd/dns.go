package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/urfave/cli/v3"
)

func DNSCommand() *cli.Command {
	return &cli.Command{
		Name:  "dns",
		Usage: "Perform DNS lookup",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "url",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			url := c.StringArg("url")
			if url == "" {
				return errors.New("url cannot be empty")
			}

			return DoDNS(url)
		},
	}
}

func DoDNS(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("URL解析错误: %v\n", err)
		return err
	}

	host := u.Hostname()

	// 获取MX记录
	mxRecords, err := net.LookupMX(host)
	if err == nil {
		fmt.Printf("\nMX记录:\n")
		for _, mx := range mxRecords {
			fmt.Printf("Host: %s, Pref: %d\n", mx.Host, mx.Pref)
		}
	} else {
		fmt.Printf("\nError : %v\n", err)
	}

	// 获取TXT记录
	txtRecords, err := net.LookupTXT(host)
	if err == nil {
		fmt.Printf("\nTXT记录:\n")
		for _, txt := range txtRecords {
			fmt.Println(txt)
		}
	} else {
		fmt.Printf("\nError : %v\n", err)
	}

	// 获取NS记录
	nsRecords, err := net.LookupNS(host)
	if err == nil {
		fmt.Printf("\nNS记录:\n")
		for _, ns := range nsRecords {
			fmt.Println(ns.Host)
		}
	} else {
		fmt.Printf("\nError : %v\n", err)
	}

	// 获取IP地址
	ips, err := net.LookupIP(host)
	if err == nil {
		fmt.Printf("\nIP记录:\n")
		for _, ip := range ips {
			fmt.Println(ip)
		}
	} else {
		fmt.Printf("\nError : %v\n", err)
	}

	return nil
}
