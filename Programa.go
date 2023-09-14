package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"bufio"
	"io"
	"math"
	"path/filepath"
	"errors"
)

func main() {
    folder := "./Files"
    var typeDistance string
    var n int
    var typeFile string

    files, err := ioutil.ReadDir(folder)
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, fileInfo := range files {
        if !fileInfo.IsDir() {
            fmt.Println("------------------------------------------")
            fmt.Println("El nombre del archivo es:", fileInfo.Name())
            // Abre el archivo utilizando os.Open
            file, err := os.Open(filepath.Join(folder, fileInfo.Name()))
            if err != nil {
                fmt.Println("Error al abrir el archivo:", err)
                continue
            }
            defer file.Close() // Asegura que el archivo se cierre al salir de la función

            n, err = NumberCities(file)
            if err != nil {
                fmt.Println("Error al obtener el número de ciudades:", err)
                continue
            }
            fmt.Println("El número de ciudades es", n)

            typeFile = TypeProblem(file)

            if typeFile == "ATSP" {
                file.Seek(0, 0) // Asegura que el archivo esté al principio antes de abrirlo nuevamente
                typeDistance = TypeDistanceATSP(file)
                switch typeDistance {
                case "Explicita":
                    ProblemaExplicit(file, n)
                }
            } else {
                file.Seek(0, 0) // Asegura que el archivo esté al principio antes de abrirlo nuevamente
                typeDistance = TypeDistance(file)
				fmt.Println("El tipo de problema es: "+typeDistance)
                switch typeDistance {
                case "Euclidiana":
                    ProblemaEuclidiano(filepath.Join(folder, fileInfo.Name()), n)
                case "Geografica":
                    ProblemaGeografico(filepath.Join(folder, fileInfo.Name()), n)
                case "Circular":
                    ProblemaCeil(file, n)
                case "ATT":
                    ProblemaATT(filepath.Join(folder, fileInfo.Name()), n)
                case "DiagonalSuperior":
                    ProblemaSuperior(filepath.Join(folder, fileInfo.Name()), n)
                case "DiagonalInferior":
                    ProblemaInferior(filepath.Join(folder, fileInfo.Name()), n)
                }
            }
        }
    }
}

func NumberCities(file *os.File) (int, error) {
    defer file.Seek(0, 0) // Asegura que el archivo esté al principio antes de leerlo

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "DIMENSION") {
            parts := strings.Split(line, ":")
            if len(parts) != 2 {
                return 0, errors.New("Formato DIMENSION incorrecto")
            }
            dimensionStr := strings.TrimSpace(parts[1])
            dimension, err := strconv.Atoi(dimensionStr)
            if err != nil {
                return 0, err
            }
            return dimension, nil
        }
    }

    return 0, errors.New("No se encontró la dimensión")
}

func TypeProblem(file *os.File) string {
    typeProblem := ""

    // Obtén el nombre y la ruta del archivo
    fileName := file.Name()

    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
        return typeProblem
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.TrimSpace(line) == "TYPE : TSP" {
            typeProblem = "TSP"
            break
        } else if strings.TrimSpace(line) == "TYPE: ATSP" {
            typeProblem = "ATSP"
            break
        }
    }

    return typeProblem
}


func TypeDistance(file *os.File) string {
	distance := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "EDGE_WEIGHT_TYPE : EUC_2D":
			distance = "Euclidiana"
		case "EDGE_WEIGHT_TYPE: GEO":
			distance = "Geografica"
		case "EDGE_WEIGHT_TYPE : ATT":
			distance = "ATT"
		case "EDGE_WEIGHT_TYPE : MATRIX":
			distance = "MATRIX"
		case "EDGE_WEIGHT_TYPE : CEIL_2D":
			distance = "Circular"
		case "EDGE_WEIGHT_TYPE: EXPLICIT":
			scanner.Scan() // Read the next line for type
			distance = "DiagonalInferior"
			if strings.Contains(scanner.Text(), "LOWER_DIAG_ROW") {
				distance = "DiagonalInferior"
			}else{
				distance = "DiagonalSuperior"
			}
		case "EDGE_WEIGHT_TYPE: EUC_2D":
			distance = "Euclidiana"
		}
	}

	defer file.Close()

	return distance
}

