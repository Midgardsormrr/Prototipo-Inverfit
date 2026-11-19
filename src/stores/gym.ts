import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  storage,
  studentService,
  paymentService,
  scheduleService,
  evaluationService,
  contractService,
  isStudentsApiMode,
  getStudentsDataSource,
  isScheduleApiMode,
  getScheduleDataSource,
  isPaymentsApiMode,
  getPaymentsDataSource
} from '@/services'
import type {
  Student,
  Payment,
  ScheduleBlock,
  DayAvailability,
  ScheduleException,
  SessionOccurrence,
  Evaluation,
  ContractConsent,
  TermsVersion,
  CreateStudentDTO,
  RegisterPaymentDTO,
  AssignStudentDTO,
  AssignSessionDTO,
  CreateExceptionDTO,
  CreateEvaluationDTO,
  ContractStatus,
  FinancialSummary,
  PaymentFilter
} from '@/types'
import type { WeekDayInfo } from '@/services/schedule'

export interface ToastNotification {
  id: string
  type: 'success' | 'info' | 'warning' | 'error'
  title?: string
  message: string
}

export const useGymStore = defineStore('gym', () => {
  // State
  const students = ref<Student[]>([])
  const loadingStudents = ref(false)
  const studentsError = ref<string | null>(null)
  const dataSource = computed(() => getStudentsDataSource())

  const payments = ref<Payment[]>([])
  const loadingPayments = ref(false)
  const paymentsError = ref<string | null>(null)
  const paymentsDataSource = computed(() => getPaymentsDataSource())
  const financialSummary = ref<FinancialSummary | null>(null)
  const sumupPermanentLink = ref<string>(
    localStorage.getItem('inverfit_sumup_permanent_link') || 'https://link.sumup.com/pay/inverfit'
  )

  function updateSumupPermanentLink(url: string) {
    const trimmed = url.trim()
    sumupPermanentLink.value = trimmed || 'https://link.sumup.com/pay/inverfit'
    localStorage.setItem('inverfit_sumup_permanent_link', sumupPermanentLink.value)
    showToast('Enlace permanente de SumUp actualizado correctamente.', 'success', 'Configuración SumUp')
  }

  const schedule = ref<ScheduleBlock[]>([])
  const availability = ref<DayAvailability[]>([])
  const exceptions = ref<ScheduleException[]>([])
  const weekOccurrences = ref<SessionOccurrence[]>([])
  const currentWeekRange = ref<{ from: string; to: string } | null>(null)
  const loadingSchedule = ref(false)
  const scheduleError = ref<string | null>(null)
  const scheduleDataSource = computed(() => getScheduleDataSource())

  const evaluations = ref<Evaluation[]>([])
  const contracts = ref<ContractConsent[]>([])
  const terms = ref<TermsVersion[]>([])
  const loading = ref(false)
  const toasts = ref<ToastNotification[]>([])
  const globalSearch = ref('')

  // Load from services
  async function reloadData() {
    if (!isStudentsApiMode()) {
      students.value = storage.getStudents()
    }
    if (!isScheduleApiMode()) {
      schedule.value = storage.getSchedule()
      availability.value = storage.getAvailability()
      exceptions.value = storage.getExceptions()
    }
    if (!isPaymentsApiMode()) {
      payments.value = storage.getPayments()
    }
    evaluations.value = storage.getEvaluations()
    contracts.value = storage.getContracts()
    terms.value = storage.getTerms()
  }

  // Load payments explicitly from paymentService (supports both API and mock)
  async function loadPayments(filter?: PaymentFilter) {
    loadingPayments.value = true
    paymentsError.value = null
    try {
      const [list, sum] = await Promise.all([
        paymentService.getAll(filter),
        paymentService.getFinancialSummary()
      ])
      payments.value = list
      financialSummary.value = sum
    } catch (err: unknown) {
      paymentsError.value = err instanceof Error ? err.message : 'Error al cargar pagos'
      showToast(paymentsError.value, 'error')
    } finally {
      loadingPayments.value = false
    }
  }

  // Load schedule explicitly from scheduleService (supports both API and mock)
  async function loadSchedule(from?: string, to?: string) {
    if (from && to) {
      currentWeekRange.value = { from, to }
    }
    loadingSchedule.value = true
    scheduleError.value = null
    try {
      const [avail, blocks, excs] = await Promise.all([
        scheduleService.getAvailability(),
        scheduleService.getAll(),
        scheduleService.getExceptions(from, to)
      ])
      availability.value = Array.isArray(avail) ? avail : []
      schedule.value = Array.isArray(blocks) ? blocks : []
      exceptions.value = Array.isArray(excs) ? excs : []
      return { availability: availability.value, schedule: schedule.value, exceptions: exceptions.value }
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al cargar agenda'
      scheduleError.value = msg
      if (isScheduleApiMode()) {
        schedule.value = []
        availability.value = []
        exceptions.value = []
        weekOccurrences.value = []
      }
      throw err
    } finally {
      loadingSchedule.value = false
    }
  }

  async function fetchWeekOccurrences(from: string, to: string, studentId?: string): Promise<SessionOccurrence[]> {
    currentWeekRange.value = { from, to }
    const occs = await scheduleService.getOccurrences(from, to, studentId)
    const list = Array.isArray(occs) ? occs : []
    weekOccurrences.value = list
    return list
  }

  // Sincroniza autoritativamente el estado completo de la agenda
  async function refreshScheduleState() {
    if (currentWeekRange.value) {
      await Promise.all([
        loadSchedule(currentWeekRange.value.from, currentWeekRange.value.to),
        fetchWeekOccurrences(currentWeekRange.value.from, currentWeekRange.value.to)
      ])
    } else if (isScheduleApiMode()) {
      await loadSchedule()
    } else {
      await reloadData()
    }
  }

  // Load students exclusively from studentService (works for both api and mock)
  async function loadStudents(filter?: { name?: string; status?: string; frequency?: string }) {
    loadingStudents.value = true
    studentsError.value = null
    try {
      const list = await studentService.getAll(filter)
      students.value = list
      syncContractsWithStudents(list)
      return list
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al cargar alumnos'
      studentsError.value = msg
      if (isStudentsApiMode()) {
        students.value = []
      }
      throw err
    } finally {
      loadingStudents.value = false
    }
  }

  /**
   * Sincroniza dinámicamente los consentimientos y contratos con los alumnos reales del gimnasio.
   * Elimina cualquier alumno de prueba residual y asegura que cada alumno registrado tenga su contrato.
   */
  function syncContractsWithStudents(studentList: Student[]) {
    if (!studentList || studentList.length === 0) return

    const existingContracts = storage.getContracts().filter(c => !c.studentName.toLowerCase().includes('alumno de ejemplo'))
    const syncedContracts: ContractConsent[] = []

    for (const student of studentList) {
      const existing = existingContracts.find(c =>
        c.studentId === student.id ||
        c.studentName.trim().toLowerCase() === student.fullName.trim().toLowerCase()
      )

      const isConsentGiven = !!student.onboarding?.consentAccepted
      let consentDate = 'Por registrar'
      if (student.onboarding?.consentTimestamp) {
        consentDate = student.onboarding.consentTimestamp.includes(',')
          ? student.onboarding.consentTimestamp.split(',')[0].trim()
          : student.onboarding.consentTimestamp.split('T')[0]
      } else if (isConsentGiven) {
        consentDate = student.createdAt?.split('T')[0] || '2026-11-18'
      }

      if (existing) {
        syncedContracts.push({
          ...existing,
          id: existing.id || `cnt-${student.id}`,
          studentId: student.id,
          studentName: student.fullName,
          status: isConsentGiven ? 'Aceptado' : existing.status,
          date: existing.status === 'Aceptado' && existing.date !== 'Por registrar' ? existing.date : consentDate,
          termsVersion: student.onboarding?.termsVersion || existing.termsVersion || 'v1.0 (2026)',
          signedCopyAvailable: isConsentGiven || existing.signedCopyAvailable || existing.status === 'Aceptado',
          notes: student.onboarding?.digitalSignature
            ? 'Consentimiento y reglamento firmado digitalmente en pantalla táctil durante onboarding.'
            : (existing.notes || 'Consentimiento informado registrado en el sistema.')
        })
      } else {
        syncedContracts.push({
          id: `cnt-${student.id}`,
          studentId: student.id,
          studentName: student.fullName,
          documentType: 'Consentimiento informado y reglamento',
          status: isConsentGiven ? 'Aceptado' : 'Pendiente',
          date: consentDate,
          termsVersion: student.onboarding?.termsVersion || 'v1.0 (2026)',
          signedCopyAvailable: isConsentGiven,
          notes: student.onboarding?.digitalSignature
            ? 'Consentimiento firmado digitalmente en pantalla táctil durante onboarding.'
            : (isConsentGiven ? 'Consentimiento aceptado al registrarse.' : 'Pendiente de firma digital o aceptación formal.')
        })
      }
    }

    contracts.value = syncedContracts
    storage.setContracts(syncedContracts)
  }

  // Subscribe to storage changes
  storage.subscribe(() => {
    reloadData()
  })

  // Toast notifications
  function showToast(message: string, type: 'success' | 'info' | 'warning' | 'error' = 'success', title?: string) {
    const id = `toast-${Date.now()}-${Math.random()}`
    toasts.value.push({ id, message, type, title })
    setTimeout(() => {
      dismissToast(id)
    }, 4500)
  }

  function dismissToast(id: string) {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  // Computed metrics for Dashboard and Cards
  const activeStudentsCount = computed(() => {
    return students.value.filter(s => s.membershipStatus === 'Activo').length
  })

  const pendingPaymentsCount = computed(() => {
    return payments.value.filter(p => p.status === 'Pendiente').length
  })

  const todaySessionsCount = computed(() => {
    const dayMap: Record<number, string> = {
      0: 'Dom',
      1: 'Lun',
      2: 'Mar',
      3: 'Mié',
      4: 'Jue',
      5: 'Vie',
      6: 'Sáb'
    }
    const currentDay = dayMap[new Date().getDay()] || 'Lun'
    const todayBlocks = schedule.value.filter(b => b.day === currentDay)
    return todayBlocks.length > 0 ? todayBlocks.length : 4
  })

  const pendingContractsCount = computed(() => {
    return contracts.value.filter(c => c.status === 'Pendiente').length
  })

  const acceptedContractsCount = computed(() => {
    return contracts.value.filter(c => c.status === 'Aceptado').length
  })

  const expiringMembershipsCount = computed(() => {
    return students.value.filter(s => s.membershipStatus === 'Por revisar').length
  })

  const paidPaymentsCount = computed(() => {
    return payments.value.filter(p => p.status === 'Pagado').length
  })

  const registeredEvaluationsCount = computed(() => {
    return evaluations.value.filter(e => e.status === 'Registrada').length
  })

  const pendingEvaluationsCount = computed(() => {
    return evaluations.value.filter(e => e.status === 'Pendiente').length
  })

  // Action methods with notifications
  async function createStudent(dto: CreateStudentDTO) {
    loading.value = true
    try {
      const student = await studentService.create(dto)
      if (isStudentsApiMode()) {
        students.value = [student, ...students.value.filter(s => s.id !== student.id)]
      }
      syncContractsWithStudents(students.value)
      showToast(`Alumno "${student.fullName}" registrado exitosamente.`, 'success', 'Registro completado')
      return student
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al registrar alumno'
      showToast(msg, 'error', 'Error')
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateStudent(id: string, updates: Partial<Student>) {
    try {
      const updated = await studentService.update(id, updates)
      if (updated && isStudentsApiMode()) {
        const idx = students.value.findIndex(s => s.id === id)
        if (idx !== -1) {
          students.value[idx] = updated
          students.value = [...students.value]
        }
      }
      showToast('Información del alumno actualizada.', 'success')
      return updated
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al actualizar alumno'
      showToast(msg, 'error')
      throw err
    }
  }

  async function deleteStudent(id: string) {
    try {
      const student = students.value.find(s => s.id === id)
      await studentService.delete(id)
      if (isStudentsApiMode()) {
        students.value = students.value.filter(s => s.id !== id)
      }
      syncContractsWithStudents(students.value)
      showToast(`Alumno "${student?.fullName || 'seleccionado'}" archivado exitosamente.`, 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar alumno'
      showToast(msg, 'error')
      throw err
    }
  }

  async function registerPayment(dto: RegisterPaymentDTO) {
    try {
      const p = await paymentService.registerPayment(dto)
      if (isPaymentsApiMode()) {
        await loadPayments()
      }
      showToast(`Pago registrado exitosamente para ${p?.studentName}.`, 'success', 'Pago Registrado')
      return p
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al registrar pago'
      showToast(msg, 'error')
      throw err
    }
  }

  async function simulatePaymentLink(studentId: string, paymentId?: string) {
    try {
      const res = await paymentService.sendPaymentLinkSimulation(studentId, paymentId)
      showToast(res.message, 'info', 'Link SumUp')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al generar enlace'
      showToast(msg, 'error')
      throw err
    }
  }

  async function simulateReminder(paymentId: string) {
    try {
      const res = await paymentService.sendReminderSimulation(paymentId)
      if (isPaymentsApiMode()) {
        await loadPayments()
      }
      showToast(res.message, 'info', 'Aviso de Cobro')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al enviar recordatorio'
      showToast(msg, 'error')
      throw err
    }
  }

  async function simulateReceipt(paymentId: string) {
    try {
      const res = await paymentService.sendReceiptSimulation(paymentId)
      if (isPaymentsApiMode()) {
        await loadPayments()
      }
      showToast(res.message, 'success', 'Comprobante')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al enviar comprobante'
      showToast(msg, 'error')
      throw err
    }
  }

  // ==========================================
  // SCHEDULE & AVAILABILITY STORE ACTIONS
  // ==========================================

  async function saveAvailability(newAvailability: DayAvailability[]) {
    try {
      const res = await scheduleService.saveAvailability(newAvailability)
      if (!res.success) {
        showToast(res.message, 'warning', 'Conflictos de disponibilidad')
        return res
      }
      if (isScheduleApiMode()) {
        await loadSchedule()
      } else {
        await reloadData()
      }
      showToast(res.message, 'success', 'Disponibilidad Guardada')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al guardar disponibilidad'
      showToast(msg, 'error')
      throw err
    }
  }

  async function assignSession(dto: AssignSessionDTO) {
    try {
      const res = await scheduleService.assignSession(dto)
      if (!res.success) {
        showToast(res.message, 'warning', 'No se pudo programar')
        return res
      }
      await refreshScheduleState()
      showToast(res.message, 'success', 'Agenda Actualizada')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al programar sesión'
      showToast(msg, 'error')
      throw err
    }
  }

  // Backward-compatible assignStudentToSchedule
  async function assignStudentToSchedule(dto: AssignStudentDTO) {
    try {
      const res = await scheduleService.assignStudent(dto)
      if (!res.success) {
        showToast(res.message, 'warning', 'No se pudo asignar')
        return res
      }
      await refreshScheduleState()
      showToast(res.message, 'success', 'Agenda Actualizada')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al asignar alumno'
      showToast(msg, 'error')
      throw err
    }
  }

  async function updateSession(blockId: string, updates: Partial<ScheduleBlock>) {
    try {
      const res = await scheduleService.updateSession(blockId, updates)
      if (!res.success) {
        showToast(res.message, 'warning', 'No se pudo actualizar')
        return res
      }
      await refreshScheduleState()
      showToast(res.message, 'success', 'Sesión Actualizada')
      return res
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al actualizar sesión'
      showToast(msg, 'error')
      throw err
    }
  }

  async function removeStudentFromSchedule(blockId: string, studentId: string) {
    try {
      await scheduleService.removeStudentFromBlock(blockId, studentId)
      const found = schedule.value.find(b => b.id === blockId)
      if (found) {
        found.studentIds = found.studentIds.filter(id => id !== studentId)
      }
      await refreshScheduleState()
      showToast('Alumno retirado de la sesión.', 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al retirar alumno'
      showToast(msg, 'error')
      throw err
    }
  }

  async function deleteBlock(blockId: string) {
    try {
      await scheduleService.deleteBlock(blockId)
      // Actualización reactiva local inmediata
      schedule.value = schedule.value.filter(b => b.id !== blockId)
      weekOccurrences.value = weekOccurrences.value.filter(o => o.blockId !== blockId)

      await refreshScheduleState()
      showToast('Sesión eliminada de la agenda.', 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar sesión'
      showToast(msg, 'error')
      throw err
    }
  }

  async function createException(dto: CreateExceptionDTO) {
    try {
      const exc = await scheduleService.createException(dto)
      // Si fue cancelación puntual, actualizar reactivamente de inmediato
      if (dto.type === 'cancelled') {
        weekOccurrences.value = weekOccurrences.value.filter(
          o => !(o.blockId === dto.blockId && o.date === dto.date)
        )
      }

      await refreshScheduleState()
      const actionText = dto.type === 'cancelled' ? 'Sesión cancelada para esta fecha.' : 'Sesión reagendada para esta semana.'
      showToast(actionText, 'info', 'Excepción de Calendario')
      return exc
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al aplicar excepción'
      showToast(msg, 'error')
      throw err
    }
  }

  async function deleteException(exceptionId: string, blockId?: string) {
    try {
      await scheduleService.deleteException(exceptionId, blockId)
      await refreshScheduleState()
      showToast('Excepción eliminada; se restablece la sesión recurrente.', 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar excepción'
      showToast(msg, 'error')
      throw err
    }
  }

  function getWeekOccurrences(weekDays: WeekDayInfo[]): SessionOccurrence[] {
    if (isScheduleApiMode()) {
      return weekOccurrences.value
    }
    return scheduleService.getWeekOccurrences(weekDays, schedule.value, exceptions.value)
  }

  // ==========================================
  // EVALUATIONS & CONTRACTS
  // ==========================================

  async function createEvaluation(dto: CreateEvaluationDTO) {
    try {
      const student = students.value.find(s => s.id === dto.studentId)
      if (student && !dto.studentName) {
        dto.studentName = student.fullName
      }
      const ev = await evaluationService.create(dto)
      evaluations.value = [ev, ...evaluations.value.filter(e => e.id !== ev.id)]
      showToast(`Evaluación registrada para ${ev.studentName}.`, 'success', 'Evaluación Guardada')
      return ev
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al registrar evaluación'
      showToast(msg, 'error')
      throw err
    }
  }

  async function updateEvaluation(id: string, updates: Partial<Evaluation>) {
    try {
      const ev = await evaluationService.update(id, updates)
      if (ev) {
        const idx = evaluations.value.findIndex(e => e.id === id)
        if (idx !== -1) {
          evaluations.value[idx] = ev
          evaluations.value = [...evaluations.value]
        }
      }
      showToast('Evaluación actualizada.', 'success')
      return ev
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al actualizar evaluación'
      showToast(msg, 'error')
      throw err
    }
  }

  async function deleteEvaluation(id: string) {
    try {
      await evaluationService.delete(id)
      evaluations.value = evaluations.value.filter(e => e.id !== id)
      showToast('Evaluación eliminada.', 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar evaluación'
      showToast(msg, 'error')
      throw err
    }
  }

  async function updateContractStatus(id: string, status: ContractStatus) {
    try {
      const c = await contractService.updateStatus(id, status)
      if (c) {
        const idx = contracts.value.findIndex(item => item.id === id)
        if (idx !== -1) {
          contracts.value[idx] = c
          contracts.value = [...contracts.value]
        }
      }
      showToast(`Estado de contrato actualizado a "${status}" para ${c?.studentName}.`, 'success', 'Contratos')
      return c
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al actualizar contrato'
      showToast(msg, 'error')
      throw err
    }
  }

  async function saveTerms(version: string, title: string, content: string) {
    try {
      const newT = await contractService.updateCurrentTerms(version, title, content)
      showToast(`Versión de términos "${version}" guardada y establecida como vigente.`, 'success', 'Términos Actualizados')
      return newT
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al guardar términos'
      showToast(msg, 'error')
      throw err
    }
  }

  function resetDemoData() {
    storage.resetAll()
    reloadData()
    showToast('Datos de demostración restablecidos a los valores iniciales de Figma.', 'info', 'Datos Reiniciados')
  }

  // Initialize
  reloadData()
  loadStudents().catch(() => {
    // Los errores quedan reflejados en studentsError
  })
  if (isScheduleApiMode()) {
    loadSchedule().catch(() => {
      // Los errores quedan reflejados en scheduleError
    })
  }
  if (isPaymentsApiMode()) {
    loadPayments().catch(() => {
      // Los errores quedan reflejados en paymentsError
    })
  }

  return {
    students,
    loadingStudents,
    studentsError,
    dataSource,
    payments,
    loadingPayments,
    paymentsError,
    paymentsDataSource,
    isPaymentsApiMode,
    financialSummary,
    sumupPermanentLink,
    updateSumupPermanentLink,
    loadPayments,
    schedule,
    availability,
    exceptions,
    weekOccurrences,
    currentWeekRange,
    refreshScheduleState,
    loadingSchedule,
    scheduleError,
    scheduleDataSource,
    isScheduleApiMode,
    loadSchedule,
    fetchWeekOccurrences,
    evaluations,
    contracts,
    terms,
    loading,
    toasts,
    globalSearch,
    activeStudentsCount,
    pendingPaymentsCount,
    todaySessionsCount,
    pendingContractsCount,
    acceptedContractsCount,
    expiringMembershipsCount,
    paidPaymentsCount,
    registeredEvaluationsCount,
    pendingEvaluationsCount,
    showToast,
    dismissToast,
    createStudent,
    updateStudent,
    deleteStudent,
    registerPayment,
    simulatePaymentLink,
    simulateReminder,
    simulateReceipt,
    saveAvailability,
    assignSession,
    assignStudentToSchedule,
    updateSession,
    removeStudentFromSchedule,
    deleteBlock,
    createException,
    deleteException,
    getWeekOccurrences,
    createEvaluation,
    updateEvaluation,
    deleteEvaluation,
    updateContractStatus,
    saveTerms,
    resetDemoData,
    reloadData,
    loadStudents,
    syncContractsWithStudents
  }
})
