package main

import ( 
  "fmt"
)

func CheckIdade(idade int) bool {
	if idade < 18 {
	  return false
	}
  return true
}

func main() {
  fmt.println(CheckIdade(18))
}