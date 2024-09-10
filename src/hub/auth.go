package hub

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	ALPHABET    string = "_-~.abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CODE_LENGTH int    = 128
)

func Random_char() string {
	return string(ALPHABET[rand.Intn(len(ALPHABET))])
}

func Random_code(length int) string {
	random_code := ""
	for i := 0; i < length; i++ {
		random_code += Random_char()
	}
	return random_code
}

// not exactly sure how this is used, come back later once that is clear. return type in python is string not a byte slice/list.
func Code_challenge(code_verifier string) string {
	sha256_hash := sha256.New()
	sha256_hash.Write([]byte(code_verifier))
	digest := sha256_hash.Sum(nil)
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(digest)
}

func Send_challenge(ip_address string, code_verifier string) string {
	auth_url := "https://" + string(ip_address) + ":8443/v1/oauth/authorize"
	params := map[string]string{
		"audience":              "homesmart.local",
		"response_type":         "code",
		"code_challenge":        Code_challenge(code_verifier),
		"code_challenge_method": "S256",
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	req, err := http.NewRequest("GET", auth_url, nil)
	if err != nil {
		fmt.Printf("error encountered making GET request: %s", err)
		return ""
	}
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Type", "application/json") // or any other headers you need
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error encountered performing request: %s", err)
		return ""
	}
	defer resp.Body.Close()
	// Do something with resp
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// handle non-2xx status codes
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return ""
	}
	// Parse the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("error decoding json response: %s", err)
		return ""
	}
	code, ok := result["code"].(string)
	if !ok {
		fmt.Print("error getting code entry from result")
		return ""
	}

	return code
}

func Get_token(ip_address string, code string, code_verifier string) string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("error encountered getting Hostname: %s", err)
		return ""
	}
	data := fmt.Sprintf("code=%s&name=%s&grant_type=authorization_code&code_verifier=%s",
		code, hostname, code_verifier)
	header := map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	token_url := fmt.Sprintf("https://%s:8443/v1/oauth/token", ip_address)
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("POST", token_url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Printf("error encountered making POST request: %s", err)
		return ""
	}
	// Add headers to the request
	for key, value := range header {
		req.Header.Set(key, value)
	}
	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		// handle error
		fmt.Printf("error encountered performing request: %s", err)
		return ""
	}
	defer resp.Body.Close()
	// Do something with resp
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// handle non-2xx status codes
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return ""
	}
	// Parse the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("error decoding json response: %s", err)
		return ""
	}
	token, ok := result["access_token"].(string)
	if !ok {
		fmt.Print("error getting access_token entry from json")
		return ""
	}
	return token
}

func Main() {
	var err error
	var ip_address string
	if len(os.Args) > 1 {
		ip_address = os.Args[1]
	} else {
		fmt.Print("Input the ip iaddress of your Dirigera then hit ENTER ...\n")
		_, err = fmt.Scan(&ip_address)
	}
	if err != nil {
		fmt.Printf("error encountered getting ip address from user input: %s", err)
		return
	}
	code_verifier := Random_code(CODE_LENGTH)
	code := Send_challenge(ip_address, code_verifier)
	fmt.Print("Press the action button on Dirigera bridge, then enter any character and hit ENTER ...")
	var stopper string
	// for some reason this does not stop for ENTER with Scanln.. fix that at some point
	_, err = fmt.Scan(&stopper)
	token := Get_token(ip_address, code, code_verifier)
	fmt.Printf("Your TOKEN: %s", token)
}
