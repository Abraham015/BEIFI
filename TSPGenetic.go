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
	"math/rand"
	"time"
)

type Individual struct {
	Chromosome []int
	Fitness    int
}

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
                    ProblemaExplicit(filepath.Join(folder, fileInfo.Name()), n)
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
                    ProblemaCeil(filepath.Join(folder, fileInfo.Name()), n)
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

func generateRandomIndividual(numCities int) Individual {
	chromosome := rand.Perm(numCities)
	return Individual{Chromosome: chromosome}
}

func selectParent(population []Individual) Individual {
	// Método de selección: selección por torneo
	tournamentSize := 5
	tournament := make([]Individual, tournamentSize)
	for i := 0; i < tournamentSize; i++ {
		tournament[i] = population[rand.Intn(len(population))]
	}
	bestIndividual := tournament[0]
	for _, individual := range tournament {
		if individual.Fitness < bestIndividual.Fitness {
			bestIndividual = individual
		}
	}
	return bestIndividual
}

func crossover(parent1, parent2 Individual) Individual {
	// Cruzamiento de un punto
	crossoverPoint := rand.Intn(len(parent1.Chromosome))
	childChromosome := make([]int, len(parent1.Chromosome))
	copy(childChromosome, parent1.Chromosome[:crossoverPoint])
	usedCities := make(map[int]bool)
	for _, city := range childChromosome {
		usedCities[city] = true
	}
	childIndex := crossoverPoint
	for _, city := range parent2.Chromosome {
		if !usedCities[city] {
			childChromosome[childIndex] = city
			childIndex++
		}
	}
	return Individual{Chromosome: childChromosome}
}

func mutate(individual Individual) {
	// Mutación de intercambio
	index1 := rand.Intn(len(individual.Chromosome))
	index2 := rand.Intn(len(individual.Chromosome))
	individual.Chromosome[index1], individual.Chromosome[index2] = individual.Chromosome[index2], individual.Chromosome[index1]
}

func DistanciaEuclidiana(x []int, y []int, chromosome []int) int {
	var x1, x2, y1, y2 int
	distance := 0

	/*for i := 0; i < len(y); i++ {
		if i < len(y)-1 {
			x1 = x[i]
			x2 = x[i+1]
			y1 = y[i]
			y2 = y[i+1]
			distance += int(math.Sqrt(math.Pow(float64(math.Abs(float64(y2-y1))), 2) + math.Pow(float64(math.Abs(float64(x2-x1))), 2)) + 0.5)
		}
	}*/

	//Se debe recorrer el cromosoma 
	for i:=0; i<len(chromosome)-1;i++{
		x1 = x[chromosome[i]]
		x2 = x[chromosome[i+1]]
		y1 = y[chromosome[i]]
		y2 = y[chromosome[i+1]]
		distance += int(math.Sqrt(math.Pow(float64(math.Abs(float64(y2-y1))), 2) + math.Pow(float64(math.Abs(float64(x2-x1))), 2)) + 0.5)
	}

	return distance
}

func GeneticEuclideano(xdistances []int, ydistances []int) {
	rand.Seed(time.Now().UnixNano())

	// Configuración del algoritmo genético
	populationSize := 1000
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(len(xdistances))
	}
	
	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		// Calcular fitness de la población
		for i := range population {
			population[i].Fitness = DistanciaEuclidiana(xdistances, ydistances, population[i].Chromosome)
		}

		/*for i:=range population{
			for j:=0; j<populationSize-1; i++{
				fmt.Printf("%d ", population[i].Chromosome[j])
			}
			fmt.Printf("Fitness: %d\n",population[i].Fitness)
		}*/

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}
		
		// Imprimir la mejor distancia en esta generación
		//fmt.Printf("Generación %d - Mejor Distancia: %d\n", generation, bestIndividual.Fitness)

		// Seleccionar padres y realizar cruzamiento
		newPopulation := make([]Individual, populationSize)
		for i := range population {
			parent1 := selectParent(population)
			parent2 := selectParent(population)
			child := crossover(parent1, parent2)
			newPopulation[i] = child
		}

		// Aplicar mutaciones
		for i := range newPopulation {
			if rand.Float64() < mutationRate {
				mutate(newPopulation[i])
			}
		}

		// Reemplazar la antigua población con la nueva generación
		population = newPopulation
	}

	// Encontrar el mejor individuo después de todas las generaciones
	for i := range population {
		population[i].Fitness = DistanciaEuclidiana(xdistances, ydistances, population[i].Chromosome)
	}
	bestIndividual := population[0]
	for _, individual := range population {
		//fmt.Println(individual.Fitness)
		if individual.Fitness < bestIndividual.Fitness {
			if individual.Fitness > 0{
				bestIndividual = individual
			}
			
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %d\n", bestIndividual.Fitness)
	//fmt.Println("Orden de las Ciudades:", bestIndividual.Chromosome)
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

	GeneticEuclideano(firstnode, secondnode)
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
				numbernode[count]=aux[0]
				firstnode[count],_=strconv.ParseFloat(aux[1], 32)
				secondnode[count],_=strconv.ParseFloat(aux[2], 32)
			}
			count++
		}
	}	
}

