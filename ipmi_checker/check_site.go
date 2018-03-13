package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func runCommand(command string, args ...string) ([]byte, error) {
	stdout, err := exec.Command(command, args...).Output()
	return stdout, err
}

func doRequest(hostname string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(fmt.Sprintf("http://ipmi.%s.stch.ru", hostname))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 402 {
		return fmt.Errorf("Bad Status code: %d", resp.StatusCode)
	}
	return nil
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	reqErr := doRequest(hostname)
	if reqErr != nil {
		fmt.Println("request error:", reqErr)

		out, err := runCommand("ipmitool", "mc", "reset", "cold")

		fmt.Println("error:", err)
		fmt.Printf("output: %q\n", out)
	}
}
