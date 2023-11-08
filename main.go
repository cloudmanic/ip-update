package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/digitalocean/godo"

	_ "github.com/jpfuentes2/go-env/autoload"
)

const checkInterval = 1 * time.Minute

func main() {
	for {
		// Get the current public IP address.
		ip, err := getPublicIP()

		if err != nil {
			fmt.Println("Error getting public IP:", err)
			continue
		}

		// Update the DNS record with the new IP address (if it has changed)
		success, err := updateDNSRecord(ip)

		if err != nil {
			fmt.Println("Error updating ip address with Digital Ocean:", err)
			continue
		}

		if success {
			fmt.Println("DNS record updated successfully! " + ip)
		}

		// Sleep and try again.
		time.Sleep(checkInterval)
	}
}

//
// getPublicIP will query ipify.org to get our current ip address.
//
func getPublicIP() (string, error) {
	// Get the public IP address from a service like 'https://api.ipify.org'
	resp, err := http.Get("https://api.ipify.org")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(ip), nil
}

//
// updateDNSRecord will send two api calls to Digital Ocean. One to get the current record and another to update it.
//
func updateDNSRecord(newIP string) (bool, error) {
	// Convert the environment variable to an int64.
	id, err := strconv.ParseInt(os.Getenv("RECORD_ID"), 10, 64)

	if err != nil {
		return false, err
	}

	recordID := int(id)

	// Setup Digital Ocean client
	client := godo.NewFromToken(os.Getenv("DO_TOKEN"))

	ctx := context.TODO()
	record, _, err := client.Domains.Record(ctx, os.Getenv("DOMAIN"), recordID)

	if err != nil {
		log.Fatalf("Unable to get DNS record: %v", err)
	}

	// If the IP address has not changed no need to do anything else.
	if newIP == record.Data {
		return false, nil
	}

	// Update the record with the new IP address
	editRequest := &godo.DomainRecordEditRequest{
		Type: record.Type,
		Name: record.Name,
		Data: newIP,
	}
	_, _, err = client.Domains.EditRecord(ctx, os.Getenv("DOMAIN"), recordID, editRequest)

	if err != nil {
		return false, err
	}

	return true, nil
}
