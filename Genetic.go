package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Individual struct {
	Chromosome []int
	Fitness    float64
}

func generateRandomIndividual(numCities int) Individual {
	chromosome := rand.Perm(numCities)
	return Individual{Chromosome: chromosome}
}

func calculateTotalDistance(individual Individual, distances [][]float64) float64 {
	totalDistance := 0.0
	for i := 0; i < len(individual.Chromosome)-1; i++ {
		fromCity := individual.Chromosome[i]
		toCity := individual.Chromosome[i+1]
		totalDistance += distances[fromCity][toCity]
	}
	// Agregar la distancia de vuelta a la primera ciudad
	totalDistance += distances[individual.Chromosome[len(individual.Chromosome)-1]][individual.Chromosome[0]]
	return totalDistance
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Definir distancias entre ciudades (en este ejemplo, distancias aleatorias)
	numCities := 100
	distances := make([][]float64, numCities)
	for i := range distances {
		distances[i] = make([]float64, numCities)
		for j := range distances[i] {
			if i != j {
				distances[i][j] = float64(rand.Intn(100)) // Distancias aleatorias entre 0 y 99
			}
		}
	}

	// Configuración del algoritmo genético
	populationSize := 1000
	numGenerations := 100
	mutationRate := 0.01

	// Inicializar población aleatoria
	population := make([]Individual, populationSize)
	for i := range population {
		population[i] = generateRandomIndividual(numCities)
	}

	// Evolución de la población
	for generation := 0; generation < numGenerations; generation++ {
		// Calcular fitness de la población
		for i := range population {
			population[i].Fitness = calculateTotalDistance(population[i], distances)
		}

		// Encontrar el mejor individuo (el de menor distancia)
		bestIndividual := population[0]
		for _, individual := range population {
			if individual.Fitness < bestIndividual.Fitness {
				bestIndividual = individual
			}
		}

		// Imprimir la mejor distancia en esta generación
		fmt.Printf("Generación %d - Mejor Distancia: %.2f\n", generation, bestIndividual.Fitness)

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
	bestIndividual := population[0]
	for _, individual := range population {
		if individual.Fitness < bestIndividual.Fitness {
			bestIndividual = individual
		}
	}

	// Imprimir la mejor distancia encontrada y el orden de las ciudades
	fmt.Printf("Mejor Distancia Final: %.2f\n", bestIndividual.Fitness)
	fmt.Println("Orden de las Ciudades:", bestIndividual.Chromosome)
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
