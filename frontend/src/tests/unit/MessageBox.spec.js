import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { shallowMount } from '@vue/test-utils'
import MessageBox from '@/components/MessageBox.vue'

describe('MessageBox.vue', () => {
  let wrapper
  const title = 'Test Title'
  const message = 'This is a test message.'
  const duration = 3000

  beforeEach(() => {
    vi.useFakeTimers()
    wrapper = shallowMount(MessageBox, {
      props: { title, message, duration },
    })
  })

  afterEach(() => {
    vi.clearAllTimers()
  })

  it('renders the title and message correctly', () => {
    expect(wrapper.text()).toContain(title)
    expect(wrapper.text()).toContain(message)
  })

  it('is visible when created', () => {
    expect(wrapper.vm.visible).toBe(true)
    expect(wrapper.isVisible()).toBe(true)
  })

  it('hides the message box after the specified duration', async () => {
    expect(wrapper.vm.visible).toBe(true)

    vi.advanceTimersByTime(duration)
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.visible).toBe(false)
    expect(wrapper.isVisible()).toBe(false)
  })

  it('hides the message box when clicked', async () => {
    await wrapper.trigger('click')
    expect(wrapper.vm.visible).toBe(false)
    expect(wrapper.isVisible()).toBe(false)
  })
})
