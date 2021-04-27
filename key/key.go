package key

import (
	"bufio"
	"fmt"
	"os"
)

var (
	Key string
	Keyword string
)

func Init() string{
	fmt.Print("Search for > ")
	//fmt.Fscan(
	//	os.Stdin,
	//	&Key,
	//)
	Key := bufio.NewScanner(os.Stdin)
	Key.Scan()
	if err := Key.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка ввода:", err)
	}
	Keyword=Key.Text()
	return Keyword
}
