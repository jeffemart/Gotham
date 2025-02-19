package main

import "fmt"

type TestT struct {
	ID int
}

func Test() (*TestT, error) {
	return &TestT{}, nil
}

func TestS(s *TestT) error {
	return nil
}

func main() {
 t, err := Test()
 if err != nil {
	panic(err)
 }
 fmt.Println(t)
 
 var T TestT
 if err := TestS(&T); err != nil {
	panic(err)
 }
}
