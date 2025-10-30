<script setup lang="ts">
import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { computed, ref, onMounted } from "vue";
import { useRouter } from "vue-router";

import diagramCollectionAPI from "@/api/diagramCollectionServices";
import { useApiFeedback } from "@/composables/useApiFeedback";
import type {DiagramCollection} from "@/types/api.types.ts";

const router = useRouter();
const searchQuery = ref("");

const diagramCollections = ref<DiagramCollection[]>([]);
const loading = ref(true);

const { withApiFeedback } = useApiFeedback();

async function loadCollections() {
  loading.value = true;
  try {
    const data = await withApiFeedback(
        diagramCollectionAPI.getList(), // returns { list, ... }
        { errorMessage: "Failed to load collections." }
    );
    diagramCollections.value = data.list ?? [];
  } finally {
    loading.value = false;
  }
}

onMounted(loadCollections);

const filteredDiagramCollections = computed(() => {
  if (!searchQuery.value) return diagramCollections.value;
  const q = searchQuery.value.toLowerCase();
  return diagramCollections.value.filter(c =>
      c.name?.toLowerCase()?.includes(q) ||
      c.description?.toLowerCase()?.includes(q)
  );
});

const go = (id: string | number) => {
  router.push({ name: "MermaidDiagramTools-ProjectShow", params: { id } });
};
</script>
<template>
  <div class="min-h-[90vh] w-full flex justify-center rounded-lg bg-muted/10">
    <!-- Responsive main container -->
    <div class="h-full w-full sm:w-5/6 md:w-3/4 lg:w-2/4 mx-auto px-4 md:px-0">
      <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="p-5 space-y-7 flex-shrink-0">
          <div class="w-full space-y-3">
            <h1 class="text-2xl font-semibold text-gray-900 text-center lg:text-left">Diagram Project</h1>
            <div class="mx-auto lg:mx-0 w-12 h-0.5 bg-primary"></div>
          </div>

          <!-- Search Input -->
          <div class="w-full">
            <Input
                v-model="searchQuery"
                placeholder="Search collections..."
                class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg w-full"
            />
          </div>
        </div>

        <!-- Scrollable Cards Container -->
        <div class="flex-1 min-h-0 p-5 pt-0">
          <ScrollArea class="h-full">
            <div class="flex flex-col items-center space-y-5">
              <!-- Loading -->
              <div
                  v-if="loading"
                  class="text-center py-12 text-sm text-gray-600"
                  aria-live="polite"
                  aria-busy="true"
              >
                Loading collections…
              </div>

              <!-- Cards -->
              <Card
                  v-else
                  v-for="collection in filteredDiagramCollections"
                  :key="collection.id"
                  role="button"
                  tabindex="0"
                  @click="go(collection.id)"
                  @keyup.enter="go(collection.id)"
                  :aria-label="`Open ${collection.name}`"
                  class="group border border-gray-200 hover:border-gray-300 hover:shadow-md transition-all duration-200 bg-white rounded-lg cursor-pointer w-full"
              >
                <CardHeader class="p-5 py-0">
                  <div class="flex items-start justify-between">
                    <div class="flex-1 space-y-2">
                      <CardTitle class="text-base font-medium text-gray-900 group-hover:text-gray-800">
                        {{ collection.name }}
                      </CardTitle>
                      <CardDescription class="text-gray-600 text-sm leading-relaxed line-clamp-2">
                        {{ collection.description || 'No description provided.' }}
                      </CardDescription>
                    </div>
                  </div>
                </CardHeader>
              </Card>

              <!-- Empty State -->
              <div
                  v-if="!loading && filteredDiagramCollections.length === 0"
                  class="text-center py-12"
              >
                <div class="text-gray-400 mb-2">
                  <svg
                      class="mx-auto h-12 w-12"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                  >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                    />
                  </svg>
                </div>
                <h3 class="text-sm font-medium text-gray-900 mb-1">
                  No collections found
                </h3>
                <p class="text-sm text-gray-500">
                  Try adjusting your search terms.
                </p>
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
