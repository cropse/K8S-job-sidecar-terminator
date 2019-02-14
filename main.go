package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// type flagSet struct {
// 	*flag.FlagSet
// }

// func myNewFlagSet(name string, errorHandling flag.ErrorHandling) *myFlagSet {
// 	f := flag.NewFlagSet(name, errorHandling)
// 	newf := &myFlagSet{f}
// 	return newf
// }

// func (f *myFlagSet) Parsed() bool {
// 	return f.parsed
// }

// Fi
func main() {
	// Init parser
	// var erHandling flag.ErrorHandling
	// myFlag := flag.NewFlagSet("flag", erHandling)
	// fmt.Println(myFlag.Parsed())
	flag.Usage = func() {
		// fmt.Fprintf(os.Stderr, "This is not helpful.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "job-terminator [option] command\n")
		flag.PrintDefaults()
	}
	killer := flag.Bool("killer", false, "The one who Kill process")
	// flag.Parse(os.Args[1:])
	flag.Parse()
	fmt.Printf("killer: %t\n", *killer)

	//
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
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get("http://127.0.0.1:8080/kill")
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
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
	cmd.Start()
	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
		go func() {
			time.Sleep(5 * time.Second)
			os.Exit(0)
		}()

	})

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
