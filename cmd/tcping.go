package cmd

import (
	"fmt"
	"net"
	"strconv"
	"time"
)



func Tcping(host string, port uint16, timeout uint8) error {
	target := net.JoinHostPort(host, strconv.Itoa(int(port)))

	total, fail := 0, 0
	for i := 0; i < 4; i++ {
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
	fmt.Printf("\t %v successful, %v failed.  (%v%% fail)\n", total-fail, fail, float32(fail*100)/float32(total))

	return nil
}

func doOneConnect(target string, timeout time.Duration) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	duration := time.Since(start)
	if err == nil {
		conn.Close()
	}

	return duration, err
}
