<template>
  <div class="mt-2 bg-gray-100 p-2 rounded-lg">
    <p class="text-gray-700">
      <span class="font-semibold">Clicks:</span> {{ analytics.click_count }}
    </p>
    <p class="text-gray-700">
      <span class="font-semibold">Visitors:</span>
    </p>
    <ul class="list-disc list-inside text-gray-600">
      <li v-for="username in analytics.visitors_usernames" :key="username">
        {{ username }}
      </li>
    </ul>
  </div>
</template>

<script setup>
defineProps({
  analytics: {
    type: Object,
    required: true
  }
})
</script>
