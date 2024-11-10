<template>
  <div class="p-4 bg-white rounded-lg shadow-md">
    <h3 class="text-lg font-semibold text-blue-600">
      <a :href="link.url" target="_blank" class="hover:underline" @click.prevent="handleClick">{{ link.title }}</a>
    </h3>
    <p class="text-gray-500 text-sm">{{ link.url }}</p>

    <LinkAnalytics v-if="isOwner" :analytics="link.analytics" class="mt-4" />
  </div>
</template>

<script setup>
import { useLinkStore } from '@/stores/link'
import LinkAnalytics from './LinkAnalytics.vue'

const props = defineProps(['link', 'isOwner'])
const linkStore = useLinkStore()

const handleClick = async () => {
  try {
    await linkStore.trackClick(props.link.id)
    window.open(props.link.url, '_blank')
  } catch (error) {
    console.error("Failed to track link click:", error)
  }
}
</script>
