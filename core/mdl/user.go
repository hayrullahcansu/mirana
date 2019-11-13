package mdl

type User struct {
	UserId    string
	Name      string
	Balance   float32
	Win       int
	Lose      int
	Push      int
	Blackjack int
}

func NewUser(id string, name string) *User {
	return &User{
		UserId:  id,
		Balance: 3000,
		Name:    name,
	}
}
