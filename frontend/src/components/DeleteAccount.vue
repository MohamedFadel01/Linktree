<template>
  <div class="p-4 mt-6 bg-red-50 rounded-lg">
    <p class="text-red-600 font-semibold">Danger Zone</p>
    <button @click="showConfirmation" class="text-red-500 underline mt-2">
      Delete Account
    </button>
  </div>

  <ConfirmationBox v-if="confirmationVisible" title="Confirm Deletion"
    message="Are you sure you want to delete your account? This action cannot be undone." @confirm="confirmDelete"
    @cancel="confirmationVisible = false" />
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import ConfirmationBox from '@/components/ConfirmationBox.vue'

const router = useRouter()
const userStore = useUserStore()
const confirmationVisible = ref(false)

const showConfirmation = () => {
  confirmationVisible.value = true
}

const confirmDelete = async () => {
  try {
    await userStore.deleteAccount()
    showMessageBox('Account Deleted', 'Your account has been successfully deleted.', 'success');
    router.push('/')
  } catch (error) {
    showMessageBox('Deletion Failed', error.message || 'There was an error deleting your account. Please try again.', 'error');
  } finally {
    confirmationVisible.value = false
  }
}

const showMessageBox = (title, message, type) => {
  console.log(`${title}: ${message}`);
}
</script>
