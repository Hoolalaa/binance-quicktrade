# binance-quicktrade
a bot to help beginner make a quick trade on binance.com

## Background

After being invited to a few trading discords thanks to my [discord stock tickers](https://github.com/rssnyder/discord-stock-ticker), I started getting spam in my DMs for people selling access to other servers, and bots to help you make money from "pump and dumps".

Interested, I started looking into the pump and dumps, and I had a good guess at what these "premium bots" were attempting to accomplish. So I thought I would try out making one myself, and of course making it free to use.

## Disclaimer

I am not a genious programmer, or a trader by any means. This code comes with NO guarentee of profit. Trade at your own risk. I am not respoinsible for losses incurred from using this code.

You should have you binance console up when using this program in case anything goes wrong, so you can cancel any unintended trades.

## General flow

The script walks you through selecting the coin you have that you would like to use to buy with. BTC is the most popular choice as it has many trading pairs on binance.

After setting the coin and amount, you are asked for a profit limit. This is where you can set your profit limit amount. So if you want to cash out after you make 20%, you would enter 20 here.

You are also asked how long you would like to wait, if at all. This means if you want to bail on the trade after 10 seconds and cut your losses, you would enter 10 here.

Lastly, you are asked for the coin to buy. As soon as you enter the coin here, trades will execute. Bailing out of the script after this point can lead to uninteded side effects.

From here, the script will walk you through buying the coin at the current market price, setting a limit order for your target profit, and watching your order. If you enter no time limit, this will go on until the status of your order is changed from NEW or PARTIALLY FILLED.

If you specify the time limit and it is reached, your limit order will be cancelled and the program will attempt to sell what you have at the current market price.

## Usage

Download the program for your operating system from the [release page](https://github.com/rssnyder/binance-quicktrade/releases).

Navigate to where you store the binary, and execute it. You should have a [binance API key](https://www.binance.com/en/support/faq/360002502072-How-to-create-API) ready to go.

## Support

If you have issues with this program or **constructive** critisism, please open a github issue or find me on discord at `jonesbooned#1111`.

Did this work for you? Maybe [buy me a coffee](https://ko-fi.com/rileysnyder)! Or send some crypto so I can buy more takis:

eth: 0x27B6896cC68838bc8adE6407C8283a214ecD4ffE

doge: DTWkUvFakt12yUEssTbdCe2R7TepExBA2G

bch: qrnmprfh5e77lzdpalczdu839uhvrravlvfr5nwupr

btc: 1N84bLSVKPZBHKYjHp8QtvPgRJfRbtNKHQ
