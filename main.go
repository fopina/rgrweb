package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fopina/rgrweb/assets"
	gpio "github.com/fopina/rgrweb/helpers"
	flag "github.com/spf13/pflag"
)

func main() {
	var bindAddress = flag.StringP("bind", "b", "127.0.0.1:8081", "address:port to bind webserver")
	var gpioIn = flag.IntP("pin-in", "i", 0, "Input GPIO (feedback) - 0 means fake it")
	var gpioOut = flag.IntP("pin-out", "o", 0, "Output GPIO (trigger) - 0 means fake it")
	var testInput = flag.Bool("test-input", false, "Reading input GPIO for 5 seconds (testing)")
	var testOutput = flag.Bool("test-output", false, "Enable output GPIO for 5 seconds (testing)")
	var highDuration = flag.DurationP("duration", "d", 5*time.Second, "Time that output GPIO pin will be HIGH")
	flag.Parse()

	err := gpio.SetUp(*gpioIn, *gpioOut)
	if err != nil {
		panic(err)
	}
	defer gpio.Cleanup()

	// TODO: capture exit signals to ensure cleanup is done as ^C will skip garbage collection
	/*
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(quit)
	*/

	if *testInput {
		log.Println("Use Ctrl-C to stop...")
		for range time.Tick(time.Second) {
			log.Printf("GPIO%d=%v\n", *gpioIn, gpio.ReadFeedback())
		}
	}

	if *testOutput {
		gpio.SetTrigger(true)
		time.Sleep(*highDuration)
		gpio.SetTrigger(false)
		return
	}

	var resetTimer *time.Timer

	http.HandleFunc("/api/open", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("open request")
		gpio.SetTrigger(true)
		fmt.Fprintf(w, "ok")
		if resetTimer != nil {
			resetTimer.Stop()
		}
		resetTimer = time.NewTimer(*highDuration)

		go func() {
			<-resetTimer.C
			log.Printf("open request - done")
			gpio.SetTrigger(false)
		}()
	})

	http.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%v", gpio.ReadFeedback())
	})

	http.Handle("/", http.FileServer(assets.Assets))
	log.Println("Listening on http://" + *bindAddress)
	log.Fatal(http.ListenAndServe(*bindAddress, nil))
}
