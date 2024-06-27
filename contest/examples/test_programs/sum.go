package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	values := strings.Split(text, " ")
	a, _ := strconv.Atoi(values[0])
	b, _ := strconv.Atoi(values[1])
	fmt.Printf("%d\n", a+b)
}
