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
                    //ProblemaExplicit(filepath.Join(folder, fileInfo.Name()), n)
                }
            } else {
                file.Seek(0, 0) // Asegura que el archivo esté al principio antes de abrirlo nuevamente
                typeDistance = TypeDistance(file)
				fmt.Println("El tipo de problema es: "+typeDistance)
                switch typeDistance {
                case "Euclidiana":
                    ProblemaEuclidiano(filepath.Join(folder, fileInfo.Name()), n)
                case "Geografica":
                    //ProblemaGeografico(filepath.Join(folder, fileInfo.Name()), n)
                case "Circular":
                    //ProblemaCeil(filepath.Join(folder, fileInfo.Name()), n)
                case "ATT":
                    //ProblemaATT(filepath.Join(folder, fileInfo.Name()), n)
                case "DiagonalSuperior":
                    //ProblemaSuperior(filepath.Join(folder, fileInfo.Name()), n)
                case "DiagonalInferior":
                    //ProblemaInferior(filepath.Join(folder, fileInfo.Name()), n)
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