import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import Login from '../../components/Login.vue'
import MessageBox from '../../components/MessageBox.vue'

const mockLogin = vi.fn()
const mockAuthStore = {
  login: mockLogin,
}

vi.mock('../../stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))

const mockRouterPush = vi.fn()
vi.mock('vue-router', async () => {
  const actual = await vi.importActual('vue-router')
  return {
    ...actual,
    useRouter: () => ({
      push: mockRouterPush,
    }),
  }
})

describe('Login.vue', () => {
  let router

  beforeEach(() => {
    setActivePinia(createPinia())

    router = createRouter({
      history: createWebHistory(),
      routes: [
        {
          path: '/',
          name: 'home',
          component: Login,
        },
        {
          path: '/profile',
          name: 'profile',
          component: { template: '<div>Profile Page</div>' },
        },
        {
          path: '/signup',
          name: 'signup',
          component: { template: '<div>Signup Page</div>' },
        },
      ],
    })

    vi.clearAllMocks()
    mockLogin.mockReset()
    mockRouterPush.mockReset()

    router.push('/')
    router.isReady()
  })

  it('renders login form with all elements', () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
    })

    expect(wrapper.find('h2').text()).toBe('Login')
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)

    const signupLink = wrapper.find('p.text-gray-600')
    expect(signupLink.text()).toContain("Don't have an account?")
    expect(signupLink.html()).toContain('to="/signup"')
  })

  it('updates v-model values when typing in inputs', async () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
    })

    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('password123')

    expect(wrapper.vm.username).toBe('testuser')
    expect(wrapper.vm.password).toBe('password123')
  })

  it('calls login method and redirects on successful login', async () => {
    mockLogin.mockResolvedValueOnce()

    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
    })

    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('password123')

    await wrapper.find('form').trigger('submit')

    expect(mockLogin).toHaveBeenCalledWith('testuser', 'password123')
    expect(mockRouterPush).toHaveBeenCalledWith('/profile')
  })

  it('shows error message on login failure', async () => {
    const errorMessage = 'Invalid credentials'
    mockLogin.mockRejectedValueOnce(new Error(errorMessage))

    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
    })

    await wrapper.find('input[type="text"]').setValue('testuser')
    await wrapper.find('input[type="password"]').setValue('wrongpass')

    await wrapper.find('form').trigger('submit')

    expect(wrapper.vm.messageBoxVisible).toBe(true)
    expect(wrapper.vm.messageTitle).toBe('Login Failed')
    expect(wrapper.vm.messageText).toBe(errorMessage)
    expect(wrapper.vm.messageType).toBe('error')
  })

  it('hides error message after 5 seconds', async () => {
    vi.useFakeTimers()

    mockLogin.mockRejectedValueOnce(new Error('Invalid credentials'))

    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
    })

    await wrapper.find('form').trigger('submit')

    expect(wrapper.vm.messageBoxVisible).toBe(true)

    await vi.advanceTimersByTime(5000)

    expect(wrapper.vm.messageBoxVisible).toBe(false)

    vi.useRealTimers()
  })

  it('requires username and password fields', async () => {
    const wrapper = mount(Login, {
      global: {
        plugins: [router],
        components: {
          MessageBox,
        },
        stubs: {
          RouterLink: true,
        },
      },
      attachTo: document.body,
    })

    const submitButton = wrapper.find('button[type="submit"]')

    await submitButton.trigger('click')

    expect(mockLogin).not.toHaveBeenCalled()

    const usernameInput = wrapper.find('input[type="text"]')
    const passwordInput = wrapper.find('input[type="password"]')
    expect(usernameInput.attributes('required')).toBeDefined()
    expect(passwordInput.attributes('required')).toBeDefined()
  })
})
