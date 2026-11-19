import { createRouter, createWebHistory } from 'vue-router'
import InicioView from '@/views/InicioView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/inicio'
    },
    {
      path: '/inicio',
      name: 'inicio',
      component: InicioView,
      meta: { title: 'Inicio — Panel del entrenador' }
    },
    {
      path: '/agenda',
      name: 'agenda',
      component: () => import('@/views/AgendaView.vue'),
      meta: { title: 'Agenda semanal' }
    },
    {
      path: '/alumnos',
      name: 'alumnos',
      component: () => import('@/views/AlumnosView.vue'),
      meta: { title: 'Alumnos' }
    },
    {
      path: '/alumnos/:id',
      name: 'alumno-detalle',
      component: () => import('@/views/FichaAlumnoView.vue'),
      meta: { title: 'Ficha del alumno' },
      props: true
    },
    {
      path: '/membresias-pagos',
      name: 'membresias-pagos',
      component: () => import('@/views/MembresiasPagosView.vue'),
      meta: { title: 'Membresías y pagos' }
    },
    {
      path: '/evaluaciones',
      name: 'evaluaciones',
      component: () => import('@/views/EvaluacionesView.vue'),
      meta: { title: 'Evaluaciones' }
    },
    {
      path: '/contratos',
      name: 'contratos',
      component: () => import('@/views/ContratosView.vue'),
      meta: { title: 'Contratos y consentimientos' }
    },
    {
      path: '/registro-alumno',
      name: 'registro-alumno',
      component: () => import('@/views/PublicOnboardingView.vue'),
      meta: { title: 'Ficha de Ingreso y Consentimiento', public: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { title: 'Acceso Coach Pablo', public: true }
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/inicio'
    }
  ]
})

router.beforeEach((to, _from, next) => {
  if (to.meta?.title) {
    document.title = `${to.meta.title} | INVERFIT`
  }

  const isPublic = !!to.meta?.public
  const isAuth = typeof window !== 'undefined' && !!localStorage.getItem('inverfit_auth_user')

  if (to.path === '/login' && isAuth) {
    return next('/inicio')
  }

  if (!isPublic && !isAuth) {
    return next({ path: '/login', query: to.path !== '/inicio' ? { redirect: to.fullPath } : undefined })
  }

  next()
})

export default router
