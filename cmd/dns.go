package cmd

import (
	"fmt"
	"net"
	"net/url"
	"os"
)


func DoDNS(rawURL string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("URL解析错误: %v\n", err)
		os.Exit(1)
	}

	host := u.Hostname()

	// 获取MX记录
	mxRecords, err := net.LookupMX(host)
	if err == nil {
		fmt.Printf("\nMX记录:\n")
		for _, mx := range mxRecords {
			fmt.Printf("Host: %s, Pref: %d\n", mx.Host, mx.Pref)
		}
	}

	// 获取TXT记录
	txtRecords, err := net.LookupTXT(host)
	if err == nil {
		fmt.Printf("\nTXT记录:\n")
		for _, txt := range txtRecords {
			fmt.Println(txt)
		}
	}

	// 获取NS记录
	nsRecords, err := net.LookupNS(host)
	if err == nil {
		fmt.Printf("\nNS记录:\n")
		for _, ns := range nsRecords {
			fmt.Println(ns.Host)
		}
	}

		// 获取IP地址
	ips, err := net.LookupIP(host)
	if err == nil {
		fmt.Printf("\nIP记录:\n")
		for _, ip := range ips {
			fmt.Println(ip)
		}
	}
}
