<template>
  <div class="p-4 bg-white rounded-lg shadow-md">
    <div class="flex justify-between items-center">
      <input v-model="linkData.title" placeholder="Link title"
        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
      <button @click="showConfirmation" class="text-red-500 ml-2">Delete</button>
    </div>
    <input v-model="linkData.url" placeholder="Link URL"
      class="w-full px-4 py-2 mt-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500" />
    <button @click="saveLink" class="bg-blue-500 text-white px-4 py-2 rounded mt-2">
      Save
    </button>

    <MessageBox v-if="messageVisible" :title="messageTitle" :message="message" :type="messageType"
      @close="messageVisible = false" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useLinkStore } from '@/stores/link'
import ConfirmationBox from '@/components/ConfirmationBox.vue'
import MessageBox from '@/components/MessageBox.vue'

const props = defineProps(['link'])
const emit = defineEmits(['delete-link'])
const linkStore = useLinkStore()

const linkData = ref({
  title: '',
  url: ''
})

const messageVisible = ref(false)
const messageTitle = ref('')
const message = ref('')
const messageType = ref('')

watch(() => props.link, (newVal) => {
  if (newVal) {
    linkData.value = { ...newVal }
  }
}, { immediate: true })

const saveLink = async () => {
  try {
    await linkStore.updateLink(props.link.id, linkData.value)
    showMessage('Link updated successfully', 'success')
  } catch (error) {
    showMessage(error.message, 'error')
  }
}

const showConfirmation = () => {
  emit('delete-link', props.link.id)
}

const confirmDeleteLink = () => {
  emit('delete-link', props.link.id)
}

const showMessage = (msg, type) => {
  messageTitle.value = type === 'success' ? 'Success' : 'Error'
  message.value = msg
  messageType.value = type
  messageVisible.value = true
}
</script>
