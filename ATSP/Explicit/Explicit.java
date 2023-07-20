package ATSP.Explicit;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

import ATSP.Distance.Distance;

public class Explicit {

    Distance distance = new Distance();

    public Explicit() {
    }

    public void readExplicit(File file, int[][] data, int n) throws IOException {
        BufferedReader reader = new BufferedReader(new FileReader(file));
        String line;

        // Ignorar las primeras 7 líneas del archivo
        for (int i = 0; i < 7; i++) {
            reader.readLine();
        }

        int row = 0;
        int col = 0;
        while ((line = reader.readLine()) != null && !line.equals("EOF")) {
            String[] values = line.trim().split("\\s+");
            for (String value : values) {
                int num = Integer.parseInt(value);
                if (num == 9999) {
                    // Saltar al siguiente elemento de la siguiente fila
                    row++;
                    col = 0;
                } else {
                    data[row][col] = num;
                    col++;
                }

                // Comprobar si hemos llenado completamente la matriz
                if (row == n) {
                    break;
                }
            }

            // Comprobar si hemos llenado completamente la matriz
            if (row == n) {
                break;
            }
        }

        reader.close();
    }

    public void problemaExplicit(File file, int n) {
        int[][] data = new int[n][n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            readExplicit(file, data, n);
        } catch (Exception e) {
            e.printStackTrace();
        }
        List<Integer> tour = new ArrayList<>();
        int totalDistance=distance.solveATSP(data, n, tour);
        System.out.println("Recorrido óptimo:");
        for (int city : tour) {
            System.out.print(city + " -> ");
        }
        System.out.println("0"); // Imprimir el punto de inicio nuevamente para completar el ciclo

        System.out.println("Distancia total recorrida: " + totalDistance);
    }
}

