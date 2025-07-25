package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func sendRequest(wg *sync.WaitGroup) {
	defer wg.Done()

	url := "http://localhost:8001/transfer"
	payload := []byte(`{
		"merchant_id": 1,
		"amount": 10000,
		"account_number": 1,
		"simulate_success": true
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Response:", resp.Status, string(body))
}

func main() {
	rps := 10
	duration := 10 * time.Second
	ticker := time.NewTicker(time.Second / time.Duration(rps))

	var wg sync.WaitGroup

	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		wg.Add(1)
		go sendRequest(&wg)
		<-ticker.C
	}

	wg.Wait()
	fmt.Println("Burst test completed")
}
