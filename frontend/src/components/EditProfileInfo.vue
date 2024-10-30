<template>
  <div class="p-4 bg-white rounded-lg shadow-md mb-4">
    <h2 class="text-2xl font-bold text-gray-800 mb-4">Edit Profile</h2>
    <form @submit.prevent="showConfirmation">
      <div class="mb-4">
        <label class="block text-gray-600">Full Name</label>
        <input type="text" v-model="formData.full_name"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
      </div>
      <div class="mb-4">
        <label class="block text-gray-600">Password</label>
        <input type="password" v-model="formData.password" placeholder="Enter new password"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
      </div>
      <div class="mb-4">
        <label class="block text-gray-600">Bio</label>
        <textarea v-model="formData.bio"
          class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          rows="3"></textarea>
      </div>
      <button type="submit"
        class="bg-primary text-white font-semibold px-6 py-2 rounded-lg hover:bg-blue-700 transition duration-200">
        Save Changes
      </button>
    </form>

    <ConfirmationBox v-if="confirmationVisible" title="Confirm Update"
      message="Are you sure you want to save these changes?" @confirm="confirmSaveProfile"
      @cancel="confirmationVisible = false" />

    <MessageBox v-if="messageVisible" :title="messageTitle" :message="message" :type="messageType"
      @close="messageVisible = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import ConfirmationBox from '@/components/ConfirmationBox.vue'
import MessageBox from '@/components/MessageBox.vue'

const props = defineProps(['fullName', 'bio'])
const userStore = useUserStore()
const authStore = useAuthStore()

const formData = ref({
  full_name: '',
  bio: '',
  password: ''
})

watch(() => props.fullName, (newVal) => {
  formData.value.full_name = newVal
}, { immediate: true })

watch(() => props.bio, (newVal) => {
  formData.value.bio = newVal
}, { immediate: true })

const confirmationVisible = ref(false)
const messageVisible = ref(false)
const messageTitle = ref('')
const message = ref('')
const messageType = ref('')

const showConfirmation = () => {
  confirmationVisible.value = true
}

const confirmSaveProfile = async () => {
  try {
    await userStore.updateProfile(formData.value)
    showMessage('Profile updated successfully', 'success')
  } catch (error) {
    showMessage(error.message || 'There was an error updating your profile. Please try again.', 'error')
  } finally {
    confirmationVisible.value = false
  }
}

const showMessage = (msg, type) => {
  messageTitle.value = type === 'success' ? 'Success' : 'Error'
  message.value = msg
  messageType.value = type
  messageVisible.value = true
}
</script>
