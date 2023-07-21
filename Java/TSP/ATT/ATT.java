package TSP.ATT;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;

public class ATT {
    public static void ReadFileATT(File file, String[] numbernode, int[] firstnode, int[] secondnode)
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
                            firstnode[count] = Integer.parseInt(datos[1].trim());
                            secondnode[count] = Integer.parseInt(datos[2].trim());
                        }
                        count++;
                    }
                }
            }
        }
    }

    public static int distanciaATT(String[] numbernode, int[] x, int[] y) {
        // Estos dos primeros arreglos son para almacenar la diferencia de los puntos
        float[] xd = new float[numbernode.length];
        float[] yd = new float[numbernode.length];
        int dij = 0, tij = 0;
        int aux = 0;
        double rij = 0;
        // Este for sera para las diferencias de las distancias
        for (int i = 0; i < numbernode.length; i++) {
            if (i < numbernode.length - 1) {
                xd[i] = x[i] - x[i + 1];
                yd[i] = y[i] - y[i + 1];
            }
        }

        // Este for sera para el calculo de la distancia
        for (int i = 0; i < yd.length; i++) {
            rij = Math.sqrt((Math.pow(xd[i], 2) + Math.pow(yd[i], 2)) / 10.0);
            tij = (int) (rij + 0.5);
            dij += (tij < rij) ? tij + 1 : tij;
        }

        return dij;
    }

    public void problemaATT(File file, int n) {
        String[] numbernode = new String[n];
        int[] firstnode = new int[n];
        int[] secondnode = new int[n];
        // Se guardaran los datos del archivo para poder calcular la distancia
        try {
            ReadFileATT(file, numbernode, firstnode, secondnode);
        } catch (Exception e) {
            e.printStackTrace();
        }
        System.out.println(
                "La distancia total para este archivo es de "
                        + distanciaATT(numbernode, firstnode, secondnode));
    }
}
