# Plataforma de Gestión de Gimnasio — INVERFIT (Panel del Entrenador)

Prototipo frontend completamente funcional y reactivo desarrollado para la gestión operativa y financiera del gimnasio **INVERFIT**, orientado al entrenador (Pablo). Realiza un port fiel de los wireframes oficiales de Figma (resolución de referencia 1366 × 850) con componentes reutilizables, formularios interactivos, validaciones, operaciones CRUD completas y persistencia local de datos.

---

## 🚀 Tecnologías Utilizadas

- **Vue 3** (Composition API, `<script setup lang="ts">`)
- **TypeScript** (tipado estricto en modelos, DTOs y servicios)
- **Vite** (bundler de alto rendimiento)
- **Vue Router 4** (navegación SPA con historial HTML5)
- **Pinia** (estado global reactivo sincronizado)
- **Tailwind CSS v4** (estilos modulares basados en los wireframes de Figma y la identidad corporativa INVERFIT)
- **Lucide Icons** (iconografía limpia y coherente)

---

## 🎨 Identidad Visual y Branding — INVERFIT

La aplicación integra la imagen corporativa oficial de **INVERFIT**, manteniendo una estética predominantemente clara, profesional y ergonómica para el uso prolongado del entrenador en su panel administrativo:

| Color | Código | Aplicación en la Plataforma |
|---|---|---|
| **Negro profundo** | `#0B0B0D` | Elementos de marca, textos de alto contraste en botones dorados |
| **Gris grafito** | `#191A1D` | Barra lateral (`Sidebar`), avatares, bordes estructurales (`#282A30`) |
| **Dorado** | `#CEA549` | Acciones principales (botones primarios, enlaces activos, badges, chevrones) |
| **Blanco cálido** | `#F5F3EE` | Superficie general y fondo de la aplicación |
| **Blanco** | `#FFFFFF` | Tarjetas de contenido, tablas, modales y formularios |

### Logotipo y Activos de Marca
- **Logotipo vectorizado:** Se implementó el logotipo oficial compuesto por la tipografía inclinada en negrita y el isotipo de dobles chevrones geométricos dorados (`INVERFIT❯❯`).
- **Componente reactivo:** `src/components/brand/BrandLogo.vue` configurable por variante (`dark`, `light`, `gold`), tamaño y modo compacto.
- **Activos en `public/`:** `logo-inverfit.svg`, `logo-inverfit-dark.svg`, `logo-inverfit-symbol.svg` y `favicon.svg` vinculado al navegador.

---

## 📱 Pantallas Implementadas (Fieles a Figma)

1. **Pantalla 1 — Inicio (`/inicio`):**
   - Resumen general de actividad con 4 tarjetas de indicadores en tiempo real: *Alumnos activos*, *Pagos pendientes*, *Sesiones de hoy* y *Contratos pendientes*.
   - Agenda de hoy con bloques programados e hipervínculo funcional para ver los alumnos asignados.
   - Panel de pendientes de atención con navegación directa a vencimientos, contratos y evaluaciones.
   - Accesos rápidos y botón "+ Registrar alumno".

2. **Pantalla 2 — Agenda Semanal (`/agenda`):**
   - **Disponibilidad flexible y editable:** El entrenador configura de forma autónoma sus intervalos de atención para cada día de la semana (0, 1 o múltiples turnos diarios, ej. 08:30–12:00 y 16:00–21:00; días de descanso desactivados).
   - **Diseño con panel contextual:** Vista semanal organizada de lunes a domingo con tarjetas de sesiones reales programadas y estado visual de disponibilidad, complementada por un panel lateral derecho para el día seleccionado con métricas de horas activas, listado de intervalos y accesos rápidos.
   - **Programación ágil de sesiones:**
     - Duración sugerida automática de 75 minutos (`término = inicio + 75m`), totalmente editable tanto en inicio como en fin sin selectores rígidos.
     - Cálculo en tiempo real de la duración efectiva.
     - **Cero descansos obligatorios:** Se permiten sesiones consecutivas o continuas sin restricciones artificiales (el descanso lo gestiona el coach a su criterio).
     - Validación preventiva: Detección estricta de solapamientos reales y restricción de sesiones fuera de los intervalos habilitados.
   - **Excepciones de calendario y recurrencia:**
     - Las sesiones son recurrentes por defecto cada semana.
     - Capacidad de reagendar o cancelar puntualmente una sesión para una fecha específica (`YYYY-MM-DD`) mediante excepciones, sin alterar la plantilla recurrente general.
   - **Modal de detalle y gestión de cupos:** Consulta de alumnos en bloque, retiro individual, adición de nuevos alumnos y control de excepciones.

