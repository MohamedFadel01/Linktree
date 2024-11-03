import { describe, beforeEach, it, expect, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'
import AxiosMockAdapter from 'axios-mock-adapter'

const mock = new AxiosMockAdapter(axios)

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  afterEach(() => {
    mock.reset()
  })

  it('should login successfully', async () => {
    const store = useAuthStore()
    const username = 'testuser'
    const password = 'password123'

    mock.onPost('http://localhost:8188/api/v1/users/login').reply(200, {
      token: 'mock_token',
    })

    const result = await store.login(username, password)

    expect(result).toBe(true)
    expect(store.token).toBe('mock_token')
    expect(store.username).toBe(username)
    expect(localStorage.getItem('token')).toBe('mock_token')
    expect(localStorage.getItem('username')).toBe(username)
  })

  it('should handle login failure', async () => {
    const store = useAuthStore()
    const username = 'wronguser'
    const password = 'wrongpassword'

    mock.onPost('http://localhost:8188/api/v1/users/login').reply(401, {
      error: 'Invalid credentials',
    })

    await expect(store.login(username, password)).rejects.toThrow(
      'Invalid credentials',
    )
    expect(store.token).toBe(null)
    expect(store.username).toBe(null)
  })

  it('should signup successfully', async () => {
    const store = useAuthStore()
    const userData = { username: 'newuser', password: 'password123' }

    mock.onPost('http://localhost:8188/api/v1/users/signup').reply(200)

    const result = await store.signup(userData)

    expect(result).toBe(true)
  })

  it('should handle signup failure', async () => {
    const store = useAuthStore()
    const userData = { username: 'existinguser', password: 'password123' }

    mock.onPost('http://localhost:8188/api/v1/users/signup').reply(400, {
      error: 'User already exists',
    })

    await expect(store.signup(userData)).rejects.toThrow('User already exists')
  })

  it('should logout successfully', () => {
    const store = useAuthStore()
    store.token = 'mock_token'
    store.username = 'testuser'
    localStorage.setItem('token', store.token)
    localStorage.setItem('username', store.username)

    store.logout()

    expect(store.token).toBe(null)
    expect(store.username).toBe(null)
    expect(localStorage.getItem('token')).toBe(null)
    expect(localStorage.getItem('username')).toBe(null)
  })
})
