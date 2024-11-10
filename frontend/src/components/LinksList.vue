<template>
  <div class="space-y-4">
    <LinkItem v-for="link in links" :key="link.id" :link="link" :isOwner="isOwner" />
  </div>
</template>

<script setup>
import { defineProps } from 'vue'
import LinkItem from './LinkItem.vue'

const props = defineProps({
  links: Array,
  isOwner: Boolean,
})
</script>
