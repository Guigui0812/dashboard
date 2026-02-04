<script setup>
import { ref, onMounted, onBeforeMount } from 'vue';
import ServiceCard from './components/ServiceCard.vue';
import yaml from 'js-yaml';

const services = ref([]);

onBeforeMount(async () => {
  try {
    const response = await fetch('/config/services.yaml');
    const yamlText = await response.text();
    const doc = yaml.load(yamlText);
    services.value = doc.services;
  } catch (e) {
    console.error("Erreur lors du chargement du YAML :", e);
  }
});
</script>

<template>
  <div class="min-h-screen dark:bg-gray-900 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-900">
    <h1 class="text-5xl font-extrabold text-center text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-emerald-500 mb-12 drop-shadow-lg">
      Mon Homelab Dashboard
    </h1>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 p-4">
      <ServiceCard
        v-for="service in services"
        :key="service.name"
        :name="service.display_name || service.name"
        :url="service.url"
        :logo="service.logo"
        :tags="service.tags || []"
      />
    </div>
  </div>
</template>
