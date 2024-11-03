import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import EditableLinksList from '@/components/EditableLinksList.vue'
import { useLinkStore } from '@/stores/link'

vi.mock('@/components/EditableLinkItem.vue', () => ({
  default: {
    name: 'EditableLinkItem',
    template: '<div class="editable-link-item"></div>',
  },
}))

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

describe('EditableLinksList', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const mockLinks = [
    { id: 1, title: 'Link 1', url: 'https://test1.com' },
    { id: 2, title: 'Link 2', url: 'https://test2.com' },
  ]

  it('renders all links', () => {
    const wrapper = mount(EditableLinksList, {
      props: {
        links: mockLinks,
      },
    })

    const linkItems = wrapper.findAllComponents({ name: 'EditableLinkItem' })
    expect(linkItems).toHaveLength(2)
  })

  it('shows confirmation dialog when deleting a link', async () => {
    const wrapper = mount(EditableLinksList, {
      props: {
        links: mockLinks,
      },
    })

    await wrapper.vm.showConfirmation(1)
    expect(wrapper.vm.confirmationVisible).toBe(true)
    expect(wrapper.vm.linkIdToDelete).toBe(1)
  })

  it('emits update-links event after adding new link', async () => {
    const linkStore = useLinkStore()
    linkStore.createLink = vi
      .fn()
      .mockResolvedValue({ id: 3, title: 'New Link', url: 'https://' })

    const wrapper = mount(EditableLinksList, {
      props: {
        links: mockLinks,
      },
    })

    await wrapper.find('button.text-blue-500').trigger('click')
    await nextTick()
    await nextTick()

    expect(linkStore.createLink).toHaveBeenCalled()
    expect(wrapper.emitted('update-links')).toBeTruthy()
  })
})
