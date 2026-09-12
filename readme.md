# Nombre del Proyecto: CINE+
# Integrantes: Berthelot Valentin; Rodriguez Giacomasso Santiago.

# Entrega TP1: Mi Primera Aplicación Web.
* Para ejecutar el proyecto: abrir una terminal en la raíz y ejecutar "go run .", luego, abrir el navegador de preferencia e ingresar a http://localhost:8080

# Entrega TP2: Mi Primera Aplicación Web.
## Para realizar la automatización y persistencia del proyecto, se debe contar con los siguientes requisitos:
* Go en una versión >= 1.25.0.
* Docker y Docker Compose, necesario para levantar el contenedor con la base de datos.
* Make, para la automatización de tareas necesarias para realizar el test.
* Atlas, utilizada para realizar las migraciones de la base de datos.
* SQLC, usado para la generación de código GO a partir de consultas SQL realizadas en el motor PostgreSQL.

## Ejecución de parte 2:
* Para ejecutar esta segunda parte, debe hacerse un clonado del repositorio, obteniendose de esa forma el link, posteriormente en el IDE de preferencia se
abre una terminal, se ejecuta el comando **git clone** junto con el link del repositorio. Por último, desde la ruta principal del proyecto debe ejecutarse **make test**.