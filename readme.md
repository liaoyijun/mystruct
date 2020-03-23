## Quick Start

```
$ git clone git@github.com:Farmyard/mystruct.git mystruct
$ cd mystruct
$ go install
```

## Usage

```
$ mystruct --help
NAME:
   mystruct - A golang convenient converter supports Database to Struct

USAGE:
   mystruct [GLOBAL OPTIONS] [DATABASE]

VERSION:
   0.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value, -h value      Host address of the database. (default: "127.0.0.1")
   --port value, -P value      Port number to use for connection. Honors $MYSQL_TCP_PORT. (default: 3306)
   --username value, -u value  User name to connect to the database. (default: "root")
   --password value, -p value  Password to connect to the database.
   --database value, -D value  Database to use.
   --help                      Show this message and exit (default: false)
   --version                   Output mycli's version. (default: false)
```

## example

```
$ mysql root@127.0.0.1:(none)>use my_database
$ mysql root@127.0.0.1:(my_database)>my_table1,my_table2
```