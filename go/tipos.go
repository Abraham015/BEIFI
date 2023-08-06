package main

import(
	"fmt"
	"strconv"
)

func main(){
	edad:=22
	//Para convertir un int a string
	edad_str:=strconv.Itoa(edad)

	fmt.Println("Mi edad es "+edad_str)
}