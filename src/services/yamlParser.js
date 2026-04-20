import { fetchUrlStatus } from './serviceChecker.js';
import yaml from 'js-yaml';

export async function yamlParser(yamlFile, services) {
  var servicesList = []

  try {
    let response = await fetch(yamlFile);
    let yamlText = await response.text();
    let doc = yaml.load(yamlText);
    servicesList= doc.services;
    console.log("Services chargés :", doc.services);
  } catch (e) {
    console.error("Erreur lors du chargement du YAML :", e);
  }

  services.value = servicesList;

  for (let i=0; i < services.value.length; i++) {
    let service = services.value[i];
    try {
      let status = await fetchUrlStatus(service.url);
      console.log(`Service ${service.name} is ${status}`);
      service.status = status;

    } catch (e) {
      console.error(`Erreur lors de la vérification du service ${service.name} :`, e);
    }
  }

  return services
}