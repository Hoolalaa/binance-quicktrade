package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
)

func main() {

	// Get users API keys
	var keyId string
	var keySecret string

	fmt.Println("Enter Binance API Key ID: ")
	fmt.Scanln(&keyId)

	fmt.Println("\nEnter Binance API Key Secret: ")
	fmt.Scanln(&keySecret)

	// Log into binance
	binance.UseTestnet = true
	client := binance.NewClient(keyId, keySecret)

	// Correct time issues
	client.NewSetServerTimeService().Do(context.Background())

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println("Error logging into binance.")
		fmt.Println(err)
		return
	}
	// fmt.Printf("%+v\n", res)

	// Get base coin to use for trade
	var baseCoin string

	fmt.Println("\nEnter base coin to trade with: ")
	fmt.Scanln(&baseCoin)
	baseCoin = strings.ToUpper(baseCoin)

	// Find amount we have avalible to trade
	var freeAmountFloat float64

	for _, coin := range res.Balances {
		if coin.Asset == baseCoin {
			if freeAmountFloat, err := strconv.ParseFloat(coin.Free, 64); err == nil {
				fmt.Printf("\n%g %s avalible to trade.", freeAmountFloat, baseCoin)
			}
		}
	}

	// Get amount to trade
	var tradeAmount string

	fmt.Println("\nEnter amount to trade with (eg. 0.01): ")
	fmt.Scanln(&tradeAmount)

	if tradeAmountFloat, err := strconv.ParseFloat(tradeAmount, 64); err == nil {
		if tradeAmountFloat > freeAmountFloat {
			fmt.Printf("Trading with %g %s\n", tradeAmountFloat, baseCoin)
		} else {
			fmt.Printf("Not enough funds, you have %g %s avalible to trade.\n", tradeAmountFloat, baseCoin)
			return
		}
	} else {
		fmt.Println("Error setting trade amount.")
		fmt.Println(err)
		return
	}

	// Get take profit percentage
	var takePercentage string
	var takePercentageMultiplier float64

	fmt.Println("\nEnter percent gain to sell at (eg. 20 for 20%): ")
	fmt.Scanln(&takePercentage)

	if takePercentageFloat, err := strconv.ParseFloat(takePercentage, 64); err == nil {
		takePercentageMultiplier = takePercentageFloat/100 + 1
	} else {
		fmt.Println("Error calculating profit target.")
		fmt.Println(err)
		return
	}

	// Get timeout
	var timeout string

	fmt.Println("\nEnter a time where you would like to bail from your trade if you sell order has not executed.")
	fmt.Println("If the specified amount of time passes and the price has not reached your gain target, you will sell at market price.")
	fmt.Println("Enter a time in seconds (eg. 10 for 10 seconds): ")
	fmt.Scanln(&timeout)

	timeoutInt, err := strconv.Atoi(timeout)
	if err != nil {
		fmt.Println("Error setting timeout (in seconds).")
		fmt.Println(err)
		return
	}

	timeoutDuration, err := time.ParseDuration(fmt.Sprintf("%ds", timeoutInt))
	if err != nil {
		fmt.Println("Error setting timeout (in seconds).")
		fmt.Println(err)
		return
	}

	// Get coin to buy
	var targetCoin string

	fmt.Println("\nWARNING: As soon as you enter your target coin, trades will be made. Proceed with caution.")
	fmt.Println("Enter coin to buy: ")
	fmt.Scanln(&targetCoin)
	targetCoin = strings.ToUpper(targetCoin)

	targetPair := targetCoin + baseCoin

	// Place market order
	buyOrder, err := client.NewCreateOrderService().Symbol(targetPair).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		Quantity(tradeAmount).
		Do(context.Background())

	if err != nil {
		fmt.Println("Error placing market order.")
		fmt.Println(err)
		return
	}

	//Mark time we bought at
	start := time.Now()

	// Find highest buy price in market fills
	highestPrice := 0.0

	for _, fill := range buyOrder.Fills {
		if fillPriceFloat, err := strconv.ParseFloat(fill.Price, 64); err == nil {
			if fillPriceFloat > highestPrice {
				highestPrice = fillPriceFloat
			}
		} else {
			fmt.Println("Error finding market buy prices.")
			fmt.Println(err)
			return
		}
	}

	// Show order to user
	fmt.Printf("\nOrder: %d\nStatus: %s\nExecuted Quantity: %s\nPrice: %g\n\n", buyOrder.OrderID, buyOrder.Status, buyOrder.ExecutedQuantity, highestPrice)

	// Get target price for profit exit
	targetPrice := highestPrice * takePercentageMultiplier

	fmt.Printf("We are targeting a price of %g for taking profit.\n", targetPrice)

	// Set profit take order
	profitOrder, err := client.NewCreateOrderService().Symbol(targetPair).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		Price(fmt.Sprintf("%f", targetPrice)).Quantity(buyOrder.ExecutedQuantity).
		TimeInForce(binance.TimeInForceTypeGTC).Do(context.Background())

	if err != nil {
		fmt.Println("Error placing sell order. YOU ARE AT RISK.")
		fmt.Println(err)
		return
	}

	// Show order to user
	fmt.Printf("\nOrder: %d\nStatus: %s\nOrigional Quantity: %s\nPrice: %s\n\n", profitOrder.OrderID, profitOrder.Status, profitOrder.OrigQuantity, profitOrder.Price)

	// Time passed
	fmt.Println("\nTime since purchase")
	fmt.Println(time.Since(start))

	// Watch sell order
	for {
		watchOrder, err := client.NewGetOrderService().Symbol(targetPair).
			OrderID(profitOrder.OrderID).Do(context.Background())

		if err != nil {
			fmt.Println("Can't get sell order status. YOU ARE AT RISK.")
			fmt.Println(err)
			return
		}

		fmt.Printf("\nOrder: %d\nStatus: %s\n\n", profitOrder.OrderID, profitOrder.Status)
		if watchOrder.Status == binance.OrderStatusTypeFilled {
			fmt.Println("Success! You have sold for your target profit! Please see the binance console for more information.")
			return
		} else if watchOrder.Status == binance.OrderStatusTypePartiallyFilled {
			fmt.Println("Part of your sell order has filled, keep waiting...")
		} else if watchOrder.Status == binance.OrderStatusTypeCanceled {
			fmt.Println("Something has gone wrong. Your sell order was cancelled.")
			return
		} else if watchOrder.Status == binance.OrderStatusTypePendingCancel {
			fmt.Println("Something has gone wrong. Your sell order was cancelled.")
			return
		} else if watchOrder.Status == binance.OrderStatusTypeRejected {
			fmt.Println("Something has gone wrong. Your sell order was rejected.")
			return
		} else if watchOrder.Status == binance.OrderStatusTypeExpired {
			fmt.Println("Something has gone wrong. Your sell order expired.")
			return
		} else if time.Since(start) > timeoutDuration {
			fmt.Println("We are out of time!")
			break
		}
	}

	// Cancel sell order
	_, err = client.NewCancelOrderService().Symbol(targetPair).
		OrderID(profitOrder.OrderID).Do(context.Background())

	if err != nil {
		fmt.Println("Error cancelling sell order. YOU ARE AT RISK.")
		fmt.Println(err)
		return
	}

	// Sell what we have at the current price
	sellOrder, err := client.NewCreateOrderService().Symbol(targetPair).
		Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).
		Quantity(buyOrder.ExecutedQuantity).Do(context.Background())

	if err != nil {
		fmt.Println("Error placing market sell order.")
		fmt.Println(err)
		return
	}

	// Show order to user
	fmt.Printf("\nOrder: %d\nStatus: %s\nExecuted Quantity: %s\nPrice: %s\n\n", sellOrder.OrderID, sellOrder.Status, sellOrder.ExecutedQuantity, sellOrder.Price)

	for _, fill := range sellOrder.Fills {
		fmt.Printf("Sold %s for %s\n", fill.Quantity, fill.Price)
	}

	if sellOrder.Status != binance.OrderStatusTypeFilled {
		fmt.Println("WARNING: It seems your market sell order failed. YOU ARE AT RISK.")
		return
	}

	fmt.Println("\n\nExiting...")

	return
}