func DistanciaGeografica(x []float64, y []float64, chromosome []int) int {
	distance:=0
	var q1, q2, q3 float64
	RRR := 6378.388 // Radio de la Tierra en km
	pi :=  3.141592

	latitude := make([]float64, len(x))
	longitude := make([]float64, len(x))
	for i:=0; i<len(chromosome)-1; i++{
		deg := int(x[chromosome[i]] + 0.5)
		min := x[chromosome[i]] - float64(deg)
		latitude[i] = pi * (float64(deg) + 5.0*min/3.0) / 180.0
		deg = int(y[chromosome[i]] + 0.5)
		min = y[chromosome[i]] - float64(deg)
		longitude[i] = pi * (float64(deg) + 5.0*min/3.0) / 180.0
	}

	for i:=0; i<len(x)-1; i++{
		q1 = math.Cos(longitude[i] - longitude[i+1])
		q2 = math.Cos(latitude[i] - latitude[i+1])
		q3 = math.Cos(latitude[i] + latitude[i+1])
		distance += int(RRR * math.Acos(0.5*((1.0+q1)*q2-(1.0-q1)*q3)) + 1.0)
	}

	return distance
}

func GeneticGeografico(xdistances []float64, ydistances []float64) {
	rand.Seed(time.Now().UnixNano())

	// Configuración del algoritmo genético
	populationSize := 1000
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(len(xdistances))
	}
	
	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		// Calcular fitness de la población
		for i := range population {
			population[i].Fitness = DistanciaGeografica(xdistances, ydistances, population[i].Chromosome)
		}

		/*for i:=range population{
			for j:=0; j<populationSize-1; i++{
				fmt.Printf("%d ", population[i].Chromosome[j])
			}
			fmt.Printf("Fitness: %d\n",population[i].Fitness)
		}*/

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}
		
		// Imprimir la mejor distancia en esta generación
		//fmt.Printf("Generación %d - Mejor Distancia: %d\n", generation, bestIndividual.Fitness)

		// Seleccionar padres y realizar cruzamiento
		newPopulation := make([]Individual, populationSize)
		for i := range population {
			parent1 := selectParent(population)
			parent2 := selectParent(population)
			child := crossover(parent1, parent2)
			newPopulation[i] = child
		}

		// Aplicar mutaciones
		for i := range newPopulation {
			if rand.Float64() < mutationRate {
				mutate(newPopulation[i])
			}
		}

		// Reemplazar la antigua población con la nueva generación
		population = newPopulation
	}

	// Encontrar el mejor individuo después de todas las generaciones
	for i := range population {
		population[i].Fitness = DistanciaGeografica(xdistances, ydistances, population[i].Chromosome)
	}
	bestIndividual := population[0]
	for _, individual := range population {
		//fmt.Println(individual.Fitness)
		if individual.Fitness < bestIndividual.Fitness {
			if individual.Fitness > 0{
				bestIndividual = individual
			}
			
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %d\n", bestIndividual.Fitness)
	//fmt.Println("Orden de las Ciudades:", bestIndividual.Chromosome)
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

	GeneticGeografico(firstnode, secondnode)
	//fmt.Println("La distancia total para este archivo es de", DistanciaGeografica(numbernode, firstnode, secondnode))
}

func ReadFileCeil(file *os.File, numbernode []string, firstnode, secondnode []float64) {
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
				numbernode[count]=aux[0]
				firstnode[count],_=strconv.ParseFloat(aux[1], 32)
				secondnode[count],_=strconv.ParseFloat(aux[2], 32)
			}
			count++
		}
	}	
}

