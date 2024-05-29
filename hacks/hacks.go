package main

import "fmt"

type User struct {
	Id       string
	Password string
}

func (u User) Format(f fmt.State, verb rune) {
	//  Put what you want to show on screen
	f.Write([]byte(fmt.Sprintf("ID: %s", u.Id)))
}

// If you don't want to display password when you print User
// Output: User ID: UUID

func main() {
	user := User{Id: "UUID", Password: "I will not ganna tell you"}
	fmt.Print("User ", user)
}
