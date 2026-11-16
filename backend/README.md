# INVERFIT — Backend API REST (Go + Gin + MongoDB)

Servicio backend modular para la plataforma de gestión del centro de entrenamiento **INVERFIT**. Desarrollado en lenguaje **Go**, utilizando el framework web **Gin** y el driver oficial de **MongoDB**.

---

## 1. Descripción del Proyecto

El backend de INVERFIT implementa una arquitectura monolítica modular por capas (**Router → Handler → Service → Repository → MongoDB**). En esta **Fase 1 (Inicialización)**, se establece la infraestructura base de la aplicación:

- Servidor HTTP de alto rendimiento con Gin.
- Conexión y verificación de conectividad con MongoDB mediante el driver oficial.
- Gestión segura de configuración mediante variables de entorno y archivo `.env`.
- Middleware de producción: CORS restrictivo por orígenes, logging estructurado sin fuga de credenciales y recuperación controlada de panics.
- Endpoints de comprobación de salud para liveness (`/health/live`) y readiness (`/health/ready`).
- Manejo uniforme de errores en formato JSON.
- Cierre ordenado (*Graceful Shutdown*) del servidor HTTP y del pool de MongoDB.

---

## 2. Tecnologías Utilizadas

| Componente | Tecnología | Versión | Rol en el Proyecto |
|---|---|---|---|
| **Lenguaje** | Go | `1.23+` (compatible con toolchain 1.25+) | Lenguaje principal de backend |
| **Framework Web** | Gin (`github.com/gin-gonic/gin`) | `v1.12.0` | Enrutamiento HTTP, binding de DTOs y middlewares |
| **Base de Datos** | MongoDB Official Go Driver | `v1.17.10` | Conexión, ping y persistencia en MongoDB |
| **Logging** | Go Standard Library (`log/slog`) | Estándar | Logging estructurado JSON sin dependencias externas |
| **Pruebas** | Go Testing Framework (`testing`) | Estándar | Pruebas unitarias de config, handlers, router y middlewares |

---

## 3. Estructura de Directorios

