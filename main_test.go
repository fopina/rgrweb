package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fopina/rgrweb/helpers"
)

func ExampleTestInput() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	helpers.SetFakeTrigger(false)
	*testInput = true
	*testOutput = false
	*highDuration = 1 * time.Millisecond
	go runIt()
	time.Sleep(2 * time.Millisecond)
	// Output:
	// Use Ctrl-C to stop...
	// GPIO0=false
	// GPIO0=false
}

func TestTestOutput(t *testing.T) {
	*testInput = false
	*testOutput = true
	*highDuration = 10 * time.Millisecond
	if helpers.ReadFakeFeedback() {
		t.Fatalf("should NOT BE TRUE!")
	}
	go runIt()
	time.Sleep(2 * time.Millisecond)
	if !helpers.ReadFakeFeedback() {
		t.Fatalf("should not be FALSE!")
	}
	time.Sleep(10 * time.Millisecond)
	if helpers.ReadFakeFeedback() {
		t.Fatalf("should not be TRUE!")
	}
}

func TestRunIt(t *testing.T) {
	*testInput = false
	*testOutput = false
	*bindAddress = "127.0.0.1:9999"
	helpers.SetFakeTrigger(false)
	go runIt()

	var resp *http.Response
	var err error

	for {
		// FIXME: how to properly wait for main() to get ready without modifying it...?
		// gief python monkeypatching power
		resp, err = http.Get("http://127.0.0.1:9999/api/check")
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "connection refused") {
			t.Fatalf("failed: %v", err)
		}
		t.Logf("main() not ready yet...")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if string(body) != "false" {
		t.Errorf("check should be false at start: %v", string(body))
	}

	resp, err = http.Get("http://127.0.0.1:9999/api/open")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if string(body) != "ok" {
		t.Fatalf("open not ok: %v", string(body))
	}

	resp, err = http.Get("http://127.0.0.1:9999/api/check")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if string(body) != "true" {
		t.Fatalf("check should be true after open: %v", string(body))
	}
}
