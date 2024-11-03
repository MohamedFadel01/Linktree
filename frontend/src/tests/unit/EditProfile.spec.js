import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import EditProfile from '@/components/EditProfile.vue'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'

const routerPushMock = vi.fn()
const useRouterMock = vi.fn(() => ({
  push: routerPushMock,
}))

vi.mock('vue-router', () => ({
  useRouter: () => useRouterMock(),
}))

vi.mock('@/components/EditProfileInfo.vue', () => ({
  default: {
    name: 'EditProfileInfo',
    template: '<div class="edit-profile-info"></div>',
  },
}))

vi.mock('@/components/EditableLinksList.vue', () => ({
  default: {
    name: 'EditableLinksList',
    template: '<div class="editable-links-list"></div>',
  },
}))

vi.mock('@/components/DeleteAccount.vue', () => ({
  default: {
    name: 'DeleteAccount',
    template: '<div class="delete-account"></div>',
  },
}))

vi.mock('axios', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    defaults: {
      headers: {
        common: {},
      },
    },
  },
}))

describe('EditProfile', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    routerPushMock.mockClear()
    localStorage.clear()
  })

  it('redirects to login if user is not authenticated', async () => {
    const authStore = useAuthStore()
    authStore.$patch({
      token: null,
    })

    mount(EditProfile)
    await nextTick()

    expect(routerPushMock).toHaveBeenCalledWith('/login')
  })

  it('fetches user profile on mount if authenticated', async () => {
    const authStore = useAuthStore()
    const userStore = useUserStore()

    authStore.$patch({
      token: 'fake-token',
      username: 'testuser',
    })

    const fetchUserProfileMock = vi.fn().mockResolvedValue({
      full_name: 'Test User',
      bio: 'Test Bio',
      links: [],
    })
    userStore.fetchUserProfile = fetchUserProfileMock

    mount(EditProfile)
    await nextTick()

    expect(fetchUserProfileMock).toHaveBeenCalledWith('testuser')
  })

  it('shows loading state when userData is not available', async () => {
    const authStore = useAuthStore()
    const userStore = useUserStore()

    authStore.$patch({
      token: 'fake-token',
      username: 'testuser',
    })

    userStore.fetchUserProfile = vi
      .fn()
      .mockImplementation(() => new Promise(() => {}))

    const wrapper = mount(EditProfile)
    await nextTick()

    expect(wrapper.text()).toContain('Loading user data')
  })

  it('shows error message when there is an error', async () => {
    const authStore = useAuthStore()
    const userStore = useUserStore()

    authStore.$patch({
      token: 'fake-token',
      username: 'testuser',
    })

    const errorResponse = {
      response: {
        data: {
          error: 'Failed to load profile',
        },
      },
    }

    userStore.fetchUserProfile = vi.fn().mockImplementation(async () => {
      userStore.error = 'Failed to load profile'
      throw errorResponse
    })

    vi.spyOn(console, 'error').mockImplementation(() => {})

    const wrapper = mount(EditProfile)

    await nextTick()
    await nextTick()

    expect(wrapper.text()).toContain('Failed to load profile')

    console.error.mockRestore()
  })
})
