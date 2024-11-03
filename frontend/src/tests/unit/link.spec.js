// src/tests/unit/link.spec.js
import { describe, it, beforeEach, expect, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useLinkStore } from '../../stores/link.js'
import { useUserStore } from '../../stores/user.js'
import axios from 'axios'

vi.mock('axios')

vi.mock('../../stores/user.js', () => ({
  useUserStore: vi.fn(),
}))

describe('Link Store', () => {
  let store
  let userStore

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useLinkStore()

    vi.clearAllMocks()

    userStore = {
      userData: {
        username: 'testuser',
      },
      fetchUserProfile: vi.fn().mockResolvedValue({}),
    }
    useUserStore.mockReturnValue(userStore)
  })

  describe('createLink', () => {
    it('should successfully create a link and refresh user profile', async () => {
      const linkData = {
        title: 'Test Link',
        url: 'https://example.com',
      }

      vi.mocked(axios.post).mockResolvedValueOnce({ data: linkData })

      await store.createLink(linkData)

      expect(axios.post).toHaveBeenCalledWith(
        'http://localhost:8188/api/v1/links',
        linkData,
      )
      expect(userStore.fetchUserProfile).toHaveBeenCalledWith('testuser')
    })

    it('should handle create link error', async () => {
      const errorMessage = 'Invalid URL'
      vi.mocked(axios.post).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      const linkData = { url: 'invalid-url' }

      await expect(store.createLink(linkData)).rejects.toThrow(errorMessage)
      expect(userStore.fetchUserProfile).not.toHaveBeenCalled()
    })
  })

  describe('updateLink', () => {
    it('should successfully update a link and refresh user profile', async () => {
      const linkId = '123'
      const linkData = {
        title: 'Updated Link',
        url: 'https://example.com/updated',
      }

      vi.mocked(axios.put).mockResolvedValueOnce({ data: linkData })

      await store.updateLink(linkId, linkData)

      expect(axios.put).toHaveBeenCalledWith(
        `http://localhost:8188/api/v1/links/${linkId}`,
        linkData,
      )
      expect(userStore.fetchUserProfile).toHaveBeenCalledWith('testuser')
    })

    it('should handle update link error', async () => {
      const linkId = '123'
      const errorMessage = 'Link not found'
      vi.mocked(axios.put).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      await expect(store.updateLink(linkId, {})).rejects.toThrow(errorMessage)
      expect(userStore.fetchUserProfile).not.toHaveBeenCalled()
    })
  })

  describe('deleteLink', () => {
    it('should successfully delete a link and refresh user profile', async () => {
      const linkId = '123'

      vi.mocked(axios.delete).mockResolvedValueOnce({ data: {} })

      await store.deleteLink(linkId)

      expect(axios.delete).toHaveBeenCalledWith(
        `http://localhost:8188/api/v1/links/${linkId}`,
      )
      expect(userStore.fetchUserProfile).toHaveBeenCalledWith('testuser')
    })

    it('should handle delete link error', async () => {
      const linkId = '123'
      const errorMessage = 'Link not found'
      vi.mocked(axios.delete).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      await expect(store.deleteLink(linkId)).rejects.toThrow(errorMessage)
      expect(userStore.fetchUserProfile).not.toHaveBeenCalled()
    })
  })

  describe('trackClick', () => {
    it('should successfully track link click and refresh user profile', async () => {
      const linkId = '123'

      vi.mocked(axios.post).mockResolvedValueOnce({ data: {} })

      await store.trackClick(linkId)

      expect(axios.post).toHaveBeenCalledWith(
        `http://localhost:8188/api/v1/analytics/${linkId}/click`,
      )
      expect(userStore.fetchUserProfile).toHaveBeenCalledWith('testuser')
    })

    it('should handle track click error gracefully', async () => {
      const linkId = '123'
      const consoleSpy = vi.spyOn(console, 'error')

      vi.mocked(axios.post).mockRejectedValueOnce(new Error('Network error'))

      await store.trackClick(linkId)

      expect(consoleSpy).toHaveBeenCalled()
      expect(userStore.fetchUserProfile).not.toHaveBeenCalled()
    })
  })
})
