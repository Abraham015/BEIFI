import matplotlib.pyplot as plt
import numpy as np

# Nombre del archivo de texto
file_path = './Files/Data/fitness_data.txt'

# Se lee el archivo y se obtienen los encabezados
data = np.genfromtxt(file_path, names=True)

# Grafica cada columna con el nombre correspondiente
for column_name in data.dtype.names:
    plt.plot(data[column_name], label=column_name)

# Añade etiquetas y leyenda
plt.xlabel('Índice')
plt.ylabel('Valores')
plt.legend()

# Muestra el gráfico
plt.show()
