<script setup lang="ts">
import {Copy, ExternalLink, Check, Plus, Clock, AlertCircle} from 'lucide-vue-next'
import {Button} from '@/components/ui/button'
import {Card, CardContent} from '@/components/ui/card'
import {ref, computed} from 'vue'
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from "@/components/ui/accordion"
import {Badge} from "@/components/ui/badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import {Checkbox} from "@/components/ui/checkbox"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"

// Types
interface Parameter {
  name: string
  type: string
  required: boolean
  description: string
}

interface Response {
  code: string
  description: string
  example?: string
}

interface Endpoint {
  method: string
  path: string
  name: string
  description: string
  parameters: Parameter[]
  responses: Response[]
}

interface ApiCategory {
  id: string
  name: string
  endpoints: Endpoint[]
}

interface Scenario {
  id: string
  endpointPath: string
  httpStatus: number
  name: string
  description: string
  responseDelay: number
  isActive: boolean
  createdAt: string
  response?: any
}

// State
const copiedItems = ref<Record<string, boolean>>({})
const activeScenarios = ref<Record<string, string>>({})

// Mock scenarios data
const scenarios = ref<Scenario[]>([
  {
    id: 'sc1',
    endpointPath: '/auth/login',
    httpStatus: 200,
    name: 'Successful Login',
    description: 'Returns valid JWT token with user data',
    responseDelay: 1,
    isActive: false,
    createdAt: '2024-01-15',
    response: {token: 'jwt_token', user: {id: 1, name: 'John Doe'}}
  },
  {
    id: 'sc2',
    endpointPath: '/auth/login',
    httpStatus: 401,
    name: 'Invalid Credentials',
    description: 'Returns error for wrong username/password',
    responseDelay: 0.5,
    isActive: false,
    createdAt: '2024-01-15',
    response: {error: 'Invalid credentials'}
  },
  {
    id: 'sc3',
    endpointPath: '/auth/login',
    httpStatus: 429,
    name: 'Rate Limited',
    description: 'Too many login attempts',
    responseDelay: 0,
    isActive: false,
    createdAt: '2024-01-16',
    response: {error: 'Too many requests'}
  },
  {
    id: 'sc4',
    endpointPath: '/auth/login',
    httpStatus: 500,
    name: 'Server Error',
    description: 'Simulates internal server error',
    responseDelay: 3,
    isActive: false,
    createdAt: '2024-01-16',
    response: {error: 'Internal server error'}
  }
])

// API configuration
const baseUrl = 'https://mock.com/test123123123123213'
const documentationUrl = 'https://docs.example.com'

