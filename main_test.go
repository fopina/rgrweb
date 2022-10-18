package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/fopina/rgrweb/helpers"
)

func assertError(t *testing.T, err error, message string) {
	if err.Error() != message {
		t.Fatalf("unexpected err: %v", err)
	}
}

func reset() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	// needed to reset HandleFunc registry across tests
	http.DefaultServeMux = http.NewServeMux()

	*tokenFile = ""
	*tokens = []string{}
	*noAuthRequired = false
	*testInput = false
	*testOutput = false
	helpers.SetFakeTrigger(false)
}

func DISABLED_Example_test_input() {
	reset()
	*testInput = true
	*highDuration = 6 * time.Millisecond
	go runIt()
	time.Sleep(9 * time.Millisecond)
	// Output:
	// Use Ctrl-C to stop...
	// GPIO0=false
}

func TestTestOutput(t *testing.T) {
	reset()
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

func TestRunIt_TokenAndAuth(t *testing.T) {
	reset()
	// hacky way to stop runIt for "succesful auth setup" == make it break after the checks
	*bindAddress = "invalid"
	err := runIt()
	assertError(t, err, "either token-file, token or no-auth need to be specified")
	reset()
	*bindAddress = "invalid"
	*noAuthRequired = true
	err = runIt()
	assertError(t, err, "listen tcp: address invalid: missing port in address")
	reset()
	f, err := ioutil.TempFile("", "testtokenfile")
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("a:12345\nb:54321"), 0644)
	*bindAddress = "invalid"
	*tokenFile = f.Name()
	err = runIt()
	assertError(t, err, "listen tcp: address invalid: missing port in address")
	reset()
	*bindAddress = "invalid"
	*tokens = []string{"asd"}
	err = runIt()
	assertError(t, err, "invalid token format asd")
	*tokens = []string{"name:token"}
	err = runIt()
	assertError(t, err, "listen tcp: address invalid: missing port in address")
}

func waitForIt(t *testing.T) {
	var err error

	for {
		// FIXME: how to properly wait for main() to get ready without modifying it...?
		// gief python monkeypatching power
		_, err = http.Get("http://127.0.0.1:9999/")
		if err == nil {
			return
		}
		if !strings.Contains(err.Error(), "connection refused") {
			t.Fatal(err)
		}
		t.Log("http server not ready yet...")
	}
}

func TestRunIt_NoAuth(t *testing.T) {
	reset()
	*bindAddress = "127.0.0.1:9999"
	*noAuthRequired = true
	go runIt()

	waitForIt(t)

	resp, err := http.Get("http://127.0.0.1:9999/api/check")
	if err != nil {
		t.Fatalf("failed: %v", err)
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
	// hacky, make sure timer runs out and log is completed
	time.Sleep(*highDuration)
}

func TestRunIt_WithBadAuth(t *testing.T) {
	reset()
	*bindAddress = "127.0.0.1:9999"
	*tokens = []string{"a:12345", "b:98765"}
	go runIt()

	waitForIt(t)

	resp, err := http.Get("http://127.0.0.1:9999/api/check")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if string(body) != "invalid token\n" {
		t.Fatalf("/api/check without token did not return error: %v", string(body))
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

	if string(body) != "invalid token\n" {
		t.Fatalf("/api/open without token did not return error: %v", string(body))
	}

	if helpers.ReadFakeFeedback() {
		t.Fatalf("gpio should not be true after invalid open request")
	}

	// try just one more time with an actual token but invalid
	req, err := http.NewRequest("GET", "http://127.0.0.1:9999/api/check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Token", "1234")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if string(body) != "invalid token\n" {
		t.Fatalf("/api/check with invalid token did not return error: %v", string(body))
	}
}

func TestRunIt_WithGoodAuth(t *testing.T) {
	reset()
	*bindAddress = "127.0.0.1:9999"
	*tokens = []string{"a:12345", "b:98765"}
	go runIt()

	waitForIt(t)

	req, err := http.NewRequest("GET", "http://127.0.0.1:9999/api/check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Token", "12345")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if string(body) != "false" {
		t.Errorf("check should be false at start: %v", string(body))
	}

	req, err = http.NewRequest("GET", "http://127.0.0.1:9999/api/open", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Token", "12345")
	resp, err = http.DefaultClient.Do(req)
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

	req, err = http.NewRequest("GET", "http://127.0.0.1:9999/api/check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Token", "12345")
	resp, err = http.DefaultClient.Do(req)
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
	// hacky, make sure timer runs out and log is completed
	time.Sleep(*highDuration)
}
