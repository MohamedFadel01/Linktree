import { describe, it, beforeEach, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '@/App.vue'
import Navbar from '@/components/Navbar.vue'
import { createPinia } from 'pinia'
import { setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'
import { createRouter, createWebHistory } from 'vue-router'

vi.mock('axios')

const router = createRouter({
  history: createWebHistory(),
  routes: [],
})

describe('App.vue', () => {
  let authStore

  beforeEach(() => {
    setActivePinia(createPinia())
    authStore = useAuthStore()
    axios.defaults.headers.common['Authorization'] = undefined
  })

  it('renders the Navbar component', () => {
    const wrapper = mount(App, {
      global: {
        plugins: [router],
      },
    })
    expect(wrapper.findComponent(Navbar).exists()).toBe(true)
  })

  it('sets Axios authorization header on mount if token exists', () => {
    authStore.token = 'mocked_token'

    mount(App, {
      global: {
        plugins: [router],
      },
    })

    expect(axios.defaults.headers.common['Authorization']).toBe(
      'Bearer mocked_token',
    )
  })

  it('does not set Axios authorization header if token does not exist', () => {
    authStore.token = ''

    mount(App, {
      global: {
        plugins: [router],
      },
    })

    expect(axios.defaults.headers.common['Authorization']).toBeUndefined()
  })
})
