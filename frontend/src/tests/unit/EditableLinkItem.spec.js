import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { nextTick } from 'vue'
import EditableLinkItem from '@/components/EditableLinkItem.vue'
import { useLinkStore } from '@/stores/link'

vi.mock('@/components/MessageBox.vue', () => ({
  default: {
    name: 'MessageBox',
    template: '<div class="message-box"></div>',
  },
}))

describe('EditableLinkItem', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  const mockLink = {
    id: 1,
    title: 'Test Link',
    url: 'https://test.com',
  }

  it('displays link data correctly', async () => {
    const linkStore = useLinkStore()
    linkStore.updateLink = vi.fn()

    const wrapper = mount(EditableLinkItem, {
      props: {
        link: mockLink,
      },
    })

    await nextTick()
    const titleInput = wrapper.find('input[placeholder="Link title"]')
    const urlInput = wrapper.find('input[placeholder="Link URL"]')

    expect(titleInput.element.value).toBe('Test Link')
    expect(urlInput.element.value).toBe('https://test.com')
  })

  it('emits delete-link event when delete button is clicked', async () => {
    const wrapper = mount(EditableLinkItem, {
      props: {
        link: mockLink,
      },
    })

    const deleteButton = wrapper.find('button.text-red-500')
    await deleteButton.trigger('click')

    expect(wrapper.emitted('delete-link')).toBeTruthy()
    expect(wrapper.emitted('delete-link')[0]).toEqual([1])
  })

  it('updates link data when inputs change', async () => {
    const wrapper = mount(EditableLinkItem, {
      props: {
        link: mockLink,
      },
    })

    const titleInput = wrapper.find('input[placeholder="Link title"]')
    const urlInput = wrapper.find('input[placeholder="Link URL"]')

    await titleInput.setValue('Updated Title')
    await urlInput.setValue('https://updated.com')
    await nextTick()

    expect(wrapper.vm.linkData.title).toBe('Updated Title')
    expect(wrapper.vm.linkData.url).toBe('https://updated.com')
  })
})
