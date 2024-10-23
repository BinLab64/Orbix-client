package main

import (
	"context"
	"fmt"

	"github.com/BinLab64/Orbix-client/api"
)

func main() {
	fmt.Printf("--- %s", api.DefaultBaseURL)
	TestApi()
}

func TestApi() {
	client := api.NewClient(api.ClientOptions{})
	res, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	// res.Symbols
	fmt.Println(res)
	// client.NewExchangeInfoService()

}
