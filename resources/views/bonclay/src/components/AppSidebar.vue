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
import {type Component, computed} from "vue";
import TeamSwitcher from "@/components/TeamSwitcher.vue";


const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: "icon",
})


const route = useRoute()
const data = computed(() => ({
  user: {
    name: "shadcn",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  teams: [
    {
      name: "Acme Inc",
      logo: Clover,
      plan: "Enterprise",
    }
  ],
  navMain: [
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
  ],
}))
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <TeamSwitcher :teams="data.teams" />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="data.navMain"/>
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user"/>
    </SidebarFooter>
    <SidebarRail/>
  </Sidebar>
</template>
