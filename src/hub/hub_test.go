package hub

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

var test_ALPHABET string = "_-~.abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Test_auth_Random_char_length(t *testing.T) {
	length := len(Random_char())
	want := 1
	if !(want == length) {
		t.Fatalf("Received random char of invalid length: %s, expected: %s", strconv.Itoa(length), strconv.Itoa(want))
	}
}

func Test_auth_Random_char_in_alphabet(t *testing.T) {
	want := test_ALPHABET
	received := Random_char()
	if !(strings.Contains(want, received)) {
		t.Fatalf("Received random char not in specified alphabet: %s", received)
	}
}

func Test_auth_Random_code_length(t *testing.T) {
	want := 32
	length := len(Random_code(want))
	if !(want == length) {
		t.Fatalf("Received random code of invalid length: %s, expected: %s", strconv.Itoa(length), strconv.Itoa(want))
	}
}

func Test_auth_Random_code_in_alphabet(t *testing.T) {
	want := test_ALPHABET
	received := Random_code(32)
	for _, c := range received {
		if !(strings.Contains(want, string(c))) {
			t.Fatalf("Received random char not in specified alphabet: %s", received)
		}
	}
}

func Test_auth_Main(t *testing.T) {
	var ip_address string
	ip_address = "192.168.178.27"
	code_verifier := Random_code(CODE_LENGTH)
	code := Send_challenge(ip_address, code_verifier)
	fmt.Print("Press the action button on Dirigera bridge, then hit ENTER ...")
	fmt.Scan()
	token := Get_token(ip_address, code, code_verifier)
	fmt.Printf("Your TOKEN: %s", token)
}
