<template>
  <div class="flex flex-col items-start text-dark gap-0.5">
    <h1 class="text-5xl font-bold leading-none">{{ time }}</h1>

    <div v-if="weather" class="flex items-center space-x-2 mt-2 text-sm sm:text-lg">
      <font-awesome-icon :icon="['fas', 'cloud']" />
      <span>{{ weather.temp }}° {{ weather.location }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import axios from 'axios'

const time = ref('00:00')
const weather = ref<null | { temp: number; location: string }>(null)

function updateTime() {
  const now = new Date()
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  time.value = `${hours}:${minutes}`
}

async function fetchWeatherFromCoords(lat: number, lon: number) {
  try {
    // wttr.in aceita lat,lon como parâmetro de cidade
    const res = await axios.get(`https://wttr.in/${lat},${lon}?format=j1`)

    const city = res.data.nearest_area[0].areaName[0].value
    const region = res.data.nearest_area[0].region[0].value
    const temp = parseInt(res.data.current_condition[0].temp_C)

    const location = `${city}, ${region}`

    weather.value = {
      temp: temp,
      location: location,
    }
  } catch (err) {
    console.error(err)
  }
}

async function fetchWeather() {
  try {
    // wttr.in aceita lat,lon como parâmetro de cidade
    const locRes = await axios.get('https://ipwho.is/')
    const { city, region, latitude, longitude } = locRes.data

    console.log(city, region)
    fetchWeatherFromCoords(latitude, longitude)
  } catch (err) {
    console.error(err)
  }
}

onMounted(() => {
  updateTime()
  const interval = setInterval(updateTime, 1000)
  fetchWeather()
  onUnmounted(() => clearInterval(interval))
})
</script>
