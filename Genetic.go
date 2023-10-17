package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type Individual struct {
	Chromosome []int
	Fitness    int
}

func generateRandomIndividual(numCities int) Individual {
	chromosome := rand.Perm(numCities)
	return Individual{Chromosome: chromosome, Fitness: 0}
}

func calculateDistance(city1, city2 [2]int) float64 {
	return float64((city1[0]-city2[0])*(city1[0]-city2[0]) + (city1[1]-city2[1])*(city1[1]-city2[1]))
}

func calculateFitness(individual Individual, cities [][2]int) int {
	totalDistance := 0.0
	for i := 0; i < len(individual.Chromosome)-1; i++ {
		city1 := cities[individual.Chromosome[i]]
		city2 := cities[individual.Chromosome[i+1]]
		totalDistance += calculateDistance(city1, city2)
	}
	// Agregar la distancia de vuelta a la primera ciudad
	firstCity := cities[individual.Chromosome[0]]
	lastCity := cities[individual.Chromosome[len(individual.Chromosome)-1]]
	totalDistance += calculateDistance(lastCity, firstCity)
	return int(totalDistance)
}

func crossover(parent1, parent2 Individual) Individual {
	crossoverPoint := rand.Intn(len(parent1.Chromosome))
	childChromosome := append(parent1.Chromosome[:crossoverPoint], parent2.Chromosome[crossoverPoint:]...)
	child := Individual{Chromosome: childChromosome, Fitness: 0}
	return child
}

func mutate(individual Individual) Individual {
	index1 := rand.Intn(len(individual.Chromosome))
	index2 := rand.Intn(len(individual.Chromosome))
	individual.Chromosome[index1], individual.Chromosome[index2] = individual.Chromosome[index2], individual.Chromosome[index1]
	return individual
}

func geneticAlgorithm(cities [][2]int, populationSize, generations int) Individual {
	rand.Seed(time.Now().UnixNano())

	// Inicializar la población
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(len(cities))
	}

	// Evolucionar la población
	for generation := 0; generation < generations; generation++ {
		// Calcular la aptitud de cada individuo en la población
		for i := range population {
			population[i].Fitness = calculateFitness(population[i], cities)
		}

		// Ordenar la población por aptitud (menor distancia es mejor)
		sort.Slice(population, func(i, j int) bool {
			return population[i].Fitness < population[j].Fitness
		})

		// Seleccionar padres para la reproducción (torneo)
		parent1 := population[rand.Intn(populationSize/2)]
		parent2 := population[rand.Intn(populationSize/2)]

		// Realizar cruce y mutación para crear nuevos individuos
		child := crossover(parent1, parent2)
		child = mutate(child)

		// Reemplazar el peor individuo con el nuevo hijo
		population[populationSize-1] = child
	}

	// Ordenar la población una última vez antes de devolver el mejor individuo
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})

	return population[0]
}

func main() {
	// Definir las ciudades como coordenadas (x, y)
	cities := [][2]int{{0, 0}, {1, 2}, {3, 1}, {5, 2}, {6, 0}}

	populationSize := 100
	generations := 1000

	bestSolution := geneticAlgorithm(cities, populationSize, generations)
	fmt.Printf("Mejor ruta encontrada: %v\n", bestSolution.Chromosome)
	fmt.Printf("Distancia total: %d\n", bestSolution.Fitness)
}
