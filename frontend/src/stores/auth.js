import { defineStore } from 'pinia'
import axios from 'axios'
import { API_BASE_URL } from '../config.js'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || null,
    username: localStorage.getItem('username') || null,
  }),

  getters: {
    isAuthenticated: state => !!state.token,
    getUsername: state => state.username,
  },

  actions: {
    async login(username, password) {
      try {
        const response = await axios.post(`${API_BASE_URL}/v1/users/login`, {
          username,
          password,
        })
        this.token = response.data.token
        this.username = username
        localStorage.setItem('token', this.token)
        localStorage.setItem('username', this.username)
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
        return true
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Login failed')
      }
    },

    async signup(userData) {
      try {
        await axios.post(`${API_BASE_URL}/v1/users/signup`, userData)
        return true
      } catch (error) {
        throw new Error(error.response?.data?.error || 'Signup failed')
      }
    },

    logout() {
      this.token = null
      this.username = null
      localStorage.removeItem('token')
      localStorage.removeItem('username')
      delete axios.defaults.headers.common['Authorization']
    },
  },
})
