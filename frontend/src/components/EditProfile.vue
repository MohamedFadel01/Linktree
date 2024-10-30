<template>
  <section class="min-h-screen bg-gradient-to-r from-teal-100 to-slate-100 p-6">
    <div class="max-w-3xl mx-auto">
      <template v-if="userStore.userData">
        <EditProfileInfo :fullName="userStore.userData.full_name" :bio="userStore.userData.bio" />

        <EditableLinksList :links="userStore.userData.links" @update-links="updateLinks" />

        <DeleteAccount />
      </template>

      <template v-else>
        <p v-if="userStore.error" class="text-red-500">{{ userStore.error }}</p>
        <p v-else>Loading user data...</p>
      </template>
    </div>
  </section>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import EditProfileInfo from './EditProfileInfo.vue'
import EditableLinksList from './EditableLinksList.vue'
import DeleteAccount from './DeleteAccount.vue'

const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()

onMounted(async () => {
  if (!authStore.isAuthenticated) {
    router.push('/login')
    return
  }
  await userStore.fetchUserProfile(authStore.username)
})

const updateLinks = () => {
  userStore.fetchUserProfile(authStore.username)
}
</script>
