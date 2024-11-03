import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import Navbar from '@/components/Navbar.vue'
import MessageBox from '@/components/MessageBox.vue'
import { createPinia, setActivePinia } from 'pinia'

const fetchMock = vi.fn()
globalThis.fetch = fetchMock

const mockAuthStore = {
  isAuthenticated: false,
  logout: vi.fn(),
}

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: {} },
    { path: '/login', component: {} },
    { path: '/profile', component: {} },
    { path: '/:username', component: {} },
  ],
})

describe('Navbar', () => {
  let wrapper

  beforeEach(async () => {
    setActivePinia(createPinia())
    mockAuthStore.isAuthenticated = false
    vi.clearAllMocks()
    fetchMock.mockClear()

    await router.push('/')

    wrapper = mount(Navbar, {
      global: {
        plugins: [router],
        components: { MessageBox },
        stubs: ['router-link'],
      },
    })
  })

  it('performs search and navigates to user profile', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({}),
    })

    const searchInput = wrapper.find('input')
    await searchInput.setValue('testuser')
    await wrapper.find('form').trigger('submit')

    expect(fetchMock).toHaveBeenCalledWith(
      'http://localhost:8188/api/v1/users/testuser',
    )
    expect(searchInput.element.value).toBe('')
  })

  it('shows error message when user is not found', async () => {
    fetchMock.mockResolvedValueOnce({
      ok: false,
    })

    await wrapper.find('input').setValue('nonexistentuser')
    await wrapper.find('form').trigger('submit')

    await wrapper.vm.$nextTick()
    expect(wrapper.findComponent(MessageBox).exists()).toBe(true)
    expect(wrapper.vm.messageText).includes('does not exist')
  })

  it('shows error message on network error', async () => {
    fetchMock.mockRejectedValueOnce(new Error('Network error'))

    await wrapper.find('input').setValue('testuser')
    await wrapper.find('form').trigger('submit')

    await wrapper.vm.$nextTick()
    expect(wrapper.findComponent(MessageBox).exists()).toBe(true)
    expect(wrapper.vm.messageText).includes('error')
  })

  it('handles logout when authenticated', async () => {
    mockAuthStore.isAuthenticated = true
    wrapper = mount(Navbar, {
      global: {
        plugins: [router],
        components: { MessageBox },
        stubs: ['router-link'],
      },
    })

    await router.isReady()

    const menuButton = wrapper.findAll('button')[0]
    await menuButton.trigger('click')

    const logoutButton = wrapper
      .findAll('button')
      .find(b => b.text().includes('Logout'))
    await logoutButton.trigger('click')

    await router.isReady()

    expect(mockAuthStore.logout).toHaveBeenCalled()
    expect(router.currentRoute.value.path).toBe('/')
  })
})
