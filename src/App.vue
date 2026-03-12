<script setup>
import { ref, onMounted, onBeforeMount } from 'vue';
import ServiceCard from './components/ServiceCard.vue';
import yaml from 'js-yaml';

async function fetchUrlStatus(url) {
  try {
    const response = await fetch("http://localhost:8080/proxy?url=" + encodeURIComponent(url), { method: 'GET' });
    if (!response.ok) {
      return 'error';
    }
    const jsonBody = await response.json();
    return jsonBody.status === 'ok' ? 'ok' : 'error';
  } catch (e) {
    return 'error';
  }
}

const services = ref([]);
onBeforeMount(async () => {
  try {
    const response = await fetch('/config/services.yaml');
    const yamlText = await response.text();
    const doc = yaml.load(yamlText);
    services.value = doc.services;
    console.log("Services chargés :", doc.services);
  } catch (e) {
    console.error("Erreur lors du chargement du YAML :", e);
  }



  for (let i=0; i < services.value.length; i++) {
    const service = services.value[i];
    try {
      const status = await fetchUrlStatus(service.url);
      console.log(`Service ${service.name} is ${status}`);
      service.status = status;

    } catch (e) {
      console.error(`Erreur lors de la vérification du service ${service.name} :`, e);
    }
  }
  });
</script>

<template>
  <div>
    <h1 class="text-5xl font-extrabold text-center text-transparent bg-clip-text bg-gradient-to-r from-green-400 to-emerald-500 mb-12 drop-shadow-lg">
      Mon Homelab Dashboard
    </h1>
    <div class="flex flex-wrap justify-center gap-4 max-w-8xl mx-auto">
      <ServiceCard
        v-for="service in services"
        :key="service.name"
        :name="service.display_name || service.name"
        :url="service.url"
        :logo="service.logo"
        :tags="service.tags || []"
        :status="service.status"
      />
    </div>
  </div>
</template>
