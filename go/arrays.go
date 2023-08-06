package main

import "fmt"

func main(){
	//en go es mas comun utilizar slices
	//La sintaxis de un array es la siguiente
	var ages[3] int
	//Para imprimir el valor de los elementos del arreglo son:
	fmt.Println(ages[0])
	fmt.Println(len(ages)) //Elementos dentro del arreglo

	//Para declarar un valor inicializado es:
	var ages2[3] int=[3]int{10,20,30}
	fmt.Println(ages2)

	//Nosotros podemos omitir el tamaño del arreglo como:
	test:=[...]int{10:12}
	fmt.Println(test)

	//Nosotros podemos comparar dos arreglos
}