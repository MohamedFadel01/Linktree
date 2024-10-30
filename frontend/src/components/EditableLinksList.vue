<template>
  <div class="space-y-4">
    <EditableLinkItem v-for="link in links" :key="link.id" :link="link" @delete-link="showConfirmation(link.id)" />
    <button @click="addLink" class="text-blue-500 hover:underline mt-4">
      + Add New Link
    </button>

    <ConfirmationBox v-if="confirmationVisible" title="Confirm Deletion"
      message="Are you sure you want to delete this link?" @confirm="confirmDeleteLink" @cancel="cancelDelete" />

    <MessageBox v-if="messageVisible" :title="messageTitle" :message="message" :type="messageType"
      @close="messageVisible = false" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useLinkStore } from '@/stores/link'
import EditableLinkItem from './EditableLinkItem.vue'
import ConfirmationBox from '@/components/ConfirmationBox.vue'
import MessageBox from '@/components/MessageBox.vue'

const props = defineProps(['links'])
const emit = defineEmits(['update-links'])
const linkStore = useLinkStore()

const confirmationVisible = ref(false)
const linkIdToDelete = ref(null)
const messageVisible = ref(false)
const messageTitle = ref('')
const message = ref('')
const messageType = ref('')

const addLink = async () => {
  try {
    await linkStore.createLink({
      title: 'New Link',
      url: 'https://'
    })
    emit('update-links')
    showMessage('Link added successfully', 'success')
  } catch (error) {
    showMessage(error.message, 'error')
  }
}

const showConfirmation = (id) => {
  if (!confirmationVisible.value) {
    confirmationVisible.value = true
    linkIdToDelete.value = id
  }
}

const confirmDeleteLink = async () => {
  if (linkIdToDelete.value !== null) {
    try {
      await linkStore.deleteLink(linkIdToDelete.value)
      emit('update-links')
      showMessage('Link deleted successfully', 'success')
    } catch (error) {
      showMessage(error.message, 'error')
    } finally {
      confirmationVisible.value = false
      linkIdToDelete.value = null
    }
  }
}

const cancelDelete = () => {
  confirmationVisible.value = false
  linkIdToDelete.value = null
}

const showMessage = (msg, type) => {
  messageTitle.value = type === 'success' ? 'Success' : 'Error'
  message.value = msg
  messageType.value = type
  messageVisible.value = true
}
</script>
