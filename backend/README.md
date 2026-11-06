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
│   └── router/
│       ├── router.go               # Registro de middlewares, handlers y respuestas 404/405
│       └── router_test.go          # Pruebas de integración del router
├── .env.example                    # Plantilla de variables de entorno (sin secretos)
├── .gitignore                      # Exclusiones de Git (archivos .env, binarios, temporales)
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

## 7. Comandos de Ejecución del Backend

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

## 8. Endpoints Disponibles en Fase 1

### `GET /health/live`
Comprueba que el proceso de la aplicación está activo y respondiendo solicitudes HTTP. No depende de MongoDB.

* **Respuesta Exitosa (HTTP 200):**
```json
{
  "status": "ok",
  "service": "inverfit-api"
}
```

---

### `GET /health/ready`
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

### Formato Uniforme de Respuestas de Error
Toda respuesta de error HTTP (400, 404, 405, 500, 503) respeta la siguiente estructura estándar:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Ruta no encontrada"
  }
}
```

---

## 9. Ejecución de Pruebas y Análisis Estático

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

## 10. Consideraciones de Seguridad y Buenas Prácticas

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