// API endpoints data
const apiCategories: ApiCategory[] = [
  {
    id: 'authentication',
    name: 'Authentication',
    endpoints: [
      {
        method: 'POST',
        path: '/auth/login',
        name: 'Login',
        description: 'Authenticate user with email and password',
        parameters: [
          {name: 'email', type: 'string', required: true, description: 'User email address'},
          {name: 'password', type: 'string', required: true, description: 'User password'}
        ],
        responses: [
          {code: '200', description: 'Authentication successful', example: '{ "token": "jwt_token", "user": {...} }'},
          {code: '401', description: 'Invalid credentials'},
          {code: '429', description: 'Too many requests'}
        ]
      },
      {
        method: 'POST',
        path: '/auth/logout',
        name: 'Logout',
        description: 'Invalidate user session',
        parameters: [],
        responses: [
          {code: '200', description: 'Logout successful'}
        ]
      },
      {
        method: 'POST',
        path: '/auth/refresh',
        name: 'Refresh Token',
        description: 'Refresh authentication token',
        parameters: [
          {name: 'refreshToken', type: 'string', required: true, description: 'Refresh token'}
        ],
        responses: [
          {code: '200', description: 'Token refreshed successfully'},
          {code: '401', description: 'Invalid refresh token'}
        ]
      }
    ]
  },
  {
    id: 'transactions',
    name: 'Transactions',
    endpoints: [
      {
        method: 'GET',
        path: '/transactions',
        name: 'List Transactions',
        description: 'Retrieve paginated list of transactions with filtering options',
        parameters: [
          {name: 'page', type: 'number', required: false, description: 'Page number (default: 1)'},
          {name: 'limit', type: 'number', required: false, description: 'Items per page (default: 20, max: 100)'},
          {
            name: 'status',
            type: 'string',
            required: false,
            description: 'Filter by status (pending, completed, failed)'
          },
          {name: 'dateFrom', type: 'string', required: false, description: 'Start date (ISO 8601 format)'},
          {name: 'dateTo', type: 'string', required: false, description: 'End date (ISO 8601 format)'}
        ],
        responses: [
          {code: '200', description: 'List retrieved successfully'}
        ]
      },
      {
        method: 'POST',
        path: '/transactions',
        name: 'Create Transaction',
        description: 'Create a new transaction',
        parameters: [
          {name: 'amount', type: 'number', required: true, description: 'Transaction amount'},
          {name: 'currency', type: 'string', required: true, description: 'Currency code (USD, EUR, etc.)'},
          {name: 'description', type: 'string', required: false, description: 'Transaction description'},
          {name: 'categoryId', type: 'string', required: false, description: 'Category ID'}
        ],
        responses: [
          {code: '201', description: 'Transaction created successfully'},
          {code: '400', description: 'Invalid request data'}
        ]
      },
      {
        method: 'GET',
        path: '/transactions/{id}',
        name: 'Get Transaction',
        description: 'Retrieve a specific transaction by ID',
        parameters: [
          {name: 'id', type: 'string', required: true, description: 'Transaction ID'}
        ],
        responses: [
          {code: '200', description: 'Transaction retrieved successfully'},
          {code: '404', description: 'Transaction not found'}
        ]
      }
    ]
  },
  {
    id: 'master-data',
    name: 'Master Data',
    endpoints: [
      {
        method: 'GET',
        path: '/master/categories',
        name: 'Get Categories',
        description: 'Retrieve all available categories',
        parameters: [],
        responses: [
          {code: '200', description: 'Categories retrieved successfully'}
        ]
      },
      {
        method: 'GET',
        path: '/master/currencies',
        name: 'Get Currencies',
        description: 'Retrieve supported currencies',
        parameters: [],
        responses: [
          {code: '200', description: 'Currencies retrieved successfully'}
        ]
      }
    ]
  }
]

// For single accordion (only one item can be open at a time)
const openAccordionItem = ref<string | null>(null)

// Function to check if an accordion item is open
const isAccordionOpen = (itemValue: string): boolean => {
  return openAccordionItem.value === itemValue
}

// Handle accordion state changes
const handleAccordionChange = (value: string | null) => {
  openAccordionItem.value = value
}

// Computed
const totalEndpoints = computed(() =>
    apiCategories.reduce((acc, cat) => acc + cat.endpoints.length, 0)
)

const activeScenarioCount = computed(() =>
    scenarios.value.filter(s => s.isActive).length
)

// Methods
const copyToClipboard = async (text: string, itemId: string = 'default'): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text)
    copiedItems.value[itemId] = true
    setTimeout(() => {
      delete copiedItems.value[itemId]
    }, 2000)
  } catch (err) {
    console.error('Failed to copy text: ', err)
  }
}

const getMethodColor = (method: string): string => {
  const colors: Record<string, string> = {
    GET: 'bg-green-500 hover:bg-green-600',
    POST: 'bg-blue-500 hover:bg-blue-600',
    PUT: 'bg-orange-500 hover:bg-orange-600',
    DELETE: 'bg-red-500 hover:bg-red-600',
    PATCH: 'bg-purple-500 hover:bg-purple-600'
  }
  return colors[method] || 'bg-gray-500 hover:bg-gray-600'
}

