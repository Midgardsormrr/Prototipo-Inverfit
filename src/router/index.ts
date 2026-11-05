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
      path: '/:pathMatch(.*)*',
      redirect: '/inicio'
    }
  ]
})

router.beforeEach((to, _from, next) => {
  if (to.meta?.title) {
    document.title = `${to.meta.title} | INVERFIT`
  }
  next()
})

export default router
