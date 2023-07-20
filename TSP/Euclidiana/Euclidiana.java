package TSP.Euclidiana;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

public class Euclidiana {

    public static void ReadFileEuclidiana(File file, String[] numbernode, int[] firstnode, int[] secondnode)
            throws IOException {
        int count = 0; // Para poder obtener las substring
        int flag = 0; // Se tendra una bandera para comenzar a guardar los datos
        int i = 1;
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                if (flag == 0) {
                    if (line.compareTo("NODE_COORD_SECTION") == 0) // Se comenzará a colocar los valores de los nodos
                        flag++;
                } else {
                    if (line.compareTo("EOF") == 0) {
                        break;
                    } else {
                        String aux = line.trim();
                        String[] datos = aux.split(" ", 3);
                        numbernode[count] = datos[0];
                        if (datos[1].contains("e")) { // True
                            String[] numbers = datos[1].split("e");
                            firstnode[count] = (int) (Double.parseDouble(numbers[0])
                                    * Math.pow(10, Double.parseDouble(numbers[1])));
                            numbers = datos[2].split("e");
                            secondnode[count] = (int) (Double.parseDouble(numbers[0])
                                    * Math.pow(10, Double.parseDouble(numbers[1])));
                        } else { // False
                            if (datos[1].length() == 0) {
                                String[] num = datos[2].trim().split(" ", 2);
                                firstnode[count] = Integer.parseInt(num[0]);
                                secondnode[count] = num[1].contains(" ") ? Integer.parseInt(num[1].trim())
                                        : Integer.parseInt(num[1]);
                            } else {
                                if (datos[1].indexOf(".") != -1) {
                                    float num1 = Float.parseFloat(datos[1].trim());
                                    float num2 = Float.parseFloat(datos[2].trim());
                                    firstnode[count] = Math.round(num1);
                                    secondnode[count] = Math.round(num2);
                                } else {
                                    firstnode[count] = Integer.parseInt(datos[1].trim());
                                    secondnode[count] = Integer.parseInt(datos[2].trim());
                                }
                            }
                        }
                        count++;
                    }
                }
            }
        }
    }

    public static int distanciaEuclidiana(String[] nodes, int[] x, int[] y) {
        int x1 = 0;
        int x2 = 0;
        int y1 = 0;
        int y2 = 0;
        int distance = 0;
        for (int i = 0; i < nodes.length; i++) {
            if (i < nodes.length - 1) {
                x1 = x[i];
                x2 = x[i + 1];
                y1 = y[i];
                y2 = y[i + 1];
                distance += (int) (Math.sqrt(Math.pow(Math.abs(y2 - y1), 2) + Math.pow(Math.abs(x2 - x1), 2))
                        + 0.5);
            }
        }
        return distance;
    }

    public void problemaEuclidiano(File file, int n) {
        String[] numbernode = new String[n];
        int[] firstnode = new int[n];
        int[] secondnode = new int[n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            ReadFileEuclidiana(file, numbernode, firstnode, secondnode);
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println("La distancia total para este archivo es de "
                + distanciaEuclidiana(numbernode, firstnode, secondnode));
    }
}
