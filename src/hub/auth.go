package hub

import (
	"crypto/sha256"
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
	return string(sha256_hash.Sum(nil))
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
	}
	req, err := http.NewRequest("GET", auth_url, nil)
	if err != nil {
		// handle error
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json") // or any other headers you need
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	// Do something with resp
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// handle non-2xx status codes
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
	// Parse the JSON response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "error decoding JSON response"
	}
	code, ok := result["code"].(string)
	if !ok {
		return "error extracting 'code' from response"
	}

	return code
}

func Get_token(ip_address string, code string, code_verifier string) string {
	fmt.Printf("to be implemented, %s %s %s", ip_address, code, code_verifier)
	return "lol nope"
}

func main() {
	var ip_address string
	if len(os.Args) > 1 {
		ip_address = os.Args[1]
	} else {
		fmt.Print("Input the ip iaddress of your Dirigera then hit ENTER ...\n")
		return
	}
	code_verifier := Random_code(CODE_LENGTH)
	code := Send_challenge(ip_address, code_verifier)
	fmt.Print("Press the action button on Dirigera bridge, then hit ENTER ...")
	fmt.Scan()
	token := Get_token(ip_address, code, code_verifier)
	fmt.Printf("Your TOKEN: %s", token)
}
