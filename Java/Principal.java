import java.io.File;

import ATSP.Explicit.Explicit;
import Tools.Tools;
import General.General;
import TSP.ATT.ATT;
import TSP.Ceil.Ceil;
import TSP.Euclidiana.Euclidiana;
import TSP.Geografica.Geografico;
import TSP.Matrix.Matrix;

public class Principal {
    public static void main(String[] args) {
        /************************************************************* */
        //Se crea la ruta para los archivos
        File folder = new File("./Files");
        /************************************************************ */
        /*Se declaran las variables del programa */
        String typeDistance = "";
        int n=0;
        String typeFile="";
        /************************************************************ */
        /*Se crean los objetos de las clases correspondientes de TSP */
        Euclidiana euclidiana=new Euclidiana();
        Geografico geografico=new Geografico();
        Ceil ceil=new Ceil();
        ATT att=new ATT();
        Matrix ma=new Matrix();
        Tools tool=new Tools();
        General g=new General();
        /*********************************************************** */
        /*Se crean los objetos de las clases correspondientes de ATSP */
        Explicit explicit=new Explicit();
        /*********************************************************** */
        for (File file : folder.listFiles()) {
            if(!file.isDirectory()){
                System.out.println("------------------------------------------");
                System.out.println("El nombre del archivo es: " + file.getName());
                n=g.numberCities(file);
                System.out.println("El número de ciudades es "+n);
                typeFile=g.typeProblem(file);
                if (typeFile.compareTo("ATSP") == 0) {
                    typeDistance=tool.typeDistanceATSP(file);
                    switch (typeDistance) {
                        case "Explicita":
                            explicit.problemaExplicit(file, n);
                            break;
                    }
                }else{
                    typeDistance=tool.typeDistance(file);
                    switch (typeDistance) {
                        case "Euclidiana":
                            euclidiana.problemaEuclidiano(file, n);
                            break;
                        case "Geografica":
                            geografico.problemaGeografico(file, n);
                            break;
                        case "Circular":
                            ceil.problemaCeil(file, n);
                            break;
                        case "ATT":
                            att.problemaATT(file, n);
                            break;
                        case "DiagonalSuperior":
                            ma.problemaSuperior(file, n);
                            break;
                        case "DiagonalInferior":
                            ma.problemaInferior(file, n);
                            break;
                    }
                }
            }
        }
    }
}