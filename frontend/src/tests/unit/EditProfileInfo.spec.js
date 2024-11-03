import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import EditProfileInfo from '@/components/EditProfileInfo.vue'
import { useUserStore } from '@/stores/user'

vi.mock('@/components/MessageBox.vue', () => ({
  default: {
    name: 'MessageBox',
    template: '<div class="message-box"></div>',
  },
}))

vi.mock('@/components/ConfirmationBox.vue', () => ({
  default: {
    name: 'ConfirmationBox',
    template: '<div class="confirmation-box"></div>',
  },
}))

describe('EditProfileInfo', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const mockProps = {
    fullName: 'John Doe',
    bio: 'Test bio',
  }

  it('initializes form with provided props', () => {
    const wrapper = mount(EditProfileInfo, {
      props: mockProps,
    })

    expect(wrapper.vm.formData.full_name).toBe('John Doe')
    expect(wrapper.vm.formData.bio).toBe('Test bio')
  })

  it('shows confirmation dialog when form is submitted', async () => {
    const wrapper = mount(EditProfileInfo, {
      props: mockProps,
    })

    await wrapper.find('form').trigger('submit')
    expect(wrapper.vm.confirmationVisible).toBe(true)
  })

  it('updates form data when inputs change', async () => {
    const wrapper = mount(EditProfileInfo, {
      props: mockProps,
    })

    const nameInput = wrapper.find('input[type="text"]')
    const bioTextarea = wrapper.find('textarea')
    const passwordInput = wrapper.find('input[type="password"]')

    await nameInput.setValue('Jane Doe')
    await bioTextarea.setValue('Updated bio')
    await passwordInput.setValue('newpassword')

    expect(wrapper.vm.formData.full_name).toBe('Jane Doe')
    expect(wrapper.vm.formData.bio).toBe('Updated bio')
    expect(wrapper.vm.formData.password).toBe('newpassword')
  })

  it('displays success message after successful profile update', async () => {
    const userStore = useUserStore()
    userStore.updateProfile = vi.fn().mockResolvedValue()

    const wrapper = mount(EditProfileInfo, {
      props: mockProps,
    })

    userStore.updateProfile.mockResolvedValueOnce()

    await wrapper.vm.confirmSaveProfile()
    await nextTick()
    await nextTick()

    expect(wrapper.vm.messageVisible).toBe(true)
    expect(wrapper.vm.messageType).toBe('success')
    expect(wrapper.vm.message).toBe('Profile updated successfully')
  })
})
