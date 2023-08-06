package main

import "fmt"

func main(){
	/*if false{ //No se cumple la condicion
		fmt.Println("Se cumple la condicion")
	}else{
		fmt.Println("No se cumple la condición")
	}*/

	var num int
	fmt.Println("Digite un numero")
	fmt.Scanln(&num)

	if num==0{
		fmt.Println("Es un numero neutro")
	}else if num%2==0{
		fmt.Println("Es un numero par")
	}else{
		fmt.Println("Es un numero impar")
	}
}