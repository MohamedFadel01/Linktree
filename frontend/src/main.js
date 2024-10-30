import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import App from './App.vue'
import './assets/main.css'
import { useAuthStore } from './stores/auth.js'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('./components/Hero.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('./components/Login.vue'),
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('./components/Signup.vue'),
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('./components/Profile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/edit-profile',
      name: 'edit-profile',
      component: () => import('./components/EditProfile.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/:username',
      name: 'user-profile',
      component: () => import('./components/Profile.vue'),
    },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (
    to.matched.some(record => record.meta.requiresAuth) &&
    !authStore.isAuthenticated
  ) {
    next('/login')
  } else {
    next()
  }
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.mount('#app')
