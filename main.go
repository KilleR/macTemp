package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	file, err := os.OpenFile(timeStamp+".txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln("Failed to open file:", err)
	}

	cpuTempRex := regexp.MustCompile("CPU die temperature: ([0-9]{2}.[0-9]{2}) C")

	tmpCmd := exec.Command("powermetrics", "--samplers", "smc,thermal")
	tmpCmd.Stdin = os.Stdin
	//tmpCmd.Stdout = os.Stdout
	tmpCmd.Stderr = os.Stderr
	stdOut, err := tmpCmd.StdoutPipe()
	if err != nil {
		log.Fatalln("failed to get stdout", err)
	}


	err = tmpCmd.Start()
	if err != nil {
		log.Fatalln("failed to start", err)
	}


	scanner := bufio.NewScanner(stdOut)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		matches := cpuTempRex.FindStringSubmatch(m)
		if len(matches) > 0 {
			fmt.Println("MATCH:", matches)
			stamp := time.Now().Format("2006-01-02 15:04:05")
			_, err = file.WriteString(matches[1] + ", " + stamp + "\n")
			if err != nil {
				log.Fatalln("Failed to write:", err)
			}
		}
		fmt.Println(m)
	}
	err = tmpCmd.Wait()
	if err != nil {
		log.Fatalln("failed to wait", err)
	}
}
