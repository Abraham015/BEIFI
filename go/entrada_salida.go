package main

import "fmt"

func main(){
	var(
		nombre string
		edad int
		pi float64
	)

	fmt.Println("Ingrese su nombre: ")
	fmt.Scanln(&nombre)

	fmt.Println("Ingrese su edad: ")
	fmt.Scanln(&edad)

	fmt.Println("Ingrese el valor de pi: ")
	fmt.Scanln(&pi)

	fmt.Printf("Nombre: %s\nEdad: %d\nValor de Pi: %f", nombre, edad, pi)
}