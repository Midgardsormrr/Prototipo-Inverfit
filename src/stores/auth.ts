import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface AuthUser {
  id: string
  name: string
  email: string
  role: string
  avatarLetter: string
}

const STORAGE_KEY = 'inverfit_auth_user'

export const useAuthStore = defineStore('auth', () => {
  // Inicialización desde localStorage
  const storedRaw = typeof window !== 'undefined' ? localStorage.getItem(STORAGE_KEY) : null
  let initialUser: AuthUser | null = null
  if (storedRaw) {
    try {
      initialUser = JSON.parse(storedRaw)
    } catch {
      initialUser = null
    }
  }

  const user = ref<AuthUser | null>(initialUser)
  const isAuthenticated = ref<boolean>(!!initialUser)
  const authError = ref<string | null>(null)

  function login(email: string, password: string): boolean {
    authError.value = null

    const cleanEmail = email.trim().toLowerCase()
    const cleanPass = password.trim()

    // Credenciales oficiales de Coach Pablo
    const isPablo =
      cleanEmail === 'pablo@inverfit.cl' ||
      cleanEmail === 'coach@inverfit.cl' ||
      cleanEmail === 'admin@inverfit.cl' ||
      cleanEmail === 'pablo'

    const isPassValid =
      cleanPass === 'inverfit2026' ||
      cleanPass === 'pablo2026' ||
      cleanPass === '123456'

    if (isPablo && isPassValid) {
      const loggedUser: AuthUser = {
        id: 'usr-pablo',
        name: 'Pablo',
        email: cleanEmail.includes('@') ? cleanEmail : 'pablo@inverfit.cl',
        role: 'Entrenador Principal',
        avatarLetter: 'P'
      }
      user.value = loggedUser
      isAuthenticated.value = true
      localStorage.setItem(STORAGE_KEY, JSON.stringify(loggedUser))
      return true
    }

    authError.value = 'Correo electrónico o contraseña incorrectos. Verifica tus credenciales de acceso.'
    return false
  }

  function logout() {
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem(STORAGE_KEY)
  }

  return {
    user,
    isAuthenticated,
    authError,
    login,
    logout
  }
})
