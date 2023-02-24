package main

import (
    "bytes"
    "encoding/json"
    "errors"
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

    //for {
    NOW := time.Now()
    fmt.Printf("START \nINFO - Current time is: %s\n", NOW.Format("2006-01-02 15:04:05 Monday"))
    // check IP
    compareCurrentAndDomainIP(GODADDY_KEY, GODADDY_SECRET, GODADDY_DOMAIN, IP_PROVIDER)
    // sleep
    time.Sleep(time.Second * time.Duration(*POLLING_INTERVAL))
    //}
}

func verifyVar(VAR string, MESSAGE string) {
    if VAR == "" {
        log.Fatalf("ERROR - %s", MESSAGE)
    }
}

func getCurrentIP(IP_PROVIDER string) (string, error) {
    response, err := http.Get(IP_PROVIDER)
    if err != nil {
        log.Fatal(err)
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
    jsonData := make([]struct {
        Data string `json:"data"`
    }, 1)
    json.NewDecoder(response.Body).Decode(&jsonData)
    return jsonData[0].Data, nil
}

func compareCurrentAndDomainIP(GODADDY_KEY string, GODADDY_SECRET string, GODADDY_DOMAIN string, IP_PROVIDER string) {
    // current IP
    CURRENT_IP, err := getCurrentIP(IP_PROVIDER)
    errMsg(err)
    fmt.Printf("INFO - Current IP is: %s\n", CURRENT_IP)

    // domain IP
    DOMAIN_IP, err := getDomainIP(GODADDY_KEY, GODADDY_SECRET, GODADDY_DOMAIN)
    errMsg(err)
    fmt.Printf("INFO - Domain IP is: %s\n", DOMAIN_IP)

    if CURRENT_IP != DOMAIN_IP {
        fmt.Println("INFO - The IP is different!")
        RESPONSE, err := putDomainIP(CURRENT_IP, GODADDY_KEY, GODADDY_SECRET, GODADDY_DOMAIN)
        errMsg(err)
        fmt.Printf(RESPONSE)
    } else {
        fmt.Println("INFO - The IP is same.")
    }
}

func errMsg(ERR error) {
    if ERR != nil {
        log.Fatal(ERR)
    }
}

func putDomainIP(CURRENT_IP string, GODADDY_KEY string, GODADDY_SECRET string, GODADDY_DOMAIN string) (string, error) {

    type Data struct {
        Data string `json:"data"`
        TTL  int    `json:"ttl"`
    }
    jsonData, _ := json.Marshal([]Data{
        {
            Data: CURRENT_IP,
            TTL:  3600,
        },
    })
    requestBody := bytes.NewBuffer(jsonData)
    //fmt.Println(requestBody)
    request, err := http.NewRequest("PUT",
        fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/@", GODADDY_DOMAIN),
        requestBody)
    // NOTE
    // JSON will be: [{"data":"11.22.33.44","ttl":3600}]

    //var buf bytes.Buffer
    //err := json.NewEncoder(&buf).Encode(&struct {
    //	Data string `json:"data"`
    //	TTL  int    `json:"ttl"`
    //} {
    //	CURRENT_IP,
    //	3600,
    //})
    //fmt.Print(&buf)
    //if err != nil {
    //	return "", err
    //}
    //request, err := http.NewRequest("PUT",
    //	fmt.Sprintf("https://api.godaddy.com/v1/domains/%s/records/A/@", GODADDY_DOMAIN),
    //	&buf)
    // NOTE
    // JSON will be: {"data":"11.22.33.44","ttl":3600}

    //fmt.Println(request)
    if err != nil {
        return "", err
    }
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", GODADDY_KEY, GODADDY_SECRET))

    client := &http.Client{}

    response, err := client.Do(request)
    if err != nil {
        return "", err
    }
    if response.StatusCode == 200 {
        return "INFO - Success!", nil
    } else {
        return "", errors.New(fmt.Sprintf("ERROR - HTTP status code %d\n", response.StatusCode))
    }
}
