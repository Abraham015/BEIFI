package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	/**************************************************************/
	// Se crea la ruta para los archivos
	folder := "./Files"
	/**************************************************************/
	/* Se declaran las variables del programa */
	var typeDistance string
	var n int
	var typeFile string
	/**************************************************************/
	/* Se crean los objetos de las clases correspondientes de TSP */
	euclidiana := Euclidiana{}
	geografico := Geografico{}
	ceil := Ceil{}
	att := ATT{}
	ma := Matrix{}
	tool := Tools{}
	g := General{}
	/**************************************************************/
	/* Se crean los objetos de las clases correspondientes de ATSP */
	explicit := Explicit{}
	/**************************************************************/
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			fmt.Println("------------------------------------------")
			fmt.Println("El nombre del archivo es:", file.Name())
			n = g.NumberCities(file)
			fmt.Println("El número de ciudades es", n)
			typeFile = g.TypeProblem(file)

			if typeFile == "ATSP" {
				typeDistance = tool.TypeDistanceATSP(file)
				switch typeDistance {
				case "Explicita":
					explicit.ProblemaExplicit(file, n)
				}
			} else {
				typeDistance = tool.TypeDistance(file)
				switch typeDistance {
				case "Euclidiana":
					euclidiana.ProblemaEuclidiano(file, n)
				case "Geografica":
					geografico.ProblemaGeografico(file, n)
				case "Circular":
					ceil.ProblemaCeil(file, n)
				case "ATT":
					att.ProblemaATT(file, n)
				case "DiagonalSuperior":
					ma.ProblemaSuperior(file, n)
				case "DiagonalInferior":
					ma.ProblemaInferior(file, n)
				}
			}
		}
	}
}

type General struct{}

func (g *General) NumberCities(file os.FileInfo) int {
	cities := 0
	name := file.Name()
	var number string

	for _, c := range name {
		if c >= '0' && c <= '9' {
			number += string(c)
		}
	}

	cities, _ = strconv.Atoi(number)
	return cities
}

func (g *General) TypeProblem(file os.FileInfo) string {
	typeProblem := ""
	data, err := ioutil.ReadFile(file.Name())
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

type General struct{}

func (g *General) NumberCities(file *os.File) int {
	cities := 0
	name := file.Name()
	var number string

	for _, c := range name {
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			number += string(c)
		}
	}

	cities, _ = strconv.Atoi(number)
	return cities
}

func (g *General) TypeProblem(file *os.File) string {
	var t string

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}

		if strings.TrimSpace(line) == "TYPE : TSP" {
			t = "TSP"
			break
		} else if strings.TrimSpace(line) == "TYPE: ATSP" {
			t = "ATSP"
			break
		}

		if err == io.EOF {
			break
		}
	}

	return t
}

type Tools struct{}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) TypeDistance(file *os.File) string {
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
				distance = "DiagonalSuperior"
			}
		case "EDGE_WEIGHT_TYPE: EUC_2D":
			distance = "Euclidiana"
		}
	}

	return distance
}

func (t *Tools) TypeDistanceATSP(file *os.File) string {
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

func (t *Tools) CompletarMatriz(matriz [][]int, size int) {
	for i := 1; i < size; i++ {
		for j := 0; j < size; j++ {
			if matriz[i][j] != 0 {
				matriz[j][i] = matriz[i][j]
			}
		}
	}
}

func (t *Tools) CalcularCostoCamino(matriz [][]int, camino []int) int {
	costoTotal := 0
	n := len(camino)

	for i := 0; i < n-1; i++ {
		nodoActual := camino[i]
		nodoSiguiente := camino[i+1]
		costoTotal += matriz[nodoSiguiente][nodoActual]
	}

	return costoTotal
}

func (t *Tools) ImprimirMatriz(matriz [][]int, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			print(matriz[i][j], " ")
		}
		println()
	}
}

