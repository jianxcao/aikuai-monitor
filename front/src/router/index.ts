import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LanViewVue from '@/views/LanView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/lan',
      name: 'Lan',
      component: LanViewVue
    },
    {
      path: '/',
      name: 'home',
      component: HomeView
    }
  ]
})

export default router
