# DOCUMENTACIÓN FUNCIONAL Y TÉCNICA DEL PROYECTO — INVERFIT

---

**Fecha de elaboración:** 4 de noviembre de 2026  
**Versión documental:** v1.1.0 (Post-Reestructuración de Agenda y Normalización de Identidad)  
**Estado de la aplicación:** Prototipo frontend interactivo 100% funcional en cliente con persistencia local en `localStorage` (sin backend ni base de datos real).  
**Propósito y vigencia:** Este documento constituye la especificación funcional, técnica y operativa exhaustiva del estado actual de la plataforma INVERFIT. Debe actualizarse de forma obligatoria cada vez que se incorporen o modifiquen funcionalidades, reglas de negocio o integraciones en el proyecto.

---

## 1. Información General

### 1.1. Nombre del Proyecto
**INVERFIT** (Centro de entrenamiento y acondicionamiento físico enfocado en fuerza, estética y rendimiento en deportes de *endurance*).

### 1.2. Objetivo de la Plataforma
El objetivo principal de la aplicación es dotar al entrenador y fundador (**Pablo**) de una herramienta web centralizada, ágil y visualmente profesional para administrar integralmente la operación del centro de entrenamiento:
1. **Control de alumnos:** Matrículas, fichas individuales, planes, datos de salud y contactos de emergencia.
2. **Agenda semanal dinámica:** Administración flexible de disponibilidad horaria por intervalos, programación de sesiones individuales y grupales (duración sugerida de 75 minutos, 100% editable), sin imposición de descansos forzados y con soporte para excepciones de calendario.
3. **Membresías y pagos:** Seguimiento de vencimientos, registro de abonos, exportación de balances a CSV y simulación de emisión de comprobantes y enlaces de cobro digital.
4. **Evaluaciones físicas:** Registro dinámico de evaluaciones iniciales y periódicas con métricas cuantitativas y cualitativas personalizables.
5. **Contratos y consentimientos informados:** Control de aceptación de términos del servicio, resguardo legal preventivo y versionado de documentos contractuales.

### 1.3. Tecnologías Utilizadas
La solución está construida sobre un stack moderno de tecnologías web frontend:

| Tecnología | Versión | Rol en el Proyecto |
|---|---|---|
| **Vue.js** | `^3.5.44` | Framework principal reactivo utilizando la **Composition API** y sintaxis `<script setup lang="ts">`. |
| **TypeScript** | `~6.0.3` | Tipado estático estricto para entidades de dominio, DTOs, interfaces de servicio y props de componentes. |
| **Vite** | `^8.3.5` | Servidor de desarrollo de alta velocidad y empaquetador optimizado para producción. |
| **Pinia** | `^4.1.0` | Store centralizado para reactividad global, gestión de estado entre vistas y emisión de notificaciones Toast. |
| **Vue Router** | `^4.6.4` | Enrutamiento SPA con historial HTML5 (`createWebHistory`) y títulos dinámicos de página. |
| **Tailwind CSS** | `^4.3.4` | Framework de diseño utilitario integrado mediante `@tailwindcss/vite` para estilos coherentes y responsivos. |
| **Lucide Vue Next** | `^1.0.0` | Colección de iconografía SVG ligera y consistente. |
| **vue-tsc** | `^3.3.13` | Comprobador de tipos TypeScript para plantillas Vue durante el proceso de build. |

### 1.4. Arquitectura General del Frontend
La aplicación sigue el patrón de diseño arquitectónico desacoplado en capas:

```
┌────────────────────────────────────────────────────────┐
│                   Vistas (src/views/)                  │
│       7 pantallas completas fieles a Figma             │
└───────────┬────────────────────────────────┬───────────┘
            │                                │
┌───────────▼────────────────────────────────▼───────────┐
│              Componentes (src/components/)             │
│   Modales de acción, Layout (TopNav, Sidebar), UI      │
└───────────┬────────────────────────────────┬───────────┘
            │                                │
┌───────────▼────────────────────────────────▼───────────┐
│               Pinia Store (src/stores/gym.ts)          │
│   Estado global reactivo, contadores y notificaciones  │
└───────────┬────────────────────────────────┬───────────┘
            │                                │
┌───────────▼────────────────────────────────▼───────────┐
│             Servicios de Negocio (src/services/)       │
│  students.ts | schedule.ts | payments.ts | etc.       │
└───────────┬────────────────────────────────┬───────────┘
            │                                │
┌───────────▼────────────────────────────────▼───────────┐
│           Almacenamiento Local (storage.ts)            │
│  Abstracción de localStorage con emisor reactivo       │
└────────────────────────────────────────────────────────┘
```

### 1.5. Estado Actual del Desarrollo y Alcance de la Versión
- **Alcance funcional:** Es un prototipo interactivo de alta fidelidad, completamente operativo en el navegador del cliente. No se trata de simples maquetas estáticas: permite crear alumnos, modificar agendas, calcular duraciones, programar excepciones, alternar estados de pago, descargar archivos CSV reales y alterar versiones legales.
- **Ausencia de backend real:** No existe un servidor backend en Go, una API REST externa ni una base de datos PostgreSQL/MySQL en esta fase. Toda la persistencia opera de forma síncrona en el `localStorage` del navegador donde se ejecute la sesión.
- **Datos iniciales:** Al iniciar por primera vez o al hacer clic en *"Reiniciar demo"*, la aplicación se precarga con un conjunto coherente de datos ficticios basados en los wireframes de Figma.

---

## 2. Arquitectura del Proyecto y Estructura de Directorios

El código fuente se organiza bajo `src/` respetando una clara separación de responsabilidades:

```
Prototipo-Inverfit/
├── docs/
│   └── FUNCIONALIDADES_IMPLEMENTADAS.md    # Documento de especificación vigente
├── public/                                 # Activos públicos estáticos (logotipos SVG, favicon)
├── src/
│   ├── assets/                             # Recursos gráficos compilados
│   ├── components/
│   │   ├── brand/                          # Logotipo vectorial oficial de INVERFIT
│   │   ├── layout/                         # Barra superior (TopNavbar) y lateral (Sidebar)
│   │   ├── modals/                         # Ventanas modales de creación, edición y detalle
│   │   └── ui/                             # Notificaciones flotantes (ToastContainer)
│   ├── router/                             # Configuración de Vue Router y guardias de navegación
│   ├── services/                           # Lógica de dominio desacoplada y validaciones
│   ├── stores/                             # Store global de Pinia (gym.ts)
│   ├── types/                              # Definiciones de tipos e interfaces TypeScript
│   ├── views/                              # Las 7 vistas principales del sistema
│   ├── App.vue                             # Componente raíz con layout y contenedor Toast
│   ├── main.ts                             # Entrada principal de la aplicación
│   └── style.css                           # Configuración de Tailwind CSS v4 y variables CSS
├── package.json                            # Dependencias y scripts de ejecución
├── tsconfig.json                           # Configuración general de TypeScript
└── vite.config.ts                          # Configuración de Vite y plugins
```