func DistanciaCeil(x []float64, y []float64, chromosome []int) int {
	var x1, x2, y1, y2 float64
	distance := 0
	for i:=0; i < len(chromosome)-1; i++{
		x1 = x[chromosome[i]]
		x2 = x[chromosome[i+1]]
		y1 = y[chromosome[i]]
		y2 = y[chromosome[i+1]]
		distance += int(math.Sqrt(math.Pow(math.Abs(y2-y1), 2) + math.Pow(math.Abs(x2-x1), 2)) + 0.5)
	}
	/*for i := 0; i < len(nodes); i++ {
		if i < len(nodes)-1 {
			x1 = x[i]
			x2 = x[i+1]
			y1 = y[i]
			y2 = y[i+1]
			distance += int(math.Sqrt(math.Pow(math.Abs(y2-y1), 2) + math.Pow(math.Abs(x2-x1), 2)) + 0.5)
		}
	}*/

	return distance
}

func GeneticCeil(xdistances []float64, ydistances []float64) {
	rand.Seed(time.Now().UnixNano())

	// Configuración del algoritmo genético
	populationSize := 1000
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(len(xdistances))
	}
	
	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		// Calcular fitness de la población
		for i := range population {
			population[i].Fitness = DistanciaCeil(xdistances, ydistances, population[i].Chromosome)
		}

		/*for i:=range population{
			for j:=0; j<populationSize-1; i++{
				fmt.Printf("%d ", population[i].Chromosome[j])
			}
			fmt.Printf("Fitness: %d\n",population[i].Fitness)
		}*/

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}
		
		// Imprimir la mejor distancia en esta generación
		//fmt.Printf("Generación %d - Mejor Distancia: %d\n", generation, bestIndividual.Fitness)

		// Seleccionar padres y realizar cruzamiento
		newPopulation := make([]Individual, populationSize)
		for i := range population {
			parent1 := selectParent(population)
			parent2 := selectParent(population)
			child := crossover(parent1, parent2)
			newPopulation[i] = child
		}

		// Aplicar mutaciones
		for i := range newPopulation {
			if rand.Float64() < mutationRate {
				mutate(newPopulation[i])
			}
		}

		// Reemplazar la antigua población con la nueva generación
		population = newPopulation
	}

	// Encontrar el mejor individuo después de todas las generaciones
	for i := range population {
		population[i].Fitness = DistanciaCeil(xdistances, ydistances, population[i].Chromosome)
	}
	bestIndividual := population[0]
	for _, individual := range population {
		//fmt.Println(individual.Fitness)
		if individual.Fitness < bestIndividual.Fitness {
			if individual.Fitness > 0{
				bestIndividual = individual
			}
			
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %d\n", bestIndividual.Fitness)
}

func ProblemaCeil(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()
	numbernode := make([]string, n)
	firstnode := make([]float64, n)
	secondnode := make([]float64, n)

	ReadFileCeil(file, numbernode, firstnode, secondnode)
	GeneticCeil(firstnode, secondnode)
	//fmt.Println("La distancia total para este archivo es de", DistanciaCeil(numbernode, firstnode, secondnode))
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

func DistanciaATT(x []int, y []int, chromosome []int) int {
	xd := make([]float64, len(x))
	yd := make([]float64, len(x))
	dij := 0
	tij := 0
	var rij float64

	for i:=0; i < len(chromosome)-1; i++{
		xd[i] = float64(x[chromosome[i]] - x[chromosome[i+1]])
		yd[i] = float64(y[chromosome[i]] - y[chromosome[i+1]])
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

func GeneticATT(xdistances []int, ydistances []int) {
	rand.Seed(time.Now().UnixNano())

	// Configuración del algoritmo genético
	populationSize := 1000
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(len(xdistances))
	}
	
	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		// Calcular fitness de la población
		for i := range population {
			population[i].Fitness = DistanciaATT(xdistances, ydistances, population[i].Chromosome)
		}

		/*for i:=range population{
			for j:=0; j<populationSize-1; i++{
				fmt.Printf("%d ", population[i].Chromosome[j])
			}
			fmt.Printf("Fitness: %d\n",population[i].Fitness)
		}*/

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}
		
		// Imprimir la mejor distancia en esta generación
		//fmt.Printf("Generación %d - Mejor Distancia: %d\n", generation, bestIndividual.Fitness)

		// Seleccionar padres y realizar cruzamiento
		newPopulation := make([]Individual, populationSize)
		for i := range population {
			parent1 := selectParent(population)
			parent2 := selectParent(population)
			child := crossover(parent1, parent2)
			newPopulation[i] = child
		}

		// Aplicar mutaciones
		for i := range newPopulation {
			if rand.Float64() < mutationRate {
				mutate(newPopulation[i])
			}
		}

		// Reemplazar la antigua población con la nueva generación
		population = newPopulation
	}

	// Encontrar el mejor individuo después de todas las generaciones
	for i := range population {
		population[i].Fitness = DistanciaATT(xdistances, ydistances, population[i].Chromosome)
	}
	bestIndividual := population[0]
	for _, individual := range population {
		//fmt.Println(individual.Fitness)
		if individual.Fitness < bestIndividual.Fitness {
			if individual.Fitness > 0{
				bestIndividual = individual
			}
			
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %d\n", bestIndividual.Fitness)
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

	GeneticATT(firstnode, secondnode)
	//fmt.Println("La distancia total para este archivo es de", DistanciaATT(numbernode, firstnode, secondnode))
}

func ReadFileSuperior(file *os.File, matriz [][]int, n int) {
	scanner := bufio.NewScanner(file)
	NoDisplay:=false
	fila := 0
	leyendoMatriz := false
	columna:=0
	// Encuentra la línea "EDGE_WEIGHT_SECTION"
	for scanner.Scan() {
		lineaActual := scanner.Text()
		if strings.HasPrefix(lineaActual, "DISPLAY_DATA_TYPE: NO_DISPLAY"){
			NoDisplay=true
		}
		if lineaActual == "EDGE_WEIGHT_SECTION" {
			leyendoMatriz = true
			break
		}
	}

	if NoDisplay{
		for scanner.Scan() {
			lineaActual := scanner.Text()
			if leyendoMatriz && lineaActual != "EOF" {
				valoresStr := strings.Fields(lineaActual)
					//fmt.Println(valoresStr)
					for _, valorStr := range valoresStr {
						valor, err := strconv.Atoi(valorStr)
						if err != nil {
							fmt.Println("Error al convertir el valor:", err)
							return
						}
			
						if valor == 0 {
							if columna < n{
								for columna < n{
									if columna<n{
										matriz[fila][columna]=0
										columna++
									}else{
										break
									}
								}
							}
						} else {
							if columna < n{
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
					if fila >= n {
						break
					}
			}
		}
	}else{
		for scanner.Scan() {
			lineaActual := scanner.Text()
			if leyendoMatriz && lineaActual != "EOF" {
				linea := strings.Fields(lineaActual)
				for columna, valorStr := range linea {
					valor, err := strconv.Atoi(valorStr)
					if err != nil {
						fmt.Println("Error al convertir el valor:", err)
						return
					}
					matriz[fila][columna] = valor
				}
			}
			fila++
		}
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

	ReadFileSuperior(file, matriz, n)	

	EscribirMatrizEnCSV("./Files/Excel/"+fileName[6:11]+".csv",matriz)	

	distanciaMATRIX(matriz, n)
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
}

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

	//Se completa la matriz
	CompletarMatriz(matriz, n)

	EscribirMatrizEnCSV("./Files/Excel/"+fileName[6:10]+".csv",matriz)	

	distanciaMATRIX(matriz, n)
}

func SolveATSP(distances [][]int, numCities int, tour []int, chromosome []int) int {
	visited := make([]bool, numCities)
	startCity := 0 
	fmt.Println("Entra a la funcion de la distancia")
	tour = append(tour, startCity)
	visited[startCity] = true
	totalDistance := 0

	// Construir el recorrido visitando el vecino más cercano disponible
	for i := 0; i < numCities; i++ {
		currentCity := tour[i]
		nextCity := findNearestNeighbor(currentCity, visited, distances)
		tour = append(tour, nextCity)
		
		fmt.Println("La ciudad es ", nextCity)
		visited[nextCity] = true
		totalDistance += distances[currentCity][nextCity]
		fmt.Println("Iteracion numero ", i)
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

func ReadExplicit(file *os.File, data [][]int, n int) {
	scanner := bufio.NewScanner(file)
	for i := 0; i < 7; i++ {
		if !scanner.Scan() {
			fmt.Println("Error: no se encontraron suficientes líneas en el archivo")
			return
		}
	}

	// Variables para rastrear la fila y columna
	row := 0
	col := 0

	// Bucle principal para leer el archivo
	for scanner.Scan() {
		lineaActual := scanner.Text()
		if lineaActual != "EOF" && row < n {
			linea := strings.Fields(lineaActual)
			for _, valorStr := range linea {
				//fmt.Println(valorStr)
				valor, err := strconv.Atoi(valorStr)
				if err != nil {
					fmt.Println("Error al convertir el valor:", err)
					return
				}
				if valor == 9999999{
					row++
					col=0
				}else if valor == 9999 {
					row++
					col = 0
				} else {
					if col < n{
						data[row][col] = valor
						col++
					}
				}
				if row == n {
					break
				}
			}
			if row == n {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
}

func GeneticATSP(distances [][]int, numCities int, tour []int) {
	rand.Seed(time.Now().UnixNano())

	// Configuración del algoritmo genético
	populationSize := 5
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(numCities)
	}
	fmt.Println(population)
	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		fmt.Println("La generacion es: %d", generation)
		// Calcular fitness de la población
		
		for i := range population {
			population[i].Fitness = SolveATSP(distances, numCities, tour,  population[i].Chromosome)
		}

		/*for i:=range population{
			for j:=0; j<populationSize-1; i++{
				fmt.Printf("%d ", population[i].Chromosome[j])
			}
			fmt.Printf("Fitness: %d\n",population[i].Fitness)
		}*/

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}
		
		// Imprimir la mejor distancia en esta generación
		//fmt.Printf("Generación %d - Mejor Distancia: %d\n", generation, bestIndividual.Fitness)

		// Seleccionar padres y realizar cruzamiento
		newPopulation := make([]Individual, populationSize)
		for i := range population {
			parent1 := selectParent(population)
			parent2 := selectParent(population)
			child := crossover(parent1, parent2)
			newPopulation[i] = child
		}

		// Aplicar mutaciones
		for i := range newPopulation {
			if rand.Float64() < mutationRate {
				mutate(newPopulation[i])
			}
		}

		// Reemplazar la antigua población con la nueva generación
		population = newPopulation
	}

	// Encontrar el mejor individuo después de todas las generaciones
	for i := range population {
		population[i].Fitness = SolveATSP(distances, numCities, tour,  population[i].Chromosome)
	}
	bestIndividual := population[0]
	for _, individual := range population {
		//fmt.Println(individual.Fitness)
		if individual.Fitness < bestIndividual.Fitness {
			if individual.Fitness > 0{
				bestIndividual = individual
			}
			
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %d\n", bestIndividual.Fitness)
}

func ProblemaExplicit(fileName string, n int) {
	file, err := os.Open(fileName)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()
    data := make([][]int, n)
    for i := range data {
        data[i] = make([]int, n)
    }

    ReadExplicit(file, data, n)

    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	var tour []int
	GeneticATSP(data, n, tour)
    //
    //totalDistance := SolveATSP(data, n, tour)

    //fmt.Println("0")

    //fmt.Println("Distancia total recorrida:", totalDistance)
}