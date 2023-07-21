package TSP.Ceil;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

public class Ceil {
    public static void ReadFileCeil(File file, String[] numbernode, float[] firstnode, float[] secondnode)
            throws IOException {
        int count = 0; // Para poder obtener las substring
        int flag = 0; // Se tendra una bandera para comenzar a guardar los datos
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                if (flag == 0) {
                    if (line.compareTo("NODE_COORD_SECTION") == 0 || line.compareTo("NODE_COORD_SECTION ") == 0) // Se
                                                                                                                 // comenzará
                                                                                                                 // a
                                                                                                                 // colocar
                                                                                                                 // los
                                                                                                                 // valores
                                                                                                                 // de
                                                                                                                 // los
                                                                                                                 // nodos
                        flag++;
                    if (line.compareTo("DISPLAY_DATA_SECTION") == 0)
                        flag++;
                } else {
                    if (line.compareTo("EOF") == 0) {
                        continue;
                    } else {
                        if (count < numbernode.length) {
                            String aux = line.trim();
                            String[] datos = aux.split(" ", 3);
                            numbernode[count] = datos[0];
                            firstnode[count] = Integer.parseInt(datos[1].strip());
                            secondnode[count] = Integer.parseInt(datos[2].trim());
                        }
                        count++;
                    }
                }
            }
        }
    }

    // Esta funcion sera para calcular al distancia con CEIL_2D
    public static int distanciaCeil(String[] nodes, float[] x, float[] y) {
        float x1 = 0, x2 = 0, y1 = 0, y2 = 0;
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

    public void problemaCeil(File file, int n) {
        String[] numbernode = new String[n];
        float[] firstnode = new float[n];
        float[] secondnode = new float[n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            ReadFileCeil(file, numbernode, firstnode, secondnode);
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println(
                "La distancia total para este archivo es de "
                        + distanciaCeil(numbernode, firstnode, secondnode));
    }
}