```text
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Punto de entrada, inyección de dependencias y graceful shutdown
├── internal/
│   ├── config/
│   │   ├── config.go               # Carga, validación de variables de entorno y sanitización de URIs
│   │   └── config_test.go          # Pruebas unitarias de configuración
│   ├── database/
│   │   ├── mongodb.go              # Conexión oficial MongoDB, ping con timeout y desconexión segura
│   │   └── mongodb_test.go         # Pruebas de conectividad y sanitización de credenciales
│   ├── handlers/
│   │   ├── health_handler.go       # Handlers /health/live y /health/ready con respuestas estandarizadas
│   │   └── health_handler_test.go  # Pruebas de liveness y readiness (mockeado)
│   ├── middleware/
│   │   ├── cors.go                 # Control de acceso CORS por orígenes configurables
│   │   ├── logging.go              # Logging estructurado por petición HTTP (método, ruta, status, latencia)
│   │   ├── recovery.go             # Captura de panics con respuesta JSON 500 uniforme
│   │   └── middleware_test.go      # Pruebas de CORS, logging y recovery
│   ├── modules/
│   │   ├── students/               # Módulo de Alumnos (Fase 2)
│   │   │   ├── model.go            # Modelo BSON de dominio y constantes de estado/frecuencia
│   │   │   ├── dto.go              # DTOs de entrada y salida JSON con mapeo de ObjectID a string
│   │   │   ├── errors.go           # Errores centinela de dominio tipados
│   │   │   ├── repository.go       # Persistencia MongoDB (CRUD, soft-delete, índices y filtros)
│   │   │   ├── repository_test.go  # Pruebas unitarias e integración en vivo con MongoDB
│   │   │   ├── service.go          # Lógica de negocio, validaciones y valores por defecto
│   │   │   ├── service_test.go     # Pruebas unitarias del servicio con mock en memoria
│   │   │   ├── handler.go          # Controladores HTTP Gin y mapeo de errores
│   │   │   └── handler_test.go     # Pruebas HTTP de handlers con httptest
│   │   ├── schedule/               # Módulo de Agenda Semanal y Disponibilidad (Fase 3)
│   │   │   ├── model.go            # Modelos de dominio BSON (DayAvailability, Session, ScheduleException)
│   │   │   ├── dto.go              # DTOs JSON (Disponibilidad, Sesiones, Ocurrencias, Excepciones)
│   │   │   ├── errors.go           # Errores centinela tipados de validación y solapamiento
│   │   │   ├── student_reader.go   # Interfaz desacoplada para resolución de alumnos participantes
│   │   │   ├── repository.go       # Repositorio MongoDB para 3 colecciones con índices únicos
│   │   │   ├── repository_test.go  # Pruebas de persistencia e índices en base de datos de pruebas
│   │   │   ├── service.go          # Lógica de negocio (intervalos, cero pausas, proyección de fechas)
│   │   │   ├── service_test.go     # Pruebas unitarias de reglas de agenda y detección de conflictos
│   │   │   ├── handler.go          # Controladores HTTP Gin y mapeo de errores de agenda
│   │   │   └── handler_test.go     # Pruebas HTTP de endpoints de agenda con Gin y httptest
│   │   └── payments/               # Módulo de Pagos y Finanzas (Fase 4A)
│   │       ├── model.go            # Modelo BSON de pagos, tipos y desinfección CSV Formula Injection
│   │       ├── dto.go              # DTOs JSON (Registro, Actualización, Resumen Financiero, Filtros)
│   │       ├── errors.go           # Errores centinela de dominio tipados (montos, fechas, estados)
│   │       ├── student_reader.go   # Interfaz desacoplada para validar alumno y activar membresías
│   │       ├── repository.go       # Persistencia MongoDB (CRUD, soft-delete, agregaciones e índices)
│   │       ├── repository_test.go  # Pruebas de persistencia e índices en inverfit_test
│   │       ├── service.go          # Lógica financiera (precisión entera CLP, activación, auditoría)
│   │       ├── service_test.go     # Pruebas unitarias de reglas financieras y resúmenes
│   │       ├── handler.go          # Controladores HTTP Gin y streaming progresivo de CSV UTF-8
│   │       └── handler_test.go     # Pruebas HTTP de endpoints de pagos con httptest
│   └── router/
│       ├── router.go               # Registro de middlewares, handlers y respuestas 404/405
│       └── router_test.go          # Pruebas de integración del router
├── .env.example                    # Plantilla de variables de entorno (sin secretos)
├── .gitignore                      # Exclusiones de Git (archivos .env, binarios, temporales)
├── compose.yaml                    # Orquestación de MongoDB en Docker
├── go.mod                          # Definición del módulo y dependencias
├── go.sum                          # Checksums de integridad de dependencias
└── README.md                       # Documentación técnica de esta infraestructura
```

---

## 4. Requisitos Previos

- **Go:** Versión 1.23 o superior instalada.
  * Para comprobar: `go version`
- **MongoDB (Opcional para Fase 1):** Instancia local de MongoDB en el puerto `27017` o cluster remoto en MongoDB Atlas.
  * *Nota:* Si MongoDB no está activo localmente, el servidor levantará igualmente: `/health/live` responderá `200 OK` y `/health/ready` responderá `503 Service Unavailable`, informando que la dependencia no está lista sin bloquear el inicio del proceso.

---

## 5. Configuración de Variables de Entorno

Copie la plantilla `.env.example` para generar su archivo de configuración local:

```bash
cd backend
cp .env.example .env
```

### Variables Soportadas

| Variable | Tipo | Obligatoria | Valor por Defecto | Descripción |
|---|---|---|---|---|
| `APP_ENV` | `string` | No | `development` | Entorno de ejecución (`development`, `test`, `production`). En `production`, Gin opera en modo `release`. |
| `HTTP_PORT` | `int` | No | `8080` | Puerto en el que escuchará el servidor HTTP (1-65535). |
| `MONGODB_URI` | `string` | **Sí** | — | URI de conexión a MongoDB. Ej: `mongodb://localhost:27017` o `mongodb+srv://...` |
| `MONGODB_DATABASE` | `string` | No | `inverfit_dev` | Nombre de la base de datos a utilizar. |
| `CORS_ALLOWED_ORIGINS` | `string` | No | `http://localhost:5173` | Orígenes web autorizados separados por coma. |
| `LOG_LEVEL` | `string` | No | `info` | Nivel mínimo de log (`debug`, `info`, `warn`, `error`). |

---

## 6. Gestión de MongoDB mediante Docker Compose (Fase 1.1)

El entorno local utiliza Docker Compose para orquestar la instancia de base de datos de desarrollo mediante `backend/compose.yaml`.

