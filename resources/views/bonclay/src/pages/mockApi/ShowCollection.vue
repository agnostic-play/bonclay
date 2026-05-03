<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Card, CardContent } from '@/components/ui/card'
import {
  ResizablePanelGroup,
  ResizablePanel,
  ResizableHandle,
} from '@/components/ui/resizable'
import CollectionHeader from '@/pages/mockApi/components/CollectionHeader.vue'
import ConfigurationSection from '@/pages/mockApi/components/ConfigurationSection.vue'
import ApiCategoriesSection from '@/pages/mockApi/components/ApiCategoriesSection.vue'
import EditCollectionDialog from '@/pages/mockApi/components/EditCollectionDialog.vue'
import CreateEndpointDialog from '@/pages/mockApi/components/CreateEndpointDialog.vue'
import EnvironmentVariablesPanel from '@/pages/mockApi/components/EnvironmentVariablesPanel.vue'
import { type CollectionDetail } from '@/types/api.types'
import { getCollectionDetail } from '@/api'
import { useRoute } from 'vue-router'

const route = useRoute()

const collection = ref<CollectionDetail | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const fetchCollection = async (slug: string) => {
  loading.value = true
  try {
    collection.value = await getCollectionDetail(slug)
    error.value = null
  } catch (err) {
    error.value = 'Failed to fetch collection detail'
    console.error(err)
  } finally {
    loading.value = false
  }
}

watch(
  () => route.params.collectionSlug,
  async (newSlug) => {
    const slug = newSlug ? String(newSlug) : undefined
    if (!slug) { collection.value = null; return }
    await fetchCollection(slug)
  },
  { immediate: true }
)

const mockBaseUrl = computed(() => {
  const apiBase: string = (import.meta.env as any).VITE_API_BASE_URL ?? ''
  const origin = apiBase ? new URL(apiBase).origin : window.location.origin
  const slug = collection.value?.slug || String(route.params.collectionSlug)
  return `${origin}/mock/${slug}`
})

const totalEndpoints = computed(() => {
  if (!collection.value) return 0
  return Object.values(collection.value.endpoints).reduce((acc, eps) => acc + eps.length, 0)
})

const categoryCount = computed(() => {
  if (!collection.value) return 0
  return Object.keys(collection.value.endpoints).length
})

const activeScenarioCount = computed(() => {
  if (!collection.value) return 0
  let count = 0
  Object.values(collection.value.endpoints).forEach(eps => {
    eps.forEach(ep => { if (ep.active_scenario) count++ })
  })
  return count
})

const editDialogOpen = ref(false)
const endpointDialogOpen = ref(false)
const variablesOpen = ref(false)

const refresh = () => {
  const slug = route.params.collectionSlug ? String(route.params.collectionSlug) : undefined
  if (slug) fetchCollection(slug)
}
</script>

<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <ResizablePanelGroup direction="horizontal" class="!h-auto gap-0">

      <!-- Collection panel -->
      <ResizablePanel :default-size="100" :min-size="40">
        <Card class="overflow-hidden border shadow-sm bg-white">
          <CardContent class="pt-5">
            <div v-if="loading" class="p-8 text-center">
              <div class="animate-pulse">Loading collection...</div>
            </div>

            <div v-else-if="error" class="p-8 text-center text-red-600">
              {{ error }}
            </div>

            <div v-else-if="collection">
              <div class="px-4 py-0">
                <div class="space-y-4">
                  <CollectionHeader
                    :collection-id="collection.id"
                    :collection-name="collection.name"
                    :collection-desc="collection.desc"
                    :total-endpoints="totalEndpoints"
                    :active-scenario-count="activeScenarioCount"
                    :category-count="categoryCount"
                    :forward-proxy-url="collection.forward_proxy_url || undefined"
                    :is-proxy-enable="collection.is_proxy_enable"
                    @edit-collection="editDialogOpen = true"
                    @create-new-endpoint="endpointDialogOpen = true"
                    @open-variables="variablesOpen = !variablesOpen"
                    @updated="refresh"
                  />
                  <ConfigurationSection
                    :collection-id="collection.id"
                    :base-url="mockBaseUrl"
                    :documentation-url="collection.docs"
                  />
                </div>
              </div>

              <ApiCategoriesSection
                :categories="collection.endpoints"
                :base-url="mockBaseUrl"
                :is-proxy-enable="collection.is_proxy_enable"
                @scenarios-updated="refresh"
              />
            </div>
          </CardContent>
        </Card>
      </ResizablePanel>

      <!-- Variables panel (toggled) -->
      <template v-if="variablesOpen && collection">
        <ResizableHandle :with-handle="true" class="mx-3 bg-transparent data-[resize-handle-state=drag]:bg-border" />
        <ResizablePanel :default-size="35" :min-size="20" :max-size="55">
          <Card class="overflow-hidden border shadow-sm bg-white h-full">
            <EnvironmentVariablesPanel
              :collection-id="collection.id"
              :collection-slug="collection.slug"
              @close="variablesOpen = false"
            />
          </Card>
        </ResizablePanel>
      </template>

    </ResizablePanelGroup>

    <!-- Dialogs -->
    <EditCollectionDialog
      v-if="collection"
      v-model:open="editDialogOpen"
      :collection-id="collection.id"
      :initial-name="collection.name"
      :initial-desc="collection.desc"
      :initial-docs="collection.docs"
      :initial-forward-proxy-url="collection.forward_proxy_url"
      @updated="refresh"
    />

    <CreateEndpointDialog
      v-if="collection"
      v-model:open="endpointDialogOpen"
      :collection-id="collection.id"
      :existing-endpoints="Object.values(collection.endpoints).flat()"
      @created="refresh"
    />
  </div>
</template>
