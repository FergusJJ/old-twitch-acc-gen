package datastructures

type Account struct {
	Email        string
	Password     string
	UserName     string
	UserPassword string
}

var AccountsMap = make(map[int]*Account)