3. **Pantalla 3 — Alumnos (`/alumnos`):**
   - Listado maestro con búsqueda reactiva por nombre.
   - Filtros por estado (*Activo*, *Por revisar*, *Inactivo*) y planes por frecuencia (*1 vez/sem.*, *2 veces/sem.*, *3 veces/sem.*, *4 veces/sem.*).
   - Acceso directo a la ficha del alumno ("Ver →").
   - Creación y eliminación de alumnos de demostración con confirmación previa.

4. **Pantalla 4 — Ficha del Alumno (`/alumnos/:id`):**
   - Resumen del alumno con avatar, plan, estado y acceso a emitir enlace de pago.
   - Navegación interna por pestañas:
     - **Resumen:** 4 tarjetas divididas en datos y plan, contrato y consentimiento, últimos pagos y evaluaciones.
     - **Datos personales:** edición completa de nombre, frecuencia, horarios, contacto de emergencia y salud.
     - **Pagos:** historial individual con registro de nuevos pagos.
     - **Evaluaciones:** historial de evaluaciones y métricas físicas del alumno.
     - **Contrato:** previsualización del documento y toggle de estado de aceptación (demostración).

5. **Pantalla 5 — Membresías y Pagos (`/membresias-pagos`):**
   - Control financiero con KPIs de membresías por vencer, pagos pendientes y pagos registrados.
   - Filtro por estado y buscador por alumno.
   - Botón "Exportar" que genera y descarga automáticamente un archivo CSV con el historial de pagos.
   - Simulación de operaciones:
     - **Enviar enlace de pago:** genera URL de checkout simulada con botón de copiado y texto para WhatsApp.
     - **Enviar comprobante:** simula la emisión y envío por correo de boleta/voucher centralizado (evitando el uso manual de Gmail/SII).
     - **Enviar recordatorios:** acción en lote o individual con registro de fecha de notificación.

6. **Pantalla 6 — Evaluaciones (`/evaluaciones`):**
   - Indicadores de evaluaciones registradas, alumnos por evaluar y última fecha de actualización.
   - Filtros por tipo de evaluación (*Evaluación inicial* y *Seguimiento*) y período.
   - Creación de evaluaciones con **campos dinámicos y métricas configurables** (peso, estatura, movilidad, fuerza, etc.), adaptables a los acuerdos futuros con el cliente.
   - Asociación inmediata a la ficha del alumno.

7. **Pantalla 7 — Contratos y Consentimientos (`/contratos`):**
   - Panel de control de documentos vigentes (*Términos y condiciones*, *Consentimiento informado*, *Contrato de servicios*).
   - KPIs de contratos aceptados, pendientes y versión vigente.
   - Visor modal del texto contractual con opción de simular la aceptación del alumno.
   - Acción "Gestionar términos" para editar la versión vigente y actualizar las cláusulas del gimnasio.

---

## 💾 Persistencia y Capa de Datos Simulada

La aplicación utiliza una capa desacoplada de servicios (`src/services/`):

- `students.ts`: CRUD de alumnos y sincronización de referencias.
- `payments.ts`: registro de pagos, simulación de enlaces/recordatorios y exportación CSV.
- `schedule.ts`: disponibilidad personalizada por intervalos, validación de sesiones libres (75 min sugerido), detección de solapamiento y gestión de excepciones de calendario.
- `evaluations.ts`: registro y métricas de seguimiento físico.
- `contracts.ts`: control de versiones de términos y estados de consentimiento.
- `storage.ts`: abstracción de `localStorage` con emisor de eventos reactivos.

### Sincronización entre módulos
Todas las operaciones repercuten de forma inmediata en las 7 pantallas:
- Al crear o editar un alumno, se actualiza el listado, el selector de la agenda, sus registros de pago, contrato y contadores del dashboard.
- Al registrar un pago o evaluación, se actualiza el historial en la ficha individual y en las tablas globales.
- Al recargar la página, todos los datos persisten (tanto en modo MongoDB como en modo LocalStorage).
- Sesión autenticada persistente para el entrenador Pablo con botón de cierre de sesión seguro en la barra lateral.

---

## 🛠️ Instalación y Puesta en Marcha Rápida (Para el Equipo)

