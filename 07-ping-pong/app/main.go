package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	podName     string
	podIndex    int
	namespace   string
	serviceName string
)

func main() {
	podName = os.Getenv("POD_NAME")
	if podName == "" {
		log.Fatal("POD_NAME environment variable not set")
	}

	namespace = os.Getenv("NAMESPACE")
	if namespace == "" {
		log.Fatal("NAMESPACE environment variable not set")
	}

	serviceName = os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		log.Fatal("SERVICE_NAME environment variable not set")
	}

	parts := strings.Split(podName, "-")
	indexStr := parts[len(parts)-1]
	var err error
	podIndex, err = strconv.Atoi(indexStr)
	if err != nil {
		log.Fatalf("Failed to parse pod index from pod name '%s': %v", podName, err)
	}

	log.Printf("Starting pod %s (index %d) in namespace %s.", podName, podIndex, namespace)

	http.HandleFunc("/ping", pingHandler)
	if podIndex == 0 {
		http.HandleFunc("/finish", finishHandler)
	}

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	time.Sleep(2 * time.Second)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}
	message := string(body)
	log.Printf("Received ping with message: '%s'", message)

	if podIndex > 0 {
		if podIndex%2 == 0 {
			message += " ping"
		} else {
			message += " pong"
		}
	} else {
		message = "ping"
	}

	nextIndex := podIndex + 1

	err = callPod(nextIndex, "/ping", message)

	if err != nil {
		log.Printf("Could not reach next pod (index %d). Assuming I am the last pod in the chain.", nextIndex)
		log.Printf("Calling back to pod-0 with final message: '%s'", message)

		// Call back to pod-0's /finish endpoint.
		// In a real app, you might add retries here.
		callPod(0, "/finish", message)
	}

	fmt.Fprintln(w, "OK")
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading finish request body: %v", err)
		return
	}
	log.Println("--- PING PONG SEQUENCE COMPLETE ---")
	log.Printf("Final message: %s\n", string(body))
	log.Println("---------------------------------")
	fmt.Fprintln(w, "Cycle complete.")
}

func callPod(targetIndex int, endpoint, message string) error {
	targetPodName := fmt.Sprintf("%s-%d", strings.Join(strings.Split(podName, "-")[:len(strings.Split(podName, "-"))-1], "-"), targetIndex)
	url := fmt.Sprintf("http://%s.%s.%s.svc.cluster.local:8080%s", targetPodName, serviceName, namespace, endpoint)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	log.Printf("Making POST request to %s", url)
	resp, err := client.Post(url, "text/plain", strings.NewReader(message))
	if err != nil {
		log.Printf("Request to %s failed: %v", url, err)
		return err // Signal failure
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		errMsg := fmt.Errorf("request to %s returned non-2xx status: %s", url, resp.Status)
		log.Println(errMsg)
		return errMsg // Signal failure
	}

	log.Printf("Successfully called pod %d, status: %s", targetIndex, resp.Status)
	return nil // Signal success
}