const getMethodBgColor = (method: string): string => {
  const colors: Record<string, string> = {
    GET: 'bg-green-300/10',
    POST: 'bg-blue-300/10',
    PUT: 'bg-orange-300/10',
    DELETE: 'bg-red-300/10',
    PATCH: 'bg-purple-300/10'
  }
  return colors[method] || 'bg-gray-300/10'
}

const getStatusColor = (status: number): string => {
  if (status >= 200 && status < 300) return 'text-green-600 bg-green-50'
  if (status >= 300 && status < 400) return 'text-blue-600 bg-blue-50'
  if (status >= 400 && status < 500) return 'text-orange-600 bg-orange-50'
  return 'text-red-600 bg-red-50'
}

const getEndpointScenarios = (path: string): Scenario[] => {
  return scenarios.value.filter(s => s.endpointPath === path)
}

const toggleScenario = (scenarioId: string, endpointPath: string): void => {
  scenarios.value = scenarios.value.map(s => ({
    ...s,
    isActive: s.endpointPath === endpointPath ? s.id === scenarioId : s.isActive
  }))
  activeScenarios.value[endpointPath] = scenarioId
}

const createNewScenario = (endpointPath: string): void => {
  console.log('Creating new scenario for:', endpointPath)
  // Implementation for creating new scenario
}

const editEndpoint = (endpointPath: string): void => {
  console.log('Editing endpoint:', endpointPath)
  // Implementation for editing endpoint
}

const editCollection = (): void => {
  console.log('Editing collection')
  // Implementation for editing collection
}

const createNewEndpoint = (): void => {
  console.log('Creating new endpoint')
  // Implementation for creating new endpoint
}

const viewScenario = (scenarioId: string): void => {
  console.log('Viewing scenario:', scenarioId)
  // Implementation for viewing scenario details
}

const formatDelay = (seconds: number): string => {
  if (seconds === 0) return '0s'
  if (seconds === 1) return '1s'
  return `${seconds}s`
}

// Prevent event bubbling for copy button
const handleCopyClick = (event: Event, text: string, itemId: string): void => {
  event.stopPropagation()
  copyToClipboard(text, itemId)
}
</script>

