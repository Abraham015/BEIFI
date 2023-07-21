package TSP.Matrix;

import Tools.Tools;
import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

public class Matrix {

    Tools tool=new Tools();

    public static void ReadFileInferior(File file, int[][] matriz, int n) {
        String flag = "";
        int i = 0;
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            // Saltar las primeras 6 líneas
            for (i = 0; i < 7; i++) {
                flag = br.readLine();
                if (flag.compareTo("DISPLAY_DATA_TYPE : TWOD_DISPLAY") == 0)
                    br.readLine();
            }

            // Leer los datos y almacenarlos en la matriz
            int fila = 0;
            int columna = 0;
            String linea;
            boolean continuarEnLaMismaFila = false;
            while ((linea = br.readLine()) != null && !linea.equals("EOF")) {
                String[] valores = linea.trim().split("\\s+");
                for (String valorStr : valores) {
                    int valor = Integer.parseInt(valorStr);
                    if (valor == 0) {
                        if (columna < n) {
                            // Rellenar el resto de la fila con ceros si se encuentra un cero antes de
                            // alcanzar la dimensión total
                            while (columna < n) {
                                matriz[fila][columna++] = 0;
                            }
                        }
                        continuarEnLaMismaFila = false; // Reiniciar la bandera
                    } else {
                        matriz[fila][columna++] = valor; // Almacenar el valor en la matriz
                        continuarEnLaMismaFila = true; // Establecer la bandera para continuar en la misma fila
                    }

                    // Pasar a la siguiente fila si no hay más espacio en la actual
                    if (columna == n) {
                        fila++;
                        columna = 0;
                    }

                    // Comprobar si se ha alcanzado la última posición de la matriz
                    if (fila == n) {
                        break;
                    }
                }

                // Comprobar si se ha alcanzado la última posición de la matriz
                if (fila == n) {
                    break;
                }
            }

            // Rellenar con ceros si la matriz no se ha llenado completamente
            while (fila < n) {
                for (columna = 0; columna < n; columna++) {
                    matriz[fila][columna] = 0;
                }
                fila++; // Pasar a la siguiente fila
            }

        } catch (IOException e) {
            System.out.println("Error al leer el archivo: " + e.getMessage());
            e.printStackTrace();
        }
    }

    public static void ReadFileSuperior(File file, int[][] matriz, int numbernode) {
        String[] linea;
        int fila = 0;
        int i = 0;
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String lineaActual;
            while ((lineaActual = br.readLine()) != null && !lineaActual.equals("EOF")) {
                if (fila >= 8) {
                    if (i < matriz.length) {
                        linea = lineaActual.trim().split("\\s+");
                        for (int columna = 0; columna < linea.length; columna++) {
                            matriz[fila - 8][columna] = Integer.parseInt(linea[columna]);
                        }
                    }

                }
                fila++;
                i++;
            }
        } catch (IOException e) {
            e.printStackTrace();
            // System.err.format("Error de E/S: %s%n", e);
        }
    }

    public static void distanciaMATRIX(int[][] matrix, int n) {
        int[] visitados = new int[n];
        int[] ruta = new int[n];
        int[] mejor_ruta = new int[n];
        int mejor_distancia = Integer.MAX_VALUE;
        int distancia_actual = 0;
        int element = -1;

        for (int i = 0; i < n; i++) {
            visitados[i] = 0;
            ruta[i] = -1;
        }

        visitados[0] = 1;
        ruta[0] = 0;

        int pos_actual = 0;

        while (true) {
            element = -1; // actualizar el valor de element en cada iteración

            for (int i = 0; i < n; i++) {
                if (visitados[i] == 0) {
                    if (element == -1) {
                        element = i;
                    }
                    distancia_actual = matrix[pos_actual][i];
                    if (distancia_actual < matrix[pos_actual][element]) {
                        element = i;
                    }
                }
            }

            if (element == -1) {
                break;
            }

            visitados[element] = 1;
            ruta[element] = pos_actual;
            pos_actual = element;
        }

        distancia_actual = matrix[pos_actual][0];
        ruta[n - 1] = 0;
        ruta[0] = pos_actual;

        for (int i = 0; i < n; i++) {
            if (i != n - 1) {
                distancia_actual += matrix[ruta[i]][ruta[i + 1]];
            }
        }

        if (distancia_actual < mejor_distancia) {
            mejor_distancia = distancia_actual;
            for (int i = 0; i < n; i++) {
                mejor_ruta[i] = ruta[i];
            }
        }

        System.out.println("La distancia total es: " + mejor_distancia);
    }

    public void problemaSuperior(File file, int n) {
        int[][] data = new int[n][n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            ReadFileSuperior(file, data, n);
            // imprimirMatriz(data, n);
            tool.escribirMatrizEnCSV("./Files/Excel/Matriz.csv", data);
        } catch (Exception e) {
            e.printStackTrace();
        }
        distanciaMATRIX(data, n);
    }

    public void problemaInferior(File file, int n) {
        int[][] data = new int[n][n];
        try {
            ReadFileInferior(file, data, n);
            tool.escribirMatrizEnCSV("./Files/Excel/MatrizDiagonal.csv", data);
            tool.completarMatriz(data, n);
            tool.escribirMatrizEnCSV("./Files/Excel/MatrizCompleta.csv", data);
        } catch (Exception e) {
            e.printStackTrace();
        }
        distanciaMATRIX(data, n);
    }
}