### Características de la configuración:
- **Imagen:** `mongo:7.0` (imagen oficial estable).
- **Red:** Enlazado exclusivamente a `127.0.0.1:27017` en el host (sin exposición pública).
- **Volumen Persistente:** `inverfit-mongo-data` montado en `/data/db`.
- **Autenticación Obligatoria:** Credenciales gestionadas mediante variables de entorno locales (`MONGO_INITDB_ROOT_USERNAME` y `MONGO_INITDB_ROOT_PASSWORD`).
- **Healthcheck Automatizado:** Utiliza `mongosh` con autenticación interna contra la base de datos `admin`.
- **Política de Reinicio:** `unless-stopped`.

### Comandos de ciclo de vida:

```bash
# Iniciar MongoDB en segundo plano
cd backend
docker compose up -d

# Consultar el estado y comprobación de salud (healthy)
docker compose ps

# Consultar logs del servicio de MongoDB
docker compose logs mongodb
# o seguir logs en tiempo real:
docker compose logs -f mongodb

# Detener el contenedor preservando los datos
docker compose down

# ADVERTENCIA: NUNCA utilice el modificador -v (volume delete)
# docker compose down -v  <-- BORRARÍA EL VOLUMEN PERSISTENTE
```

### Configuración en `.env`:
```env
MONGO_INITDB_ROOT_USERNAME=usuario_local
MONGO_INITDB_ROOT_PASSWORD=password_aleatorio_seguro
MONGODB_DATABASE=inverfit_dev
MONGODB_URI=mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@localhost:27017/inverfit_dev?authSource=admin
```

---

## 7. Despliegue Dockerizado Unificado del Backend (compose.stack.yaml)

Para ejecutar tanto la **API REST en Go** como **MongoDB 7.0** en contenedores independientes y aislados, se dispone del archivo `compose.stack.yaml` en la raíz del repositorio.

### 7.1. Características de la Arquitectura Dockerizada

- **Construcción Multi-stage (`backend/Dockerfile`):**
  - **Etapa 1 (Builder):** Utiliza `golang:1.25-alpine`, descarga dependencias cacheadas y compila un binario estático puro con `CGO_ENABLED=0`, `-trimpath` y `-ldflags="-w -s"`. Soporta arquitecturas cruzadas mediante `TARGETOS` y `TARGETARCH`.
  - **Etapa 2 (Runtime):** Utiliza `alpine:3.20`, incluye certificados CA (`ca-certificates`), zona horaria chilena (`tzdata` con `TZ=America/Santiago`) y herramienta de diagnóstico `curl`.
  - **Seguridad y Privilegios Mínimos:** Se ejecuta bajo el usuario no privilegiado `appuser:appgroup` (UID:GID `10001:10001`). No contiene compiladores ni archivos `.env`.
- **Aislamiento de Red:**
  - Los servicios se comunican a través de una red bridge interna privada (`inverfit-stack-net`).
  - La API Go expone su puerto HTTP mapeado exclusivamente en `127.0.0.1:8080` del host.
  - **MongoDB no expone ningún puerto al host** de manera predeterminada, evitando colisiones con el entorno de desarrollo local (puerto 27017).
- **Persistencia Segura y Volumen Aislado:**
  - Los datos residen en el volumen nombrado `inverfit-stack-mongo-data`, completamente separado del volumen de desarrollo `inverfit-mongo-data`.
- **Credenciales de Aplicación (No Root):**
  - Un script de inicialización (`docker/init-mongo.sh`) crea automáticamente el usuario de aplicación (`inverfit_app`) con rol `readWrite` restringido a la base de datos `inverfit_container`.
  - La API Go se conecta usando este usuario de aplicación, nunca las credenciales de superusuario root.

### 7.2. Configuración de Entorno (.env.stack)

Copie la plantilla `.env.stack.example` en la raíz del proyecto para crear su archivo de configuración:

```bash
cp .env.stack.example .env.stack
```

Variables requeridas en `.env.stack`:
```env
APP_ENV=production
HTTP_PORT=8080
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://127.0.0.1:5173

MONGODB_DATABASE=inverfit_container

# Credenciales de superusuario (solo para inicialización interna de MongoDB)
MONGO_INITDB_ROOT_USERNAME=inverfit_root_admin
MONGO_INITDB_ROOT_PASSWORD=super_secreto_root_password_2026

# Credenciales de aplicación (usadas por la API Go)
MONGO_APP_USERNAME=inverfit_app
MONGO_APP_PASSWORD=super_secreto_app_password_2026

# URI de conexión interna de la API Go hacia el servicio mongodb
MONGODB_URI=mongodb://inverfit_app:super_secreto_app_password_2026@mongodb:27017/inverfit_container?authSource=inverfit_container
```

