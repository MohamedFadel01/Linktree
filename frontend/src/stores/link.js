import { defineStore } from 'pinia'
import axios from 'axios'
import { useUserStore } from './user.js'
import { API_BASE_URL } from '../config.js'

export const useLinkStore = defineStore('link', {
  state: () => ({
    links: [],
  }),

  actions: {
    async createLink(linkData) {
      try {
        await axios.post(`${API_BASE_URL}/v1/links`, linkData)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to create link')
      }
    },

    async updateLink(id, linkData) {
      try {
        await axios.put(`${API_BASE_URL}/v1/links/${id}`, linkData)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to update link')
      }
    },

    async deleteLink(id) {
      try {
        await axios.delete(`${API_BASE_URL}/v1/links/${id}`)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to delete link')
      }
    },

    async trackClick(id) {
      try {
        await axios.post(`${API_BASE_URL}/v1/analytics/${id}/click`)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        console.error('Failed to track link click:', error)
      }
    },
  },
})
