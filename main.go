package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fopina/rgrweb/assets"
	gpio "github.com/fopina/rgrweb/helpers"
	flag "github.com/spf13/pflag"
)

// top level for easier testing
var bindAddress = flag.StringP("bind", "b", "127.0.0.1:8081", "address:port to bind webserver")
var gpioIn = flag.IntP("pin-in", "i", 0, "Input GPIO (feedback) - 0 means fake it")
var gpioOut = flag.IntP("pin-out", "o", 0, "Output GPIO (trigger) - 0 means fake it")
var testInput = flag.Bool("test-input", false, "Reading input GPIO for 5 seconds (testing)")
var testOutput = flag.Bool("test-output", false, "Enable output GPIO for 5 seconds (testing)")
var noAuthRequired = flag.Bool("no-auth", false, "Start without any tokens configured (hopefully testing only!)")
var tokens = flag.StringArrayP("token", "t", nil, "Tokens required for authentication, format LOGGING_ID:TOKEN")
var tokenFile = flag.String("token-file", "", "File containing list of tokens, format LOGGING_ID:TOKEN, one per line")
var highDuration = flag.DurationP("duration", "d", 5*time.Second, "Time that output GPIO pin will be HIGH")
var generateToken = flag.Bool("generate-token", false, "Helper to generate random token")

func runIt() error {
	err := gpio.SetUp(*gpioIn, *gpioOut)
	if err != nil {
		return err
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
		for range time.Tick(*highDuration) {
			log.Printf("GPIO%d=%v\n", *gpioIn, gpio.ReadFeedback())
		}
	}

	if *testOutput {
		gpio.SetTrigger(true)
		time.Sleep(*highDuration)
		gpio.SetTrigger(false)
		return nil
	}

	if *generateToken {
		const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
		rand.Seed(time.Now().UnixNano())
		b := make([]byte, 16)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}
		fmt.Println(string(b))
		return nil
	}

	if *tokenFile != "" {
		file, err := os.Open(*tokenFile)
		if err != nil {
			return err
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			*tokens = append(*tokens, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	if !*noAuthRequired && len(*tokens) == 0 {
		return fmt.Errorf("either token-file, token or no-auth need to be specified")
	}

	var tokenEntries = make(map[string]string)
	var tokenValue []string

	for _, token := range *tokens {
		tokenValue = strings.Split(token, ":")
		if len(tokenValue) != 2 {
			return fmt.Errorf("invalid token format %v", token)
		}
		tokenEntries[tokenValue[1]] = tokenValue[0]
	}

	var resetTimer *time.Timer

	http.HandleFunc("/api/open", func(w http.ResponseWriter, r *http.Request) {
		if *noAuthRequired {
			log.Printf("open request")
		} else {
			if tokenName, ok := tokenEntries[r.Header.Get("X-Token")]; ok {
				log.Printf("open request for %v", tokenName)
			} else {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
		}
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
		if !*noAuthRequired {
			if _, ok := tokenEntries[r.Header.Get("X-Token")]; !ok {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
		}
		fmt.Fprintf(w, "%v", gpio.ReadFeedback())
	})

	http.Handle("/", http.FileServer(assets.Assets))
	log.Println("Listening on http://" + *bindAddress)
	return http.ListenAndServe(*bindAddress, nil)
}

func main() {
	flag.Parse()
	err := runIt()
	if err != nil {
		log.Fatal(err)
	}
}
