package ATSP.Distance;

import java.util.List;

public class Distance {
    public Distance(){}

    // Función para resolver el ATSP usando la heurística del vecino más cercano
    public int solveATSP(int[][] distances, int numCities, List<Integer> tour) {
        boolean[] visited = new boolean[numCities];
        int startCity = 0; // Puedes elegir cualquier ciudad como punto de inicio

        tour.add(startCity);
        visited[startCity] = true;
        int totalDistance = 0;

        // Construir el recorrido visitando el vecino más cercano disponible
        for (int i = 0; i < numCities - 1; i++) {
            int currentCity = tour.get(i);
            int nextCity = findNearestNeighbor(currentCity, visited, distances);
            tour.add(nextCity);
            visited[nextCity] = true;
            totalDistance += distances[currentCity][nextCity];
        }

        // Regresar al punto de inicio para completar el ciclo
        tour.add(startCity);
        totalDistance += distances[tour.get(numCities - 1)][startCity];

        return totalDistance;
    }

    // Función para encontrar el vecino más cercano disponible desde la ciudad actual
    private static int findNearestNeighbor(int city, boolean[] visited, int[][] distances) {
        int nearestNeighbor = -1;
        int minDistance = Integer.MAX_VALUE;

        for (int i = 0; i < visited.length; i++) {
            if (!visited[i] && distances[city][i] < minDistance) {
                minDistance = distances[city][i];
                nearestNeighbor = i;
            }
        }

        return nearestNeighbor;
    }
}