### 7.3. Comandos de Gestión del Stack Dockerizado

#### 1. Construir la imagen de la API de forma independiente
```bash
docker build -t prototipo-inverfit-api:latest -f backend/Dockerfile backend
```

#### 2. Levantar el stack completo en segundo plano
```bash
docker compose -f compose.stack.yaml --env-file .env.stack up -d --build
```

#### 3. Verificar el estado y comprobaciones de salud (Healthchecks)
```bash
docker compose -f compose.stack.yaml --env-file .env.stack ps
```
Ambos servicios deben reportar estado `healthy`:
- `inverfit-stack-mongodb`: healthcheck nativo con `mongosh` y comando ping.
- `inverfit-stack-api`: healthcheck nativo consultando `/health/ready` mediante `curl`.

Verificación manual de endpoints de salud desde el host:
```bash
# Liveness (proceso Go activo)
curl -i http://127.0.0.1:8080/health/live

# Readiness (conectividad verificada con MongoDB)
curl -i http://127.0.0.1:8080/health/ready
```

#### 4. Consultar logs de los servicios
```bash
# Logs del backend Go
docker compose -f compose.stack.yaml --env-file .env.stack logs -f api

# Logs de MongoDB
docker compose -f compose.stack.yaml --env-file .env.stack logs -f mongodb

# Logs de ambos servicios simultáneamente
docker compose -f compose.stack.yaml --env-file .env.stack logs -f
```

#### 5. Detener el stack (`stop` vs `down`)
```bash
# Pausar contenedores sin removerlos (reanudable rápidamente con start):
docker compose -f compose.stack.yaml --env-file .env.stack stop

# Detener y remover contenedores y red, preservando intacto el volumen de datos:
docker compose -f compose.stack.yaml --env-file .env.stack down

# ¡ADVERTENCIA CRÍTICA!
# NUNCA ejecute 'down -v' a menos que desee eliminar irrevocablemente los datos:
# docker compose -f compose.stack.yaml down -v  <-- BORRARÍA EL VOLUMEN PERSISTENTE
```

#### 6. Respaldo y Restauración del Volumen Dockerizado

**Crear un respaldo comprimido (Backup):**
```bash
# Crear directorio local de backups si no existe
mkdir -p backups

# Ejecutar un contenedor efímero para respaldar el contenido del volumen
docker run --rm \
  -v inverfit-stack-mongo-data:/data/db:ro \
  -v "$(pwd)/backups:/backup" \
  alpine tar czvf /backup/inverfit-stack-mongo-backup.tar.gz -C /data/db .
```

**Restaurar un respaldo en el volumen (Restore):**
```bash
# Detener el stack antes de restaurar
docker compose -f compose.stack.yaml --env-file .env.stack stop mongodb

# Restaurar datos desde el archivo de backup
docker run --rm \
  -v inverfit-stack-mongo-data:/data/db \
  -v "$(pwd)/backups:/backup" \
  alpine sh -c "rm -rf /data/db/* && tar xzvf /backup/inverfit-stack-mongo-backup.tar.gz -C /data/db"

# Volver a iniciar el servicio
docker compose -f compose.stack.yaml --env-file .env.stack up -d mongodb
```

### 7.4. Comparativa de Entornos: Desarrollo Local vs. Stack Dockerizado

| Característica | Entorno Desarrollo Local (`backend/compose.yaml`) | Entorno Stack Completo (`compose.stack.yaml`) |
|---|---|---|
| **Ejecución Backend** | Nativo en Windows (`go run ./cmd/api`) | Contenedor Linux Alpine aislado (`prototipo-inverfit-api`) |
| **Ejecución MongoDB** | Contenedor Docker (`inverfit-mongodb`) | Contenedor Docker (`inverfit-stack-mongodb`) |
| **Puerto MongoDB** | `127.0.0.1:27017` expuesto al host | No expuesto al host (solo visible en `inverfit-stack-net`) |
| **Puerto API** | `8080` directo en Windows | `127.0.0.1:8080` mapeado desde contenedor |
| **Volumen de Datos** | `inverfit-mongo-data` | `inverfit-stack-mongo-data` |
| **Base de Datos** | `inverfit_dev` | `inverfit_container` |
| **Usuario MongoDB** | Root `admin` | Usuario de app con rol `readWrite` (`inverfit_app`) |
| **Propósito Principal** | Desarrollo interactivo, depuración con breakpoint y recarga veloz | Validación reproducible de despliegue, pruebas de integración y sandbox previo a producción |

