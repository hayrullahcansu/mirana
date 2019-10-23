package xxx

import "fmt"

type BasePlayer struct {
	Listener
}

func (b *BasePlayer) Test() {

}

func (b *BasePlayer) TestMessage(input string) string {
	return input
}

func (b *BasePlayer) Reis() {
	fmt.Println("HAHAHAHAHAHA")
}
