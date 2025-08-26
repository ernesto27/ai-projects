# Claude tutorial, tips

## Tabla de contenidos

- [Introducción, definicion](#introducción-definicion)
- [Instalacion](#instalacion)
- [Repositorio para seguir el tutorial](#repositorio-para-seguir-el-tutorial)
- [Init claude command](#init-claude-command)
- [Explicacion, analisis del proyecto](#explicacion-analisis-del-proyecto)
- [Implementar feature, tarea](#implementar-feature-tarea)
- [Revisar cambios](#revisar-cambios)
- [Use de @ para archivos](#use-de--para-archivos)
- [Fix BUG](#fix-bug)
- [Modo plan](#modo-plan)
- [Modo Thinking](#modo-thinking) - **No implementado aún**
- [Claude Code Tools](#claude-code-tools)
- [Agregar imagenes al contexto](#agregar-imagenes-al-contexto)
- [Uso de comandos](#uso-de-comandos)
- [Retomar conversacion](#retomar-conversacion)
- [Uso de sub-agentes](#uso-de-sub-agentes)
- **Use MCP - playwright** - **No implementado aún**
- **Use memory** - **No implementado aún**

## Introducción, definicion

Claude code desde hace unos meses se convirtio en mi herramienta favorita de AI para tareas relacionadas con la tecnologia, como previo usuario de Github Copilot y eventual de Cursor note estas. Si bien copilot funciona excelente y todavia lo utlizo en algunas ocasiones, note diferencias en particular que me llevaron a utilizar mas claude code:

Claude Code es una agente de IA que se ejecuta desde la terminal que permite entre otras cosas investigar proyecto, agregar funcionalidades, arreglar bugs, etc. salio en preview en febrero de 2025 y en Mayo estuvo disponible para todos los usuarios.

En los ultimos meses otras empresas como google y openAI sacaron herramientas similares:

- [Gemini CLI](https://blog.google/technology/developers/introducing-gemini-cli-open-source-ai-agent/)
- [OpenAI Codex CLI](https://help.openai.com/en/articles/11096431-openai-codex-cli-getting-started)

### Ventajas de Claude Code

- **Ejecucion en terminal**: en mi caso al utilizar constantemente la terminal se convierte en algo mas natural

- **Separacion de la interfaz de AI por fuera del editor**: algo que no me gusta mucho con respecto a la integracion de IA en los editores es que se vuelve un poco complejo la UX al tener muchos paneles al mismo tiempo y se pierde foco en el editor de codigo.

![Claude tutorial overview](./image01.png)

- **Agnostico e independiente de editores**: aunque en caso de requerirlo Claude Code se puede integrar con editores como VSCode, Cursor.

## Instalacion

### Requerimientos

Al momento se requiere tener una subscripcion a Anthropic para utilizar claude code:

- **Pro**: Permite utilizar el modelo Claude code 4, el cual es excelente para las tareas diaras de desarrollo.
- **Max**: Permite utilizar el model Claude Opus, el cual es el model mas avanzado de Anthropic el cual esta preparado para tareas mas complejas.

O utilizar un Api key, es este caso se va a pagar por uso, el cual se puede incrementar rapidamente dependiendo del tipo de tarea, no olvidad que al ser un agente de IA, cada pregunta, consultar puede generar multiples consultas a la API de Anthropic.

**NodeJS**: 18 o superior

### Instalacion de claude code

```bash
npm install -g @anthropic-ai/claude-code
```

Los ejemplos de este tutorial estan realizando con el model Pro.

**Otros requerimientos**:
- Docker 

## Repositorio para seguir el tutorial

Los ejemplos van a estar basados en este repositorio:

> TODO - CREAR PROYECTO ESTILO REDDIT O TWITTER.

## Init claude command

Lo primero que se debe hacer es ejecutar el comando claude dentro del reposotirio o del proyecto que deseen utilizar.

Acepten los permisos para que Claude puede modificar los archivos del proyecto.

El primer comando a ejecutar es el siguiente:

```bash
/init 
```

Este comando va a inicializar el agente, el cual va a analizar el proyecto acutal, estructra, archivos, dependencias, etc y va a generar un archivo de configuracion llamado `CLAUDE.md` en la raiz del proyecto con la documentacion del proyecto, la utilidad de este archivo es para que Claude envie estas instrucciones cada vez que se inicie una nueva conversacion y puede tener "memoria" de las practicas y estructura del proyecto.

### Tipos de archivos CLAUDE.md

Hay tres tipos de archivos CLAUDE.md que se pueden generar:

1. **Por proyecto**: el cual se va a subir al repositorio.
2. **Para uso local**: con reglas especificas del proyecto, el cual no se va a subir al repositorio.
3. **Para uso general**: en todos los proyectos en los que se utilize Claude Code.

En nuestro caso vamos a utilizar la opcion 1.


## Explicacion, analisis del proyecto
Uno de los usos mas efectivos de Claude Code es el analisis, explicacion de un proyecto, esto puede ser utili en alguno de estos casos:

- Cuando nos incorporamos a un proyecto nuevo que ya tiene un tiempo considerable de desarrollo.
- Para entender un feature.
- Para revision de un pull request.

Para esto podriamos usar un promp como el siguiente:

```bash
Haz un analisis del proyecto, codigo, estructura, dependencias, integraciones, etc.
Genera un diagrama de flujo base y tambien en formato mermaid, guarda el contenido en un archivo llamado RESEARCH.md
```

Aca vemos a Claude code en accion,  generando un TODO list y utilizando diferentes tools para investigar el proyecto.

![Claude tutorial overview](./explain.png)

Resultado generado:

- [RESEARCH.md](RESEARCH.md)




## Implementar feature, tarea

A la hora de implementar una tarea, lo recomendable es ser lo mas especifico posible en el prompt, esto es importante ya que de otra manera Cluade Code va a inferir, suponer y puede darse que el resultado este lejos de lo que se espera, como toda herramienta si se utiliza de manera incorrecta puede generar mas trabajo que beneficio.

Tomemos como ejemplo que tenemos una API backend en el cual queremos agregar una conexion a una base de datos Mysql.

### Ejemplo de un mal prompt

```bash
Agrega una base de datos mysql a la API.
```

Si bien esto puede dar un resultado funcional, lo que va a suceder es que claude code va a tomar desiciones que tal vez no sean las esperadas, como utlizar una librearia que no queremos, hacer algo mas complejo de lonecesario, etc.

### Un mejor prompt seria

```bash
Quiero agregar a la api existente una conexion a la base de datos mysql.

- Agregar servicio mysql en docker compose, utiliza la version 8.
- Utiliza la libreria "goose" para las migraciones de la base de datos.
- Utiliza la libreria GORM para la conexion a la base de datos.
- Crea una tabla llamada "posts" en un archivo de migracion con los siguientes campos: id, title, content
- Todos los archivo asociados a la base de datos deben estar en una carpeta llamada "db" 
```

Al detallar en especifo las versiones, librerias, rutas se puede generar un resultado mas acorde a lo que se espera, asi que siempre es recomendable tomarse unos minutos para pensar el detalle del prompt.

## Revisar cambios

A medida que Claude code va generando cambios en el proyecto, nos da la posibilidad de configurar el "auto-accept" ya sea seleccionando esa opcion cuando Claude code necesite editar un archivo o ejecutando el comando o presionando dos veces shift-tab 1 vez.

![Claude tutorial overview](./auto-accept.png)

Si bien esta opcion puede ser util en ciertas circunstancias,  como por ejemplo veremos mas adelante utilizando el agente en background,  lo recomendable es ir revisando los cambios linea por linea, para asegurarnos de que todo este correcto e ir ajustando, indicando a Claude que queremos cambiar,  de manera similar a lo que hariamos en un code review.

Recordar que los LLM son no deterministicos, es decir que pueden generar diferentes resultados para el mismo prompt,  ademas de que nosotros vamos a ser los responsables si el codigo se sube a produccion y tiene errores o bugs.


# Use de @ para archivos.

Si bien Claude Code tiene la capacidad de entender el contexto del proyecto y buscar el archivo en que se debe aplicar un cambio,  podemos ser mas especificos y decirle a Claude cual es el archivo especifico que necesitamos editar,  esto va a permitir que Claude sea mas rapido y no malgastar tokens utilizando las tools.

![File](./file.png)


# Fix BUG
Generalmente cuando tenemos un bug y acceso al texto de este en particular, copiar el mensaje en el Claude puede ser suficiente.

![Fix Bug](./debug.png)

Aca puede sumar comentar que fue lo que intentamos previamente hacer para solucionar el bug, agregar archivos,  contexto, etc.


# Modo plan
Claude tiene un model plan, el cual se puede habilitar presionando shift+tab 2 veces.

![Plan](./plan.png)

Este modo es util para tareas medianas/complejas, ya que permite que Claude pueda planificar los pasos a seguir de una manera estructurada, iterando y validando cada paso.

Por ejemplo.

```bash
Integrate our existing Golang API with AWS services: S3 for file uploads and SQS for async task processing
```

Ejemplo de resultado:

![Plan result](./plan2.png)

Despues de esto podemos aceptar el plan o seguir ajustando el plan.


# Modo Thinking


# Claude Code Tools

Los agentes de AI como Claude Code, copilot, cursor, etc tienen la capacidad de utilizar herramientas para realizar tareas especificas en vias de solucionar un problema.

En el caso de Claude Code, algunas de las tools disponible al momento son:

### Herramientas de Archivos y Proyecto

- **Read**: Leer archivos del proyecto
- **Write**: Crear nuevos archivos o en caso de ser necesario editar
- **Edit/MultiEdit**: Modificar archivos existentes con operaciones precisas de buscar y reemplazar

### Herramientas de Búsqueda

- **Glob**: Buscar patrones de archivos (ej: `**/*.ts`, `src/**/*.tsx`)
- **Grep**: Búsqueda basada en el comando grep
- **LS**: Listar contenidos de directorios

### Herramientas de Sistema

- **Bash**: Ejecutar comandos de shell con comillas adecuadas y soporte de ejecución en segundo plano

### Herramientas Web

- **WebFetch**: Recuperar y analizar contenido web
- **WebSearch**: Buscar en la web información actual



# Agregar imagenes al contexto.

Para agregar imagenes a la conversacion, se debe copiar el path de la imagen ( 
ejemplo ubuntu,  click derecho sobre la imagen y seleccionar "copiar") y  posteriormente pegar con  ctrl+shift+v en la terminal.

Esto nos puede servir tanto para analizar diagramas, diseños, errores, etc.

![Image from clipboard](./image-command.png)


# Uso de comandos

Algo que puede suceder es que en algun momento nos encontremos repitiendo una y otra vez el mismo prompt o tarea,  ya sea para explicar un proyecto,  actualizar , revision, etc.
Para evitar esto podemos crear archivos en formato Markdown en la carpeta `.claude/commands` 

Por ejemplo veamos un comando para hacer un review de seguridad.

securty-audit.md
```markdown
Perform a security audit on this code. Look for:
- SQL injection vulnerabilities
- XSS risks
- Authentication/authorization issues
- Sensitive data exposure
- Input validation problems
```

Otro ejemplo de comando para actualizar el archivo README.md con los ultimos cambios.

update-readme.md
```markdown
Actualiza README.md con los cambios recientes:

- Analiza archivos modificados y detecta: nuevas features, dependencias, configuraciones, estructura
- Actualiza secciones afectadas: Features, Installation, Usage, Configuration, API docs
- Añade/actualiza Changelog con formato: ### Added/Changed/Fixed/Removed
- Mantén formato y idioma existente, no elimines contenido válido
- Lista los cambios propuestos antes de aplicarlos
```

Estos comando se utiliza desde claude code de la siguiente manera:

```bash
/security-audit
```


```bash
/update-readme
```


# Retomar conversacion.

Si necesitamos retomar una conversacion previa,  debido a que no terminamos una tarea que no terminamos o presionamos ctrl+c por error,  podemos utilizar el comando:

```bash
/resume
```

esto nos va a mostar un listado con las ultimas conversaciones,  seleccionamos la que queremos retomar y podemos proseguir.

![Resume](./resume.png)


# Uso de sub-agentes
Claude nos permiter poder crear sub-agentes, para realizar diversas tareas,  esto puede resultar algo similar a lo que vimos anteriormente con los comandos, pero hacer de esta manera tiene algunas diferencias:

- El agente tiene su propia ventana de contexto, o sea no comparte el contexto con Claude como pasa con los comandos. 
- Tiene su propio promp de sistem
- Se pueden habilitar o deshabilitar diferentes tools.

Para crear un sub-agente,  se utiliza el comando:

```bash
/agent
```
Create new aget

![Sub-agent](./agent1.png)

Aca podemos generar un agente especificamente para el proyecto desde el que estemos ejecutando Claude o a nivel general del sistema.

En la opcion siguiente seleccionamos la opcion.

1. Generate with Claude (recommended)

En la siguiente seccion tenemos que definir el rol del agente.

![Sub-agent description](./agent2.png)

En este caso vamos a crear un agente para que genere un code-review del codigo.

Ejemplo

```markdown
You are a senior Golang specialist performing a comprehensive code review. Analyze this Go code with focus on:

GOLANG BEST PRACTICES:
- Idiomatic Go patterns and conventions
- Effective use of goroutines and channels
- Proper error handling (not just if err != nil)
- Context usage and cancellation
- Defer statements placement and usage

PERFORMANCE & MEMORY:
- Memory allocations and potential leaks
- Goroutine leaks
- Efficient use of slices vs arrays
- String concatenation optimization
- Sync.Pool usage where appropriate
- Benchmark-worthy code sections

CONCURRENCY SAFETY:
- Race conditions
- Proper mutex usage (sync.Mutex vs sync.RWMutex)
- Channel deadlocks
- WaitGroup patterns
- Atomic operations where needed

CODE STRUCTURE:
- Interface design and composition
- Package organization and naming
- Exported vs unexported identifiers
- Embedded types usage
- Method receivers (pointer vs value)

TESTING & RELIABILITY:
- Table-driven test completeness
- Benchmark tests for critical paths
- Race detector compatibility
- Mock interfaces design
- Test coverage gaps

SECURITY:
- SQL injection in database/sql usage
- Command injection in os/exec
- Path traversal vulnerabilities
- Proper crypto package usage
- Sensitive data in logs

Provide specific line-by-line feedback with severity levels (Critical/Major/Minor) and code examples for improvements. Focus on Go-specific issues that linters might miss.
``` 

Podemos seleccionar las tools que queremos habilitar, en este caso vamos a dejar todo habilitado como viene por defecto.

![Sub-agent tools](./agent4.png)

Si vemos el contenido de nuestro proyecto,  podemos observar que Claude genero un nuevo archivo en la carpeta `.claude/agents` con el nombre golang-code-reviewer.

Podemos invocar el agente de la siguiente manera:

```bash
use golang-code-reviewer
```

Resultado 

![Sub-agent result](./agent-result.png)