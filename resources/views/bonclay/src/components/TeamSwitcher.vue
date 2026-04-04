<script setup lang="ts">
import type { Component } from "vue"

import { ChevronsUpDown, Plus } from "lucide-vue-next"
import { ref, watch } from "vue"
import CreateSquadSheet from '@/components/CreateSquadSheet.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useRouter } from 'vue-router'
import { useSquadSession } from '@/composables/useSquadSession'

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"

const props = defineProps<{
  teams: {
    id?: string
    name: string
    logo: Component
    plan: string
    slug?: string
  }[]
}>()

const router = useRouter()
const { isMobile } = useSidebar()
const { activeSquad, setActiveSquad } = useSquadSession()
const activeTeam = ref(props.teams[0] ?? null)
const createSquadRef = ref<InstanceType<typeof CreateSquadSheet> | null>(null)

const emit = defineEmits<{
  squadCreated: []
}>()

// When teams are loaded/updated from parent, restore selection from session
// or default to the first team
watch(
  () => props.teams,
  (newVal) => {
    if (!Array.isArray(newVal) || newVal.length === 0) return

    // Try to restore the previously active squad from session
    if (activeSquad.value?.slug) {
      const match = newVal.find(t => t.slug === activeSquad.value!.slug)
      if (match) {
        activeTeam.value = match
        return
      }
    }

    // Default to first team if nothing in session
    if (!activeTeam.value || activeTeam.value === null) {
      activeTeam.value = newVal[0]
    }
  },
  { immediate: true }
)

const selectTeam = (team: any) => {
  activeTeam.value = team
  // Persist to session store
  if (team.id || team.slug) {
    setActiveSquad({
      id: team.id || '',
      slug: team.slug || '',
      name: team.name,
      plan: team.plan,
    })
  }
  if (team.slug) {
    router.push({ name: 'MockApiTools-Index', params: { slug: team.slug } })
  }
}

const handleAddTeam = () => {
  createSquadRef.value?.openDialog()
}

const handleSquadCreated = () => {
  emit('squadCreated')
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
              size="lg"
              class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <div class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
              <component :is="activeTeam.logo" class="size-4" />
            </div>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">
                {{ activeTeam.name }}
              </span>
              <span class="truncate text-xs">{{ activeTeam.plan }}</span>
            </div>
            <ChevronsUpDown class="ml-auto" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
            class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            align="start"
            :side="isMobile ? 'bottom' : 'right'"
            :side-offset="4"
        >
          <DropdownMenuLabel class="text-xs text-muted-foreground">
            Squads
          </DropdownMenuLabel>
          <DropdownMenuItem
            v-for="(team, index) in teams"
            :key="team.name"
            class="gap-2 p-2"
            :class="{ 'bg-accent': activeTeam?.slug === team.slug }"
            @click="selectTeam(team)"
          >
            <div class="flex size-6 items-center justify-center rounded-sm border">
              <component :is="team.logo" class="size-4 shrink-0" />
            </div>
            {{ team.name }}
            <DropdownMenuShortcut>⌘{{ index + 1 }}</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem class="gap-2 p-2" @click="handleAddTeam">
            <div class="flex size-6 items-center justify-center rounded-md border bg-background">
              <Plus class="size-4" />
            </div>
            <div class="font-medium text-muted-foreground">
              Add squad
            </div>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <CreateSquadSheet ref="createSquadRef" @created="handleSquadCreated" />
    </SidebarMenuItem>
  </SidebarMenu>
</template>
