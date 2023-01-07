# mail-indexer
##### English [README.md](https://github.com/mata649/mail_indexer/blob/master/README.md)
Esta aplicación se divide en **Indexer** y **Efinder**. El **Indexer** obtiene correos electrónicos de una ruta, luego analiza los correos electrónicos para proporcionarles una estructura y finalmente se ingesta un motor de zinc con ellos.
**Efinder** es solo una aplicación web para acceder a correos electrónicos indexados, la interfaz se creó en **Vue.js** con **TailwindCs** y la API se desarrolló en Go usando **Chi** como enrutador , también aproveché la **API** para servir la distribución generada por Vue, usando la URL principal como **Servidor de archivos estáticos**

## Indexer
El Indexer es el núcleo de la aplicación, recibe una ruta de correo electrónico como argumento, para ingerir el motor de zinc con los correos electrónicos..
### Pasos
1. Recibe una ruta donde están los correos
2. Obtiene un slice con las rutas de todos los archivos en el directorio
3. Divide las rutas en segmentos más pequeños, la cantidad de correos electrónicos por segmento se puede configurar en el archivo **config.json**, en la propiedad **emailsPerFile**
4. Itera las rutas divididas al mismo tiempo, la cantidad de goroutines que se ejecutan al mismo tiempo se puede configurar en **config.json**, en la propiedad **nWorkers**
5. Itera las rutas para obtener el correo electrónico de cada archivo, luego el correo electrónico se agrega a un slice de correos electrónicos.
6. Cuando se han obtenido todos los correos electrónicos, el slice de correos electrónicos se convierte en un búfer de bytes en formato NDJSON.
7. Realiza la solicitud a Zinc Engine enviando el búfer de bytes como binario, el usuario, la contraseña y el host de Zinc Engine se pueden configurar en **config.json**
### Ejecuntando el codigo
#### Importante
En Linux el programa arroja un error por la cantidad de archivos que puede abrir un proceso, puedes evitar este error cambiando en el **config.json** el **emailsPerFile**, por defecto está en **1000**, pero si tiene problemas con esto, puede configurar **100** correos electrónicos por archivo.

La segunda opción para resolver esto (**no recomendada**) es aumentar la cantidad de archivos que un proceso permite abrir. Puedes hacerlo con el comando: `ulimit -n 8000`

#### VSCode
Por alguna razón, el error anterior solo aparece si está ejecutando el programa directamente desde una terminal, si ejecuta el código desde la terminal integrada en **VSCode**, el programa no arroja un error incluso con **1000** emailsPerFiel .

Para ejecutar el código puedes hacerlo con el siguiente comando

    ./indexer -emailpath /example/path/to/enron_mail_20110402
O puede ejecutar el archivo main.go con el siguiente comando

    go run cmd/main.go -emailpath /path/example/to/enron_mail_20110402 
Hay algunas banderas que puede agregar si desea ver el profiling del código

    go run cmd/main.go -emailpath /path/example/to/enron_mail_20110402 -cpuprofile profile/cpu_profile.out -memprofile profile/mem_profile.out 

## Test
El programa tiene algunas pruebas unitarias para probar cada función, esto es útil cuando estamos haciendo cambios en el código y queremos comprobar si la función tiene el comportamiento esperado, para ejecutar todas las pruebas del programa puedes hacerlo con la siguiente comando.

    go test ./...
  
## Optimización
La **v1** del programa tiene algunos problemas de optimización, esto se podría comprobar en la carpeta [profile/v1](https://github.com/mata649/mail_indexer/tree/master/indexer/profile/v2).
En la **v1** del programa, todos es relativamente similar hasta los pasos 6 y 7 en la **v2** que son:

6. Cuando se han obtenido todos los correos electrónicos, el slice de correos electrónicos se convierte en un búfer de bytes en formato NDJSON.
7. Realiza la solicitud a Zinc Engine enviando el búfer de bytes como binario, el usuario, la contraseña y el host de Zinc Engine se pueden configurar en **config.json**

En la **v1** después de haber obtenido el slice de correos electrónicos, se guardan en un archivo en un directorio de datos, con una carpeta con la fecha y hora de ejecución como nombre: `2023-01-07 14:13 :56`. En esa carpeta, los archivos se enumeraron con un número, por ejemplo: **file1.ndjson**, **file2.ndjson**, ...., con la cantidad de correos electrónicos establecidos en **emailsPerFile** en cada archivo. Luego, cuando vamos a ingerir el **Zinc Engine** con los correos electrónicos, tenemos que obtener las rutas de los archivos y leerlos para obtener los bytes del búfer para enviar al Zinc Engine.

Esto tiene algunas ventajas y desventajas,
#### Ventajas

 - Disponemos de un registro histórico de los correos electrónicos utilizados en cada ejecución.
  - Si ocurre algún error en la ingesta, podemos registrar el archivo que obtuvimos el error para verificarlo.
  - Si ocurre algún error en la ingesta y no podemos seguir ingestando más archivos, podemos ver cuál fue el último archivo en ser ingerido y luego continuar la ingesta desde ese archivo.

#### Desventajas
  - Para cada archivo, tenemos que abrir dos buffers, el primero cuando estamos creando el archivo y después cuando tenemos que leerlo para ingerir el Zinc Engine.

La diferencia entre crear el archivo para leerlo después y enviar el búfer directamente fue de **2m 55s** a **2m 21s**. Para ser honesto, depende de las necesidades de un proyecto real, si necesita un **Indexer** "tolerante a fallas" cuando si ocurre algún problema tiene un punto de control para continuar después de la ingestión, puede usar **v1** , pero si quieres un **Indexer** que ingiera el **Zinc Engine** rápidamente, puedes usar el **v2**.

Quizás en la **v3** del programa se podría enviar el búfer de bytes directamente para realizar la solicitud y también guardar el archivo, por lo que estamos usando el mismo búfer para ambos propósitos.
## Efinder
Efinder es una aplicación web para buscar un término en correos electrónicos indexados. Se construyó con **Vue.js** y **Tailwind.css** para desarrollar la interfaz, y Go con el enrutador **Chi** para solicitar el **Zinc Engine** y también para servir la distribución de la aplicación Vue, que funciona como un **servidor de archivos estáticos**

Puede ejecutar la aplicación con el siguiente comando:

    ./efinder -host http://localhost:4080 -user admin -password Complexpass#123
### Flags

- **host**: host del **Zinc Engine** para realizar las solicitudes
 - **user**: usuario necesario para realizar la autenticación básica en el **Zinc Engine**
 - **password**: contraseña necesaria para realizar la autenticación básica en el **Zinc Engine**
