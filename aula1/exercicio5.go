package main

import "fmt"

func main() {
  a := [10]int{1, 20, 3, 4, 5, 6, 7, 8, 9, 0}

    m := make(map[int]int)

    for i := 0; i < len(a); i++ {
        m[i] = a[i]
    }

        fmt.Println(m)
}