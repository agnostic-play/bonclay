<script setup lang="ts">
import type {SidebarProps} from '@/components/ui/sidebar'


import {
  GalleryVerticalEnd,
  SquareTerminal,
} from "lucide-vue-next"
import NavMain from '@/components/NavMain.vue'
import NavUser from '@/components/NavUser.vue'
import TeamSwitcher from '@/components/TeamSwitcher.vue'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'
import {useRoute} from "vue-router";
import {computed} from "vue";

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
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    }
  ],
  navMain: [
      {
        title: "Mock API",
        url: "#",
        icon: SquareTerminal,
        isActive:true,
        items: [
          {
            title: "API Collection",
            isActive: route.name == "tesst",
            url: "/test",
          },
          {
            title: "Histories",
            isActive: route.name == "tesst",
            url: "/test",
          },
        ],
      },
    ],
}))


//
// // This is sample data.
// const data = {
//   user: {
//     name: "shadcn",
//     email: "m@example.com",
//     avatar: "/avatars/shadcn.jpg",
//   },
//   teams: [
//     {
//       name: "Acme Inc",
//       logo: GalleryVerticalEnd,
//       plan: "Enterprise",
//     }
//   ],
//   navMain: [
//     {
//       title: "Playground",
//       url: "#",
//       icon: SquareTerminal,
//       isActive: route.path == ,
//       items: [
//         {
//           title: "Home",
//           url: "/",
//         },
//
//       ],
//     },
//   ],
//   projects: [
//     {
//       name: "Design Engineering",
//       url: "#",
//       icon: Frame,
//     },
//     {
//       name: "Sales & Marketing",
//       url: "#",
//       icon: PieChart,
//     },
//     {
//       name: "Travel",
//       url: "#",
//       icon: Map,
//     },
//   ],
// }
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <TeamSwitcher :teams="data.teams"/>
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="data.navMain"/>
      <!--      <NavProjects :projects="data.projects" />-->
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user"/>
    </SidebarFooter>
    <SidebarRail/>
  </Sidebar>
</template>
