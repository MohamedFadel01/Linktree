import { defineStore } from 'pinia'
import axios from 'axios'
import { useUserStore } from './user.js'

export const useLinkStore = defineStore('link', {
  state: () => ({
    links: [],
  }),

  actions: {
    async createLink(linkData) {
      try {
        await axios.post('http://localhost:8188/api/v1/links', linkData)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to create link')
      }
    },

    async updateLink(id, linkData) {
      try {
        await axios.put(`http://localhost:8188/api/v1/links/${id}`, linkData)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to update link')
      }
    },

    async deleteLink(id) {
      try {
        await axios.delete(`http://localhost:8188/api/v1/links/${id}`)
        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Failed to delete link')
      }
    },

    async trackClick(id) {
      try {
        await axios.post(`http://localhost:8188/api/v1/analytics/${id}/click`)

        const userStore = useUserStore()
        await userStore.fetchUserProfile(userStore.userData.username)
      } catch (error) {
        console.error('Failed to track link click:', error)
      }
    },
  },
})