---

## 8. Comandos de Ejecución Local Directa del Backend

### 1. Descargar dependencias
```bash
cd backend
go mod tidy
```

### 2. Ejecutar en modo desarrollo
```bash
go run ./cmd/api
```

### 3. Compilar binario de producción
```bash
# Windows
go build -o bin/api.exe ./cmd/api

# Linux / macOS
go build -o bin/api ./cmd/api
```

---

## 9. Endpoints de la API REST

### 9.1. Endpoints Técnicos y de Monitoreo

#### `GET /health/live`
Comprueba que el proceso de la aplicación está activo y respondiendo solicitudes HTTP. No depende de MongoDB.

* **Respuesta Exitosa (HTTP 200):**
```json
{
  "status": "ok",
  "service": "inverfit-api"
}
```

---

#### `GET /health/ready`
Comprueba que el backend puede atender solicitudes de negocio, verificando que la base de datos MongoDB responda activamente a un `Ping`.

* **Respuesta Exitosa (HTTP 200 - Base de datos disponible):**
```json
{
  "status": "ready",
  "service": "inverfit-api",
  "dependencies": {
    "mongodb": "connected"
  }
}
```

* **Respuesta Degradada (HTTP 503 - Base de datos no disponible):**
```json
{
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "El servicio no se encuentra disponible"
  }
}
```

---

### 9.2. Módulo de Alumnos (`/api/v1/students`) — Fase 2

Colección MongoDB: `students`. Implementa borrado lógico (*soft delete*) mediante el campo técnico `archivedAt`. Los registros archivados quedan excluidos automáticamente de las consultas y búsquedas sin generar cascadas destructivas sobre pagos, contratos o agendas históricas.

#### Endpoints Disponibles:

| Método | Ruta | Descripción | Parámetros / Payload | Respuestas |
|---|---|---|---|---|
| `GET` | `/api/v1/students` | Listado de alumnos activos con filtros opcionales | Query params: `name` (búsqueda parcial insensible a mayúsculas), `status` (`Activo`, `Por revisar`, `Inactivo`), `frequency` (`1 vez/sem.`, `2 veces/sem.`, `3 veces/sem.`, `4 veces/sem.`) | `200 OK` (Array de `StudentResponse`) |
| `GET` | `/api/v1/students/:id` | Obtener detalle completo de un alumno | Path param: `id` (hexadecimal de 24 caracteres correspondiente al ObjectID) | `200 OK`, `400 Bad Request` (`INVALID_ID`), `404 Not Found` (`STUDENT_NOT_FOUND`) |
| `POST` | `/api/v1/students` | Registrar nuevo alumno en el sistema | JSON Body con `fullName` obligatorio y campos opcionales | `201 Created` (`StudentResponse`), `400 Bad Request` |
| `PUT` | `/api/v1/students/:id` | Actualización parcial de ficha de alumno | Path param: `id`. JSON Body con campos a modificar | `200 OK` (`StudentResponse`), `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/api/v1/students/:id` | Baja lógica / archivado de un alumno | Path param: `id` | `200 OK` (`{ message, id }`), `400 Bad Request`, `404 Not Found` |

#### Estructura de Datos (`StudentResponse` JSON):
```json
{
  "id": "67cb137b00388b15d9cf9e52",
  "fullName": "Carlos Mendoza",
  "email": "carlos.mendoza@email.com",
  "phone": "+56 9 9123 4567",
  "planName": "Plan estándar",
  "frequency": "3 veces/sem.",
  "membershipStatus": "Activo",
  "fixedSchedule": "Lun-Mie-Vie 08:30",
  "membershipStartDate": "2026-11-06",
  "membershipDueDate": "Por definir",
  "notes": "Objetivo hipertrofia",
  "healthConditions": "Sin observaciones de salud registradas",
  "emergencyContact": "+56 9 8888 7777",
  "createdAt": "2026-11-06T18:00:00Z",
  "updatedAt": "2026-11-06T18:00:00Z"
}
```

