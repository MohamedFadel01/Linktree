import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import DeleteAccount from '@/components/DeleteAccount.vue'
import { useUserStore } from '@/stores/user'

const routerPushMock = vi.fn()
const useRouterMock = vi.fn(() => ({
  push: routerPushMock,
}))

vi.mock('vue-router', () => ({
  useRouter: () => useRouterMock(),
}))

vi.mock('@/components/ConfirmationBox.vue', () => ({
  default: {
    name: 'ConfirmationBox',
    template: '<div class="confirmation-box"><slot></slot></div>',
    props: ['title', 'message'],
    emits: ['confirm', 'cancel'],
  },
}))

const consoleLogMock = vi.fn()
vi.spyOn(console, 'log').mockImplementation(consoleLogMock)

describe('DeleteAccount', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    routerPushMock.mockClear()
    consoleLogMock.mockClear()
  })

  it('renders the delete account button', () => {
    const wrapper = mount(DeleteAccount)
    expect(wrapper.text()).toContain('Delete Account')
    expect(wrapper.text()).toContain('Danger Zone')
  })

  it('shows confirmation dialog when delete button is clicked', async () => {
    const wrapper = mount(DeleteAccount)

    expect(wrapper.findComponent({ name: 'ConfirmationBox' }).exists()).toBe(
      false,
    )

    await wrapper.find('button').trigger('click')

    expect(wrapper.findComponent({ name: 'ConfirmationBox' }).exists()).toBe(
      true,
    )
  })

  it('hides confirmation dialog when cancel is clicked', async () => {
    const wrapper = mount(DeleteAccount)

    await wrapper.find('button').trigger('click')

    await wrapper.findComponent({ name: 'ConfirmationBox' }).vm.$emit('cancel')

    expect(wrapper.findComponent({ name: 'ConfirmationBox' }).exists()).toBe(
      false,
    )
  })

  it('successfully deletes account and redirects', async () => {
    const wrapper = mount(DeleteAccount)
    const userStore = useUserStore()

    userStore.deleteAccount = vi.fn().mockResolvedValue()

    await wrapper.find('button').trigger('click')
    await wrapper.findComponent({ name: 'ConfirmationBox' }).vm.$emit('confirm')

    await wrapper.vm.$nextTick()

    expect(userStore.deleteAccount).toHaveBeenCalled()

    expect(consoleLogMock).toHaveBeenCalledWith(
      'Account Deleted: Your account has been successfully deleted.',
    )

    expect(routerPushMock).toHaveBeenCalledWith('/')

    expect(wrapper.findComponent({ name: 'ConfirmationBox' }).exists()).toBe(
      false,
    )
  })

  it('handles deletion failure', async () => {
    const wrapper = mount(DeleteAccount)
    const userStore = useUserStore()

    const errorMessage = 'Failed to delete account'
    userStore.deleteAccount = vi.fn().mockRejectedValue(new Error(errorMessage))

    await wrapper.find('button').trigger('click')
    await wrapper.findComponent({ name: 'ConfirmationBox' }).vm.$emit('confirm')

    await wrapper.vm.$nextTick()

    expect(consoleLogMock).toHaveBeenCalledWith(
      `Deletion Failed: ${errorMessage}`,
    )

    expect(routerPushMock).not.toHaveBeenCalled()

    expect(wrapper.findComponent({ name: 'ConfirmationBox' }).exists()).toBe(
      false,
    )
  })

  it('handles deletion failure with default error message', async () => {
    const wrapper = mount(DeleteAccount)
    const userStore = useUserStore()

    userStore.deleteAccount = vi.fn().mockRejectedValue(new Error())

    await wrapper.find('button').trigger('click')
    await wrapper.findComponent({ name: 'ConfirmationBox' }).vm.$emit('confirm')

    await wrapper.vm.$nextTick()

    expect(consoleLogMock).toHaveBeenCalledWith(
      'Deletion Failed: There was an error deleting your account. Please try again.',
    )
  })
})
