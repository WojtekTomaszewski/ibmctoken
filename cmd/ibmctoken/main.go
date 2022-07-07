package main

import (
	"fmt"
	"github.com/WojtekTomaszewski/ibmctoken"
	"os"
)

func main() {
	// Expects api key to be argument
	if len(os.Args) != 2 {
		fmt.Println("Usage: ibmctoken <api-key>")
		os.Exit(1)
	}

	apikey := os.Args[1]

	token := ibmctoken.NewToken(apikey)

	err := token.RequestToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(token.AccessToken)
}
