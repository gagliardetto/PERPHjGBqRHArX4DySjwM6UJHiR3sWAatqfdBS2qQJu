package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	program "github.com/gagliardetto/PERPHjGBqRHArX4DySjwM6UJHiR3sWAatqfdBS2qQJu"
)

func main() {
	var input string
	flag.StringVar(&input, "input", "", "Input file containing JSON lines of accounts")
	flag.Parse()
	if input == "" {
		panic("Input file must be specified with -input flag")
	}
	// Open the input file
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		account, err := decodeLine(line)
		if err != nil {
			panic(err)
		}
		dataBytes, err := account.GetDataBytes()
		if err != nil {
			panic(err)
		}
		parsedAccount, err := program.ParseAnyAccount(dataBytes)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to parse account %s : %w", account.Pubkey, err))
			continue
		}
		println("Account:", account.Pubkey, "Data Length:", len(dataBytes))
		spew.Dump(parsedAccount)
		// Process the data bytes as needed
		// For example, print the length of the data
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

type Account struct {
	Pubkey     string `json:"pubkey"`
	Owner      string `json:"owner"`
	Executable bool   `json:"executable"`
	DataLen    uint64 `json:"data_len"`
	Lamports   uint64 `json:"lamports"`
	Data       string `json:"data"`
}

// GetDataBytes
func (a *Account) GetDataBytes() ([]byte, error) {
	// decode as base64
	return base64.StdEncoding.DecodeString(a.Data)
}

func decodeLine(line []byte) (*Account, error) {
	var account Account
	err := json.Unmarshal(line, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
