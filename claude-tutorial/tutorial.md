# Claude tutorial, tips



- Intriduccion, de copilot a claude code 
- Historia,  definicion de claude
- Otras herramientas similares
- Instalacion,  setup.
- Repositorio para seguir tutorai,  descarga
- Init claude command 
- Explain project prompt
- Add a feature  
- Fix a bug
- Use @ to change specific files 
- Tools description
- Web tool documentation
- Copy image from clipboard
- Use commands 
- Use subagents.
- Use MCP - playwright
- use memory


# Introducción, definicion

Claude code desde hace unos meses se convirtio en mi herramienta favorita de AI para el desarrollo, como previo usuario de Github Copilot y eventual de Cursor note estas 
Si bien copilot funciona excelente y todavia lo utlizo en algunas ocasiones,  note  diferencias en particular que me llevaron a utilizar mas claude code:

Claude Code es una agente de IA que se ejecuta desde la terminal que permite entre otras cosas investigar proyecto, agregar funcionalidades, arreglar bugs, etc.
salio en preview en febrero de 2025 y en Mayo estuvo disponible para todos los usuarios.


En los ultimos meses otras empresas como google y openAI sacaron herramientas similares.

https://blog.google/technology/developers/introducing-gemini-cli-open-source-ai-agent/

https://help.openai.com/en/articles/11096431-openai-codex-cli-getting-started



- Ejecucion en terminal,  en mi caso al utilizar constantemente la terminal se convierte en algo mas natural

- Separacion de la interfaz de AI por fuera del editor,  algo que no me gusta mucho con respecto a la integracion de IA en los editores es  que se vuelve un poco complejo la UX al tener muchos paneles al mismo tiempo y se pierde foco en el editor de codigo. 

![Claude tutorial overview](./image01.png)

- Agnostico e independiente de editores,  aunque en caso de requerirlo Claude Code se puede integrar con editores como VSCode, Cursor.



# Instalacion.

Requerimientos: 
Al momento se requiere tener una subscripcion a Anthropic para utilizar claude code.

Pro: Permite utilizar el modelo Claude code 4, el cual es excelente para las tareas diaras de desarrollo.
Max: Permite utilizar el model Claude Opus, el cual es el model mas avanzado de Anthropic el cual esta preparado para tareas mas complejas.

O utilizar un Api key, es este caso se va a pagar por uso,  el cual se puede incrementar rapidamente dependiendo del tipo de tarea,  no olvidad que al ser un agente de IA,  cada pregunta, consultar puede generar multiples consultas a la API de Anthropic.

NodeJS: 18 o superior

Instalacion de claude code:

```bash
npm install -g @anthropic-ai/claude-code
```

Los ejemplos de este tutorial estan realizando con el model Pro.

- Docker 


# Repositorio para seguir el tutorial

Los ejemplos van a estar basados en este repositorio, 

TODO  - CREAR PROYECTO ESTILO REDDIT O TWITTER.



# Init claude command

Lo primero que se debe hacer es ejecutar el comando claude dentro del reposotirio o del proyecto que deseen utilizar, 

Acepten los permisos para que Claude puede modificar los archivos del proyecto.

El primer comando a ejecutar es el siguiente:

```bash
/init 
```

Este comando va a inicializar el agente, el cual va a analizar el proyecto acutal, estructra, archivos, dependencias, etc y va a generar un archivo de configuracion llamado CLAUDE.md en la raiz del proyecto con la documentacion del proyecto,  la utilidad de este archivo es para que Claude envie estas instrucciones cada vez que se inicie una nueva conversacion y puede tener "memoria" de las practicas y estructura del proyecto.

Hay tres tipos de archivos CLAUDE.md que se pueden generar:
1. Por proyecto,  el cual se va a subir al repositorio.
2. Para uso local, con reglas especificas del proyecto,  el cual no se va a subir al repositorio.
3. Para uso general,  en todos los proyectos en los que se utilize Claude Code.

En nuestro caso vamos a utilizar la opcion 1. 


Servicios.

API - golang -  guardar datos en memoria
Frontend - React - Listado de posts,  detalle
Docker compose   


- Feature claude API


- Agregar mysql8 a docker compose 
- Agregar migraciones con golang using goose,  tabla posts.
- Agregar conexion en API 
- Actualizar endpoint para utilizar guardar datos en mysql y obtenerlos desde mysql


Separar esto en diferentes tareas.


Feature frontend.

- Agregar pantalla para agregar un post
- Actualizar pantalla de listado y detalle de post.

