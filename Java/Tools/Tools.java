package Tools;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;

public class Tools {

    public Tools(){}

    public String typeDistance(File file) {
        String type = "";
        String distance = "";
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                switch (line) {
                    case "EDGE_WEIGHT_TYPE : EUC_2D":
                        distance = "Euclidiana";
                        break;
                    case "EDGE_WEIGHT_TYPE: GEO":
                        distance = "Geografica";
                        break;
                    case "EDGE_WEIGHT_TYPE : ATT":
                        distance = "ATT";
                        break;
                    case "EDGE_WEIGHT_TYPE : MATRIX":
                        distance = "MATRIX";
                        break;
                    case "EDGE_WEIGHT_TYPE : CEIL_2D":
                        distance = "Circular";
                        break;
                    case "EDGE_WEIGHT_TYPE: EXPLICIT":
                        type = br.readLine();
                        distance = (type.compareTo("EDGE_WEIGHT_FORMAT : LOWER_DIAG_ROW") == 0
                                || type.compareTo("EDGE_WEIGHT_FORMAT: LOWER_DIAG_ROW") == 0) ? "DiagonalSuperior"
                                        : "DiagonalInferior";
                        break;
                    case "EDGE_WEIGHT_TYPE: EUC_2D":
                        distance = "Euclidiana";
                        break;
                    case "EDGE_WEIGHT_TYPE : EXPLICIT":
                        type = br.readLine();
                        distance = (type.compareTo("EDGE_WEIGHT_FORMAT : LOWER_DIAG_ROW") == 0
                                || type.compareTo("EDGE_WEIGHT_FORMAT: LOWER_DIAG_ROW") == 0) ? "DiagonalSuperior"
                                        : "DiagonalInferior";
                        break;
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return distance;
    }

    public String typeDistanceATSP(File file) {
        String distance = "";
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                switch (line) {
                    case "EDGE_WEIGHT_TYPE: EXPLICIT":
                        distance = "Explicita";
                        break;
                }
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return distance;
    }

    public void completarMatriz(int[][] m, int size) {
        for (int i = 1; i < size; i++) {
            for (int j = 0; j < size; j++) {
                if (m[i][j] != 0)
                    m[j][i] = m[i][j];
            }
        }
    }

    public int calcularCostoCamino(int[][] matriz, int[] camino) {
        int costoTotal = 0;
        int n = camino.length;

        for (int i = 0; i < n - 1; i++) {
            int nodoActual = camino[i];
            int nodoSiguiente = camino[i + 1];
            costoTotal += matriz[nodoSiguiente][nodoActual];
        }

        return costoTotal;
    }

    public void imprimirMatriz(int[][] matriz, int n) {
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                System.out.print(matriz[i][j] + " ");
            }
            System.out.println();
        }
    }

    public void escribirMatrizEnCSV(String nombreArchivo, int[][] matriz) {
        try (BufferedWriter writer = new BufferedWriter(new FileWriter(nombreArchivo))) {
            for (int fila = 0; fila < matriz.length; fila++) {
                for (int columna = 0; columna < matriz[fila].length; columna++) {
                    writer.write(String.valueOf(matriz[fila][columna]));

                    if (columna < matriz[fila].length - 1) {
                        writer.write(" ");
                    }
                }

                writer.newLine();
            }
        } catch (IOException e) {
            System.out.println("Error al escribir en el archivo CSV: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