### 2.1. Responsabilidad de los Directorios

#### `src/views/`
Contiene los componentes de vista que corresponden a las rutas de primer nivel:
- `InicioView.vue`: Dashboard general con métricas en tiempo real, agenda del día y alertas de atención.
- `AgendaView.vue`: Vista semanal de sesiones y disponibilidad con panel lateral contextual para el día seleccionado.
- `AlumnosView.vue`: Maestro de alumnos con búsqueda, filtros y CRUD.
- `FichaAlumnoView.vue`: Perfil detallado del alumno organizado en 5 pestañas temáticas.
- `MembresiasPagosView.vue`: Control financiero, vencimientos, cobros y exportación CSV.
- `EvaluacionesView.vue`: Historial y registro de valoraciones físicas dinámicas.
- `ContratosView.vue`: Control de documentos, estado de firmas y gestión de términos.

#### `src/components/`
- **`brand/`:** `BrandLogo.vue` centraliza la representación vectorial oficial del logotipo con opciones de variantes (`dark`, `light`, `gold`), tamaño y modo compacto.
- **`layout/`:** `Sidebar.vue` (navegación vertical con selector activo, perfil del coach Pablo y acceso a reinicio demo) y `TopNavbar.vue` (barra superior fija con búsqueda global, acceso rápido a ayuda y perfil).
- **`modals/`:** Modales independientes para cada operación de negocio:
  - `AssignScheduleModal.vue`: Programación de sesiones con sugerencia de 75 min y cálculo en vivo.
  - `BlockDetailsModal.vue`: Consulta de participantes y gestión de excepciones (reagendar/cancelar solo una fecha).
  - `ConfigureAvailabilityModal.vue`: Configuración de intervalos de atención semanal del entrenador.
  - `RegisterStudentModal.vue`: Formulario de alta rápida de nuevos alumnos.
  - `RegisterPaymentModal.vue`: Registro manual de abonos y actualización de estados financieros.
  - `CreateEvaluationModal.vue`: Formulario con métricas dinámicas para seguimiento físico.
  - `ManageTermsModal.vue`: Editor de cláusulas y versionado del contrato marco.
  - `SimulatedActionModal.vue`: Previsualización interactiva de enlaces de pago, WhatsApp, correos y comprobantes.
- **`ui/`:** `ToastContainer.vue` muestra mensajes flotantes no invasivos de éxito, información, advertencia y error.

#### `src/services/`
Aisla la lógica pura de negocio y no contiene dependencias de interfaz ni de Vue:
- `students.ts`: Operaciones de alumnos y sincronización de referencias cruzadas.
- `schedule.ts`: Algoritmos de conversión de tiempo, validación de traslapes, inserción en disponibilidad, gestión de excepciones y cálculo de ocurrencias semanales.
- `payments.ts`: Registro de transacciones, generación de enlaces simulados y exportación CSV.
- `evaluations.ts`: Creación y consulta de métricas de rendimiento y salud física.
- `contracts.ts`: Estado de aceptaciones y control de versiones de cláusulas legales.
- `storage.ts`: Capa de persistencia en `localStorage`, inicialización de datos, migraciones automáticas y patrón Observer/PubSub reactivo.
- `mockData.ts`: Colección de datos semilla iniciales fieles al prototipo de Figma.

#### `src/stores/`
- `gym.ts`: Implementa el store principal de Pinia utilizando la Composition API (`defineStore('gym', () => { ... })`). Conecta los servicios con la interfaz gráfica, mantiene el estado reactivo, computa contadores globales y dispara toasts informativos ante cada acción.

#### `src/types/`
Tipos TypeScript fuertemente tipados divididos por módulo (`student.ts`, `schedule.ts`, `payment.ts`, `evaluation.ts`, `contract.ts`, `index.ts`).

#### `src/router/`
Define las 7 rutas protegidas de la aplicación, redirecciones automáticas y actualización reactiva de la etiqueta `<title>` del navegador.

---

## 3. Funcionalidades Detalladas por Módulo

### 3.1. Pantalla 1 — Inicio (`/inicio`)