---

### 9.3. Módulo de Agenda Semanal y Disponibilidad (`/api/v1/schedule`) — Fase 3

Gestiona la disponibilidad del entrenador por intervalos libres, la programación de sesiones recurrentes individuales y grupales, y las excepciones puntuales de calendario.

- **Colecciones MongoDB e Índices:**
  - `availabilities`: Configuración horaria por día (`Lun`..`Dom`). Índice único en `{ day: 1 }`.
  - `sessions`: Bloques de entrenamiento recurrentes. Índices en `{ day: 1, start_time: 1 }`, `{ student_ids: 1 }`, `{ status: 1 }`.
  - `schedule_exceptions`: Modificaciones puntuales a una fecha específica (`YYYY-MM-DD`). Índice único compuesto en `{ session_id: 1, original_date: 1 }`, índice en `{ original_date: 1 }`.
- **Desacoplamiento Modular:** El módulo no accede a la base de datos de alumnos. Consume la interfaz `StudentReader` inyectada mediante un adaptador desde `cmd/api/main.go`.
- **Reglas de Negocio Implementadas:**
  - **Agenda de administración exclusiva:** Administrada únicamente por el entrenador.
  - **Disponibilidad por intervalos:** Soporta de 0 a múltiples turnos diarios dentro del rango operativo del centro (06:30 a 22:00).
  - **Duración sugerida y libre:** 75 minutos por defecto, editable con libertad matemática estricta (`startTime < endTime`).
  - **Sin descansos forzados:** Sesiones consecutivas (ej. 08:00–09:15 y 09:15–10:30) son 100% válidas. El algoritmo de solapamiento (`startA < endB && startB < endA`) no penaliza bordes coincidentes.
  - **Sesiones grupales ilimitadas:** Sin cupo artificial arbitrario; validación estricta contra alumnos duplicados en una misma sesión.
  - **Excepciones de calendario:** Ocurrencias puntuales canceladas (`cancelled`) o reprogramadas (`rescheduled`).
  - **Detección de conflicto al reducir disponibilidad:** Retorna `409 Conflict` con código `AVAILABILITY_CONFLICT` y la lista de sesiones afectadas si el entrenador intenta eliminar o acotar un intervalo que contiene sesiones activas.
  - **Zona Horaria Oficial:** `America/Santiago` integrada vía `time/tzdata`. Fechas en `YYYY-MM-DD`, horas en `HH:mm`, timestamps del sistema en UTC.

#### Endpoints Disponibles:

| Método | Ruta | Descripción | Parámetros / Payload | Respuestas |
|---|---|---|---|---|
| `GET` | `/api/v1/schedule/availability` | Obtener configuración semanal completa de disponibilidad | Ninguno | `200 OK` (Array de 7 días) |
| `GET` | `/api/v1/schedule/availability/:day` | Obtener disponibilidad de un día específico (`Lun`..`Dom`) | Path param: `day` | `200 OK`, `400 Bad Request` |
| `PUT` | `/api/v1/schedule/availability/:day` | Configurar intervalos de disponibilidad para un día | Path param: `day`. Body: `{ enabled, intervals: [{ startTime, endTime }] }` | `200 OK`, `400 Bad Request`, `409 Conflict` (`AVAILABILITY_CONFLICT`) |
| `DELETE` | `/api/v1/schedule/availability/:day` | Deshabilitar disponibilidad de un día | Path param: `day` | `200 OK`, `409 Conflict` |
| `GET` | `/api/v1/schedule/sessions` | Proyección de ocurrencias de calendario o listado de sesiones | Query params opcionales: `from`, `to` (`YYYY-MM-DD`), `studentId` | `200 OK` (Array de ocurrencias resueltas o sesiones) |
| `POST` | `/api/v1/schedule/sessions` | Crear nueva sesión recurrente | Body: `{ day, startTime, endTime, studentIds, startDate, notes }` | `201 Created` (`SessionResponse`), `400 Bad Request`, `409 Conflict` (`SESSION_OVERLAP`) |
| `GET` | `/api/v1/schedule/sessions/:id` | Detalle de sesión recurrente | Path param: `id` | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `PUT` | `/api/v1/schedule/sessions/:id` | Actualizar sesión recurrente | Path param: `id`. Body con campos a actualizar | `200 OK`, `400 Bad Request`, `409 Conflict`, `404 Not Found` |
| `DELETE` | `/api/v1/schedule/sessions/:id` | Archivar / cancelar sesión recurrente | Path param: `id` | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `POST` | `/api/v1/schedule/sessions/:id/exceptions` | Crear excepción de calendario (cancelación o reagendamiento) | Path param: `id`. Body: `{ date, type, newDate, newStartTime, newEndTime, notes }` | `201 Created` (`ExceptionResponse`), `400 Bad Request`, `409 Conflict` |
| `GET` | `/api/v1/schedule/exceptions` | Listar excepciones en un rango de fechas | Query params: `from`, `to` | `200 OK` (Array de `ExceptionResponse`) |
| `DELETE` | `/api/v1/schedule/sessions/:id/exceptions/:exceptionId` | Eliminar excepción y restaurar recurrencia normal | Path params: `id`, `exceptionId` | `200 OK`, `400 Bad Request`, `404 Not Found` |

