<template>
    <div class="p-4 bg-white rounded-lg shadow-md mb-4">
        <h2 class="text-2xl font-bold text-gray-800">{{ fullName }}</h2>
        <p class="text-gray-600">@{{ username }}</p>
        <p class="text-gray-700 mt-2">{{ bio }}</p>
    </div>
</template>

<script>
export default {
    props: {
        fullName: String,
        username: String,
        bio: String,
    },
};
</script>
