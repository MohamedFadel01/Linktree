<template>
    <div class="fixed inset-0 flex items-center justify-center z-50">
        <div class="bg-white p-6 rounded-lg shadow-lg max-w-sm w-full">
            <h3 class="text-lg font-semibold text-gray-800">{{ title }}</h3>
            <p class="mt-2 text-gray-600">{{ message }}</p>
            <div class="flex justify-end mt-4">
                <button @click="confirm" class="bg-red-600 text-white px-4 py-2 rounded-lg mr-2">Confirm</button>
                <button @click="cancel" class="bg-gray-300 text-gray-700 px-4 py-2 rounded-lg">Cancel</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { defineEmits, defineProps } from 'vue'

const props = defineProps({
    title: String,
    message: String,
})

const emit = defineEmits(['confirm', 'cancel'])

const confirm = () => {
    emit('confirm')
}

const cancel = () => {
    emit('cancel')
}
</script>

<style scoped>
.fixed {
    position: fixed;
}
</style>
