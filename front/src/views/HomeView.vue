<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { normalizeKey } from '@/utils/index'
import InterfaceView from '../components/InterfaceView.vue'
const data = ref({})
let timeId: number | undefined
const fetchData = async () => {
  try {
    const res = await fetch("/api/interface")
    let d = await res.json()
    d = normalizeKey.underlineToCamelCase(d);
    data.value = d.data || {};
    return d;
  } catch (error) {
    console.error(error)
    return {}
  }
}
onMounted(async () => {
  await fetchData()
  timeId = setInterval(() => {
    fetchData()
  }, 3000)
})
onUnmounted(() => {
  clearInterval(timeId)
})

</script>

<template>
  <main>
    <InterfaceView :data="data"></InterfaceView>
  </main>
</template>
