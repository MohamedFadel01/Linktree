<template>
  <nav class="bg-teal-800 shadow-lg p-4 flex items-center justify-between">
    <div class="flex items-center space-x-2">
      <img src="@/assets/logo.svg" alt="Logo" class="w-8 h-8" />
      <span class="text-2xl font-semibold text-white">Linktree</span>
    </div>

    <div class="flex-1 max-w-md mx-4">
      <form @submit.prevent="handleSearch" class="relative">
        <input data-test="search-input" v-model="searchUsername" type="text" placeholder="Search username..."
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-teal-500" />
      </form>
    </div>

    <div class="flex items-center space-x-4">
      <router-link to="/" class="text-1xl font-semibold text-white hover:text-teal-300">
        Home
      </router-link>

      <div class="relative">
        <button data-test="menu-button" class="text-1xl font-semibold text-white hover:text-teal-300 flex items-center" @click="toggleDropdown">
          Menu
          <svg class="ml-1 w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        <div v-if="isDropdownOpen" class="absolute right-0 mt-2 w-48 bg-white border rounded shadow-lg">
          <router-link v-if="authStore.isAuthenticated" :to="'/profile'"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-100 w-full text-left" @click="isDropdownOpen = false">
            Profile
          </router-link>

          <button v-if="!authStore.isAuthenticated" @click="navigateToAuth"
            class="block px-4 py-2 text-gray-700 hover:bg-gray-100 w-full text-left">
            Login / Signup
          </button>
          <button v-else @click="handleLogout" class="block px-4 py-2 text-red-600 hover:bg-gray-100 w-full text-left">
            Logout
          </button>
        </div>
      </div>
    </div>

    <MessageBox v-if="messageBoxVisible" :title="messageTitle" :message="messageText" :type="messageType" />
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import MessageBox from '@/components/MessageBox.vue'

const router = useRouter()
const authStore = useAuthStore()
const searchUsername = ref('')
const isDropdownOpen = ref(false)
const messageBoxVisible = ref(false)
const messageTitle = ref('')
const messageText = ref('')
const messageType = ref('info')

const handleSearch = async () => {
  if (searchUsername.value) {
    try {
      const response = await fetch(`http://localhost:8188/api/v1/users/${searchUsername.value}`);

      if (response.ok) {
        router.push(`/${searchUsername.value}`);
      } else {
        showMessageBox('User Not Found', `The user "${searchUsername.value}" does not exist.`, 'error');
      }
    } catch (error) {
      showMessageBox('Error', `There was an error searching for the user. ${error}`, 'error');
    }
    searchUsername.value = ''
  }
}

const showMessageBox = (title, message, type) => {
  messageTitle.value = title;
  messageText.value = message;
  messageType.value = type;
  messageBoxVisible.value = true;

  setTimeout(() => {
    messageBoxVisible.value = false;
  }, 5000);
}

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
}

const handleLogout = () => {
  authStore.logout()
  router.push('/')
  isDropdownOpen.value = false
}

const navigateToAuth = () => {
  router.push('/login')
  isDropdownOpen.value = false
}
</script>
