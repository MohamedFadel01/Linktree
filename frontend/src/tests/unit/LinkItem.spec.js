import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'
import LinkItem from '../../components/LinkItem.vue'
import { useLinkStore } from '@/stores/link'

vi.mock('../../components/LinkAnalytics.vue', () => ({
  default: {
    name: 'LinkAnalytics',
    props: ['analytics'],
    template: '<div data-testid="link-analytics"></div>',
  },
}))

describe('LinkItem', () => {
  let wrapper
  let linkStore

  const windowOpenSpy = vi.spyOn(window, 'open')
  windowOpenSpy.mockImplementation(() => {})

  const mockLink = {
    id: '123',
    title: 'Test Link',
    url: 'https://example.com',
    analytics: {
      click_count: 10,
      last_clicked: '2024-01-01T00:00:00Z',
    },
  }

  beforeEach(() => {
    const pinia = createTestingPinia({
      createSpy: vi.fn,
    })

    linkStore = useLinkStore(pinia)

    wrapper = mount(LinkItem, {
      props: {
        link: mockLink,
      },
      global: {
        plugins: [pinia],
      },
    })
  })

  it('renders link title and url correctly', () => {
    expect(wrapper.find('h3').text()).toBe(mockLink.title)
    expect(wrapper.find('p').text()).toBe(mockLink.url)
  })

  it('renders LinkAnalytics component with correct props', () => {
    const analytics = wrapper.findComponent({ name: 'LinkAnalytics' })
    expect(analytics.exists()).toBe(true)
    expect(analytics.props('analytics')).toEqual(mockLink.analytics)
  })

  it('calls trackClick and opens link in new window when clicked', async () => {
    linkStore.trackClick.mockResolvedValue()

    const link = wrapper.find('a')
    await link.trigger('click')

    expect(linkStore.trackClick).toHaveBeenCalledWith(mockLink.id)
    expect(windowOpenSpy).toHaveBeenCalledWith(mockLink.url, '_blank')
  })

  it('handles tracking error gracefully', async () => {
    const error = new Error('Tracking failed')
    linkStore.trackClick.mockRejectedValue(error)

    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    const link = wrapper.find('a')
    await link.trigger('click')

    expect(consoleSpy).toHaveBeenCalledWith(
      'Failed to track link click:',
      error,
    )

    expect(windowOpenSpy).toHaveBeenCalledWith(mockLink.url, '_blank')

    consoleSpy.mockRestore()
  })

  it('prevents default link behavior', async () => {
    const link = wrapper.find('a')
    const event = {
      preventDefault: vi.fn(),
    }

    await link.trigger('click', event)

    expect(event.preventDefault).toHaveBeenCalled()
  })

  it('has correct href attribute', () => {
    const link = wrapper.find('a')
    expect(link.attributes('href')).toBe(mockLink.url)
  })

  it('has target="_blank" for external links', () => {
    const link = wrapper.find('a')
    expect(link.attributes('target')).toBe('_blank')
  })
})
