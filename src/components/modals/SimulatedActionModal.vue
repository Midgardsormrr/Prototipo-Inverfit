<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Copy, Check, Send, Mail, MessageSquare, Printer, CheckCircle2, CreditCard, ShieldCheck, FileCheck, FileText } from 'lucide-vue-next'
import { useGymStore } from '@/stores/gym'

const props = defineProps<{
  open: boolean
  mode: 'link' | 'receipt' | 'reminder'
  studentName: string
  studentEmail?: string
  studentPhone?: string
  linkUrl?: string
  paymentId?: string
  amount?: number | null
  frequency?: string
  dueDate?: string
  method?: string
  date?: string
  transactionReference?: string
  voucherUrl?: string
  voucherFileName?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const gymStore = useGymStore()

const copied = ref(false)
const reminderLogged = ref(false)
const receiptLogged = ref(false)

const emailTarget = ref('')
const customLink = ref('')

watch(() => props.open, (isOpen) => {
  if (isOpen) {
    copied.value = false
    reminderLogged.value = false
    receiptLogged.value = false
    emailTarget.value = props.studentEmail || ''
    customLink.value = props.linkUrl || 'https://link.sumup.com/pay/inverfit'
  }
})

const effectiveFrequency = computed(() => props.frequency || '2 veces/sem.')
const formattedAmount = computed(() => {
  return (props.amount && props.amount > 0)
    ? `$${props.amount.toLocaleString('es-CL')}`
    : 'Por definir'
})

function copyLink() {
  if (customLink.value) {
    navigator.clipboard.writeText(customLink.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2500)
    gymStore.showToast('Enlace de pago SumUp copiado al portapapeles.', 'success')
  }
}

// ----------------------------------------------------
// Mode 1: Payment Link (SumUp)
// ----------------------------------------------------
const linkEmailSubject = computed(() => {
  return `INVERFIT — Enlace de Pago SumUp: Membresía Mensual (${props.studentName})`
})

const linkEmailBody = computed(() => {
  return `Hola ${props.studentName},

Te comparto tu enlace oficial de SumUp para el pago de tu Membresía Mensual INVERFIT (${effectiveFrequency.value}):

${customLink.value}

Puedes realizar el pago directamente con cualquier tarjeta de crédito o débito.
Recuerda que nuestro ciclo oficial de renovación se realiza del 1 al 5 de cada mes.

Una vez realizado, avísame para confirmar tus bloques en la agenda. ¡A seguir entrenando fuerte!

Atentamente,
Coach Pablo — INVERFIT
Centro de Entrenamiento de Fuerza y Rendimiento`
})

function sendLinkByEmail() {
  const mailto = `mailto:${encodeURIComponent(emailTarget.value || '')}?subject=${encodeURIComponent(linkEmailSubject.value)}&body=${encodeURIComponent(linkEmailBody.value)}`
  window.location.href = mailto
  if (props.paymentId) {
    gymStore.simulatePaymentLink(props.studentName, props.paymentId)
  }
}

function sendLinkByWhatsApp() {
  const cleanPhone = (props.studentPhone || '').replace(/\D/g, '')
  const msg = `Hola ${props.studentName}! Te comparto tu enlace de pago SumUp para la Membresía Mensual en INVERFIT: ${customLink.value} (pago seguro con tarjeta débito/crédito). Saludos! — Coach Pablo`
  const url = cleanPhone
    ? `https://wa.me/${cleanPhone}?text=${encodeURIComponent(msg)}`
    : `https://wa.me/?text=${encodeURIComponent(msg)}`
  window.open(url, '_blank')
}

// ----------------------------------------------------
// Mode 2: Reminder Notice (Cobro del 1 al 5)
// ----------------------------------------------------
const reminderEmailSubject = computed(() => {
  return `INVERFIT — Recordatorio de Pago Mensualidad (Período 1 al 5) - ${props.studentName}`
})

const reminderEmailBody = computed(() => {
  return `Estimado(a) ${props.studentName},

Te recordamos que nos encontramos en la ventana oficial de cobro y renovación de tu Membresía Mensual INVERFIT (${effectiveFrequency.value}), correspondiente al período del 1 al 5 de este mes.

Detalle de tu Membresía:
- Alumno: ${props.studentName}
- Membresía: Membresía Mensual (${effectiveFrequency.value})
- Período de Cobro: 1 al 5 de cada mes
- Monto: ${formattedAmount.value}

CANALES OFICIALES DE PAGO:

1. Tarjetas de Débito / Crédito (SumUp):
${customLink.value}

2. Transferencia Bancaria Directa:
Banco: Banco de Chile
Tipo de Cuenta: Cuenta Corriente
N° de Cuenta: 00-12345678-01
Nombre: INVERFIT - Pablo
RUT: 12.345.678-9
Correo para comprobantes: pagos@inverfit.cl

3. Efectivo en Box de Entrenamiento

Por favor envíanos el comprobante tras realizar la transferencia o pago con tarjeta para validar tu membresía y confirmar tus horarios en la agenda.

Atentamente,
Coach Pablo — INVERFIT
Centro de Entrenamiento de Fuerza y Rendimiento`
})

async function openReminderInEmail() {
  const mailto = `mailto:${encodeURIComponent(emailTarget.value || '')}?subject=${encodeURIComponent(reminderEmailSubject.value)}&body=${encodeURIComponent(reminderEmailBody.value)}`
  window.location.href = mailto

  if (props.paymentId) {
    await gymStore.simulateReminder(props.paymentId)
    reminderLogged.value = true
  }
}

function sendReminderByWhatsApp() {
  const cleanPhone = (props.studentPhone || '').replace(/\D/g, '')
  const msg = `Hola ${props.studentName}! Te recuerdo que estamos en la ventana oficial de cobro (1 al 5 de cada mes) para tu Membresía Mensual INVERFIT (${effectiveFrequency.value}). Puedes pagar con tarjeta vía SumUp aquí: ${customLink.value} o por transferencia. Saludos! — Coach Pablo`
  const url = cleanPhone
    ? `https://wa.me/${cleanPhone}?text=${encodeURIComponent(msg)}`
    : `https://wa.me/?text=${encodeURIComponent(msg)}`
  window.open(url, '_blank')
}

async function markReminderLogged() {
  if (props.paymentId) {
    await gymStore.simulateReminder(props.paymentId)
    reminderLogged.value = true
  }
}

// ----------------------------------------------------
// Mode 3: Receipt (Comprobante Oficial)
// ----------------------------------------------------
const receiptEmailSubject = computed(() => {
  return `INVERFIT — Comprobante Oficial de Pago Membresía Mensual - ${props.studentName}`
})

const receiptEmailBody = computed(() => {
  return `Hola ${props.studentName},

Adjuntamos el detalle de confirmación de tu pago de Membresía Mensual en INVERFIT:

- Alumno: ${props.studentName}
- Membresía: Membresía Mensual (${effectiveFrequency.value})
- Medio de Pago: ${props.method || 'SumUp (Link de pago)'}
- Fecha de Registro: ${props.date || new Date().toISOString().split('T')[0]}
- Monto Pagado: ${formattedAmount.value}
- Estado: PAGADO / AL DÍA

¡Muchas gracias por tu compromiso y por entrenar con nosotros en INVERFIT!

Atentamente,
Coach Pablo — INVERFIT`
})

function sendReceiptByEmail() {
  const mailto = `mailto:${encodeURIComponent(emailTarget.value || '')}?subject=${encodeURIComponent(receiptEmailSubject.value)}&body=${encodeURIComponent(receiptEmailBody.value)}`
  window.location.href = mailto
  markReceiptLogged()
}

async function markReceiptLogged() {
  if (props.paymentId) {
    await gymStore.simulateReceipt(props.paymentId)
    receiptLogged.value = true
  }
}

function printReceipt() {
  window.print()
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-xs overflow-y-auto"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-2xl max-w-lg w-full border border-[#E2DFD8] overflow-hidden animate-in fade-in zoom-in-95 duration-150 my-6">
      <!-- Header con Identidad INVERFIT -->
      <div class="px-6 py-4 border-b border-[#E2DFD8] flex items-center justify-between bg-[#FAF9F6]">
        <div>
          <div class="flex items-center gap-2">
            <span class="text-[10px] font-bold uppercase tracking-wider text-[#CEA549] bg-[#191A1D] px-2 py-0.5 rounded">
              INVERFIT · Pagos
            </span>
            <span class="text-xs text-[#8E9097] font-semibold">
              Coach Pablo
            </span>
          </div>
          <h2 class="text-lg font-bold text-[#191A1D] mt-1.5 flex items-center gap-2">
            <template v-if="mode === 'link'">
              <CreditCard class="w-5 h-5 text-[#CEA549]" /> Enlace de Pago SumUp
            </template>
            <template v-else-if="mode === 'receipt'">
              <CheckCircle2 class="w-5 h-5 text-emerald-600" /> Comprobante Oficial de Pago
            </template>
            <template v-else>
              <Mail class="w-5 h-5 text-[#CEA549]" /> Aviso de Cobro (Período 1 al 5)
            </template>
          </h2>
        </div>
        <button
          @click="emit('close')"
          class="text-[#8E9097] hover:text-[#191A1D] text-xl font-light w-8 h-8 rounded-full flex items-center justify-center hover:bg-[#F5F3EE] transition-colors cursor-pointer"
        >
          ✕
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 space-y-4 max-h-[calc(85vh-120px)] overflow-y-auto">
        <!-- ============================================== -->
        <!-- MODO 1: ENLACE DE PAGO SUMUP                   -->
        <!-- ============================================== -->
        <div v-if="mode === 'link'" class="space-y-4">
          <p class="text-xs text-[#4A4C52]">
            Genera y envía el enlace directo de <strong>SumUp</strong> para que <strong>{{ studentName }}</strong> pueda pagar su membresía mensual con tarjeta de débito o crédito.
          </p>

          <!-- Input enlace SumUp -->
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Enlace Oficial SumUp
            </label>
            <div class="p-2.5 bg-[#FAF9F6] rounded-lg border border-[#E2DFD8] flex items-center justify-between gap-2">
              <input
                v-model="customLink"
                type="text"
                class="text-xs text-[#191A1D] font-mono w-full bg-transparent border-0 focus:outline-none"
                placeholder="https://link.sumup.com/pay/..."
              />
              <button
                @click="copyLink"
                type="button"
                class="px-2.5 py-1.5 bg-white border border-[#CEA549] rounded text-xs font-bold text-[#191A1D] hover:bg-[#CEA549]/10 flex items-center gap-1 transition-colors flex-shrink-0 cursor-pointer shadow-2xs"
              >
                <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-600" />
                <Copy v-else class="w-3.5 h-3.5 text-[#8C6B1B]" />
                {{ copied ? 'Copiado' : 'Copiar' }}
              </button>
            </div>
          </div>

          <!-- Correo destinatario -->
          <div>
            <label class="block text-xs font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Correo del Alumno para Enviar Enlace
            </label>
            <input
              v-model="emailTarget"
              type="email"
              placeholder="correo@ejemplo.com"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] focus:ring-1 focus:ring-[#CEA549]/30 transition-all font-mono"
            />
          </div>

          <!-- Previsualización del mensaje -->
          <div class="p-3 bg-[#FAF9F6] text-[#4A4C52] border border-[#E2DFD8] rounded-lg text-xs space-y-1.5">
            <div class="flex items-center justify-between">
              <span class="font-bold text-[#191A1D] text-[11px] uppercase tracking-wide flex items-center gap-1.5">
                <ShieldCheck class="w-4 h-4 text-[#CEA549]" /> Mensaje de Cobro SumUp
              </span>
              <span class="text-[11px] text-[#8E9097]">Membresía Mensual ({{ effectiveFrequency }})</span>
            </div>
            <p class="italic text-[12px] text-[#555861] whitespace-pre-wrap leading-relaxed">
              "Hola {{ studentName }}, te comparto tu enlace de SumUp para tu Membresía Mensual INVERFIT (cobro del 1 al 5): {{ customLink }}. ¡Cualquier consulta me avisas!"
            </p>
          </div>

          <!-- Botones de Acción Inmediata -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 pt-1">
            <button
              @click="sendLinkByEmail"
              type="button"
              class="w-full py-2.5 px-3 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <Mail class="w-4 h-4" />
              Enviar Link por Correo
            </button>
            <button
              @click="sendLinkByWhatsApp"
              type="button"
              class="w-full py-2.5 px-3 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <MessageSquare class="w-4 h-4" />
              Enviar por WhatsApp
            </button>
          </div>
        </div>

        <!-- ============================================== -->
        <!-- MODO 2: RECORDATORIO DE COBRO DEL 1 AL 5       -->
        <!-- ============================================== -->
        <div v-else-if="mode === 'reminder'" class="space-y-4">
          <!-- Banner período oficial -->
          <div class="p-3 bg-[#FEF7E0] border border-[#F0D878] rounded-md flex items-center justify-between text-xs text-[#8C6B1B]">
            <div class="flex items-center gap-2 font-semibold">
              <span class="w-2 h-2 rounded-full bg-[#CEA549] animate-pulse"></span>
              <span>Ventana Oficial de Cobro: <strong>Días 1 al 5 de cada mes</strong></span>
            </div>
            <span class="font-bold text-[#191A1D]">{{ formattedAmount }}</span>
          </div>

          <!-- Datos de destinatario -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-xs">
            <div>
              <label class="block text-[11px] font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
                Alumno
              </label>
              <div class="px-3 py-2 bg-[#FAF9F6] border border-[#E2DFD8] rounded-md font-bold text-[#191A1D]">
                {{ studentName }}
              </div>
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
                Correo Electrónico
              </label>
              <input
                v-model="emailTarget"
                type="email"
                placeholder="correo@ejemplo.com"
                class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] font-mono"
              />
            </div>
          </div>

          <!-- Previsualización del correo formal -->
          <div>
            <div class="flex items-center justify-between mb-1">
              <label class="text-[11px] font-semibold text-[#6B6D74] uppercase tracking-wider">
                Previsualización Correo Oficial INVERFIT
              </label>
              <span class="text-[11px] text-[#8E9097]">Asunto: {{ reminderEmailSubject }}</span>
            </div>
            <div class="p-3.5 bg-[#FAF9F6] rounded-lg text-[12px] text-[#4A4C52] border border-[#E2DFD8] font-mono whitespace-pre-wrap max-h-48 overflow-y-auto leading-relaxed select-all">
{{ reminderEmailBody }}
            </div>
          </div>

          <!-- Acciones de aviso -->
          <div class="space-y-2 pt-1">
            <button
              @click="openReminderInEmail"
              type="button"
              class="w-full py-2.5 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <Mail class="w-4 h-4" />
              Abrir Aviso en Cliente de Correo (Gmail / Outlook)
            </button>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
              <button
                @click="sendReminderByWhatsApp"
                type="button"
                class="py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
              >
                <MessageSquare class="w-4 h-4" />
                Enviar por WhatsApp
              </button>
              <button
                @click="markReminderLogged"
                type="button"
                class="py-2 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors cursor-pointer"
              >
                <Check class="w-4 h-4 text-[#8C6B1B]" />
                Registrar aviso enviado
              </button>
            </div>
          </div>

          <div v-if="reminderLogged" class="p-2.5 bg-emerald-50 text-emerald-800 rounded-md text-xs font-semibold text-center border border-emerald-200">
            ✓ Aviso registrado en el historial de cobros del alumno.
          </div>
        </div>

        <!-- ============================================== -->
        <!-- MODO 3: COMPROBANTE DE PAGO                    -->
        <!-- ============================================== -->
        <div v-else class="space-y-4">
          <!-- Recibo gráfico INVERFIT -->
          <div class="p-5 bg-[#FAF9F6] border-2 border-[#191A1D] rounded-lg space-y-3 text-xs shadow-xs">
            <div class="flex justify-between items-center border-b pb-3 border-[#E2DFD8]">
              <div class="flex items-center gap-2">
                <span class="text-sm font-black text-[#191A1D] tracking-tight">INVERFIT</span>
                <span class="text-sm font-black text-[#CEA549]">❯❯</span>
              </div>
              <span class="text-[11px] font-bold text-[#CEA549] bg-[#191A1D] px-2.5 py-0.5 rounded">
                COMPROBANTE DE PAGO
              </span>
            </div>

            <div class="grid grid-cols-2 gap-2 text-[#6B6D74]">
              <span>Alumno:</span>
              <span class="font-bold text-[#191A1D] text-right">{{ studentName }}</span>

              <span>Membresía:</span>
              <span class="font-semibold text-[#191A1D] text-right">Membresía Mensual ({{ effectiveFrequency }})</span>

              <span>Medio de Pago:</span>
              <span class="font-semibold text-[#191A1D] text-right">{{ method || 'SumUp (Link de pago)' }}</span>

              <span>N° Operación SumUp:</span>
              <span class="font-mono text-[#191A1D] text-right font-bold">{{ transactionReference || 'SUM-OFICIAL-01' }}</span>

              <span>Fecha de Pago:</span>
              <span class="font-semibold text-[#191A1D] text-right">{{ date || new Date().toISOString().split('T')[0] }}</span>

              <span>Monto Recaudado:</span>
              <span class="font-bold text-emerald-700 text-right text-sm">{{ formattedAmount }}</span>

              <span>Estado:</span>
              <span class="font-bold text-emerald-700 text-right">✓ PAGADO / AL DÍA</span>
            </div>

            <div class="border-t pt-2 border-[#E2DFD8] text-[10px] text-[#8E9097] text-center">
              INVERFIT · Centro de Entrenamiento de Fuerza y Rendimiento · Coach Pablo
            </div>
          </div>

          <!-- Voucher Original Adjuntado de SumUp o Transferencia -->
          <div v-if="voucherUrl" class="p-4 bg-[#FAF9F6] border border-[#E2DFD8] rounded-lg space-y-2">
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-[#191A1D] flex items-center gap-1.5">
                <FileCheck class="w-4 h-4 text-emerald-700" />
                Voucher / Comprobante Original Adjunto
              </span>
              <span class="text-[10px] text-[#8E9097] font-mono">{{ voucherFileName || 'voucher_sumup.png' }}</span>
            </div>

            <div class="border border-[#DCD8CF] rounded-lg overflow-hidden bg-white p-2 flex items-center justify-center">
              <img
                v-if="voucherUrl.startsWith('data:image') || voucherUrl.includes('image')"
                :src="voucherUrl"
                class="max-h-56 w-auto object-contain rounded"
                alt="Comprobante SumUp adjunto"
              />
              <div v-else class="p-6 text-center text-xs text-[#6B6D74]">
                <FileText class="w-8 h-8 text-[#CEA549] mx-auto mb-1" />
                <span>Documento PDF / Comprobante digital adjunto</span>
              </div>
            </div>
          </div>

          <!-- Correo de envío de comprobante -->
          <div>
            <label class="block text-[11px] font-semibold text-[#6B6D74] uppercase tracking-wider mb-1">
              Enviar copia al correo del alumno
            </label>
            <input
              v-model="emailTarget"
              type="email"
              placeholder="correo@ejemplo.com"
              class="w-full px-3 py-2 border border-[#DCD8CF] rounded-md text-xs text-[#191A1D] bg-white focus:outline-none focus:border-[#CEA549] font-mono"
            />
          </div>

          <!-- Acciones de Comprobante -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 pt-1">
            <button
              @click="sendReceiptByEmail"
              type="button"
              class="w-full py-2 bg-[#CEA549] hover:bg-[#BA923B] text-[#0B0B0D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors shadow-xs cursor-pointer"
            >
              <Send class="w-4 h-4" />
              Enviar por Correo
            </button>
            <button
              @click="printReceipt"
              type="button"
              class="w-full py-2 bg-white border border-[#DCD8CF] hover:bg-[#FAF9F6] text-[#191A1D] rounded-md text-xs font-bold flex items-center justify-center gap-2 transition-colors cursor-pointer"
            >
              <Printer class="w-4 h-4 text-[#8C6B1B]" />
              Imprimir Recibo
            </button>
          </div>

          <div v-if="receiptLogged" class="p-2.5 bg-emerald-50 text-emerald-800 rounded-md text-xs font-semibold text-center border border-emerald-200">
            ✓ Comprobante registrado y enviado.
          </div>
        </div>

        <!-- Footer -->
        <div class="pt-3 border-t border-[#E2DFD8] flex items-center justify-between">
          <span class="text-[11px] text-[#8E9097]">
            Cobranza Oficial INVERFIT (1 al 5 de cada mes)
          </span>
          <button
            @click="emit('close')"
            class="px-4 py-2 bg-[#F5F3EE] hover:bg-[#EAE6DD] text-[#191A1D] rounded-md text-xs font-semibold transition-colors cursor-pointer border border-[#DDD8CE]"
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>