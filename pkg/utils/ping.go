package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Ping struct {
	average time.Duration
}

func tes() {
	cmd := exec.Command("ping", "8.8.8.8")
	// Linux version
	//cmd := exec.Command("ping", "-c 4", "8.8.8.8")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	printCommand(cmd)
	err := cmd.Run()
	printError(err)
	output := cmdOutput.Bytes()
	printOutput(output)
	ping := Ping{}
	parseOutput(output, &ping)

	fmt.Println(ping)
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func parseOutput(outs []byte, ping *Ping) {
	var average = regexp.MustCompile(`Average = (\d+ms)`)
	result := average.FindStringSubmatch(string(outs))

	if len(result) > 0 {
		ping.average, _ = time.ParseDuration(result[1])
	}
	// Linux version
	/*var average = regexp.MustCompile(`min\/avg\/max\/mdev = (0\.\d+)\/(0\.\d+)\/(0\.\d+)\/(0\.\d+) ms`)
	  result := average.FindAllStringSubmatch(string(outs), -1)

	  if len(result) > 0 {
	  		ping.average, _ = time.ParseDuration(result[0][2] + "ms")
	  }*/
}
