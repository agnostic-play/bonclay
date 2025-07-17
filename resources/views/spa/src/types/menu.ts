export interface MenuItemType {
    id: string
    label: string
    icon: string
    route: string                    // Add this - direct route path
    iconColor?: 'purple' | 'orange' | 'green'
    badge?: {
        value: string
        severity?: 'danger' | 'success' | 'warning' | 'info'
    }
}

export interface MenuSectionType {
    id: string
    title: string
    isExpanded: boolean
    items: MenuItemType[]
}