El panel de inicio consolida la actividad diaria del gimnasio en un solo vistazo:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  [ Alumnos activos: 6 ]  [ Pagos pendientes: 2 ]  [ Sesiones hoy: 4 ]  ...   │
├────────────────────────────────────────┬─────────────────────────────────────┤
│  AGENDA DE HOY                         │  PENDIENTES DE ATENCIÓN             │
│  - 08:30 Bloque de entrenamiento (2 al)│  - Membresías por vencer (1 pend.)  │
│  - 10:00 Bloque de entrenamiento (1 al)│  - Contratos sin aceptar (2 pend.)  │
│  - 17:30 Bloque de entrenamiento (2 al)│  - Evaluaciones por reg. (2 pend.)  │
│  ... [Ver completa →]                  │                                     │
├────────────────────────────────────────┴─────────────────────────────────────┤
│  ACCESOS RÁPIDOS: [Ver agenda]  [Buscar alumno]  [Gestionar pagos]           │
└──────────────────────────────────────────────────────────────────────────────┘
```

#### Indicadores Disponibles (KPIs)
Calculados dinámicamente mediante propiedades computadas de Pinia:
1. **Alumnos activos:** Total de alumnos cuyo estado de membresía es `Activo`.
2. **Pagos pendientes:** Cantidad de pagos en estado `Pendiente`.
3. **Sesiones de hoy:** Número de bloques asignados para el día actual de la semana (Lunes a Domingo). Si el día actual no tiene sesiones cargadas, muestra los bloques de Lunes como vista previa informativa.
4. **Contratos pendientes:** Total de alumnos que aún no han aceptado los términos del servicio.

#### Componentes del Inicio
- **Agenda de hoy:** Lista los bloques programados para la jornada, indicando hora de inicio, nombre del bloque y cantidad de alumnos participantes. Al hacer clic en un bloque, se abre directamente el modal de detalle del bloque.
- **Pendientes de atención:** Tres tarjetas de alerta reactivas con contadores automáticos:
  - *Membresías por vencer:* Lleva directamente a `/membresias-pagos`.
  - *Contratos sin aceptar:* Lleva directamente a `/contratos`.
  - *Evaluaciones por registrar:* Lleva directamente a `/evaluaciones`.
- **Accesos rápidos:** Botones directos a las áreas principales y botón superior destacado `+ Registrar alumno`.

---

### 3.2. Pantalla 2 — Agenda Semanal (`/agenda`)

El módulo de agenda semanal representa la herramienta de trabajo diario más crítica del entrenador. Ha sido completamente rediseñado para garantizar una **flexibilidad total y un control absoluto por parte de Pablo**.

```
┌───────────────────────────────────────────────────────────────────┬──────────────────────────────┐
│  Lunes 3    Martes 4    Miércoles 5   Jueves 6   Viernes 7  ...       │ PANEL DEL DÍA (Lunes 3)      │
│  ┌────────┐ ┌────────┐  ┌────────┐   ┌────────┐ ┌────────┐        │ - 3 sesiones programadas     │
│  │ 08:30  │ │ 08:30  │  │ 08:30  │   │ 09:00  │ │ 08:30  │        │ - 4 alumnos atendidos        │
│  │ 09:45  │ │ 09:45  │  │ 09:45  │   │ 10:15  │ │ 09:45  │        │ - 3.75 hrs activas           │
│  │ (2 al.)│ │ (2 al.)│  │ (1 al.)│   │ (1 al.)│ │ (2 al.)│        ├──────────────────────────────┤
│  └────────┘ └────────┘  └────────┘   └────────┘ └────────┘        │ TURNOS CONFIGURADOS:         │
│  ┌────────┐             ┌────────┐              ┌────────┐        │ • 08:30 — 12:00              │
│  │ 10:00  │             │ 10:00  │              │ 10:00  │        │ • 16:00 — 21:00              │
│  │ 11:15  │             │ 11:15  │              │ 11:15  │        ├──────────────────────────────┤
│  └────────┘             └────────┘              └────────┘        │ [+ Nueva sesión]             │
│  [+ Espacio libre: 16:00 a 21:00]                                 │ [Configurar disponibilidad]  │
└───────────────────────────────────────────────────────────────────┴──────────────────────────────┘
```

#### Reglas de Negocio Definitivas de la Agenda
1. **Disponibilidad por intervalos configurables:**
   - La disponibilidad y las sesiones son entidades independientes. El entrenador define sus franjas horarias de apertura para cada día de la semana (Lunes a Domingo).
   - Soporta **0, 1 o múltiples intervalos por jornada** (ej. Mañana de `08:30 a 12:00` y Tarde de `16:00 a 21:00`).
   - Si un día no se trabaja (ej. Domingo), el switch del día se apaga (`enabled: false`) y el día se clasifica visualmente como *"Sin disponibilidad configurada"*.
2. **Duración de sesiones libre y sugerencia automática de 75 minutos:**
   - **No existe un selector forzado de duración.**
   - Al tipear o seleccionar una hora de inicio (ej. `08:00`), el sistema propone automáticamente como hora de término `09:15` (+75 minutos).
   - Ambos campos son editables libremente por el entrenador (puede modificar el término a `09:00`, `09:30`, `10:00`, etc.).
   - La duración efectiva en minutos se calcula dinámicamente: $\text{duración} = \text{hora\_fin} - \text{hora\_inicio}$.
   - Solo se exige una duración mínima técnica de 15 minutos y que la hora de término sea posterior a la de inicio.
3. **Cero descansos obligatorios:**
   - **No existen descansos automáticos de 15 minutos ni advertencias por bloques contiguos.**
   - El entrenador administra su fatiga y descansos según su propio criterio.
   - Dos sesiones consecutivas son 100% válidas (ej. Bloque A de `08:30 a 09:30` y Bloque B de `09:30 a 10:45`).
4. **Detección estricta de solapamiento real:**
   - Una nueva sesión solo se rechaza si existe un cruce temporal efectivo con otra sesión existente:
     $$\text{inicio}_A < \text{fin}_B \quad \text{y} \quad \text{fin}_A > \text{inicio}_B$$
5. **Restricción dentro de la disponibilidad:**
   - Toda sesión debe quedar **completamente contenida** dentro de al menos uno de los intervalos configurados para ese día. No se permite programar sesiones fuera del horario laboral establecido.
6. **Protección preventiva ante reducción de turnos:**
   - Si el entrenador intenta recortar o eliminar un intervalo de disponibilidad en un día que ya tiene sesiones agendadas que quedarían fuera del nuevo rango, el sistema **bloquea el guardado** y muestra la lista de sesiones en conflicto para que sean reubicadas previamente.
7. **Sesiones grupales y asignación de alumnos:**
   - Cada bloque horario admite de 1 a 4 alumnos simultáneos.
   - Si tiene 1 alumno se clasifica como `Individual`; con 2 o más se etiqueta como `Sesión grupal`.
   - Asignar o retirar alumnos de un bloque actualiza en tiempo real el campo de resumen `fixedSchedule` en la ficha individual de cada alumno.
8. **Recurrencia semanal y excepciones de calendario (`ScheduleException`):**
   - Las sesiones son recurrentes por defecto cada semana.
   - **Excepción puntual por fecha:** Desde el modal de detalle, el entrenador puede:
     - **Reagendar solo esta fecha:** Cambia el horario o día de una sesión únicamente para una fecha determinada (`YYYY-MM-DD`). La sesión base en las semanas anteriores y posteriores no se ve afectada.
     - **Cancelar solo esta fecha:** Oculta la sesión para ese día específico sin eliminar la recurrencia habitual.
     - **Restaurar:** En cualquier momento puede revertir la excepción y recuperar el horario recurrente base.
9. **No autoagendamiento:**
   - Los alumnos no tienen acceso a la plataforma para agendar, cancelar o modificar sesiones. El agendamiento es 100% manual y supervisado por el entrenador.

---

### 3.3. Pantalla 3 — Alumnos (`/alumnos`)

Módulo de administración maestro del padrón de atletas de INVERFIT:
- **Buscador reactivo:** Campo de búsqueda por nombre del alumno en tiempo real.
- **Filtros cruzados:**
  - *Por estado:* `Todos los estados`, `Activo`, `Por revisar`, `Inactivo`.
  - *Por frecuencia semanal:* `Todos los planes`, `1 vez/sem.`, `2 veces/sem.`, `3 veces/sem.`, `4 veces/sem.`
- **Tabla de alumnos:** Renderiza avatar con inicial y colores de marca, nombre completo, horarios fijos asignados, tipo de plan, frecuencia, fecha del próximo vencimiento, badge de estado y acciones.
- **Creación de alumnos (`RegisterStudentModal`):** Alta con validación de campos obligatorios, notas, contacto de emergencia y afecciones de salud.
  - *Sincronización automática:* Al crear un alumno se genera automáticamente su primer registro de pago pendiente, su primera evaluación física pendiente y su contrato pendiente de aceptación.
- **Eliminación con confirmación:** Acción con ventana de confirmación previa que limpia en cascada la presencia del alumno en la agenda, pagos, contratos y evaluaciones.
- **Navegación:** Enlace directo `"Ver →"` para acceder a la ficha individual.

---

### 3.4. Pantalla 4 — Ficha Individual del Alumno (`/alumnos/:id`)

Proporciona la vista consolidada de 360 grados de un alumno particular. Se divide en una cabecera con avatar, badges y accesos de cobro, y 5 pestañas de contenido:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  [P] Juan Pérez   Plan 3 sesiones   [Activo]   Horario: Lun, Mié 18:00       │
│  [Enviar enlace de pago]  [Editar información]                              │
├──────────────┬──────────────┬──────────────┬──────────────┬──────────────────┤
│ 1. Resumen   │ 2. Datos     │ 3. Pagos     │ 4. Evaluac.  │ 5. Contrato      │
├──────────────┴──────────────┴──────────────┴──────────────┴──────────────────┤
│                                                                              │
│  [Pestaña 1: Resumen]                                                        │
│  ┌───────────────────────┐ ┌───────────────────────┐                         │
│  │ Datos y plan          │ │ Contrato y consent.   │                         │
│  │ Membresía activa      │ │ Estado: Aceptado      │                         │
│  └───────────────────────┘ └───────────────────────┘                         │
│  ┌───────────────────────┐ ┌───────────────────────┐                         │
│  │ Últimos pagos         │ │ Evaluaciones          │                         │
│  │ Al día (BOL-2026-...) │ │ Inicial registrada    │                         │
│  └───────────────────────┘ └───────────────────────┘                         │
└──────────────────────────────────────────────────────────────────────────────┘
```

