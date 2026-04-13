<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import { Card, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import CreateCollectionSheet from '@/pages/mockApi/components/CreateCollectionSheet.vue'
import { getSquadDetail, getSquads } from '@/api'
import type { SquadDetail, Collection } from '@/api'
import { useSquadSession } from '@/composables/useSquadSession'

const router = useRouter()
const route = useRoute()
const { activeSquadId, activeSquadSlug } = useSquadSession()

const searchQuery = ref('')
const squad = ref<SquadDetail | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const effectiveSlug = computed(() => {
  if (route.params.slug) return String(route.params.slug)
  return activeSquadSlug.value || activeSquadId.value || undefined
})

const go = (collection: Collection) => {
  const parentSlug = effectiveSlug.value || ''
  const collectionSlug = collection.slug || collection.id
  router.push({ name: 'MockApiTools-CollectionShow', params: { slug: parentSlug, collectionSlug } })
}

const fetchSquadDetail = async (slug: string) => {
  loading.value = true
  try {
    squad.value = await getSquadDetail(slug)
    error.value = null
  } catch (err) {
    error.value = 'Failed to fetch squad detail'
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  let slug = effectiveSlug.value
  if (!slug) {
    const res = await getSquads().catch(() => null)
    slug = res?.list[0]?.slug || res?.list[0]?.id
  }
  if (slug) await fetchSquadDetail(slug)
})

watch(effectiveSlug, (slug) => {
  if (slug) fetchSquadDetail(slug)
})

const handleCollectionCreated = () => {
  if (effectiveSlug.value) fetchSquadDetail(effectiveSlug.value)
}

const filteredCollections = computed(() => {
  const collections = squad.value?.collections ?? []
  if (!searchQuery.value) return collections
  const q = searchQuery.value.toLowerCase()
  return collections.filter(c =>
    c.name.toLowerCase().includes(q) || c.desc.toLowerCase().includes(q)
  )
})
</script>

<template>
  <div class="h-[90vh] w-full flex justify-center rounded-lg bg-muted/10">
    <div class="h-full w-2/4">
      <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="p-5 space-y-7 flex-shrink-0">
          <div class="w-full space-y-3">
            <div class="flex items-center justify-between">
              <h1 class="text-2xl font-semibold text-gray-900">API Collections</h1>
              <CreateCollectionSheet
                :squad-id="squad?.id || activeSquadId || undefined"
                :squad-slug="effectiveSlug"
                @created="handleCollectionCreated"
              />
            </div>
            <div class="w-12 h-0.5 bg-primary"></div>
          </div>

          <Input
            v-model="searchQuery"
            placeholder="Search Collections..."
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <!-- Collection List -->
        <div class="flex-1 min-h-0 p-5 pt-0">
          <ScrollArea class="h-full">
            <div class="space-y-5 pr-2">
              <Card
                v-for="collection in filteredCollections"
                :key="collection.id"
                role="button"
                tabindex="0"
                :aria-label="`Open ${collection.name}`"
                class="group border border-gray-200 hover:border-gray-300 hover:shadow-md transition-all duration-200 bg-white rounded-lg cursor-pointer"
                @click="go(collection)"
                @keyup.enter="go(collection)"
              >
                <CardHeader class="p-5 py-0">
                  <CardTitle class="text-base font-medium text-gray-900 group-hover:text-gray-800">
                    {{ collection.name }}
                  </CardTitle>
                  <CardDescription class="text-gray-600 text-sm leading-relaxed line-clamp-2">
                    {{ collection.desc }}
                  </CardDescription>
                </CardHeader>
              </Card>

              <!-- Empty State -->
              <div v-if="filteredCollections.length === 0" class="text-center py-12">
                <div class="text-gray-400 mb-3">
                  <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <h3 class="text-sm font-medium text-gray-900 mb-1">No collection found</h3>
                <p class="text-sm text-gray-500 mb-4">
                  Have you selected a squad? Or try adjusting your search terms.
                </p>
                <CreateCollectionSheet
                  :squad-id="squad?.id || activeSquadId || undefined"
                  :squad-slug="effectiveSlug"
                  @created="handleCollectionCreated"
                />
              </div>
            </div>
          </ScrollArea>
        </div>
      </div>
    </div>
  </div>
</template>
