<script setup lang="ts">
import { type SidebarProps} from '@/components/ui/sidebar'


import {
  Clover,
  SquareTerminal,
} from "lucide-vue-next"
import NavMain from '@/components/NavMain.vue'
import NavUser from '@/components/NavUser.vue'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'
import {useRoute} from "vue-router";
import {type Component, computed, ref, onMounted} from "vue";
import client from "@/api/bonClayHttpClient";
import TeamSwitcher from "@/components/TeamSwitcher.vue";


const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: "icon",
})


const route = useRoute()
const user = {
  name: "shadcn",
  email: "m@example.com",
  avatar: "/avatars/shadcn.jpg",
}

const teams = ref<{
  id?: string
  name: string
  logo: Component
  plan: string
  slug?: string
}[]>([
  {
    name: "Acme Inc",
    logo: Clover,
    plan: "Enterprise",
  }
])

const navMain = computed(() => [
  {
    title: "Mock API",
    url: "#",
    icon: SquareTerminal,
    isActive: route.path.startsWith("/tools/mock-api"),
    items: [
      {
        title: "API Collection",
        routeName: "MockApiTools-Index",
      },
      {
        title: "Histories",
        routeName: "MockApiTools-Index",
      },
    ],
  },
])

const getSquads = async () => {
  try {
    const res = await client.get<any[]>('/api/v2/squads/all')
    teams.value = (res || []).map((s: any) => ({
      id: s.id || '',
      name: s.name,
      logo: Clover,
      plan: s.desc || '',
      slug: s.slug || '',
    }))
  } catch (err) {
    console.error('Failed to fetch squads for sidebar', err)
  }
}

onMounted(() => {
  void getSquads()
})
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <TeamSwitcher :teams="teams" @squad-created="getSquads" />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="navMain"/>
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="user"/>
    </SidebarFooter>
    <SidebarRail/>
  </Sidebar>
</template>