#### Detalle de las Pestañas
1. **Resumen:** Cuatro tarjetas clave que resumen la situación general (Membresía, Documento legal vigente, Último pago registrado y Última evaluación realizada).
2. **Datos personales:** Formulario con modo de edición en línea. Permite modificar nombre, teléfono, email, plan, frecuencia, estado de membresía, fecha de inicio, fecha de vencimiento, contacto de emergencia, antecedentes médicos/lesiones y notas del entrenador.
3. **Pagos:** Historial exclusivo de los pagos asociados a este alumno. Permite registrar nuevos pagos o emitir enlaces simulados de cobro específicos.
4. **Evaluaciones:** Historial cronológico de pruebas físicas del alumno, visualización detallada de métricas y acceso a crear nuevas evaluaciones vinculadas directamente a su ID.
5. **Contrato:** Permite auditar el estado de aceptación de los términos del servicio (`Aceptado` o `Pendiente`), la versión de términos aplicable, la fecha de consentimiento y un interruptor para simular la aceptación o revocación del contrato.

---

### 3.5. Pantalla 5 — Membresías y Pagos (`/membresias-pagos`)

Módulo de administración y trazabilidad de ingresos:
- **KPIs Financieros:**
  - *Membresías por vencer:* Alumnos con fecha límite próxima o estado `Por revisar`.
  - *Pagos pendientes:* Total de cobros mensuales pendientes de liquidación.
  - *Pagos registrados:* Transacciones confirmadas y archivadas.
- **Tabla de pagos:** Muestra alumno, plan, fecha de vencimiento pactada, monto en pesos chilenos ($), estado (`Pagado` / `Pendiente`), medio de pago registrado, fecha de abono y estado de envío de comprobante.
- **Registro manual de pago (`RegisterPaymentModal`):**
  - Permite ingresar monto, método de pago (*Transferencia bancaria*, *Enlace de pago*, *Efectivo / Presencial*), fecha de pago y glosa.
  - Si el alumno tenía estado *"Por revisar"*, al registrar su pago el sistema actualiza automáticamente su estado a *"Activo"*.
- **Operaciones de simulación externa:**
  - *Enviar enlace de pago:* Genera una URL de checkout ficticia (`https://demo.inverfit.invalid/checkout/{token}`) con botón de copia en un clic y plantilla formateada para WhatsApp.
  - *Enviar recordatorios:* Notificación en lote o individual. Registra la fecha y hora exacta del envío e incrementa el contador `reminderSentCount`.
  - *Enviar comprobante:* Simula la emisión de una boleta electrónica o comprobante PDF para ser enviado por correo al alumno.
- **Exportación de datos a CSV:** Botón *"Exportar"* que genera y descarga en el acto un archivo `.csv` real con todo el historial de cobros procesado en memoria mediante objetos `Blob`.

---

### 3.6. Pantalla 6 — Evaluaciones Físicas (`/evaluaciones`)

Espacio dedicado al seguimiento del progreso atlético y composición corporal de los alumnos:
- **KPIs de Seguimiento:** Evaluaciones totales registradas, alumnos con evaluación pendiente y fecha de última actualización.
- **Filtros:** Buscador por nombre de alumno y filtro por tipo (*Evaluación inicial* y *Seguimiento*).
- **Métricas dinámicas y configurables (`CreateEvaluationModal`):**
  - No limita al entrenador a un esquema rígido de campos.
  - Incluye métricas sugeridas base: *Peso corporal (kg)*, *Estatura (cm)*, *Movilidad articular*, *Sentadilla / Fuerza base*.
  - El entrenador puede añadir campos personalizados adicionales al vuelo (definiendo nombre de la métrica, valor cuantitativo/cualitativo y unidad de medida opcional) o eliminar los que no requiera.
  - Registro de observaciones técnicas generales y fecha de aplicación.
- **Trazabilidad:** Cada registro queda vinculado de inmediato a la ficha del alumno en `/alumnos/:id`.

---

### 3.7. Pantalla 7 — Contratos y Consentimientos (`/contratos`)

Módulo de resguardo y control legal de la relación entre el centro y los deportistas:
- **KPIs Documentales:** Solicitudes pendientes de aceptación, contratos aceptados y versión del marco legal vigente (ej. `v1.0 (2026)`).
- **Tabla de consentimientos:** Listado por alumno con tipo de documento (*Términos y condiciones*, *Consentimiento informado*), estado, versión asociada y fecha de firma.
- **Visor contractual modal:** Permite consultar el texto íntegro de las cláusulas operativas de INVERFIT (naturaleza del servicio, duración de sesiones, política de ausencias y reagendamientos, condición de salud, términos de pago y cláusula de exoneración).
- **Simulación de aceptación:** Botón para conmutar el estado del documento entre `Pendiente` y `Aceptado` con marca temporal automática.
- **Gestión de términos (`ManageTermsModal`):** Permite al administrador crear una nueva versión de los términos (ej. `v1.1 (2026)`), editar su texto legal y marcarla como la versión vigente aplicable a futuros registros.

