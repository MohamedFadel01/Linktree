import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import Hero from '../../components/Hero.vue'
import { createPinia, setActivePinia } from 'pinia'

const mockAuthStore = {
  isAuthenticated: false,
}

vi.mock('../../stores/auth', () => ({
  useAuthStore: () => mockAuthStore,
}))

describe('Hero.vue', () => {
  let router

  beforeEach(() => {
    mockAuthStore.isAuthenticated = false

    setActivePinia(createPinia())

    router = createRouter({
      history: createWebHistory(),
      routes: [
        {
          path: '/',
          name: 'home',
          component: Hero,
        },
        {
          path: '/signup',
          name: 'signup',
          component: { template: '<div>Signup Page</div>' },
        },
        {
          path: '/profile',
          name: 'profile',
          component: { template: '<div>Profile Page</div>' },
        },
      ],
    })

    router.push('/')
    router.isReady()
  })

  it('renders the component with correct text content', () => {
    const wrapper = mount(Hero, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.text()).toContain('Manage All Your Links in One Place')
    expect(wrapper.text()).toContain(
      'Linktree allows you to share multiple links',
    )
  })

  it('shows "Get Started" button when user is not authenticated', () => {
    const wrapper = mount(Hero, {
      global: {
        plugins: [router],
      },
    })

    const button = wrapper.find('a')
    expect(button.text()).toBe('Get Started')
    expect(button.attributes('href')).toBe('/signup')
  })

  it('shows "Go to Profile" button when user is authenticated', async () => {
    mockAuthStore.isAuthenticated = true

    const wrapper = mount(Hero, {
      global: {
        plugins: [router],
      },
    })

    const button = wrapper.find('a')
    expect(button.text()).toBe('Go to Profile')
    expect(button.attributes('href')).toBe('/profile')
  })

  it('handles button click and navigates to correct route when not authenticated', async () => {
    const wrapper = mount(Hero, {
      global: {
        plugins: [router],
      },
    })

    await wrapper.find('a').trigger('click')
    await router.isReady()

    expect(wrapper.find('a').attributes('href')).toBe('/signup')
  })

  it('handles button click and navigates to correct route when authenticated', async () => {
    mockAuthStore.isAuthenticated = true

    const wrapper = mount(Hero, {
      global: {
        plugins: [router],
      },
    })

    await wrapper.find('a').trigger('click')
    await router.isReady()

    expect(wrapper.find('a').attributes('href')).toBe('/profile')
  })
})
