<template>
  <div v-if="visible" :class="`fixed top-15 right-5 p-4 rounded shadow-md ${bgColor} ${textColor}`"
    @click="visible = false">
    <span class="font-semibold">{{ title }}</span>
    <p>{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: 'MessageBox',
  props: {
    title: {
      type: String,
      default: 'Notification',
    },
    message: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      default: 'info',
    },
    duration: {
      type: Number,
      default: 5000,
    },
  },
  data() {
    return {
      visible: true,
    };
  },
  computed: {
    bgColor() {
      switch (this.type) {
        case 'success':
          return 'bg-green-100 border-green-500 text-green-700';
        case 'error':
          return 'bg-red-100 border-red-500 text-red-700';
        case 'warning':
          return 'bg-yellow-100 border-yellow-500 text-yellow-700';
        default:
          return 'bg-blue-100 border-blue-500 text-blue-700';
      }
    },
    textColor() {
      return 'border-l-4';
    },
  },
  mounted() {
    setTimeout(() => {
      this.visible = false;
    }, this.duration);
  },
};
</script>

<style scoped>
.fixed {
  z-index: 9999;
}
</style>