---

## 4. Reglas de Negocio Implementadas

Para facilitar la futura auditoría y el traspaso a un backend real, las reglas de la aplicación se dividen en cuatro categorías:

### 4.1. Reglas Confirmadas por el Cliente (INVERFIT / Pablo)
1. **Identidad unificada:** La marca oficial, única y definitiva es **INVERFIT**.
2. **Control exclusivo del entrenador:** Todo el agendamiento y coordinación es potestad del entrenador. No existe autoatención ni reservas directas por parte de los alumnos.
3. **Membresías por frecuencia fija:** Los planes se estructuran en 1, 2, 3 o 4 sesiones semanales fijas.
4. **Política de ausencias y recuperaciones:** Las sesiones no realizadas por inasistencia del alumno no se reembolsan ni aplazan el vencimiento mensual; solo pueden reagendarse excepcionalmente dentro de la misma semana y previa autorización del coach.
5. **Autonomía horaria del coach:** No existen descansos obligatorios de 15 minutos entre sesiones ni duraciones rígidas de 60 o 75 minutos impuestas por el sistema.
6. **Cupo máximo por sesión:** Capacidad operativa de hasta 4 alumnos simultáneos en bloque de entrenamiento.

### 4.2. Decisiones de Implementación Técnica
1. **Sugerencia de 75 minutos:** La duración predeterminada sugerida para una sesión es de 75 minutos (`inicio + 75m`), editable en tiempo real.
2. **Duración mínima de sesión:** 15 minutos (para evitar errores accidentales de digitación como `09:00 a 09:00`).
3. **Solapamiento real estricto:** Solo se bloquea la programación si los rangos de minutos se intersectan efectivamente. Las sesiones consecutivas (`fin_A == inicio_B`) son válidas.
4. **Contención en disponibilidad:** Toda sesión debe estar dentro de un intervalo habilitado del día.
5. **Sincronización en cascada:** Crear o eliminar un alumno propaga automáticamente las modificaciones en la agenda, pagos, contratos y evaluaciones.
6. **Tokens ficticios seguros:** Los enlaces de pago simulados utilizan el dominio reservado para pruebas `https://demo.inverfit.invalid/checkout/{TOKEN}` según RFC 2606 / RFC 6761, evitando colisiones con dominios reales de internet.

### 4.3. Valores Predeterminados Configurables
- **Horario operativo referencial del gimnasio:** 06:30 a 22:00.
- **Disponibilidad predeterminada de Pablo:**
  - Lunes a Viernes: Turno mañana (`08:30 a 12:00`) y Turno tarde (`16:00 a 21:00`).
  - Sábado: Turno único mañana (`09:00 a 13:00`).
  - Domingo: Deshabilitado / Cerrado (`enabled: false`).
- **Versión legal base:** `v1.0 (2026)` — *Términos y Condiciones del Servicio de Entrenamiento — INVERFIT*.

### 4.4. Limitaciones y Aspectos Pendientes de Definición
- **Precios formales de los planes:** En esta etapa los montos son ingresados manualmente en cada cobro o dejados en estado *"Por definir"*.
- **Contratos con valor probatorio legal:** La aceptación en el prototipo es una simulación en cliente; no implementa Firma Electrónica Avanzada (FEA) bajo la Ley 19.799 de Chile.
- **Múltiples entrenadores:** La arquitectura actual asume un único perfil de coach (**Pablo**).

---

## 5. Persistencia y Almacenamiento Local

