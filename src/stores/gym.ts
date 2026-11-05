import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  storage,
  studentService,
  paymentService,
  scheduleService,
  evaluationService,
  contractService
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
  ContractStatus
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
  const payments = ref<Payment[]>([])
  const schedule = ref<ScheduleBlock[]>([])
  const availability = ref<DayAvailability[]>([])
  const exceptions = ref<ScheduleException[]>([])
  const evaluations = ref<Evaluation[]>([])
  const contracts = ref<ContractConsent[]>([])
  const terms = ref<TermsVersion[]>([])
  const loading = ref(false)
  const toasts = ref<ToastNotification[]>([])
  const globalSearch = ref('')

  // Load from services
  async function reloadData() {
    students.value = storage.getStudents()
    payments.value = storage.getPayments()
    schedule.value = storage.getSchedule()
    availability.value = storage.getAvailability()
    exceptions.value = storage.getExceptions()
    evaluations.value = storage.getEvaluations()
    contracts.value = storage.getContracts()
    terms.value = storage.getTerms()
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
      showToast(`Alumno "${student?.fullName || 'seleccionado'}" eliminado.`, 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar alumno'
      showToast(msg, 'error')
      throw err
    }
  }

  async function registerPayment(dto: RegisterPaymentDTO) {
    try {
      const p = await paymentService.registerPayment(dto)
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
      showToast(res.message, 'info', 'Simulación: Enlace Generado')
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
      showToast(res.message, 'info', 'Simulación: Recordatorio')
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
      showToast(res.message, 'success', 'Simulación: Comprobante')
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
      await reloadData()
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
      await reloadData()
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
      await reloadData()
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
      await reloadData()
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
      await reloadData()
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
      await reloadData()
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
      await reloadData()
      const actionText = dto.type === 'cancelled' ? 'Sesión cancelada para esta fecha.' : 'Sesión reagendada para esta semana.'
      showToast(actionText, 'info', 'Excepción de Calendario')
      return exc
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al aplicar excepción'
      showToast(msg, 'error')
      throw err
    }
  }

  async function deleteException(exceptionId: string) {
    try {
      await scheduleService.deleteException(exceptionId)
      await reloadData()
      showToast('Excepción eliminada; se restablece la sesión recurrente.', 'info')
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Error al eliminar excepción'
      showToast(msg, 'error')
      throw err
    }
  }

  function getWeekOccurrences(weekDays: WeekDayInfo[]): SessionOccurrence[] {
    return scheduleService.getWeekOccurrences(weekDays, schedule.value, exceptions.value)
  }

  // ==========================================
  // EVALUATIONS & CONTRACTS
  // ==========================================

  async function createEvaluation(dto: CreateEvaluationDTO) {
    try {
      const ev = await evaluationService.create(dto)
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

  return {
    students,
    payments,
    schedule,
    availability,
    exceptions,
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
    reloadData
  }
})
