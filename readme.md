# Servidor websocket con autentiacion en GO

Este proyecto es un servidor que emplea el protocolo de intercambio de mensajes de nivel de aplicacion websocket.

## Caracteristicas

Algunas caracteristicas de este servidor son:
  - Registro/Loggin/Eliminacion de cuentas de usuario
  - Intercambio de mensajes publicos entre todo los usuarios
  - Intercambio de mensajes privados entre usuarios
  - Listado de sesiones de usuario activas y de clientes en tiempo real

## Iniciar proyecto

Para iniciar el proyecto en primer lugar es necesario disponer una maquina con la ultima version de GO.

Teniendo el compilador de GO instalado nos dirigiremos sobre la carpeta del proyecto con `cd gowebsocketauth` y compilaremos el programa ejecutando `go build -o gowebsocketauth cmd/app/*.go`. Esto ultimo compilara todos los ficheros de go del paquete main. Es importante tener en cuenta que si estamos en windows el nombre del programa debera tener el sufijo `.exe` para la extension del programa.

Para ejectar el programa generado le daremos permisos de ejecucion con el comando `chmod +x gowebsocketauth`.

Teniendo los permisos necesario lo ejecutaremos de la siguiente manera `./gowebsocketauth`.

### Ficheros estaticos de cliente

Por defecto el programa tratara de obtener los ficheros estaticos de cliente del directorio `./internal/fileutils/views` cuya ruta es relativa al working directory desde donde se ejecuto el programa; no obstante, podemos indicar la ruta del directorio donde se encontraran los ficheros estaticos del cliente colocandolo como un segundo argumento en la ejecucion del programa, de la siguiente manera `./gowebsocketauth ./ruta/a/estaticos`.

Nuestra carpeta de estaticos debe tener unos ficheros muy concretos. El fichero principal para la pagina web del cliente sera `index.html`; por otra parte, el fichero `main.js` tambien es necesario y es el encargado de implementar la logica de cliente con el servidor websocket.

## Acceder al servidor

Con todo esto podremos acceder al servidor a traves del puerto 3000 en el dominio local en un navegador cualquiera para cargar el fichero principal de la pagina (`index.html`)

![Imagen de cliente websocket predeterminado para el proyecto](https://github.com/panprogramadorgh/gowebsocketauth/blob/master/screenshots/websocket-client.PNG)

## Comandos de interaccion con el servidor

Para que el servidor websocket interprete un mensaje como un comando hemos de colocar una barra ascendente al principio de cualquier comando.

- /register `usuario` `contraseña` : registra un nuevo usuario y hace login automaticamente

- /login `usuario` `contraseña` : hace login a nombre de un usuario especifico

- /logout : cierra la sesion actual (si es que la hay)

- /murder `usuario` `contraseña` : elimina una cuenta de usuario (debemos indicar la contraseña por claros motivos de seguridad)

- /whoami : el servidor te informa sobre a nombre de que usuario estas logeado

- /sessions : muestra una lista con informacion sobre todas las sesiones activas

- /list : muestra una lista con todos los clientes (conexiones con el servidor)

- /tell `usuario` `mensaje` : enviar mensajes privados a otro usuario

- /exit : cierra la conexion (y sesion si la hay) con el sevidor (recargar el cliente para obtener conexion de nuevo)