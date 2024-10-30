package main

import (
	"context"
	"fmt"

	"github.com/BinLab64/Orbix-client/pkg/api"
)

func main() {
	fmt.Println("--- ", api.DefaultBaseURL)
	TestApi()
}

func TestApi() {
	client := api.NewClient(api.ClientOptions{
		// ClientAuth: api.NewClientAuth(),
		Logger: nil,
	})
	// res, err := client.NewExchangeInfoService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// res.Symbols
	// fmt.Println(res)
	// client.NewExchangeInfoService()

	// res, err := client.NewOrderbookDepthService("USDT_THB").
	// 	Limit(1).
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("%+v", res)

	res, err := client.NewOrderbookService("USDT_THB").
		// Side(api.SideTypeBuy).
		Side(api.SideTypeSell).
		Do(context.Background())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("res.Asks:%+v\n", res.Asks)
	fmt.Printf("res.Bids:%+v\n", res.Bids)

	// res, err := client.NewListCurrentOrdersService("usdt_thb", 10, 0).
	// 	Do(context.Background())

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("res.:%+v\n", res)

	// res, err := client.NewListBalanceAddressService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("res.:%+v\n", res.Wallets["usdt"].AvailableBalance)

	// res, err := client.NewCreateOrderService("usdt_thb",
	// 	api.SideTypeBuy,
	// 	api.OrderTypeLimit,
	// 	"36.11",
	// 	"10").
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("res.:%+v\n", res)

	// res, err := client.NewCreateOrderService("usdt_thb",
	// 	api.SideTypeSell,
	// 	api.OrderTypeLimit,
	// 	"34",
	// 	"5").
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("res.:%+v\n", res)

	// err := client.NewCancelOrderService("52896594", "usdt_thb").
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println("Yoo", err)
	// }

	// res, err := client.NewCancelAllOrdersService("usdt_thb").
	// 	Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("res.:%+v\n", res)

}
