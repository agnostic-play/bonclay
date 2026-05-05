<template>
  <div class="min-h-[90vh] w-full flex justify-center rounded-lg bg-muted/10">
    <!-- Responsive main container -->
    <div class="h-full w-full sm:w-5/6 md:w-3/4 lg:w-2/4 mx-auto px-4 md:px-0">
      <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="p-5 space-y-7 flex-shrink-0">
          <div class="w-full space-y-3">
            <div class="flex items-center justify-between">
              <div>
                <h1 class="text-2xl font-semibold text-gray-900 text-center lg:text-left">Diagram Project</h1>
                <div class="w-12 h-0.5 bg-primary mx-auto lg:mx-0"></div>
              </div>

              <!-- Create Collection Button -->
              <Button class="h-10" @click="openCreate">
                <FilePlus2 class="mr-2 h-4 w-4" />
                Create Collection
              </Button>
            </div>
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
                  v-for="collection in diagramCollections"
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
                  v-if="!loading && diagramCollections.length === 0"
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

    <!-- Create Collection Modal -->
    <Dialog v-model:open="showCreate">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Create Collection</DialogTitle>
          <DialogDescription>Provide a name and an optional description.</DialogDescription>
        </DialogHeader>

        <div class="grid gap-3 py-2">
          <div class="grid gap-1.5">
            <Label for="cc-name">Name</Label>
            <Input id="cc-name" v-model="createForm.name" placeholder="My Collection" />
          </div>
          <div class="grid gap-1.5">
            <Label for="cc-desc">Description</Label>
            <Textarea id="cc-desc" v-model="createForm.description" rows="3" placeholder="Optional details…" />
          </div>
        </div>

        <DialogFooter class="justify-end gap-2">
          <Button variant="outline" @click="showCreate = false">Cancel</Button>
          <Button :disabled="!createForm.name.trim() || isCreating" @click="submitCreate">
            <FilePlus2 class="mr-2 h-4 w-4" />
            <span v-if="!isCreating">Create</span>
            <span v-else>Creating…</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
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

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { toast } from "vue-sonner";

import { ScrollArea } from "@/components/ui/scroll-area";
import { Input } from "@/components/ui/input";
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { FilePlus2 } from "lucide-vue-next";

import diagramCollectionAPI from "@/api/diagramCollectionServices";
import { useApiFeedback } from "@/composables/useApiFeedback";
import type { DiagramCollection } from "@/types/entities";

const router = useRouter();
const searchQuery = ref("");

const diagramCollections = ref<DiagramCollection[]>([]);
const loading = ref(true);

const showCreate = ref(false);
const isCreating = ref(false);
const createForm = ref({ name: "", description: "" });

const { withApiFeedback } = useApiFeedback();

async function loadCollections(search = '') {
  loading.value = true;
  try {
    const data = await withApiFeedback(
        diagramCollectionAPI.getList(search ? { search } : undefined),
        { errorMessage: "Failed to load collections." }
    );
    diagramCollections.value = data.data?.list ?? [];
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  createForm.value = { name: "", description: "" };
  showCreate.value = true;
}

async function submitCreate() {
  if (!createForm.value.name.trim()) return;

  isCreating.value = true;
  try {
    const payload = {
      name: createForm.value.name.trim(),
      description: createForm.value.description ?? "",
    };

    // try actCreate first; fallback to create
    const svc: any =
        (diagramCollectionAPI as any).actCreate
            ? await withApiFeedback((diagramCollectionAPI as any).actCreate(payload), {
              successMessage: "Collection created",
              errorMessage: "Failed to create collection.",
            })
            : await withApiFeedback((diagramCollectionAPI as any).create(payload), {
              successMessage: "Collection created",
              errorMessage: "Failed to create collection.",
            });

    // normalize the created object (id may be in svc.data or svc directly)
    const created = svc?.data ?? svc;
    showCreate.value = false;
    await loadCollections();

    // navigate to detail if id is available
    const newId = String(created?.id ?? created?.data?.id ?? "");
    if (newId) {
      router.push({ name: "MermaidDiagramTools-ProjectShow", params: { id: newId } });
    } else {
      toast.success("Collection created");
    }
  } catch (e: any) {
    toast.error(e?.message || "Failed to create collection.");
  } finally {
    isCreating.value = false;
  }
}

let searchTimer: ReturnType<typeof setTimeout> | undefined

watch(searchQuery, (val) => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => void loadCollections(val), 300)
})

onMounted(() => loadCollections());

const go = (id: string | number) => {
  router.push({ name: "MermaidDiagramTools-ProjectShow", params: { id } });
};
</script>