func (t *Tools) EscribirMatrizEnCSV(nombreArchivo string, matriz [][]int) {
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

type ATT struct{}

func (a *ATT) ReadFileATT(file *os.File, numbernode []string, firstnode, secondnode []int) error {
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

func (a *ATT) DistanciaATT(numbernode []string, x, y []int) int {
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

func (a *ATT) ProblemaATT(file *os.File, n int) {
	numbernode := make([]string, n)
	firstnode := make([]int, n)
	secondnode := make([]int, n)

	err := a.ReadFileATT(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", a.DistanciaATT(numbernode, firstnode, secondnode))
}

type Ceil struct{}

func (c *Ceil) ReadFileCeil(file *os.File, numbernode []string, firstnode, secondnode []float64) error {
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

func (c *Ceil) DistanciaCeil(nodes []string, x, y []float64) int {
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

func (c *Ceil) ProblemaCeil(file *os.File, n int) {
	numbernode := make([]string, n)
	firstnode := make([]float64, n)
	secondnode := make([]float64, n)

	err := c.ReadFileCeil(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", c.DistanciaCeil(numbernode, firstnode, secondnode))
}

type Euclidiana struct{}

func (e *Euclidiana) ReadFileEuclidiana(file *os.File, numbernode []string, firstnode, secondnode []int) error {
	count := 0
	flag := 0
	i := 1

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

func (e *Euclidiana) DistanciaEuclidiana(nodes []string, x, y []int) int {
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

func (e *Euclidiana) ProblemaEuclidiano(file *os.File, n int) {
	numbernode := make([]string, n)
	firstnode := make([]int, n)
	secondnode := make([]int, n)

	err := e.ReadFileEuclidiana(file, numbernode, firstnode, secondnode)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("La distancia total para este archivo es de", e.DistanciaEuclidiana(numbernode, firstnode, secondnode))
}

type Geografica struct{}

func (g *Geografica) DistanciaGeografica(numbernode []string, x, y []float64) int {
	var distance int
	var q1, q2, q3 float64
	RRR := 6378.388 // Radio de la Tierra en km
	pi := math.Pi

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

func (g *Geografica) ProblemaGeografico(file *os.File, n int) {
	numbernode := make([]string, n)
	firstnode := make([]float64, n)
	secondnode := make([]float64, n)

	// Aquí deberías implementar la función ReadFileGeografica en Go para leer los datos del archivo

	fmt.Println("La distancia total para este archivo es de", g.DistanciaGeografica(numbernode, firstnode, secondnode))
}

type Matrix struct {
	Tool Tools
}

func (m *Matrix) ReadFileInferior(file *os.File, matriz [][]int, n int) {
	// Implementa la función ReadFileInferior aquí
}

func (m *Matrix) ReadFileSuperior(file *os.File, matriz [][]int, numbernode int) {
	// Implementa la función ReadFileSuperior aquí
}

func (m *Matrix) DistanciaMATRIX(matrix [][]int, n int) {
	// Implementa la función DistanciaMATRIX aquí
}

func (m *Matrix) ProblemaSuperior(file *os.File, n int) {
	data := make([][]int, n)
	for i := 0; i < n; i++ {
		data[i] = make([]int, n)
	}

	// Implementa la lectura del archivo y llamada a ReadFileSuperior aquí

	fmt.Println("Matriz Superior:")
	m.Tool.ImprimirMatriz(data, n)

	// Implementa el llamado a Tool.EscribirMatrizEnCSV aquí

	m.DistanciaMATRIX(data, n)
}

func (m *Matrix) ProblemaInferior(file *os.File, n int) {
	data := make([][]int, n)
	for i := 0; i < n; i++ {
		data[i] = make([]int, n)
	}

	// Implementa la lectura del archivo y llamada a ReadFileInferior aquí

	fmt.Println("Matriz Inferior:")
	m.Tool.ImprimirMatriz(data, n)

	// Implementa el llamado a Tool.EscribirMatrizEnCSV aquí

	// Completa la matriz con ceros usando Tool.CompletarMatriz

	fmt.Println("Matriz Completa:")
	m.Tool.ImprimirMatriz(data, n)

	// Implementa el llamado a Tool.EscribirMatrizEnCSV aquí

	m.DistanciaMATRIX(data, n)
}

type Distance struct{}

func (d *Distance) SolveATSP(distances [][]int, numCities int, tour []int) int {
	visited := make([]bool, numCities)
	startCity := 0 // Puedes elegir cualquier ciudad como punto de inicio

	tour = append(tour, startCity)
	visited[startCity] = true
	totalDistance := 0

	// Construir el recorrido visitando el vecino más cercano disponible
	for i := 0; i < numCities-1; i++ {
		currentCity := tour[i]
		nextCity := d.findNearestNeighbor(currentCity, visited, distances)
		tour = append(tour, nextCity)
		visited[nextCity] = true
		totalDistance += distances[currentCity][nextCity]
	}

	// Regresar al punto de inicio para completar el ciclo
	tour = append(tour, startCity)
	totalDistance += distances[tour[numCities-1]][startCity]

	return totalDistance
}

func (d *Distance) findNearestNeighbor(city int, visited []bool, distances [][]int) int {
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

type Distance struct{}

func (d *Distance) SolveATSP(distances [][]int, numCities int, tour []int) int {
	// Implementa la lógica de SolveATSP aquí (código previamente proporcionado)
}

type Explicit struct {
	distance Distance
}

func (e *Explicit) ReadExplicit(file string, data [][]int, n int) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	row := 0
	col := 0

	for _, line := range lines[7:] {
		if line == "EOF" {
			break
		}

		values := strings.Fields(line)
		for _, value := range values {
			num := 0
			if value != "9999" {
				num = parseValue(value)
			}

			if num == 9999 {
				row++
				col = 0
			} else {
				data[row][col] = num
				col++
			}

			if row == n {
				break
			}
		}

		if row == n {
			break
		}
	}

	return nil
}

func parseValue(value string) int {
	// Implementa la conversión del valor de string a int aquí
	// Puedes usar strconv.Atoi() u otra función de conversión adecuada
}

func (e *Explicit) ProblemaExplicit(file string, n int) {
	data := make([][]int, n)
	for i := range data {
		data[i] = make([]int, n)
	}

	err := e.ReadExplicit(file, data, n)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var tour []int
	totalDistance := e.distance.SolveATSP(data, n, tour)

	fmt.Println("Recorrido óptimo:")
	for _, city := range tour {
		fmt.Printf("%d -> ", city)
	}
	fmt.Println("0")

	fmt.Println("Distancia total recorrida:", totalDistance)
}
