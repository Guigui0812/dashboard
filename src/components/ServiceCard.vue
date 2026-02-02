<script setup>
defineProps({
  name: String,
  url: String,
  logo: String,
  description: String
})
</script>

<template>
  <a href="https://nodejs.org/fr/download" target="_blank" rel="noopener noreferrer"
     class="flex-col flex justify-center rounded-lg bg-gray-700 m-8 shadow-lg transition-all duration-200 hover:scale-105 hover:bg-gray-600 w-1/6 h-80">
    <div class="text-center justify-items-center items-center">
      <img src="https://images.icon-icons.com/2699/PNG/512/grafana_logo_icon_171048.png" alt="logo" class="h-25 w-25">
      <div>
        <h3 class="text-xl font-bold text-white mt-6">Grafana</h3>
        <div class="flex text-white mt-4 justify-center font-bold">
            <p class="p-1 bg-orange-400 rounded mx-2">monitoring</p>
            <p class="p-1 bg-red-400 rounded mx-2">alerting</p>
            
        </div>
      </div>
    </div>
  </a>
</template>