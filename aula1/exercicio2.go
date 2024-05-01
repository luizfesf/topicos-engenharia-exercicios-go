package main

import 
(
    "errors"
    "fmt"
)


func main()  {
    str := "Hello world, Luiz Felipe!"
    lenght, err := CountString(str)

        if err != nil {
            fmt.Println("Erro: ", err)  
        } else {
                fmt.Println("Tamanho da String: ", lenght)
            }
}

func CountString(s string) (int, error) {
    if s == " " {
        return 0, errors.New("String está vazia")
    }
    return len(s), nil
}