### 5.1. Mecanismo de Almacenamiento
La persistencia de datos se gestiona a través del servicio [`src/services/storage.ts`](file:///c:/Users/dguerra/Documents/GitHub/Prototipo-Inverfit/src/services/storage.ts) utilizando la API de `localStorage` del navegador.

Para entornos donde `localStorage` no está presente (por ejemplo, durante la ejecución de pruebas unitarias en Node.js o compilación SSR), el servicio incluye un mecanismo de respaldo automático en memoria (`memoryStore`).

### 5.2. Claves de Almacenamiento (`STORAGE_KEYS`)

| Clave en LocalStorage | Tipo de Datos Persistido | Descripción |
|---|---|---|
| `inverfit_students_v1` | `Student[]` | Catálogo maestro de alumnos registrados con sus datos de contacto y salud. |
| `inverfit_payments_v1` | `Payment[]` | Histórico de pagos, vencimientos, estados y marcas de recordatorios. |
| `inverfit_schedule_v1` | `ScheduleBlock[]` | Bloques de entrenamiento recurrentes y participantes asignados. |
| `inverfit_availability_v1` | `DayAvailability[]` | Intervalos horarios configurados para cada día de la semana. |
| `inverfit_schedule_exceptions_v1` | `ScheduleException[]` | Reagendamientos y cancelaciones para fechas puntuales (`YYYY-MM-DD`). |
| `inverfit_evaluations_v1` | `Evaluation[]` | Registros de valoraciones físicas y baterías de métricas dinámicas. |
| `inverfit_contracts_v1` | `ContractConsent[]` | Registro de aceptaciones y consentimientos informados por alumno. |
| `inverfit_terms_v1` | `TermsVersion[]` | Historial de versiones del texto contractual y versión vigente. |

### 5.3. Sincronización Reactiva (Patrón Observer)
El servicio `LocalStorageService` implementa una lista interna de escuchas (`listeners: Set<() => void>`). Cada vez que se ejecuta una mutación (`set()`), se notifica a todos los suscriptores. El store de Pinia (`useGymStore`) está suscrito de forma permanente a estos cambios, lo que asegura que cualquier vista o pestaña abierta refleje las mutaciones de inmediato sin requerir recargar la página.

### 5.4. Migración Automática y Retrocompatibilidad
El método privado `getMigratedOrInit()` examina si existen datos almacenados con prefijos o denominaciones anteriores del prototipo. Si detecta claves previas con el mismo sufijo, migra su contenido de forma transparente a la nueva clave `inverfit_*` y elimina la entrada obsoleta.

### 5.5. Restauración de Datos de Demostración
La función `resetDemoData()` (accesible desde la barra lateral y el encabezado) sobrescribe todas las claves con los datos semilla iniciales definidos en `src/services/mockData.ts`, permitiendo reiniciar la aplicación a su estado de fábrica en cualquier momento.

### 5.6. Limitaciones del Almacenamiento Local
- **Aislamiento por navegador:** Los datos residen en el equipo y navegador específico del usuario; no se sincronizan entre diferentes dispositivos.
- **Límite de capacidad:** Típicamente restringido a 5 MB por origen por las políticas estándar de los navegadores web.
- **Sin concurrencia multiusuario:** No apto para edición simultánea por múltiples administradores al mismo tiempo.

---

## 6. Componentes Reutilizables

| Componente | Ruta en el Código | Responsabilidad Principal |
|---|---|---|
| **BrandLogo** | `src/components/brand/BrandLogo.vue` | Renderizado del isotipo oficial con chevrons dorados `❯❯` y texto `INVERFIT`. |
| **Sidebar** | `src/components/layout/Sidebar.vue` | Menú lateral con indicador activo dorado, perfil de Pablo y botón de reinicio de demo. |
| **TopNavbar** | `src/components/layout/TopNavbar.vue` | Barra superior con buscador global sincronizado, botón de ayuda interactiva y avatar. |
| **ToastContainer** | `src/components/ui/ToastContainer.vue` | Contenedor de notificaciones flotantes con auto-cierre y animaciones de entrada. |
| **AssignScheduleModal** | `src/components/modals/AssignScheduleModal.vue` | Asignación de sesiones con sugerencia automática de 75 min y selector de alumnos. |
| **BlockDetailsModal** | `src/components/modals/BlockDetailsModal.vue` | Visor de participantes de bloque y gestor de excepciones de calendario. |
| **ConfigureAvailabilityModal** | `src/components/modals/ConfigureAvailabilityModal.vue` | Editor de turnos por día (Lunes a Domingo) con detector de conflictos. |
| **RegisterStudentModal** | `src/components/modals/RegisterStudentModal.vue` | Alta de nuevo alumno con propagación en pagos, contratos y evaluaciones. |
| **RegisterPaymentModal** | `src/components/modals/RegisterPaymentModal.vue` | Registro de abonos y actualización del estado de membresía. |
| **CreateEvaluationModal** | `src/components/modals/CreateEvaluationModal.vue` | Registro de métricas físicas cuantitativas y cualitativas dinámicas. |
| **ManageTermsModal** | `src/components/modals/ManageTermsModal.vue` | Editor de cláusulas contractuales y control de versión vigente. |
| **SimulatedActionModal** | `src/components/modals/SimulatedActionModal.vue` | Previsualizador de links de pago, mensajes de WhatsApp y comprobantes. |

---

## 7. Rutas de Navegación (Vue Router)

Configuradas en [`src/router/index.ts`](file:///c:/Users/dguerra/Documents/GitHub/Prototipo-Inverfit/src/router/index.ts):

| Ruta | Nombre | Componente de Vista | Título del Navegador | Propósito |
|---|---|---|---|---|
| `/` | — | *(Redirección)* | — | Redirige automáticamente a `/inicio`. |
| `/inicio` | `inicio` | `InicioView.vue` | `Inicio — Panel del entrenador \| INVERFIT` | Dashboard principal y resumen operativo. |
| `/agenda` | `agenda` | `AgendaView.vue` | `Agenda semanal \| INVERFIT` | Planificador semanal y disponibilidad. |
| `/alumnos` | `alumnos` | `AlumnosView.vue` | `Alumnos \| INVERFIT` | Maestro de alumnos y filtros. |
| `/alumnos/:id` | `alumno-detalle` | `FichaAlumnoView.vue` | `Ficha del alumno \| INVERFIT` | Perfil individual y 5 pestañas de detalle. |
| `/membresias-pagos` | `membresias-pagos` | `MembresiasPagosView.vue` | `Membresías y pagos \| INVERFIT` | Gestión financiera, cobros y exportación CSV. |
| `/evaluaciones` | `evaluaciones` | `EvaluacionesView.vue` | `Evaluaciones \| INVERFIT` | Registro y seguimiento de métricas físicas. |
| `/contratos` | `contratos` | `ContratosView.vue` | `Contratos y consentimientos \| INVERFIT` | Visor de documentos y control de firmas. |
| `/:pathMatch(.*)*` | — | *(Catch-all)* | — | Redirige rutas inexistentes a `/inicio`. |

---

## 8. Funcionalidades Simuladas y Limitaciones Actuales

Para mantener total transparencia con el cliente y el equipo técnico, a continuación se desglosan las operaciones que funcionan como **simulaciones de demostración**:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        MAPA DE REALIDAD OPERATIVA                      │
├────────────────────────────────┬───────────────────────────────────────┤
│ FUNCIONALIDAD REAL EN CLIENTE  │ OPERACIÓN SIMULADA (SIN BACKEND/APIS) │
├────────────────────────────────┼───────────────────────────────────────┤
│ • Lógica de agenda y cálculo   │ • Pasarelas de pago bancarias         │
│ • Validaciones y conflictos    │ • Envío real de correos / WhatsApp    │
│ • Persistencia en localStorage │ • Emisión de boletas ante el SII      │
│ • Exportación descargable CSV  │ • Firma electrónica bajo ley chilena  │
│ • Sincronización entre módulos │ • Portal de acceso para alumnos       │
│ • Edición de datos y métricas  │ • Autenticación y control de roles    │
└────────────────────────────────┴───────────────────────────────────────┘
```

1. **Pasarelas de Pago:** La aplicación no está conectada a Transbank, Webpay, Flow, Mercado Pago ni Stripe. El botón *"Enviar enlace de pago"* genera una URL demostrativa bajo el dominio ficticio `https://demo.inverfit.invalid/checkout/...`.
2. **Mensajería y Notificaciones:** No existen integraciones con servidores SMTP ni APIs de mensajería (Twilio, Meta WhatsApp Business API). Los envíos de recordatorios son confirmaciones visuales en pantalla que actualizan contadores locales.
3. **Facturación Electrónica:** No existe conexión con el Servicio de Impuestos Internos (SII) ni con proveedores de facturación electrónica. Los comprobantes son simulaciones para verificar el flujo de trabajo del entrenador.
4. **Firma de Contratos:** La aceptación de contratos es un estado booleano almacenado en el navegador; no constituye una firma digital con validez legal.
5. **Autenticación y Sesiones:** La interfaz asume que el usuario autenticado es siempre el entrenador **Pablo**. No existen formularios de inicio de sesión, recuperación de contraseña ni tokens JWT.
6. **Portal del Alumno:** Los alumnos no disponen de una interfaz para consultar sus horarios o pagar de forma autónoma.

---

## 9. Instalación y Ejecución del Proyecto

### 9.1. Prerrequisitos
- **Node.js:** Versión 18.0 o superior (recomendado Node.js 20+ LTS).
- **npm:** Versión 9.0 o superior.

### 9.2. Instrucciones Paso a Paso

1. **Ingresar al directorio del repositorio:**
   ```bash
   cd "Prototipo-Inverfit"
   ```

2. **Instalar dependencias:**
   ```bash
   npm install
   ```

3. **Iniciar el servidor de desarrollo local:**
   ```bash
   npm run dev
   ```
   La aplicación se iniciará en `http://localhost:5173/`.

4. **Compilar para producción (validación de TypeScript y empaquetado):**
   ```bash
   npm run build
   ```
   Ejecuta `vue-tsc -b` para verificar tipos y genera los activos optimizados en el directorio `dist/`.

5. **Previsualizar la versión compilada de producción:**
   ```bash
   npm run preview
   ```

---

## 10. Implementación del Backend (Go + Gin + MongoDB) y Conexión Frontend (Fases 1, 2, 3 y 4A)

A partir de las Fases 1, 2, 3 y 4A, el proyecto cuenta con infraestructura backend real, CRUD de alumnos, el motor completo de agenda semanal y el módulo financiero de pagos y membresías:

### 10.1. Infraestructura Backend (Go 1.23+ / Gin / MongoDB 7.0)
- **Monolito modular por capas:** Router → Handler → Service → Repository → MongoDB.
- **Seguridad:** CORS estricto, ofuscación de credenciales en logs (`RedactedMongoURI`), variables de entorno locales protegidas de git.
- **Docker Compose:** Orquestación local de MongoDB 7.0 oficial en `backend/compose.yaml` con volumen persistente `inverfit-mongo-data`.
- **Healthchecks:** `/health/live` (liveness) y `/health/ready` (readiness con ping a MongoDB).

### 10.2. Módulo de Alumnos (Fase 2)
- **Endpoints REST (`/api/v1/students`):**
  - `GET /api/v1/students`: listado con filtros reactivos (`name`, `status`, `frequency`).
  - `GET /api/v1/students/:id`: detalle por identificador hexadecimal (ObjectID).
  - `POST /api/v1/students`: creación con validaciones y valores por defecto.
  - `PUT /api/v1/students/:id`: actualización parcial de datos de contacto y condiciones.
  - `DELETE /api/v1/students/:id`: baja lógica (*soft delete*) preservando historial mediante `archivedAt`.
- **Modo Dual en Frontend (Vue 3 / Pinia):**
  - `VITE_STUDENTS_DATA_SOURCE=mock` (por defecto): Mantiene funcionamiento autónomo en `localStorage` para Cloudflare Pages y demostraciones sin servidor.
  - `VITE_STUDENTS_DATA_SOURCE=api`: Delega en `StudentApiService` llamando al backend Go en `VITE_API_BASE_URL`. Sin fallback silencioso a mock.
  - Interfaz de usuario enriquecida con indicador de fuente de datos (badge Mock/API), botón de refresco, banner de reintento, estados de carga y pantalla de alumno no encontrado (404) con retorno asistido.

### 10.3. Módulo de Agenda Semanal y Disponibilidad (Fase 3)
- **Endpoints REST (`/api/v1/schedule`):**
  - Disponibilidad semanal: `GET /availability`, `GET /availability/:day`, `PUT /availability/:day`, `DELETE /availability/:day`.
  - Sesiones y Proyección de Calendario: `GET /sessions?from=&to=&studentId=`, `POST /sessions`, `GET /sessions/:id`, `PUT /sessions/:id`, `DELETE /sessions/:id`.
  - Excepciones de Calendario: `POST /sessions/:id/exceptions`, `GET /exceptions`, `DELETE /sessions/:id/exceptions/:exceptionId`.
- **Reglas de Negocio en Go:**
  - Administración exclusiva del entrenador.
  - Disponibilidad diaria por intervalos (`06:30` - `22:00`).
  - Duración sugerida 75 min (libremente editable).
  - Cero descansos forzados (sesiones continuas perfectamente válidas).
  - Verificación estricta de solapamiento temporal (`startA < endB && startB < endA`) sensible a rangos de vigencia de fechas (`startDate` y `endDate`) y excepciones activas de calendario (las fechas canceladas o reprogramadas no bloquean el horario original, y las sesiones movidas validan conflictos en destino).
  - Detección de conflicto al reducir disponibilidad (`409 AVAILABILITY_CONFLICT`).
  - Excepciones de calendario canceladas y reprogramadas con proyección sobre zona horaria `America/Santiago`.
  - Desacoplamiento modular estricto mediante interfaz `StudentReader` y adaptador en `main.go`.
- **Modo Dual en Frontend:**
  - `VITE_SCHEDULE_DATA_SOURCE=mock` (por defecto) | `api`.
  - Delegación transparente mediante `scheduleService` (`schedule.mock.ts` y `schedule.api.ts`).
  - Ficha del alumno obtiene el horario asignado dinámicamente desde las sesiones activas de la agenda.
  - **Garantía de serialización no nula y resiliencia:** Todas las respuestas de colecciones (`/sessions`, `/availability`, `/exceptions`, `intervals`, `studentIds`, `conflicts`) garantizan la serialización de arreglos vacíos `[]` en el backend Go (evitando serializar `null`), y el cliente frontend normaliza defensivamente los arreglos para prevenir excepciones de tipo `Cannot read properties of null (reading 'map')`.

### 10.4. Módulo de Pagos y Membresías (Fase 4A)
- **Endpoints REST (`/api/v1/payments`):**
  - Listado y filtros: `GET /payments?studentId=&status=&search=&from=&to=`.
  - Registro de abonos: `POST /payments`.
  - Balance consolidado: `GET /payments/summary` (cálculo en vivo mediante agregaciones en MongoDB).
  - Exportación streaming: `GET /payments/export` (CSV en UTF-8 con BOM y desinfección contra inyección de fórmulas `=, +, -, @`).
  - Historial por alumno: `GET /payments/student/:studentId`.
  - Detalle, actualización y soft-delete: `GET /payments/:id`, `PUT /payments/:id`, `DELETE /payments/:id`.
  - Simulaciones operativas: `POST /payments/:id/simulate-reminder`, `POST /payments/:id/simulate-receipt`.
- **Reglas de Negocio en Go:**
  - Precisión monetaria entera en pesos chilenos (`int64` CLP), evitando imprecisiones de coma flotante.
  - Preservación irrestricta del historial financiero mediante borrado lógico (`archivedAt`).
  - Activación automática de membresía (`Por revisar` → `Activo`) vía `StudentReader` al confirmar un pago.
  - Desacoplamiento modular absoluto respecto a la colección de alumnos.
- **Modo Dual en Frontend:**
  - `VITE_PAYMENTS_DATA_SOURCE=mock` (por defecto) | `api`.
  - Delegador `paymentService` (`payments.mock.ts` y `payments.api.ts`).
  - Vista `MembresiasPagosView.vue` enriquecida con badge de fuente de datos, botón de recarga, estados de carga y exportación reactiva.

### 10.5. Despliegue Dockerizado Unificado del Backend (compose.stack.yaml)
- **Construcción Multi-stage (`backend/Dockerfile`):**
  - Compilación estática pura (`CGO_ENABLED=0`, `-trimpath`, `-ldflags="-w -s"`) en `golang:1.25-alpine`.
  - Imagen final mínima en `alpine:3.20` con certificados CA (`ca-certificates`), zona horaria `tzdata` (`TZ=America/Santiago`) y herramienta `curl` para comprobaciones de salud.
  - Usuario no privilegiado `appuser:appgroup` (UID:GID `10001:10001`), sin compiladores ni secretos en la imagen final.
- **Orquestación y Red Aislada (`compose.stack.yaml`):**
  - Red interna privada tipo bridge (`inverfit-stack-net`).
  - Servicio `api`: mapeado exclusivamente a `127.0.0.1:8080` en el host, con comprobación periódica de salud (`/health/ready`).
  - Servicio `mongodb`: imagen oficial `mongo:7.0`, sin puertos publicados al host para prevenir colisiones con el entorno de desarrollo local (puerto 27017).
- **Persistencia Segura y Credenciales de Aplicación:**
  - Volumen persistente dedicado e independiente: `inverfit-stack-mongo-data` (aislado de `inverfit-mongo-data`).
  - Script de inicialización `docker/init-mongo.sh` para crear el usuario de aplicación `inverfit_app` con permisos `readWrite` sobre `inverfit_container`.
  - La API Go opera conectada con el usuario de aplicación, nunca con credenciales root.
  - Persistencia de datos validada ante reinicio de API, reinicio de base de datos y recreación de contenedor de MongoDB sin pérdida de volumen.

---

## 11. Matriz de Estado General del Proyecto

| Funcionalidad / Módulo | Estado | Observaciones Técnicas |
|---|---|---|
| **Normalización de marca INVERFIT** | **Implementada y comprobada** | Logotipo SVG oficial, paleta corporativa y textos normalizados al 100%. |
| **Navegación SPA y Rutas** | **Implementada y comprobada** | Las 7 rutas funcionan con historial HTML5 y títulos de pestaña dinámicos. |
| **Dashboard de Inicio** | **Implementada y comprobada** | 4 KPIs en tiempo real, agenda de hoy y alertas de pendientes operativas. |
| **Disponibilidad por intervalos (Agenda)** | **Implementada y comprobada** | Configuración de 0 a $N$ turnos por día y detección preventiva de conflictos. |
| **Duración sugerida 75 min (Agenda)** | **Implementada y comprobada** | Sugerencia automática y duración libre editable con cálculo en vivo. |
| **Cero descansos forzados (Agenda)** | **Implementada y comprobada** | Eliminada la regla de 15 min; se admiten sesiones consecutivas sin trabas. |
| **Detección de solapamiento real (Agenda)** | **Implementada y comprobada** | Verificación estricta de traslapes temporales sin falsos positivos. |
| **Excepciones de calendario (Agenda)** | **Implementada y comprobada** | Reagendamiento y cancelación por fecha específica (`YYYY-MM-DD`). |
| **Sesiones grupales (Agenda)** | **Implementada y comprobada** | Soporte de hasta 4 alumnos y sincronización de horarios fijos. |
| **Listado y filtros de Alumnos** | **Implementada y comprobada** | Búsqueda reactiva y filtrado por estado y frecuencia semanal. |
| **CRUD de Alumnos** | **Implementada y comprobada** | Alta, edición, baja con confirmación y cascada sobre datos relacionados. |
| **Ficha del Alumno (5 pestañas)** | **Implementada y comprobada** | Resumen, datos personales, pagos, evaluaciones y contrato accesibles. |
| **Control de Membresías y Pagos** | **Implementada y comprobada** | Registro manual de pagos, actualización de estados y trazabilidad. |
| **Exportación a CSV** | **Implementada y comprobada** | Generación y descarga real de archivo CSV con el historial financiero. |
| **Simulación de enlaces de cobro** | **Implementada y comprobada** | Tokenización ficticia segura y plantilla preformateada para WhatsApp. |
| **Evaluaciones con métricas dinámicas** | **Implementada y comprobada** | Soporte para añadir/quitar campos cuantitativos y cualitativos al vuelo. |
| **Visor y Versionado de Contratos** | **Implementada y comprobada** | Visualización de cláusulas, cambio de vigencia y simulación de firma. |
| **Persistencia reactiva en LocalStorage** | **Implementada y comprobada** | Almacenamiento síncrono con notificación a Pinia y botón de reinicio. |
| **Validación de tipos y compilación** | **Implementada y comprobada** | `vue-tsc` y `vite build` finalizan con código 0 y sin errores. |
| **Criterios de evaluación física de Pablo** | *Implementada, pendiente de validación* | Esquema dinámico listo; pendiente definir pruebas físicas estándar del centro. |
| **Planes tarifarios oficiales** | *Implementada, pendiente de validación* | Funciona con montos libres; requiere fijar aranceles oficiales por plan. |
| **Backend en Go + MongoDB (Fases 1, 2, 3 y 4A)** | **Implementada y comprobada** | Servidor Gin, Docker MongoDB 7.0, Alumnos, Agenda, Pagos y Balances (/api/v1/payments) y modo dual en Vue. |
| **Stack Dockerizado Unificado del Backend** | **Implementada y comprobada** | Dockerfile multi-stage (Alpine, non-root, tzdata), compose.stack.yaml, persistencia en inverfit-stack-mongo-data y credenciales de app no-root. |
| **Pasarela real de pagos (Webpay/Stripe)** | *Pendiente de desarrollo* | Requiere contrato bancario comercial y credenciales API del proveedor. |
| **Portal web de autoatención para Alumnos** | *Pendiente de desarrollo* | Fuera del alcance del panel del entrenador en la presente versión. |
| **Firma Electrónica Avanzada (FEA)** | *Pendiente de desarrollo* | Requiere integración con prestador acreditado de certificación digital. |

---

*Fin del documento de especificación funcional y técnica.*
