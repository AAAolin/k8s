package config

import "time"

const (
	Username = "admin"
	Password = "123456"
)

type Student struct {
	PersonValue Person
	Metha       meth
	age         int
}
type Person struct {
	name string
}

type meth interface {
	get()
	post()
}

func test() {
	var s Student
	s.PersonValue.name = "aaa"
	s.Metha.get()
}
