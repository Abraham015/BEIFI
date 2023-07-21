package TSP.Geografica;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

public class Geografico {
    public static void ReadFileGeografica(File file, String[] numbernode, float[] firstnode, float[] secondnode)
            throws IOException {
        int count = 0; // Para poder obtener las substring
        int flag = 0; // Se tendra una bandera para comenzar a guardar los datos
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                if (flag == 0) {
                    if (line.compareTo("NODE_COORD_SECTION") == 0) // Se comenzará a colocar los valores de los nodos
                        flag++;
                    if (line.compareTo("DISPLAY_DATA_SECTION") == 0)
                        flag++;
                } else {
                    if (line.compareTo("EOF") == 0) {
                        continue;
                    } else {
                        if (count < firstnode.length) {
                            String aux = line.trim();
                            String[] datos = aux.split(" ", 3);
                            if (datos[1].length() == 0) {
                                String[] num = datos[2].trim().split(" ", 2);
                                firstnode[count] = num[0].contains(" ") ? Float.parseFloat(num[0].trim())
                                        : Float.parseFloat(num[0]);
                                secondnode[count] = num[1].contains(" ") ? Float.parseFloat(num[1].trim())
                                        : Float.parseFloat(num[1]);
                            } else {
                                firstnode[count] = Float.parseFloat(datos[1].trim());
                                secondnode[count] = Float.parseFloat(datos[2].trim());
                            }
                        }
                        count++;
                    }
                }
            }
        }
    }

    public static int distanciaGeometrica(String[] numbernode, float[] x, float[] y) {
        double pi = 3.141592;
        double RRR = 6378.388;
        int distance = 0;
        int deg;
        double latitude[] = new double[numbernode.length];
        double longitude[] = new double[numbernode.length];
        double q1, q2, q3;
        float min;

        // En este primer for se realizarán los primeros calculos
        for (int i = 0; i < numbernode.length; i++) {
            deg = (int) (x[i] + 0.5);
            min = x[i] - deg;
            latitude[i] = pi * (deg + 5.0 * min / 3.0) / 180.00;
            deg = (int) (y[i] + 0.5);
            min = y[i] - deg;
            longitude[i] = pi * (deg + 5.0 * min / 3.0) / 180.00;
        }

        // Este for será para calcular la distancia
        for (int i = 0; i < numbernode.length; i++) {
            if (i < latitude.length - 1) {
                q1 = Math.cos(longitude[i] - longitude[i + 1]);
                q2 = Math.cos(latitude[i] - latitude[i + 1]);
                q3 = Math.cos(latitude[i] + latitude[i + 1]);
                distance += (int) (RRR * Math.acos(0.5 * ((1.0 + q1) * q2 - (1.0 - q1) * q3)) + 1.0);
            }
        }
        return distance;
    }

    public void problemaGeografico(File file, int n) {
        String[] numbernode = new String[n];
        float[] firstnode = new float[n];
        float[] secondnode = new float[n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            ReadFileGeografica(file, numbernode, firstnode, secondnode);
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println(
                "La distancia total para este archivo es de " + distanciaGeometrica(numbernode, firstnode, secondnode));
    }
}
