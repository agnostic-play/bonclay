<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { computed, ref, onMounted } from 'vue'
import client from "@/api/bonClayHttpClient";

interface Squad {
  id: string
  name: string
  slug: string
  desc: string
  collections: Collection[]
}

interface Collection {
  id: string
  name: string
  docs: string
  slug: string
  desc: string
}

const searchQuery = ref('')

import { useRouter, useRoute } from 'vue-router'
const router = useRouter()
const route = useRoute()
const go = (id: string | number) => {
  router.push({ name: 'MockApiTools-CollectionShow', params: { id } })
}

const getSquadDetail = async (slug: string): Promise<Squad> => {
  const res = await client.get<Squad>(`/api/squad/detail/${slug}`)
  return res
}

const squads = ref<Squad | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

onMounted(async () => {
  loading.value = true;
  try {
    let slug: string | undefined = undefined

    if (route.query.slug) slug = String(route.query.slug)
    // if no slug provided (selected from sidebar), do nothing
    if (!slug) {
      squads.value = null
      return
    }

    const detail = await getSquadDetail(slug)
    squads.value = detail
  } catch (err) {
    error.value = "Failed to fetch squad detail";
    console.error(err);
  } finally {
    loading.value = false;
  }
});

const filteredCollections = computed(() => {
  if (!searchQuery.value) return squads.value?.collections || []

  return squads.value?.collections.filter(collection =>
    collection.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    collection.desc.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

</script>

<template>
  <div class="h-[90vh] w-full flex justify-center rounded-lg  bg-muted/10">
    <!--    <div class="min-h-[90vh] flex-1 rounded-xl bg-muted/50 md:min-h-min" >-->

    <div class="h-full w-2/4 ">
      <div class="flex flex-col h-full">
        <!-- Fixed Header Section -->
        <div class="p-5 space-y-7 flex-shrink-0 ">
          <!-- Header -->
          <div class="w-full space-y-3">
            <h1 class="text-2xl font-semibold text-gray-900">Api Collection's</h1>
            <div class="w-12 h-0.5 bg-primary"></div>
          </div>

          <!-- Search Input -->
          <div class="w-full">
            <Input v-model="searchQuery" placeholder="Search Collections..."
              class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg" />
          </div>
        </div>

        <!-- Scrollable Cards Container -->
        <div class="flex-1 min-h-0 p-5 pt-0">
          <ScrollArea class="h-full">
            <div class="space-y-5 pr-2">
              <Card v-for="collection in filteredCollections" :key="collection.id" role="button" tabindex="0" @click="go(collection.id)"
                @keyup.enter="go(collection.id)" :aria-label="`Open ${collection.name}`"
                class="group border border-gray-200 hover:border-gray-300 hover:shadow-md transition-all duration-200 bg-white rounded-lg cursor-pointer">
                <CardHeader class="p-5 py-0">
                  <div class="flex items-start justify-between">
                    <div class="flex-1 space-y-2">
                      <CardTitle class="text-base font-medium text-gray-900 group-hover:text-gray-800">
                        {{ collection.name }}
                      </CardTitle>
                      <CardDescription class="text-gray-600 text-sm leading-relaxed line-clamp-2">
                        {{ collection.desc }}
                      </CardDescription>
                    </div>
                    <div class="ml-4 flex-shrink-0">

                    </div>
                  </div>
                </CardHeader>
              </Card>

              <!-- Empty State -->
              <div v-if="filteredCollections?.length === 0" class="text-center py-12">
                <div class="text-gray-400 mb-2">
                  <svg class="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                </div>
                <h3 class="text-sm font-medium text-gray-900 mb-1">No collection found</h3>
                <p class="text-sm text-gray-500">Have you select the squad? or try adjusting your search terms.</p>
              </div>
            </div>
          </ScrollArea>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>