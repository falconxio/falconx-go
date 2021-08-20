package main

import (
	"flag"
	"github.com/falconxio/falconx-go/client_examples"
	"log"
	"strings"
)

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {

	// Mandatory Arguments
	var APIKey *string = flag.String("api_key", "XXX", "API Key provided by FalxonX")
	var Secret *string = flag.String("secret", "XXX", "Secret provided by FalxonX")
	var Passphrase *string = flag.String("passphrase", "XXX", "Pass Phrase provided by FalxonX")

	// Optional Arguments
	var example_set *string = flag.String("example_set", "rest", "Example Set to run")
	var Host *string = flag.String("host", "", "Host To Query")
	flag.Parse()

	isAPIKeyPassed := isFlagPassed("api_key")
	isSecretPassed := isFlagPassed("secret")
	isPassphrasePassed := isFlagPassed("passphrase")

	if !isAPIKeyPassed {
		panic("api_key argument missing! ( eg. go run run_examples.go -api_key=<api_key>")
	} else if !isSecretPassed {
		panic("secret argument missing! ( eg. go run run_examples.go -secret=<secret>")
	} else if !isPassphrasePassed {
		panic("passphrase argument missing! ( eg. go run run_examples.go -passphrase=<passphrase>")
	}
	possibleExampleSets := []string{"websocket", "rest"}

	log.Printf("example_set: %s", *example_set)
	if *example_set == "websocket" {
		client_examples.RunWebSocketExamples(*APIKey, *Secret, *Passphrase, *Host)
	} else if *example_set == "rest" {
		client_examples.RunRestExamples(*APIKey, *Secret, *Passphrase, *Host)
	} else {
		log.Fatalf("Example Set Not Found! Example sets available: %s", strings.Join(possibleExampleSets[:], ","))
	}

}
