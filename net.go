package main

import (
	"fmt"
	"strings"
)

func reserver(str *string) {
	strlist := strings.Split(*str, " ")
	for i := 0; i < len(strlist)/2; i++ {
		strlist[i], strlist[len(strlist)-1-i] = strlist[len(strlist)-1-i], strlist[i]
	}
	*str = strings.Join(strlist, " ")
	return
}

func main() {
	str := "hello world"
	reserver(&str)
	fmt.Println(str)
}
