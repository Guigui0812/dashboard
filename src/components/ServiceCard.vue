<script setup>
defineProps({
  name: String,
  url: String,
  logo: String,
  tags: Array,
});
</script>

<template>
  <a
    :href="url"
    target="_blank"
    rel="noopener noreferrer"
    class="flex flex-col items-center justify-between bg-gray-800 rounded-xl p-6 m-2 shadow-xl hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-1 hover:scale-105 border border-gray-700 hover:border-green-500"
  >
    <div class="flex flex-col items-center">
      <div class="w-20 h-20 mb-4 flex items-center justify-center rounded-full bg-gray-700 p-2">
        <img :src="logo" alt="Logo" class="max-h-full max-w-full object-contain" />
      </div>
      <h3 class="text-2xl font-bold text-white mb-2 text-center">{{ name }}</h3>
    </div>
    <div class="flex flex-wrap justify-center gap-2 mt-4">
      <span
        v-for="tag in tags"
        class="px-3 py-1 text-sm font-semibold text-gray-100 bg-gradient-to-r from-orange-500 to-amber-500 rounded-full shadow-sm"
      >
        {{ tag }}
      </span>
    </div>
  </a>
</template>