---

### 9.4. Módulo de Pagos y Membresías (`/api/v1/payments`) — Fase 4A

Gestiona el registro de pagos, trazabilidad financiera, cálculo de balances y resúmenes de cobranza, así como la exportación streaming de informes a CSV.

- **Colección MongoDB e Índices:**
  - `payments`: Documentos individuales de pagos confirmados y pendientes.
  - Índices: `{ "student_id": 1 }`, `{ "status": 1 }`, `{ "created_at": -1 }`, `{ "archived_at": 1 }`.
- **Desacoplamiento Modular:** El módulo no accede directamente a la colección de alumnos. Utiliza la interfaz `StudentReader` inyectada mediante un adaptador en `cmd/api/main.go`, la cual valida la existencia del alumno y activa su membresía (`Por revisar` → `Activo`) al registrar un pago exitoso.
- **Reglas de Negocio Implementadas:**
  - **Gestión exclusiva del entrenador:** El entrenador Pablo administra y registra los pagos; no hay pasarelas ni cobros automáticos en esta fase.
  - **Precisión Monetaria Estricta:** Los montos se manejan en pesos chilenos (CLP) como enteros de 64 bits (`int64`), evitando cualquier error de coma flotante (`float32`/`float64`).
  - **Preservación Irrestricta del Historial:** Los registros nunca se destruyen físicamente; las bajas o anulaciones se gestionan mediante baja lógica (`archivedAt`).
  - **Inmutabilidad Histórica:** Los pagos pasados conservan la información pactada en su fecha de emisión (`planName`, `frequency`, `amount`); no se sobreescriben retrospectivamente ante cambios en la ficha del alumno.
  - **Resumen Financiero en Vivo (`/summary`):** Métricas calculadas en tiempo real directamente sobre los registros persistidos de MongoDB mediante agregaciones, evitando totales desactualizados.
  - **Exportación Streaming a CSV (`/export`):** Generación progresiva con codificación UTF-8 y BOM (`\xef\xbb\xbf`) para compatibilidad directa con Microsoft Excel en Windows. Incluye protección automática contra inyección de fórmulas CSV (*CSV Formula / DDE Injection*) ante prefijos `=`, `+`, `-`, `@`.

#### Endpoints Disponibles:

| Método | Ruta | Descripción | Parámetros / Payload | Respuestas |
|---|---|---|---|---|
| `GET` | `/api/v1/payments` | Listar y filtrar pagos | Query params: `studentId`, `status`, `search`, `from`, `to` | `200 OK` (Array de `PaymentResponse`) |
| `POST` | `/api/v1/payments` | Registrar nuevo pago o confirmar pendiente | Body: `{ studentId, amount, dueDate, paymentDate, paymentMethod, notes, transactionReference }` | `201 Created` (`PaymentResponse`), `400 Bad Request`, `404 Not Found` |
| `GET` | `/api/v1/payments/summary` | Consultar balance y KPIs financieros consolidados | Ninguno | `200 OK` (`FinancialSummaryResponse`) |
| `GET` | `/api/v1/payments/export` | Exportar historial de pagos en CSV progresivo | Query params opcionales equivalentes al listado | `200 OK` (`text/csv; charset=utf-8`) |
| `GET` | `/api/v1/payments/student/:studentId` | Consultar historial de pagos de un alumno | Path param: `studentId` | `200 OK` (Array), `400 Bad Request` |
| `GET` | `/api/v1/payments/:id` | Detalle completo de un pago | Path param: `id` | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `PUT` | `/api/v1/payments/:id` | Actualizar datos de un pago | Path param: `id`. Body con campos a modificar | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/api/v1/payments/:id` | Anulación o baja lógica de pago (*soft delete*) | Path param: `id` | `200 OK` (`{ message, id }`), `400`, `404` |
| `POST` | `/api/v1/payments/:id/simulate-reminder` | Registrar envío de recordatorio de cobro | Path param: `id` | `200 OK` (`PaymentResponse`), `400`, `404` |
| `POST` | `/api/v1/payments/:id/simulate-receipt` | Registrar emisión de comprobante de pago | Path param: `id` | `200 OK` (`PaymentResponse`), `400`, `404` |

---

### 9.5. Formato Uniforme de Respuestas de Error
Toda respuesta de error HTTP (400, 404, 405, 409, 500, 503) respeta la siguiente estructura estándar:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Ruta no encontrada"
  }
}
```

