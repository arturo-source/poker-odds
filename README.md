# Calculate poker odds from your terminal

This is a simple program to calculate your equity, if you know all the hands in the table. You should simply enter the board, and the hands from all the people in the table, and you get the equity from each one. Example of usage:

```bash
poker-odds --board 7c8dQd Qc9h KdQs 3h5h 7d6h 2d3d
```

![Result of running the last command](https://github.com/arturo-source/poker-odds/assets/59207995/4bf28114-9355-4611-871e-842baf190db1)

If you want to calculate the probability of getting a **Flush** with your hand, you should not enter the other hands, just your own. If you play multiple hands, you will get the probability of winning other hands with a **Spin**, not the chance of getting that combination.

![Result of running the last command with just one hand](https://github.com/arturo-source/poker-odds/assets/59207995/17b2ceb2-16d7-477e-9f7d-f64d6c37a041)

The program is **BLAZINGLY FAST**, because I use [poker engine](https://github.com/arturo-source/poker-engine), which I made using uint64, instead of arrays to calculate combinations faster (using bit operations). But obviously it gets slower (up to 4 seconds in my computer) if you don't use `--board`, because it gets tons of board combinations ([nCr of 52 5](https://en.wikipedia.org/wiki/Poker_probability#5-card_poker_hands)).

![Result of running hands without board](https://github.com/arturo-source/poker-odds/assets/59207995/c7b781e1-d501-4b1b-aa20-1e8ac82d62b1)

## Program options

- Use `-h` to know all options.
- Up to 5 cards in board. If you don't want to have initial board, do not use `--board` option.
- Default is coloring the terminal, that works for Linux and MacOS (it is disabled for Windows), but if you want to save the output, or process it, you can use `--no-color` to disable.

## Build the program your own

If you don't want to download the binary [from releases](https://github.com/arturo-source/poker-odds/releases), you can build easily with the Go compiler.

```bash
git clone https://github.com/arturo-source/poker-odds.git
cd poker-odds
go build
```

Or you can install the program in your GOPATH, to run it as a system program.

```bash
go install github.com/arturo-source/poker-odds@latest
```