<template>
  <div class="min-h-screen bg-white p-1 pt-0">
    <div class="mx-auto ">
      <!-- Main Container -->
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="pt-5">
          <div class="p-8 px-4 py-0">
            <div class="space-y-4">
              <!-- Title Section -->
              <div class="flex items-start justify-between">
                <div>
                  <h1 class="text-4xl font-bold mb-2 text-slate-900">
                    Rintis Collection API
                  </h1>
                  <p class="leading-relaxed max-w-4xl text-slate-600">
                    The Rintis Collection API provides a comprehensive mock service for testing payment processing,
                    transaction management, and authentication flows. Configure different response scenarios to simulate
                    various conditions including success states, errors, and edge cases for robust application testing.
                  </p>
                </div>
                <div class="flex gap-2">
                  <Button
                      variant="outline"
                      size="sm"
                      @click="editCollection"
                      class="h-9 px-4 font-normal"
                  >
                    Edit Collection
                  </Button>
                  <Button
                      size="sm"
                      @click="createNewEndpoint"
                      class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal"
                  >
                    <Plus class="w-4 h-4 mr-2"/>
                    New Endpoint
                  </Button>
                </div>
              </div>

              <!-- Stats Row -->
              <div class="flex flex-wrap gap-6 pt-4">
                <div class="flex items-center gap-2">
                  <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                  <span class="text-sm text-slate-700">{{ totalEndpoints }} Endpoints</span>
                </div>
                <div class="flex items-center gap-2">
                  <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
                  <span class="text-sm text-slate-700">{{ activeScenarioCount }} Active Scenarios</span>
                </div>
                <div class="flex items-center gap-2">
                  <div class="w-2 h-2 bg-orange-500 rounded-full"></div>
                  <span class="text-sm text-slate-700">{{ apiCategories.length }} Categories</span>
                </div>
              </div>

              <!-- Configuration Section -->
              <div class="bg-muted rounded-lg p-4">
                <div class="flex flex-col md:flex-row md:items-center gap-4">
                  <!-- Base URL -->
                  <div class="flex items-center gap-2 flex-1">
                    <span class="text-slate-700 font-medium text-sm">Base URL:</span>
                    <div class="flex items-center gap-2 bg-background rounded-md px-3 py-1.5 flex-1 border">
                      <code class="text-slate-900 font-mono text-xs">
                        {{ baseUrl }}
                      </code>
                      <Tooltip>
                        <TooltipTrigger asChild>
                          <Button
                              variant="ghost"
                              size="sm"
                              @click="copyToClipboard(baseUrl, 'baseUrl')"
                              class="h-6 w-6 p-0 hover:bg-slate-100"
                          >
                            <Check v-if="copiedItems['baseUrl']" class="w-3 h-3 text-green-600"/>
                            <Copy v-else class="w-3 h-3 text-slate-600"/>
                          </Button>
                        </TooltipTrigger>
                        <TooltipContent>
                          <p>Copy base URL</p>
                        </TooltipContent>
                      </Tooltip>
                    </div>
                  </div>

                  <!-- Documentation Link -->
                  <Button
                      variant="outline"
                      size="sm"
                      asChild
                  >
                    <a :href="documentationUrl" target="_blank" class="inline-flex items-center gap-2">
                      <ExternalLink class="w-3 h-3"/>
                      Documentation
                    </a>
                  </Button>
                </div>
              </div>
            </div>
          </div>

          <!-- API Categories -->
          <div class="p-8 pt-0 pb-0 space-y-7 mt-7">
            <div v-for="category in apiCategories" :key="category.id" class="space-y-1">
              <div class="flex items-center gap-3 mb-2 mt-4">
                <h2 class="text-2xl font-semibold text-slate-800">
                  {{ category.name }}
                </h2>
                <Badge variant="outline" class="text-xs">
                  {{ category.endpoints.length }} endpoints
                </Badge>
              </div>

              <Accordion
                  type="single"
                  class="w-full space-y-2"
                  collapsible
                  :model-value="openAccordionItem"
                  @update:model-value="handleAccordionChange"
              >
                <AccordionItem
                    v-for="endpoint in category.endpoints"
                    :key="`${category.id}-${endpoint.path}`"
                    :value="`${category.id}-${endpoint.path}`"
                    :class="[
                      'rounded-sm border transition-all duration-200',
                      isAccordionOpen(`${category.id}-${endpoint.path}`)
                        ? 'bg-white border-gray-200 shadow-sm'
                        : getMethodBgColor(endpoint.method)
                    ]"
                >
                  <AccordionTrigger class="no-underline hover:no-underline px-6  hover:bg-white/50 ">
                    <div class="flex items-center w-full space-x-3">
                      <!-- HTTP Method Badge -->
                      <Badge
                          :class="`text-white font-bold text-center min-w-[80px] justify-center ${getMethodColor(endpoint.method)}`"
                      >
                        {{ endpoint.method }}
                      </Badge>

                      <div
                          class="flex items-center gap-2 text-slate-700 font-mono text-sm bg-slate-100 px-3 rounded-lg">
                        <code>{{ endpoint.path }}</code>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                                variant="ghost"
                                size="sm"
                                @click="(event) => handleCopyClick(event, baseUrl + endpoint.path, baseUrl + endpoint.path)"
                                class="h-6 w-6 p-0 hover:bg-slate-200"
                            >
                              <Check v-if="copiedItems[baseUrl + endpoint.path]" class="w-3 h-3 text-green-600"/>
                              <Copy v-else class="w-3 h-3 text-slate-600"/>
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Copy full URL</p>
                          </TooltipContent>
                        </Tooltip>
                      </div>

                      <!-- Endpoint Name -->
                      <span class="ml-auto font-medium text-gray-800">
                        {{ endpoint.description }}
                      </span>
                    </div>
                  </AccordionTrigger>

                  <AccordionContent class="px-6 pb-6 bg-white rounded-b-xl">

                    <!-- Scenarios Section -->
                    <div class="space-y-4">
                      <div class="flex justify-between items-center">
                        <h4 class="text-sm font-medium text-gray-800 flex items-center gap-3">
                          <Clock class="w-4 h-4 text-gray-500"/>
                          Mock Scenarios
                        </h4>
                        <div class="flex gap-3">
                          <Button
                              variant="outline"
                              size="sm"
                              @click="editEndpoint(endpoint.path)"
                              class="h-8 px-3 font-normal text-xs"
                          >
                            Edit Endpoint
                          </Button>
                          <Button
                              size="sm"
                              @click="createNewScenario(endpoint.path)"
                              class="h-8 px-3 bg-gray-900 hover:bg-gray-800 font-normal text-xs"
                          >
                            <Plus class="w-3.5 h-3.5 mr-1.5"/>
                            New Scenario
                          </Button>
                        </div>
                      </div>

                      <div class="bg-gray-50/60 rounded-xl border border-gray-100 overflow-hidden">
                        <Table>
                          <TableHeader>
                            <TableRow class="bg-slate-50">
                              <TableHead class="w-[10%] text-slate-900 text-center font-semibold">HTTP Status
                              </TableHead>
                              <TableHead class="w-[55%] text-slate-900 font-semibold">Scenario</TableHead>
                              <TableHead class="w-[15%] text-slate-900 text-center font-semibold">Delay</TableHead>
                              <TableHead class="w-[10%] text-slate-900 text-center font-semibold">Active</TableHead>
                              <TableHead class="w-[10%] text-slate-900 text-center font-semibold">Actions</TableHead>
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            <TableRow
                                v-for="scenario in getEndpointScenarios(endpoint.path)"
                                :key="scenario.id"
                                :class="scenario.isActive ? 'bg-blue-50/50' : ''"
                            >
                              <TableCell class="text-center">
                                <Badge
                                    variant="secondary"
                                    :class="getStatusColor(scenario.httpStatus)"
                                >
                                  {{ scenario.httpStatus }}
                                </Badge>
                              </TableCell>
                              <TableCell>
                                <div>
                                  <div class="font-medium text-slate-800">{{ scenario.name }}</div>
                                  <div class="text-xs text-slate-500 mt-0.5">{{ scenario.description }}</div>
                                </div>
                              </TableCell>
                              <TableCell class="text-center">
                                <div class="flex items-center justify-center gap-1">
                                  <span class="text-sm">{{ formatDelay(scenario.responseDelay) }}</span>
                                </div>
                              </TableCell>
                              <TableCell class="text-center">
                                <Checkbox
                                    :checked="scenario.isActive"
                                    @update:checked="() => toggleScenario(scenario.id, endpoint.path)"
                                />
                              </TableCell>
                              <TableCell class="text-center">
                                <Button
                                    variant="ghost"
                                    size="sm"
                                    @click="viewScenario(scenario.id)"
                                    class="h-8 px-3 font-normal text-xs text-gray-600 hover:text-gray-800"
                                >
                                  View
                                </Button>
                              </TableCell>
                            </TableRow>
                            <TableRow v-if="getEndpointScenarios(endpoint.path).length === 0">
                              <TableCell colspan="5" class="text-center text-slate-500 py-8">
                                No scenarios configured. Click "New Scenario" to create one.
                              </TableCell>
                            </TableRow>
                          </TableBody>
                        </Table>
                      </div>
                    </div>
                  </AccordionContent>
                </AccordionItem>

              </Accordion>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<style scoped>
.accordion-trigger {
  text-decoration: none !important;
}

.accordion-trigger:hover,
.accordion-trigger:focus,
.accordion-trigger:active {
  text-decoration: none !important;
}

/* Custom animations */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>