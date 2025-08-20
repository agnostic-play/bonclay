<script setup lang="ts">
import {ref, computed} from 'vue'
import {Card, CardContent} from '@/components/ui/card'
import CollectionHeader from '@/pages/mockApi/components/CollectionHeader.vue'
import ConfigurationSection from '@/pages/mockApi/components/ConfigurationSection.vue'
import ApiCategoriesSection from '@/pages/mockApi/components/ApiCategoriesSection.vue'
import type {ApiCategory, Scenario} from '@/types/api.types'

// Mock data
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
        responses: [{code: '200', description: 'Logout successful'}]
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
        responses: [{code: '200', description: 'List retrieved successfully'}]
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
        responses: [{code: '200', description: 'Categories retrieved successfully'}]
      },
      {
        method: 'GET',
        path: '/master/currencies',
        name: 'Get Currencies',
        description: 'Retrieve supported currencies',
        parameters: [],
        responses: [{code: '200', description: 'Currencies retrieved successfully'}]
      }
    ]
  }
]

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

// Configuration
const baseUrl = 'https://mock.com/test123123123123213'
const documentationUrl = 'https://docs.example.com'

// Computed properties
const totalEndpoints = computed(() =>
    apiCategories.reduce((acc, cat) => acc + cat.endpoints.length, 0)
)

const activeScenarioCount = computed(() =>
    scenarios.value.filter(s => s.isActive).length
)

// Event handlers
const handleEditCollection = () => {
  console.log('Editing collection')
}

const handleCreateNewEndpoint = () => {
  console.log('Creating new endpoint')
}

const handleScenarioUpdate = (updatedScenarios: Scenario[]) => {
  scenarios.value = updatedScenarios
}
</script>

<template>
  <div class="min-h-screen bg-white p-5 pt-0 ">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="pt-5">
          <div class="p-8 px-4 py-0">
            <div class="space-y-4">
              <CollectionHeader
                  :total-endpoints="totalEndpoints"
                  :active-scenario-count="activeScenarioCount"
                  :category-count="apiCategories.length"
                  @edit-collection="handleEditCollection"
                  @create-new-endpoint="handleCreateNewEndpoint"
              />

              <ConfigurationSection
                  :base-url="baseUrl"
                  :documentation-url="documentationUrl"
              />
            </div>
          </div>

          <ApiCategoriesSection
              :categories="apiCategories"
              :scenarios="scenarios"
              @scenarios-updated="handleScenarioUpdate"
          />
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<style scoped>
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