import { mount } from '@vue/test-utils'
import ConfirmationBox from '@/components/ConfirmationBox.vue'
import { describe, it, beforeEach, expect } from 'vitest'

describe('ConfirmationBox.vue', () => {
  let wrapper
  const title = 'Confirm Action'
  const message = 'Are you sure you want to proceed?'

  beforeEach(() => {
    wrapper = mount(ConfirmationBox, {
      props: {
        title,
        message,
      },
    })
  })

  it('renders the title and message correctly', () => {
    expect(wrapper.find('h3').text()).toBe(title)
    expect(wrapper.find('p').text()).toBe(message)
  })

  it('emits confirm event when Confirm button is clicked', async () => {
    await wrapper.find('button.bg-red-600').trigger('click')
    expect(wrapper.emitted().confirm).toBeTruthy()
  })

  it('emits cancel event when Cancel button is clicked', async () => {
    await wrapper.find('button.bg-gray-300').trigger('click')
    expect(wrapper.emitted().cancel).toBeTruthy()
  })
})
