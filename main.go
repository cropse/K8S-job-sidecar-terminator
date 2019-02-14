package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Init parser
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "job-terminator [option] command\n")
		flag.PrintDefaults()
	}
	killer := flag.Bool("killer", false, "The one who Kill process")
	flag.Parse()
	fmt.Printf("killer: %t\n", *killer)

	subArgs := flag.Args()
	cmd := exec.Command(subArgs[0], subArgs[1:]...)
	log.Printf("Running command and waiting hook to terminated")
	if *killer == false {
		Receiver(cmd)
	} else {
		Sender(cmd)
	}
}

// Sender is the process who going to trigger the webhook to kill it self when it done.
func Sender(cmd *exec.Cmd) {
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Running result:\n %s\n", out)

	res, err := http.Get("http://127.0.0.1:8080/kill")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", response)
	if string(response) == "OK" {
		fmt.Println("")
	}
}

// Receiver is main terminated process running http server to exit process
func Receiver(cmd *exec.Cmd) {
	// handle command stdout and running
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	// http handler
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
		// exit asyncly for response
		go func() {
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}()

	})
	// running simple http server
	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
	}()

	// shouldn't reach here without error
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
