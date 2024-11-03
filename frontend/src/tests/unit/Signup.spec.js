import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import Signup from '../../components/Signup.vue'
import MessageBox from '../../components/MessageBox.vue'
import { flushPromises } from '@vue/test-utils'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: { template: '<div>Home</div>' },
    },
    {
      path: '/login',
      name: 'login',
      component: { template: '<div>Login</div>' },
    },
  ],
})

const mockAuthStore = {
  signup: vi.fn(),
}

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))

describe('Signup.vue', () => {
  beforeEach(async () => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    await router.push('/')
    await router.isReady()
  })

  it('renders the signup form correctly', () => {
    const wrapper = mount(Signup, {
      global: {
        plugins: [router],
        components: { MessageBox },
      },
    })

    expect(wrapper.find('h2').text()).toBe('Sign Up')
    expect(wrapper.findAll('input[type="text"]')).toHaveLength(2)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
    expect(wrapper.find('textarea').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('updates form data when inputs change', async () => {
    const wrapper = mount(Signup, {
      global: {
        plugins: [router],
        components: { MessageBox },
      },
    })

    const [fullNameInput, usernameInput] = wrapper.findAll('input[type="text"]')
    const bioInput = wrapper.find('textarea')
    const passwordInput = wrapper.find('input[type="password"]')

    await fullNameInput.setValue('John Doe')
    await usernameInput.setValue('johndoe')
    await bioInput.setValue('My bio')
    await passwordInput.setValue('password123')

    expect(wrapper.vm.fullName).toBe('John Doe')
    expect(wrapper.vm.username).toBe('johndoe')
    expect(wrapper.vm.bio).toBe('My bio')
    expect(wrapper.vm.password).toBe('password123')
  })

  it('calls signup method and redirects on successful submission', async () => {
    const wrapper = mount(Signup, {
      global: {
        plugins: [router],
        components: { MessageBox },
      },
    })

    mockAuthStore.signup.mockResolvedValueOnce()

    const [fullNameInput, usernameInput] = wrapper.findAll('input[type="text"]')
    const bioInput = wrapper.find('textarea')
    const passwordInput = wrapper.find('input[type="password"]')

    await fullNameInput.setValue('John Doe')
    await usernameInput.setValue('johndoe')
    await bioInput.setValue('My bio')
    await passwordInput.setValue('password123')

    await wrapper.find('form').trigger('submit')

    expect(mockAuthStore.signup).toHaveBeenCalledWith({
      full_name: 'John Doe',
      username: 'johndoe',
      bio: 'My bio',
      password: 'password123',
    })

    await flushPromises()

    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('shows error message on failed signup', async () => {
    const wrapper = mount(Signup, {
      global: {
        plugins: [router],
        components: { MessageBox },
      },
    })

    const error = new Error('Signup failed')
    mockAuthStore.signup.mockRejectedValueOnce(error)

    await wrapper.find('form').trigger('submit')

    expect(wrapper.vm.messageBoxVisible).toBe(true)
    expect(wrapper.vm.messageTitle).toBe('Signup Failed')
    expect(wrapper.vm.messageType).toBe('error')

    expect(wrapper.findComponent(MessageBox).exists()).toBe(true)
  })

  it('error message includes MessageBox component when visible', async () => {
    const wrapper = mount(Signup, {
      global: {
        plugins: [router],
        components: { MessageBox },
      },
    })

    wrapper.vm.messageBoxVisible = true
    wrapper.vm.messageTitle = 'Test Title'
    wrapper.vm.messageText = 'Test Message'
    wrapper.vm.messageType = 'error'

    await wrapper.vm.$nextTick()

    const messageBox = wrapper.findComponent(MessageBox)
    expect(messageBox.exists()).toBe(true)
    expect(messageBox.props('title')).toBe('Test Title')
    expect(messageBox.props('message')).toBe('Test Message')
    expect(messageBox.props('type')).toBe('error')
  })
})
