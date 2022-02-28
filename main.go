package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	cid "github.com/ipfs/go-cid"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
)

func main() {
	var hash string
	flag.StringVar(&hash, "hash", "", "Hash to search")

	flag.Parse()

	if hash == "" {
		// valid  hash "QmTp2hEo8eXRp6wg7jXv1BLCMh5a4F3B7buAUZNZUu772j"
		hash = generateCID()
	}

	// find hash from IPFS network
	if !find(hash) {
		fmt.Printf("Hash %s not found ", hash)
		return
	}

	fmt.Printf("Congrats!... Hash  %s found on ipfs network", hash)
}

func generateCID() string {

	var builder cid.V0Builder

	msg := randomString(35)

	c, _ := cid.V0Builder.Sum(builder, []byte(msg))
	s := c.Hash().B58String()
	fmt.Println("Created CID: ", s)
	return s
}

func find(hash string) bool {

	out := fmt.Sprintf("%s.txt", hash)
	url := loadConfig("IPFS_URL")
	sh := shell.NewShell(url)
	err := sh.Get(hash, out)
	if err != nil {
		fmt.Printf("unable to find hash %s", err)
		return false
	}

	data, err := sh.Cat(hash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return false
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)
	newStr := buf.String()
	fmt.Printf("data %s", newStr)

	return true
}

func loadConfig(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)

}

func generateHash() string {
	noRandomCharacters := 32
	randString := randomString(noRandomCharacters)

	hash := sha256.New()
	hash.Write([]byte(randString))
	bs := hash.Sum(nil)
	fmt.Println(len(bs))
	return fmt.Sprintf("%x", bs)
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}
