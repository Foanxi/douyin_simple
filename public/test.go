package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	str := "./Data/photo/"
	str1 := str[1:]
	fmt.Print(str1)
}
