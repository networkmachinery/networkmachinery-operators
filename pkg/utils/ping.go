package utils

import (
	"regexp"
	"time"
)

type Ping struct {
	min, average, max time.Duration
}

func (p *Ping) Min() string {
	return p.min.String()
}

func (p *Ping) Average() string {
	return p.average.String()
}

func (p *Ping) Max() string {
	return p.max.String()
}

func ParsePingOutput(outs []byte, ping *Ping) {
	//var average = regexp.MustCompile(`min\/avg\/max = (0\.\d+)\/(0\.\d+)\/(0\.\d+) ms`)
	var average = regexp.MustCompile(`min\/avg\/max = (\d+\.\d+)\/(\d+.\d+)\/(\d+.\d+) ms`)
	result := average.FindAllStringSubmatch(string(outs), -1)
	if len(result) > 0 {
		ping.min, _ = time.ParseDuration(result[0][1] + "ms")
		ping.average, _ = time.ParseDuration(result[0][2] + "ms")
		ping.max, _ = time.ParseDuration(result[0][3] + "ms")
	}
}