En caso de conflicto de disponibilidad (`AVAILABILITY_CONFLICT`), se anexa la lista detallada de sesiones en conflicto:
```json
{
  "error": {
    "code": "AVAILABILITY_CONFLICT",
    "message": "no es posible reducir la disponibilidad porque existen sesiones activas programadas"
  },
  "conflicts": [ ... ]
}
```

---

### 9.6. Integración Frontend (Modo Dual Mock / API)

El frontend en Vue 3 dispone de una capa de delegación transparente configurada mediante variables de entorno:

| Variable Vite | Valores | Por Defecto | Comportamiento |
|---|---|---|---|
| `VITE_STUDENTS_DATA_SOURCE` | `mock` \| `api` | `mock` | `mock`: Almacena en `localStorage` (ideal para Cloudflare Pages y demo offline).<br>`api`: Consume directamente los endpoints REST de Alumnos en Go en tiempo real. |
| `VITE_SCHEDULE_DATA_SOURCE` | `mock` \| `api` | `mock` | `mock`: Gestiona agenda y disponibilidad en `localStorage`.<br>`api`: Consume endpoints de agenda en Go con proyección de ocurrencias y validaciones en servidor. |
| `VITE_PAYMENTS_DATA_SOURCE` | `mock` \| `api` | `mock` | `mock`: Gestiona pagos y comprobantes en `localStorage`.<br>`api`: Consume endpoints de pagos en Go con cálculo en vivo de balances y exportación streaming. |
| `VITE_API_BASE_URL` | URL HTTP | `http://localhost:8080/api/v1` | URL base del backend Go. |

---

## 10. Ejecución de Pruebas y Análisis Estático

El proyecto incluye pruebas unitarias completas para configuración, base de datos (con timeouts y sanitización), middlewares y handlers:

```bash
# Ejecutar todas las pruebas unitarias
go test -v ./...

# Análisis estático de código
go vet ./...

# Compilación de verificación
go build ./cmd/api
```

---

## 11. Consideraciones de Seguridad y Buenas Prácticas

1. **Protección Absoluta de Secretos:**
   - Los archivos `.env` reales están excluidos tanto en el `.gitignore` de la raíz del repositorio como en `backend/.gitignore`.
   - `RedactedMongoURI()` ofusca automáticamente contraseñas de las URIs antes de imprimirlas en logs o mensajes de error.
2. **CORS Restrictivo:**
   - No se utiliza comodín `*`. Solo los orígenes listados en `CORS_ALLOWED_ORIGINS` reciben cabeceras de autorización de acceso cruzado.
3. **Control de Fugas en Logging:**
   - El middleware `Logging` registra metadata controlada (método, path estático, status, latencia, IP del cliente). Omite intencionadamente query strings para evitar registrar tokens, parámetros sensibles o identificadores, y nunca registra encabezados `Authorization`, contraseñas ni cuerpos JSON.
4. **Recuperación ante Panics (`Recovery`):**
   - En caso de pánico inesperado, el servidor no se cae ni envía la traza interna al cliente; responde con código `500` uniforme y registra el stack trace únicamente en el log del servidor.
5. **Apagado Ordenado (*Graceful Shutdown*):**
   - Ante señales `SIGINT` o `SIGTERM`, el servidor HTTP dispone de un plazo de 10 segundos para completar las peticiones en vuelo antes de cerrar el pool de MongoDB y detener el proceso.
