package main

import (
	"flag"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/jamesandariese/reflux"
)

var influxStatName = flag.String("measurement", "", "Name of measurement to be submitted")
var influxStatTranslation = flag.Bool("negative-unknown", true, "Unknown is -1 instead of 3")

func translateStatus(status int) int {
	if *influxStatTranslation && status == 3 {
		return -1
	}
	return status
}

func main() {
	reflux.PrepareFlags("nagios")
	flag.Parse()

	if *influxStatName == "" {
		log.Fatalln("A measurement is required")
	}

	if len(flag.Args()) == 0 {
		log.Println("A command to run is required")
		os.Exit(3)
	}
	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	out, err := cmd.CombinedOutput()
	// Use up the error if it's *not* an ExitError
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			log.Println("Couldn't run nagios command:", err)
			os.Exit(3)
		}
	}
	outReader := bytes.NewBuffer(out)
	if n, err := io.Copy(os.Stdout, outReader); err != nil {
		log.Printf("Failed copying command output at byte %d: %v", n, err)
	}

	if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		reflux.SendPointUsingFlags(*influxStatName, map[string]interface{}{"status": translateStatus(ws.ExitStatus())})
		fmt.Printf("Exit code: %d\n", ws.ExitStatus())
	} else {
		if err != nil {
			// If err is nil, exit code was the equivalent of 0
			log.Println("Could not determine exit code")
			os.Exit(3)
		}
	}
}