func TypeDistanceATSP(file *os.File) string {
	distance := ""
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		switch line {
		case "EDGE_WEIGHT_TYPE: EXPLICIT":
			distance = "Explicita"
		}
	}

	return distance
}

func CompletarMatriz(matriz [][]int, size int) {
	for i := 1; i < size; i++ {
		for j := 0; j < size; j++ {
			if matriz[i][j] != 0 {
				matriz[j][i] = matriz[i][j]
			}
		}
	}
}

func ImprimirMatriz(matriz [][]int, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			print(matriz[i][j], " ")
		}
		println()
	}
}

func EscribirMatrizEnCSV(nombreArchivo string, matriz [][]int) {
	file, err := os.Create(nombreArchivo)
	if err != nil {
		println("Error al crear el archivo CSV:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, fila := range matriz {
		for columna, valor := range fila {
			writer.WriteString(strconv.Itoa(valor))
			if columna < len(fila)-1 {
				writer.WriteString(" ")
			}
		}
		writer.WriteString("\n")
	}
}

func ReadFileATT(file *os.File, numbernode []string, firstnode, secondnode []int) error {
	count := 0
	flag := 0

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if flag == 0 {
			if strings.TrimSpace(line) == "NODE_COORD_SECTION" || strings.TrimSpace(line) == "NODE_COORD_SECTION " {
				flag++
			}
			if strings.TrimSpace(line) == "DISPLAY_DATA_SECTION" {
				flag++
			}
		} else {
			if strings.TrimSpace(line) == "EOF" {
				continue
			} else {
				if count < len(numbernode) {
					aux := strings.TrimSpace(line)
					datos := strings.Fields(aux)
					numbernode[count] = datos[0]
					firstnode[count], _ = strconv.Atoi(datos[1])
					secondnode[count], _ = strconv.Atoi(datos[2])
				}
				count++
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func DistanciaATT(numbernode []string, x, y []int) int {
	xd := make([]float64, len(numbernode))
	yd := make([]float64, len(numbernode))
	dij := 0
	tij := 0
	var rij float64

	for i := 0; i < len(numbernode); i++ {
		if i < len(numbernode)-1 {
			xd[i] = float64(x[i] - x[i+1])
			yd[i] = float64(y[i] - y[i+1])
		}
	}

	for i := 0; i < len(yd); i++ {
		rij = math.Sqrt((math.Pow(xd[i], 2) + math.Pow(yd[i], 2)) / 10.0)
		tij = int(rij + 0.5)
		if tij < int(rij) {
			dij += tij + 1
		} else {
			dij += tij
		}
	}

	return dij
}

func ProblemaATT(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()
	numbernode := make([]string, n)
	firstnode := make([]int, n)
	secondnode := make([]int, n)

	ReadFileATT(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", DistanciaATT(numbernode, firstnode, secondnode))
}

func ReadFileCeil(file *os.File, numbernode []string, firstnode, secondnode []float64) error {
	count := 0
	flag := 0

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if flag == 0 {
			if strings.TrimSpace(line) == "NODE_COORD_SECTION" || strings.TrimSpace(line) == "NODE_COORD_SECTION " {
				flag++
			}
			if strings.TrimSpace(line) == "DISPLAY_DATA_SECTION" {
				flag++
			}
		} else {
			if strings.TrimSpace(line) == "EOF" {
				continue
			} else {
				if count < len(numbernode) {
					aux := strings.TrimSpace(line)
					datos := strings.Fields(aux)
					numbernode[count] = datos[0]
					firstnode[count], _ = strconv.ParseFloat(datos[1], 64)
					secondnode[count], _ = strconv.ParseFloat(datos[2], 64)
				}
				count++
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func DistanciaCeil(nodes []string, x, y []float64) int {
	var x1, x2, y1, y2 float64
	distance := 0

	for i := 0; i < len(nodes); i++ {
		if i < len(nodes)-1 {
			x1 = x[i]
			x2 = x[i+1]
			y1 = y[i]
			y2 = y[i+1]
			distance += int(math.Sqrt(math.Pow(math.Abs(y2-y1), 2) + math.Pow(math.Abs(x2-x1), 2)) + 0.5)
		}
	}

	return distance
}

func ProblemaCeil(file *os.File, n int) {
	numbernode := make([]string, n)
	firstnode := make([]float64, n)
	secondnode := make([]float64, n)

	err := ReadFileCeil(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", DistanciaCeil(numbernode, firstnode, secondnode))
}

func ReadFileEuclidiana(file *os.File, numbernode []string, firstnode, secondnode []int) error {
	count := 0
	flag := 0
	//i := 1

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if flag == 0 {
			if strings.TrimSpace(line) == "NODE_COORD_SECTION" {
				flag++
			}
		} else {
			if strings.TrimSpace(line) == "EOF" {
				break
			} else {
				aux := strings.TrimSpace(line)
				datos := strings.Fields(aux)
				numbernode[count] = datos[0]
				if strings.Contains(datos[1], "e") {
					numbers := strings.Split(datos[1], "e")
					num, _ := strconv.ParseFloat(numbers[0], 64)
					exp, _ := strconv.ParseFloat(numbers[1], 64)
					firstnode[count] = int(num * math.Pow(10, exp))
					numbers = strings.Split(datos[2], "e")
					num, _ = strconv.ParseFloat(numbers[0], 64)
					exp, _ = strconv.ParseFloat(numbers[1], 64)
					secondnode[count] = int(num * math.Pow(10, exp))
				} else {
					if len(datos[1]) == 0 {
						nums := strings.Fields(datos[2])
						firstnode[count], _ = strconv.Atoi(nums[0])
						if len(nums) > 1 {
							secondnode[count], _ = strconv.Atoi(nums[1])
						} else {
							secondnode[count], _ = strconv.Atoi(nums[0])
						}
					} else {
						if strings.Contains(datos[1], ".") {
							num1, _ := strconv.ParseFloat(datos[1], 64)
							num2, _ := strconv.ParseFloat(datos[2], 64)
							firstnode[count] = int(math.Round(num1))
							secondnode[count] = int(math.Round(num2))
						} else {
							firstnode[count], _ = strconv.Atoi(datos[1])
							secondnode[count], _ = strconv.Atoi(datos[2])
						}
					}
				}
				count++
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func DistanciaEuclidiana(nodes []string, x, y []int) int {
	var x1, x2, y1, y2 int
	distance := 0

	for i := 0; i < len(nodes); i++ {
		if i < len(nodes)-1 {
			x1 = x[i]
			x2 = x[i+1]
			y1 = y[i]
			y2 = y[i+1]
			distance += int(math.Sqrt(math.Pow(float64(math.Abs(float64(y2-y1))), 2) + math.Pow(float64(math.Abs(float64(x2-x1))), 2)) + 0.5)
		}
	}

	return distance
}

func ProblemaEuclidiano(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()
	numbernode := make([]string, n)
	firstnode := make([]int, n)
	secondnode := make([]int, n)

	ReadFileEuclidiana(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", DistanciaEuclidiana(numbernode, firstnode, secondnode))
}

func DistanciaGeografica(numbernode []string, x, y []float64) int {
	distance:=0
	var q1, q2, q3 	
	RRR := 6378.388 // Radio de la Tierra en km
	pi :=  3.141592

	latitude := make([]float64, len(numbernode))
	longitude := make([]float64, len(numbernode))

	for i := 0; i < len(numbernode); i++ {
		deg := int(x[i] + 0.5)
		min := x[i] - float64(deg)
		latitude[i] = pi * (float64(deg) + 5.0*min/3.0) / 180.0
		deg = int(y[i] + 0.5)
		min = y[i] - float64(deg)
		longitude[i] = pi * (float64(deg) + 5.0*min/3.0) / 180.0
	}

	for i := 0; i < len(numbernode); i++ {
		if i < len(latitude)-1 {
			q1 = math.Cos(longitude[i] - longitude[i+1])
			q2 = math.Cos(latitude[i] - latitude[i+1])
			q3 = math.Cos(latitude[i] + latitude[i+1])
			distance += int(RRR * math.Acos(0.5*((1.0+q1)*q2-(1.0-q1)*q3)) + 1.0)
		}
	}

	return distance
}

func ProblemaGeografico(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()
	numbernode := make([]string, n)
	firstnode := make([]float64, n)
	secondnode := make([]float64, n)

	ReadFileGeografico(file, numbernode, firstnode, secondnode, n)

	fmt.Println("La distancia total para este archivo es de", DistanciaGeografica(numbernode, firstnode, secondnode))
}

func ReadFileGeografico(file *os.File, numbernode[] string, firstnode[] float64, secondnode[] float64, n int){
	count:=0
	flag:=0

	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		linea:=scanner.Text()
		if flag==0{
			if strings.HasPrefix(linea, "NODE_COORD_SECTION") {
				flag++
			}

			if strings.HasPrefix(linea, "DISPLAY_DATA_SECTION") {
				flag++
			}
		} else{
			if strings.HasPrefix(linea, "EOF") {
				break
			}

			if(count<len(firstnode)){
				aux := strings.Fields(linea)
				//fmt.Println("La línea de datos es ",aux)
				numbernode[count]=aux[0]
				firstnode[count],_=strconv.ParseFloat(aux[1], 32)
				secondnode[count],_=strconv.ParseFloat(aux[2], 32)
				/*fmt.Println("El número de ciudad es ", numbernode[count])
				fmt.Println("El primer nodo es ",firstnode[count])
				fmt.Println("El segundo nodo es ",secondnode[count])*/
			}
			count++
		}
	}	
}

func ReadFileInferior(file *os.File, matriz [][]int, n int) {
	leyendoMatriz := false
	columna:=0
	fila:=0

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		linea := scanner.Text()

		if strings.HasPrefix(linea, "EDGE_WEIGHT_SECTION") {
			leyendoMatriz = true
			continue
		}

		if leyendoMatriz && !strings.HasPrefix(linea, "EOF") {
			valoresStr := strings.Fields(linea)
				for _, valorStr := range valoresStr {
						//fmt.Println("El valor es "+valorStr)
						valor, err := strconv.Atoi(valorStr)
						if err != nil {
							fmt.Println("Error al convertir el valor:", err)
							return
						}
						if valor==0{
							if columna<n{
								for columna < n{
									if columna<n{
										matriz[fila][columna]=0
										columna++
									}else{
										break
									}
								}
							}
						}else{
							if columna<n{
								matriz[fila][columna]=valor
								columna++
							}
							
						}

						if columna == n{
							fila++
							columna=0
						}

						if fila==n{
							break
						}
				}
				if fila==n{
					break
				}
		}

		for fila<n{
			if fila<n{
				break
			}
			for columna=0; columna < n; columna++{
				matriz[fila][columna]=0
			}
			fila++
		}

		if strings.HasPrefix(linea, "EOF") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	/*scanner := bufio.NewScanner(file)

    // Variables para rastrear la fila y la columna en la matriz
    fila := 0
    columna := 0

    // Variable para determinar si debemos continuar en la misma fila
    continuarEnLaMismaFila := false

    // Variable para almacenar la línea anterior
    lineaAnterior := ""

    for scanner.Scan() {
        linea := scanner.Text()

        if linea == "EOF" {
            break
        }

        // Dividir la línea en valores
        valoresStr := strings.Fields(linea)

        for _, valorStr := range valoresStr {
            valor, err := strconv.Atoi(valorStr)
            if err != nil {
                fmt.Println("Error al convertir el valor:", err)
                return
            }

            // Si encontramos un cero, llenamos la fila con ceros y avanzamos a la siguiente fila
            if valor == 0 {
                continuarEnLaMismaFila = false
                fila++ // Avanzar a la siguiente fila
                columna = 0 // Reiniciar la columna
            } else {
                // Almacenar el valor en la matriz
                matriz[fila][columna] = valor
                columna++
                continuarEnLaMismaFila = true
            }

            // Comprobar si debemos continuar en la misma fila
            if continuarEnLaMismaFila {
                // Si la fila anterior tenía valores, llenar con ceros la parte inferior
                if len(valoresStr) < n {
                    valoresAnterior := strings.Fields(lineaAnterior)
                    for columna < n && len(valoresAnterior) > columna {
                        matriz[fila][columna] = 0
                        columna++
                    }
                }
            }
        }

        // Almacenar la línea actual como línea anterior
        lineaAnterior = linea
    }*/
}

func ReadFileSuperior(file *os.File, matriz [][]int, numbernode int) {
	scanner := bufio.NewScanner(file)

	// Variables para rastrear la fila
	fila := 0

	// Bucle principal para leer el archivo
	for scanner.Scan() {
		lineaActual := scanner.Text()
		lineaActual = strings.TrimSpace(lineaActual)

		if fila >= 8 && lineaActual != "EOF" {
			linea := strings.Fields(lineaActual)
			for columna, valorStr := range linea {
				valor, err := strconv.Atoi(valorStr)
				if err != nil {
					fmt.Println("Error al convertir el valor:", err)
					return
				}
				matriz[fila-8][columna] = valor
			}
		}

		fila++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}

// The function calculates the shortest distance between nodes in a matrix using the nearest neighbor
// algorithm.
func distanciaMATRIX(matrix [][]int, n int) {
	visitados := make([]int, n)
	ruta := make([]int, n)
	mejorRuta := make([]int, n)
	mejorDistancia := math.MaxInt
	distanciaActual:=0
	element:=-1

	for i := 0; i < n; i++ {
		visitados[i] = 0
		ruta[i] = -1
	}

	visitados[0] = 1
	ruta[0] = 0
	posActual:= 0

	for true{
		element = -1

		for i := 0; i < n; i++ {
			if visitados[i] == 0 {
				if element == -1 {
					element = i
				}
				distanciaActual = matrix[posActual][i]
				if distanciaActual < matrix[posActual][element] {
					element = i
				}
			}
		}

		if element == -1 {
			break
		}

		visitados[element] = 1
		ruta[element] = posActual
		posActual = element
	}

	distanciaActual = matrix[posActual][0]
	//fmt.Println("La distancia actual es: ",distanciaActual)
	ruta[n-1] = 0
	ruta[0] = posActual

	for i := 0; i < n; i++ {
		if i != n-1 {
			distanciaActual += matrix[ruta[i]][ruta[i+1]]
		}
	}

	if distanciaActual < mejorDistancia {
		mejorDistancia = distanciaActual
		for i:=0; i < n; i++{
			mejorRuta[i]=ruta[i]
		}
	}

	fmt.Println("La distancia total es: ", mejorDistancia)
}

func ProblemaSuperior(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()

	matriz := make([][]int, n)
	for i := 0; i < n; i++ {
		matriz[i] = make([]int, n)
	}

	fmt.Println("Entra a esta funcion de superior")

	ReadFileSuperior(file, matriz, n)	

	fmt.Println("Matriz Superior:")
	ImprimirMatriz(matriz, n)

	EscribirMatrizEnCSV("test.csv",matriz)	

	distanciaMATRIX(matriz, n)
}

// The function "ProblemaInferior" reads a matrix from a file, prints the lower triangular part of the
// matrix, completes the matrix, prints the complete matrix, and calculates the distance between each
// pair of elements in the matrix.
func ProblemaInferior(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()

	matriz := make([][]int, n)
	for i := 0; i < n; i++ {
		matriz[i] = make([]int, n)
	}

	ReadFileInferior(file,matriz,n)

	/*fmt.Println("Matriz Inferior:")
	ImprimirMatriz(matriz, n)*/

	EscribirMatrizEnCSV("./Files/Excel/LeerMatriz.csv",matriz)	

	//Se completa la matriz
	CompletarMatriz(matriz, n)
	/*fmt.Println("Matriz Completa:")
	ImprimirMatriz(matriz, n)*/

	EscribirMatrizEnCSV("./Files/Excel/MatrizCompleta.csv",matriz)	

	distanciaMATRIX(matriz, n)
}

func SolveATSP(distances [][]int, numCities int, tour []int) int {
	visited := make([]bool, numCities)
	startCity := 0 // Puedes elegir cualquier ciudad como punto de inicio

	tour = append(tour, startCity)
	visited[startCity] = true
	totalDistance := 0

	// Construir el recorrido visitando el vecino más cercano disponible
	for i := 0; i < numCities-1; i++ {
		currentCity := tour[i]
		nextCity := findNearestNeighbor(currentCity, visited, distances)
		tour = append(tour, nextCity)
		visited[nextCity] = true
		totalDistance += distances[currentCity][nextCity]
	}

	// Regresar al punto de inicio para completar el ciclo
	tour = append(tour, startCity)
	totalDistance += distances[tour[numCities-1]][startCity]

	return totalDistance
}

func findNearestNeighbor(city int, visited []bool, distances [][]int) int {
	nearestNeighbor := -1
	minDistance := math.MaxInt32

	for i := 0; i < len(visited); i++ {
		if !visited[i] && distances[city][i] < minDistance {
			minDistance = distances[city][i]
			nearestNeighbor = i
		}
	}

	return nearestNeighbor
}

func ReadExplicit(file *os.File, data [][]int, n int) error {
    reader := bufio.NewReader(file)

    // Ignorar las primeras 7 líneas del archivo
    for i := 0; i < 7; i++ {
        _, err := reader.ReadString('\n')
        if err != nil {
            return err
        }
    }

    row := 0
    col := 0
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }

        line = strings.TrimSpace(line)
        if line == "EOF" {
            break
        }

        values := strings.Fields(line)
        for _, valueStr := range values {
            value, err := strconv.Atoi(valueStr)
            if err != nil {
                return err
            }

            if value == 9999 {
                // Saltar al siguiente elemento de la siguiente fila
                row++
                col = 0
            } else {
                data[row][col] = value
                col++
            }

            // Comprobar si hemos llenado completamente la matriz
            if row == n {
                break
            }
        }

        // Comprobar si hemos llenado completamente la matriz
        if row == n {
            break
        }
    }

    // Rellenar con ceros si la matriz no se ha llenado completamente
    for row < n {
        for col = 0; col < n; col++ {
            data[row][col] = 0
        }
        row++
    }

    return nil
}

func ProblemaExplicit(file *os.File, n int) {
    data := make([][]int, n)
    for i := range data {
        data[i] = make([]int, n)
    }

    err := ReadExplicit(file, data, n)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    var tour []int
    totalDistance := SolveATSP(data, n, tour)

    fmt.Println("Recorrido óptimo:")
    for _, city := range tour {
        fmt.Printf("%d -> ", city)
    }
    fmt.Println("0")

    fmt.Println("Distancia total recorrida:", totalDistance)
}
