package General;
import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;

public class General {
    public General(){}
    public int numberCities(File file) {
        int cities = 0;
        String name = file.getName();
        String number = "";
        for (int i = 0; i < name.length(); i++) {
            char c = name.charAt(i);
            switch (c) {
                case '1':
                    number += "1";
                    break;
                case '2':
                    number += "2";
                    break;
                case '3':
                    number += "3";
                    break;
                case '4':
                    number += "4";
                    break;
                case '5':
                    number += "5";
                    break;
                case '6':
                    number += "6";
                    break;
                case '7':
                    number += "7";
                    break;
                case '8':
                    number += "8";
                    break;
                case '9':
                    number += "9";
                    break;
                case '0':
                    number += "0";
                    break;
            }
        }
        return cities = Integer.parseInt(number);
    }

    public String typeProblem(File file) {
        String type = "";
        try (BufferedReader br = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = br.readLine()) != null) {
                if(line.compareTo("TYPE : TSP")==0)
                    type="TSP";
                else if(line.compareTo("TYPE: ATSP")==0)
                    type="ATSP";
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        return type;
    }
}
