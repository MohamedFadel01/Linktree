// user.store.spec.js
import { describe, it, beforeEach, expect, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '../../stores/user.js'
import axios from 'axios'
import { useAuthStore } from '../../stores/auth.js'

vi.mock('axios')

vi.mock('../../stores/auth.js', () => ({
  useAuthStore: vi.fn(),
}))

describe('User Store', () => {
  let store
  let authStore

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useUserStore()

    vi.clearAllMocks()

    authStore = {
      username: 'testuser',
      logout: vi.fn(),
    }
    useAuthStore.mockReturnValue(authStore)
  })

  describe('fetchUserProfile', () => {
    it('should successfully fetch user profile', async () => {
      const mockUserData = { id: 1, username: 'testuser' }
      vi.mocked(axios.get).mockResolvedValueOnce({ data: mockUserData })

      const result = await store.fetchUserProfile('testuser')

      expect(axios.get).toHaveBeenCalledWith(
        'http://localhost:8188/api/v1/users/testuser',
      )
      expect(store.userData).toEqual(mockUserData)
      expect(result).toEqual(mockUserData)
      expect(store.error).toBeNull()
    })

    it('should handle fetch error and set error state', async () => {
      const errorMessage = 'User not found'
      vi.mocked(axios.get).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      await expect(store.fetchUserProfile('testuser')).rejects.toThrow()
      expect(store.error).toBe(errorMessage)
      expect(store.userData).toBeNull()
    })
  })

  describe('updateProfile', () => {
    it('should successfully update profile and fetch updated data', async () => {
      const userData = { name: 'Updated Name' }
      const mockUserData = { id: 1, username: 'testuser', ...userData }

      vi.mocked(axios.put).mockResolvedValueOnce({ data: mockUserData })
      vi.mocked(axios.get).mockResolvedValueOnce({ data: mockUserData })

      await store.updateProfile(userData)

      expect(axios.put).toHaveBeenCalledWith(
        'http://localhost:8188/api/v1/users',
        userData,
      )
      expect(axios.get).toHaveBeenCalledWith(
        'http://localhost:8188/api/v1/users/testuser',
      )
      expect(store.userData).toEqual(mockUserData)
    })

    it('should handle update error with proper error message', async () => {
      const errorMessage = 'Update failed'
      vi.mocked(axios.put).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      await expect(store.updateProfile({})).rejects.toThrow(errorMessage)
    })
  })

  describe('deleteAccount', () => {
    it('should successfully delete account and logout', async () => {
      vi.mocked(axios.delete).mockResolvedValueOnce({ data: {} })

      await store.deleteAccount()

      expect(axios.delete).toHaveBeenCalledWith(
        'http://localhost:8188/api/v1/users',
      )
      expect(authStore.logout).toHaveBeenCalledOnce()
    })

    it('should handle delete error with proper error message', async () => {
      const errorMessage = 'Delete failed'
      vi.mocked(axios.delete).mockRejectedValueOnce({
        response: { data: { error: errorMessage } },
      })

      await expect(store.deleteAccount()).rejects.toThrow(errorMessage)
    })
  })
})
