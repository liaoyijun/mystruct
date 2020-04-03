package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Farmyard/mystruct/gormat"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "mystruct"
	app.Usage = "A golang convenient converter supports Database to Struct"
	app.Version = "0.0.0"
	app.UsageText = "mystruct [GLOBAL OPTIONS] [DATABASE]"
	app.UseShortOptionHandling = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "host",
			Aliases: []string{"h"},
			Value:   "127.0.0.1",
			Usage:   "Host address of the database.",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"P"},
			Value:   3306,
			Usage:   "Port number to use for connection. Honors $MYSQL_TCP_PORT.",
		},
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Value:   "root",
			Usage:   "User name to connect to the database.",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Password to connect to the database.",
		},
		&cli.StringFlag{
			Name:    "database",
			Aliases: []string{"D"},
			Usage:   "Database to use.",
		},
	}
	app.Action = func(c *cli.Context) error {
		opt := &gormat.Option{
			Host:     c.String("host"),
			Port:     c.Int("port"),
			Username: c.String("username"),
			Password: c.String("password"),
			Database: c.String("database"),
		}
		if opt.Database == "" {
			opt.Database = c.Args().First()
		}
		if err := handle(opt); err != nil {
			fmt.Printf(err.Error())
		}
		return nil
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:  "help",
		Usage: "Show this message and exit",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "version",
		Usage: "Output mycli's version.",
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handle(opt *gormat.Option) error {
	f := bufio.NewReader(os.Stdin)
	if opt.Password == "" {
		fmt.Printf("Password:")
		input, _ := f.ReadString('\n')
		input = strings.TrimRight(input, "\n")
		opt.Password = input
	}
	// Open mysql
	con, err := gormat.Open(opt)
	if err != nil {
		return err
	}
	printStd(opt)
	for {
		input, _ := f.ReadString('\n')
		input = strings.TrimRight(input, "\n")
		if input == "quit" {
			break
		}
		a := strings.Split(input, " ")
		// use or change database
		if a[0] == "use" {
			opt.Database = a[len(a)-1]

			var err error
			err = con.Use(opt.Database)
			if err != nil {
				fmt.Printf("(1049, \"Unknown database '" + opt.Database + "'\")\n")
				opt.Database = ""
			}
		} else {
			if opt.Database == "" {
				fmt.Printf("(1046, 'No database selected')\n")
			} else {
				tables := strings.Split(input, ",")
				std := con.Make(tables...)
				fmt.Println(std)
			}
		}
		printStd(opt)
	}
	return nil
}

func printStd(opt *gormat.Option) {
	var database string
	if opt.Database == "" {
		database = "(none)"
	} else {
		database = opt.Database
	}
	fmt.Printf("mysql " + opt.Username + "@" + opt.Host + ":" + database + ">")
}