### Prerrequisitos
- **Node.js** v18+ (recomendado v20+) y **npm** v9+
- **Docker & Docker Compose** (para ejecutar la API en Go y MongoDB en contenedores)
- *Opcional:* Go 1.22+ si se desea ejecutar el backend de manera nativa sin Docker

---

### Pasos para Inicializar el Proyecto

#### 1. Clonar el repositorio
```bash
git clone https://github.com/Midgardsormrr/Prototipo-Inverfit.git
cd Prototipo-Inverfit
```

#### 2. Configurar las Variables de Entorno del Frontend
Copie la plantilla de configuración `.env.example` como `.env`:
```bash
# En Windows (PowerShell):
copy .env.example .env

# En Linux / macOS / Bash:
cp .env.example .env
```

El archivo `.env` viene preconfigurado para conectarse a la API Dockerizada (`http://localhost:8081/api/v1`):
```env
VITE_STUDENTS_DATA_SOURCE=api
VITE_SCHEDULE_DATA_SOURCE=api
VITE_PAYMENTS_DATA_SOURCE=api
VITE_API_BASE_URL=http://localhost:8081/api/v1
```
*(Nota: Si desea ejecutar el frontend 100% desconectado o en modo demo offline, simplemente cambie los valores a `mock`).*

#### 3. Levantar la Base de Datos y Backend (Docker Stack)
El proyecto incluye un entorno Docker multi-contenedor listo para producción/desarrollo (`compose.stack.yaml`):

```bash
# 1. Crear el archivo de variables del stack desde la plantilla
cp .env.stack.example .env.stack   # (o en PowerShell: copy .env.stack.example .env.stack)

# 2. Levantar MongoDB y la API REST en Go
docker compose -f compose.stack.yaml --env-file .env.stack up -d --build

# 3. Verificar que los contenedores estén saludables (healthy)
docker compose -f compose.stack.yaml --env-file .env.stack ps
```
- La API de Go quedará disponible en: `http://localhost:8081/api/v1`
- Verificación de salud: `http://localhost:8081/health/ready`

#### 4. Instalar Dependencias del Frontend y Correr en Desarrollo
```bash
# Instalar paquetes npm
npm install

# Iniciar servidor Vite
npm run dev
```
La aplicación web quedará disponible en: **`http://localhost:5173/`**

---

## 🔐 Credenciales de Acceso (Coach Pablo)

Al ingresar a `http://localhost:5173/login`, el sistema presenta una pantalla de inicio de sesión con branding oficial INVERFIT:

| Campo | Valor |
|---|---|
| **Correo electrónico** | `pablo@inverfit.cl` |
| **Contraseña** | `inverfit2026` |
| **Acceso Rápido** | Se incluye un botón de 1 clic *"Pablo (Coach & Admin)"* para autocompletar credenciales |

### 📝 Portal Público de Onboarding y Firma Digital
Cualquier alumno nuevo puede completar su ficha médica y firmar su consentimiento informado en pantalla táctil/mouse accediendo a la ruta pública:
- **`http://localhost:5173/registro-alumno`** (No requiere login de entrenador).

---

## 📦 Compilación para Producción

Para compilar el proyecto como una aplicación web estática optimizada para hosting (Cloudflare Pages, Netlify, Vercel, S3):

```bash
npm run build
```

Para previsualizar localmente la compilación de producción:
```bash
npm run preview
```

Los artefactos listos para despliegue se generarán en el directorio `dist/`.

---

## 🔌 Arquitectura del Backend Go + MongoDB

La plataforma cuenta con un backend modular implementado en **Go + Gin** conectado a **MongoDB 7.0**:
- **Módulos:** Alumnos (`/api/v1/students`), Agenda y Disponibilidad (`/api/v1/schedule`), Pagos y Finanzas (`/api/v1/payments`).
- **Control de concurrencia y cupos:** Validación estricta del límite oficial de **3 alumnos por bloque de 75 minutos**.
- Para más detalles sobre la arquitectura interna, endpoints y modelos, consulte [`backend/README.md`](backend/README.md).

---

## ⚖️ Versión y Respaldo Operativo
- **Versión de la plataforma:** `INVERFIT ❯❯ v1.1`
- **Ciclo de Cobranza:** Oficial del 1 al 5 de cada mes mediante pasarela e integraciones SumUp.
- **Consentimientos:** Respaldo digital de firmas conforme al reglamento interno oficial de 4 cláusulas.
