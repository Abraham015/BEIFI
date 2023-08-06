package main

import "fmt"

func main(){
	//Un slice es una estructura dinamica, como sabemos
	//los arreglos son de un tamaño definido y no se
	//puede modificar dicho tamaño 
	//La forma de declarar un slice es
	matriz:=[]int{1,2}
	//Para el null en go es nil
	if matriz==nil{
		fmt.Println("Esta vacio")
	}else{
		fmt.Println(len(matriz))
	}

	arreglo:=[]int{1,2,3}
	fmt.Println(arreglo)

	slice:=arreglo[:2]
	fmt.Println(slice)
}