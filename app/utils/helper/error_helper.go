package helper

import "fmt"

func Recover(location string) {
	if r := recover(); r != nil {
		fmt.Printf("recover panic action from %s : %s\n", location, r)
	}
}
