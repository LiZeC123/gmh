package cmd

import (
	"fmt"
	"net"
	"time"
)

/*



Probing 43.159.60.42:1234/tcp - No response - time=2011.608ms
Probing 43.159.60.42:1234/tcp - No response - time=2015.934ms
Probing 43.159.60.42:1234/tcp - No response - time=2006.609ms
Probing 43.159.60.42:1234/tcp - No response - time=2000.199ms

Ping statistics for 43.159.60.42:1234
     4 probes sent.
     0 successful, 4 failed.  (100.00% fail)
Was unable to connect, cannot provide trip statistics.


Probing 43.159.60.42:443/tcp - Port is open - time=70.494ms
Probing 43.159.60.42:443/tcp - Port is open - time=80.277ms
Probing 43.159.60.42:443/tcp - Port is open - time=1080.888ms
Probing 43.159.60.42:443/tcp - Port is open - time=78.111ms

Ping statistics for 43.159.60.42:443
     4 probes sent.
     4 successful, 0 failed.  (0.00% fail)
Approximate trip times in milli-seconds:
     Minimum = 70.494ms, Maximum = 1080.888ms, Average = 327.443ms



*/

func Tcping(host string, port string, timeout int) {
	target := net.JoinHostPort(host, port)

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
