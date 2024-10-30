<template>
  <section class="min-h-screen bg-gradient-to-r from-teal-100 p-6">
    <div class="max-w-3xl mx-auto">
      <div class="flex justify-end mb-4">
        <router-link v-if="isOwner" to="/edit-profile">
          <button class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition duration-200">
            Edit Profile
          </button>
        </router-link>
      </div>

      <p v-if="userStore.error" class="text-red-500">{{ userStore.error }}</p>

      <div v-if="!userStore.userData && !userStore.error" class="text-gray-500">
        Loading...
      </div>

      <ProfileInfo v-if="userStore.userData" :fullName="userStore.userData.full_name"
        :username="userStore.userData.username" :bio="userStore.userData.bio" />

      <LinksList v-if="userStore.userData && userStore.userData.links" :links="userStore.userData.links" />
    </div>
  </section>
</template>

<script setup>
import { onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import ProfileInfo from './ProfileInfo.vue'
import LinksList from './LinksList.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()

const fetchProfile = async (username) => {
  try {
    await userStore.fetchUserProfile(username)
  } catch (error) {
    console.error(error)
    router.push('/')
  }
}

onMounted(() => {
  const username = route.params.username || authStore.username
  fetchProfile(username)
})

watch(() => route.params.username, (newUsername) => {
  fetchProfile(newUsername)
})

const isOwner = computed(() => {
  return userStore.userData && authStore.username === userStore.userData.username
})
</script>
