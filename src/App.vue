<script setup>
import { ref, onMounted, onBeforeMount } from 'vue';
import ServiceCard from './components/ServiceCard.vue';
import { fetchUrlStatus } from './services/serviceChecker.js';
import { yamlParser } from './services/yamlParser.js';

const yamlFile = '/config/services.yaml'

const services = ref([]);
onBeforeMount(async () => {
  
  services = await yamlParser(yamlFile, services)
  
  });
</script>

<template>
  <div>
    <div class="flex text-4xl justify-items-start content-center font-bistrom bg-gradient-to-r text-green-300 my-6 ml-4 drop-shadow-lg">   
      <img class="w-15 self-center mr-2" src="/src/assets/avocado.png">
      <h1 class="self-center" >Avocado Dashboard</h1>
  </div>

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
