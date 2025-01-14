import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth.js'
import { API_BASE_URL } from '../config.js'

export const useUserStore = defineStore('user', {
  state: () => ({
    userData: null,
    error: null,
  }),

  actions: {
    async fetchUserProfile(username) {
      try {
        const response = await axios.get(`${API_BASE_URL}/v1/users/${username}`)
        this.userData = response.data
        return response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch profile'
        console.error(this.error)
        throw error
      }
    },

    async updateProfile(userData) {
      try {
        await axios.put(`${API_BASE_URL}/v1/users`, userData)
        const authStore = useAuthStore()
        await this.fetchUserProfile(authStore.username)
      } catch (error) {
        throw new Error(
          error.response?.data?.error || 'Failed to update profile',
        )
      }
    },

    async deleteAccount() {
      try {
        await axios.delete(`${API_BASE_URL}/v1/users`)
        const authStore = useAuthStore()
        authStore.logout()
      } catch (error) {
        throw new Error(
          error.response?.data?.error || 'Failed to delete account',
        )
      }
    },
  },
})
