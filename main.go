package main

import (
	"fmt"
	"os"
	"stripe-checker/src"

	"github.com/urfave/cli/v2"
)

var (
	app        = &cli.App{}
	cfg        src.Cfg
	result     src.Result
	card       src.Card
	filename   string
	output     string
	configPath string
	separator  string
)

func init() {
	app.Name = "stripe-checker"
	app.Usage = "credit card checker using stripe payment gateway."
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"cfg", "c"},
			Value:       "./config.cfg",
			Usage:       "config path to checker",
			Destination: &configPath,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"out", "o"},
			Value:       "./lives.txt",
			Usage:       "output name",
			Destination: &output,
		},
		&cli.StringFlag{
			Name:        "separator",
			Aliases:     []string{"sep", "s"},
			Value:       "|",
			Usage:       "separator that separate the credit card.",
			Destination: &separator,
		},
	}
	if len(os.Args) < 2 {
		app.RunAndExitOnError()
		os.Exit(0)
	}
	app.Action = func(ctx *cli.Context) error {
		filename = ctx.Args().First()
		cfg = src.LoadCfg(configPath)

		src.OpenFileByName(filename, func(line string) {
			card = src.GetCardByLine(line, separator)
			result = src.CheckCard(card, cfg)

			if result.Valid {
				fmt.Printf("[live] %s, %s/%s, %s (%s) \n", card.CardNumber, card.ExpMonth, card.ExpYear, card.Cvv, result.Code)
				src.SaveCard(card, output, result)
			} else {
				fmt.Printf("[die] %s, %s/%s, %s (%s) \n", card.CardNumber, card.ExpMonth, card.ExpYear, card.Cvv, result.Code)
			}

		})
		return nil
	}
}

func main() {
	err := app.Run(os.Args)
	src.HandleError(err)
}