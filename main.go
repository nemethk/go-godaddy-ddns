package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// globals
var GODADDY_KEY = ""
var GODADDY_SECRET = ""
var GODADDY_DOMAIN = ""
var IP_PROVIDER = ""

func main() {

	// required
	ptrGodaddyKey := flag.String("godaddy-key", "", "Godaddy API key")
	ptrGodaddySecret := flag.String("godaddy-secret", "", "Godaddy API secret")
	ptrGodaddyDomain := flag.String("godaddy-domain", "", "Registered domain name")

	// optional
	POLLING_INTERVAL := flag.Int("polling-interval", 10, "Polling interval in seconds")
	ptrIpProvider := flag.String("ip-provider", "https://v4.ident.me/", "IP provider API")

	flag.Parse()
	GODADDY_KEY = *ptrGodaddyKey
	GODADDY_SECRET = *ptrGodaddySecret
	GODADDY_DOMAIN = *ptrGodaddyDomain
	IP_PROVIDER = *ptrIpProvider

	// is it specified
	verifyVar(GODADDY_KEY, "Specify the GoDaddy API key")
	verifyVar(GODADDY_SECRET, "Specify the GoDaddy API secret")
	verifyVar(GODADDY_DOMAIN, "Specify the GoDaddy domain")

	// current IP
	currentIP, err := getCurrentIP(IP_PROVIDER)
	if err != nil {
		log.Fatal(err)
	}
	println(currentIP)

	// domain IP
	domainIP, err := getDomainIP(GODADDY_KEY, GODADDY_SECRET, GODADDY_DOMAIN)
	if err != nil {
		log.Fatal(err)
	}
	println(domainIP)

	// sleep
	time.Sleep(time.Second * time.Duration(*POLLING_INTERVAL))
}

func verifyVar(VAR string, MESSAGE string) {
	if VAR == "" {
		log.Fatalf("ERROR: %s", MESSAGE)
	}
}

func getCurrentIP(IP_PROVIDER string) (string, error) {
	response, err := http.Get(IP_PROVIDER)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.String(), nil
}

func getDomainIP(GODADDY_KEY string, GODADDY_SECRET string, GODADDY_DOMAIN string) (string, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A", GODADDY_DOMAIN), nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", fmt.Sprintf("sso-key %s:%s", GODADDY_KEY, GODADDY_SECRET))
	client := new(http.Client)
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	jsondata := make([]struct {
		Data string `json:"data"`
	}, 1)
	json.NewDecoder(response.Body).Decode(&jsondata)
	return jsondata[0].Data, nil
}

