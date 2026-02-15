import { createRouter, createWebHistory } from 'vue-router'
import MainLayout from '@/views/MainLayout.vue' // Если создали в layouts
import LoginView from '@/views/LoginView.vue'
import MyDeposits from '@/views/MyDeposits.vue'
import MainPage from "@/views/MainPage.vue";
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: MainPage,
    },
    {
      path: '/dashboard',
      name: 'MyAccount',
      component: MainLayout
      ,
      children: [

        {
          path: 'deposits',
          name: 'MyDeposits',
          component: MyDeposits
        },
        {
          path: 'open',
          name: 'open-deposit',
          component: () => import('@/views/OpenDeposit.vue')
        }

      ]
    }
  ]
})

export default router