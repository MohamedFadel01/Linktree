import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import { createTestingPinia } from '@pinia/testing'
import Profile from '../../components/Profile.vue'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'

vi.mock('../../components/ProfileInfo.vue', () => ({
  default: {
    name: 'ProfileInfo',
    props: ['fullName', 'username', 'bio'],
    template: '<div data-testid="profile-info"></div>',
  },
}))

vi.mock('../../components/LinksList.vue', () => ({
  default: {
    name: 'LinksList',
    props: ['links'],
    template: '<div data-testid="links-list"></div>',
  },
}))

vi.mock('../../components/LinkAnalytics.vue', () => ({
  default: {
    name: 'LinkAnalytics',
    props: ['analytics'],
    template: '<div data-testid="link-analytics"></div>',
  },
}))

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/profile/:username',
      name: 'Profile',
      component: Profile,
    },
    {
      path: '/edit-profile',
      name: 'EditProfile',
      component: { template: '<div></div>' },
    },
    {
      path: '/',
      name: 'Home',
      component: { template: '<div></div>' },
    },
  ],
})

describe('Profile', () => {
  let wrapper
  let userStore
  let authStore
  let pinia

  beforeEach(async () => {
    pinia = createTestingPinia({
      createSpy: vi.fn,
      initialState: {
        user: { userData: null, error: null },
        auth: { username: null },
      },
    })

    userStore = useUserStore(pinia)
    authStore = useAuthStore(pinia)

    await router.push('/')
    await router.isReady()

    wrapper = mount(Profile, {
      global: {
        plugins: [pinia, router],
        stubs: {
          RouterLink: {
            template: '<a><slot></slot></a>',
          },
        },
      },
    })
  })

  it('shows loading state when no userData and no error', () => {
    userStore.userData = null
    userStore.error = null

    expect(wrapper.text()).toContain('Loading...')
  })

  it('displays error message when there is an error', async () => {
    userStore.error = 'Failed to fetch profile'
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Failed to fetch profile')
  })

  it('shows ProfileInfo and LinksList when userData is available', async () => {
    userStore.userData = {
      username: 'testuser',
      full_name: 'Test User',
      bio: 'Test bio',
      links: [
        {
          url: 'https://example.com',
          analytics: {
            click_count: 0,
            last_clicked: null,
          },
        },
      ],
    }
    await wrapper.vm.$nextTick()

    expect(wrapper.findComponent({ name: 'ProfileInfo' }).exists()).toBe(true)
    expect(wrapper.findComponent({ name: 'LinksList' }).exists()).toBe(true)
  })

  it('shows edit button when user is viewing their own profile', async () => {
    userStore.userData = {
      username: 'testuser',
    }
    authStore.username = 'testuser'
    await wrapper.vm.$nextTick()

    const editButton = wrapper.get('button')
    expect(editButton.exists()).toBe(true)
    expect(editButton.text()).toBe('Edit Profile')
  })

  it("hides edit button when user is viewing someone else's profile", async () => {
    userStore.userData = {
      username: 'otheruser',
    }
    authStore.username = 'testuser'
    await wrapper.vm.$nextTick()

    const editButton = wrapper.find('button')
    expect(editButton.exists()).toBe(false)
  })

  it('calls fetchUserProfile on mount with username from route params', async () => {
    const pinia = createTestingPinia({ createSpy: vi.fn })
    const store = useUserStore(pinia)

    await router.push('/profile/testuser')
    await router.isReady()

    mount(Profile, {
      global: {
        plugins: [router, pinia],
        stubs: {
          RouterLink: {
            template: '<a><slot></slot></a>',
          },
        },
      },
    })

    await vi.waitFor(() => {
      expect(store.fetchUserProfile).toHaveBeenCalledWith('testuser')
    })
  })

  it('redirects to home page on fetch error', async () => {
    const pinia = createTestingPinia({ createSpy: vi.fn })
    const store = useUserStore(pinia)
    store.fetchUserProfile.mockRejectedValueOnce(new Error('Fetch failed'))

    await router.push('/profile/testuser')
    await router.isReady()

    mount(Profile, {
      global: {
        plugins: [router, pinia],
        stubs: {
          RouterLink: {
            template: '<a><slot></slot></a>',
          },
        },
      },
    })

    await vi.waitFor(() => {
      expect(router.currentRoute.value.path).toBe('/')
    })
  })

  it('updates profile when route params change', async () => {
    const pinia = createTestingPinia({ createSpy: vi.fn })
    const store = useUserStore(pinia)

    const wrapper = mount(Profile, {
      global: {
        plugins: [router, pinia],
        stubs: {
          RouterLink: {
            template: '<a><slot></slot></a>',
          },
        },
      },
    })

    await router.push('/profile/user1')
    await router.isReady()
    await wrapper.vm.$nextTick()

    await router.push('/profile/user2')
    await router.isReady()
    await wrapper.vm.$nextTick()

    await vi.waitFor(() => {
      expect(store.fetchUserProfile).toHaveBeenCalledWith('user2')
    })
  })
})
