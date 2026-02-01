<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import CollectionHeader from '@/pages/mockApi/components/CollectionHeader.vue'
import ConfigurationSection from '@/pages/mockApi/components/ConfigurationSection.vue'
import ApiCategoriesSection from '@/pages/mockApi/components/ApiCategoriesSection.vue'
import { type EndpointByCategory, type Scenario } from '@/types/api.types'
import client from "@/api/bonClayHttpClient";
import { useRouter, useRoute } from 'vue-router'

const route = useRoute()

const getCollectionDetail = async (slug: string): Promise<EndpointByCategory> => {
  const res = await client.get<EndpointByCategory>(`/api/v2/collection/${slug}/endpoint_scenario`)
  return res
}

const collectionDetail = ref<EndpointByCategory | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

watch(
  () => route.params.collectionSlug,
  async (newSlug) => {
    const slug = newSlug ? String(newSlug) : undefined
    if (!slug) {
      collectionDetail.value = null
      return
    }

    loading.value = true
    try {
      const detail = await getCollectionDetail(slug)
      collectionDetail.value = detail
      error.value = null
    } catch (err) {
      error.value = 'Failed to fetch collection detail'
      console.error(err)
    } finally {
      loading.value = false
    }
  },
  { immediate: true }
)

// Configuration - these could come from collection metadata in the future
const baseUrl = 'https://mock.com/test123123123123213'
const documentationUrl = 'https://docs.example.com'

// Computed properties
const totalEndpoints = computed(() => {
  if (!collectionDetail.value) return 0
  return Object.values(collectionDetail.value).reduce((acc, endpoints) => acc + endpoints.length, 0)
})

const categoryCount = computed(() => {
  if (!collectionDetail.value) return 0
  return Object.keys(collectionDetail.value).length
})

const activeScenarioCount = computed(() => {
  if (!collectionDetail.value) return 0
  let count = 0
  Object.values(collectionDetail.value).forEach(endpoints => {
    endpoints.forEach(endpoint => {
      if (endpoint.active_scenario) count++
    })
  })
  return count
})

// Event handlers
const handleEditCollection = () => {
  console.log('Editing collection')
}

const handleCreateNewEndpoint = () => {
  console.log('Creating new endpoint')
}

const handleScenarioUpdate = () => {
  // Refresh collection data after scenario update
  const slug = route.params.collectionSlug ? String(route.params.collectionSlug) : undefined
  if (slug) {
    getCollectionDetail(slug).then(detail => {
      collectionDetail.value = detail
    })
  }
}
</script>

<template>
  <div class="min-h-screen bg-white p-5 pt-0 ">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="pt-5">
          <div v-if="loading" class="p-8 text-center">
            <div class="animate-pulse">Loading collection...</div>
          </div>
          
          <div v-else-if="error" class="p-8 text-center text-red-600">
            {{ error }}
          </div>
          
          <div v-else-if="collectionDetail">
            <div class="p-8 px-4 py-0">
              <div class="space-y-4">
                <CollectionHeader :total-endpoints="totalEndpoints" :active-scenario-count="activeScenarioCount"
                  :category-count="categoryCount" @edit-collection="handleEditCollection"
                  @create-new-endpoint="handleCreateNewEndpoint" />

                <ConfigurationSection :base-url="baseUrl" :documentation-url="documentationUrl" />
              </div>
            </div>

            <ApiCategoriesSection :categories="collectionDetail"
              @scenarios-updated="handleScenarioUpdate" />
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<style scoped>
@keyframes pulse {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>