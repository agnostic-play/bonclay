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

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"

const props = defineProps<{
  teams: {
    name: string
    logo: Component
    plan: string
    slug?: string
  }[]
}>()

const router = useRouter()
const { isMobile } = useSidebar()
const activeTeam = ref(props.teams[0] ?? null)
const createSquadRef = ref<InstanceType<typeof CreateSquadSheet> | null>(null)

const emit = defineEmits<{
  squadCreated: []
}>()

// when teams are loaded/updated from parent, set activeTeam if not set
watch(
  () => props.teams,
  (newVal) => {
    if ((!activeTeam.value || activeTeam.value === null) && Array.isArray(newVal) && newVal.length) {
      activeTeam.value = newVal[0]
    }
  },
  { immediate: true }
)

const selectTeam = (team: any) => {
  activeTeam.value = team
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
            Teams
          </DropdownMenuLabel>
          <DropdownMenuItem
            v-for="(team, index) in teams"
            :key="team.name"
            class="gap-2 p-2"
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
              Add team
            </div>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <CreateSquadSheet ref="createSquadRef" @created="handleSquadCreated" />
    </SidebarMenuItem>
  </SidebarMenu>
</template>
