<template>
  <div class="encryption-tools-page">
    <div class="tool-card">
      <div class="tool-form">
        <div class="tool-header">
          <h2>
            <i class="pi pi-shield"></i>
            Encryption Tools
          </h2>
          <!-- Action Buttons moved to header -->
        </div>
        <!-- Key Input -->
        <div class="key-section">
          <label for="encryptionKey" class="key-label">
            <i class="pi pi-key"></i>
            Encryption Key
          </label>
          <div class="key-input-group">
            <InputText
                id="encryptionKey"
                v-model="encryptionKey"
                placeholder="Enter your encryption key"
                :type="showKey ? 'text' : 'password'"
                class="key-input"
            />
            <Button
                :icon="showKey ? 'pi pi-eye-slash' : 'pi pi-eye'"
                @click="showKey = !showKey"
                class="p-button-text key-toggle"
                type="button"
                v-tooltip="showKey ? 'Hide key' : 'Show key'"
            />
          </div>
        </div>

        <!-- Data Section -->
        <div class="data-section">
          <!-- Input Data -->
          <div class="data-column">
            <div class="data-header">
              <label>
                <i class="pi pi-file-edit"></i>
                Input Data
              </label>
              <div class="data-actions">
                <Button
                    icon="pi pi-copy"
                    @click="copyToClipboard(inputData)"
                    class="p-button-text p-button-sm"
                    v-tooltip="'Copy input'"
                    :disabled="!inputData"
                />
                <Button
                    icon="pi pi-trash"
                    @click="inputData = ''"
                    class="p-button-text p-button-sm"
                    v-tooltip="'Clear input'"
                    :disabled="!inputData"
                />
                <Button
                    icon="pi pi-code"
                    @click="beautifyInputJson"
                    class="p-button-text p-button-sm"
                    v-tooltip="'Beautify JSON'"
                    :disabled="!isValidJson(inputData)"
                />
              </div>
            </div>
            <Textarea
                v-model="inputData"
                rows="12"
                placeholder="Enter the data you want to encrypt or decrypt..."
                class="data-textarea input-textarea"
            />
            <div v-if="isValidJson(inputData)" class="json-indicator">
              <i class="pi pi-check-circle"></i>
              Valid JSON detected
            </div>
          </div>

          <!-- Output Data -->
          <div class="data-column">
            <div>
              <div class="data-header">
                <label>
                  <i class="pi pi-file"></i>
                  Output Data
                </label>
                <div class="data-actions">
                  <Button
                      icon="pi pi-copy"
                      @click="copyToClipboard(outputData)"
                      class="p-button-text p-button-sm"
                      v-tooltip="'Copy output'"
                      :disabled="!outputData"
                  />
                  <Button
                      icon="pi pi-code"
                      @click="beautifyOutputJson"
                      class="p-button-text p-button-sm"
                      v-tooltip="'Beautify JSON'"
                      :disabled="!isValidJson(outputData)"
                  />
                </div>
              </div>
              <Textarea
                  v-model="outputData"
                  rows="12"
                  readonly
                  placeholder="Encrypted or decrypted data will appear here..."
                  class="data-textarea output-textarea"
              />

            </div>
            <div v-if="isValidJson(outputData)" class="json-indicator">
              <i class="pi pi-check-circle"></i>
              Valid JSON detected
            </div>
            <div class="header-actions" style="margin-top: 3rem;">
              <Button
                  label="Encrypt"
                  icon="pi pi-lock"
                  @click="encryptData"
                  :disabled="!encryptionKey || !inputData"
                  :loading="isProcessing"
                  severity="success"
              />
              <Button
                  label="Decrypt"
                  icon="pi pi-unlock"
                  @click="decryptData"
                  :disabled="!encryptionKey || !inputData"
                  :loading="isProcessing"
                  severity="info"
              />
              <Button
                  label="Swap"
                  icon="pi pi-arrow-right-arrow-left"
                  @click="swapInputOutput"
                  :disabled="!outputData"
                  outlined
              />
            </div>

          </div>
        </div>

        <!-- API Status -->

      </div>


    </div>

    <!-- Toast for notifications -->
    <Toast/>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {useToast} from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Toast from 'primevue/toast'

const toast = useToast()

// State
const encryptionKey = ref('')
const showKey = ref(false)
const inputData = ref('')
const outputData = ref('')
const isProcessing = ref(false)
const lastApiCall = ref<{ success: boolean; duration: number } | null>(null)

const API_CONFIG = {
  baseUrl: 'https://your-api-domain.com/api',
  endpoints: {
    encrypt: '/encrypt',
    decrypt: '/decrypt'
  }
}

// API Functions
const callEncryptionApi = async (endpoint: string, payload: any) => {
  const startTime = Date.now()
  try {
    const response = await fetch(`${API_CONFIG.baseUrl}${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${getAuthToken()}`
      },
      body: JSON.stringify(payload)
    })

    const duration = Date.now() - startTime

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({message: 'Unknown error'}))
      throw new Error(errorData.message || `HTTP ${response.status}`)
    }

    const data = await response.json()
    lastApiCall.value = {success: true, duration}

    return data
  } catch (error) {
    const duration = Date.now() - startTime
    lastApiCall.value = {success: false, duration}
    throw error
  }
}

const getAuthToken = (): string => {
  return localStorage.getItem('auth-token') || ''
}

// Encryption Functions
const encryptData = async () => {
  if (!encryptionKey.value || !inputData.value) return

  isProcessing.value = true
  try {
    const payload = {
      key: encryptionKey.value,
      data: inputData.value
    }

    const result = await callEncryptionApi(API_CONFIG.endpoints.encrypt, payload)

    outputData.value = result.encryptedData || result.data
    toast.add({
      severity: 'success',
      summary: 'Encryption Successful',
      detail: 'Data has been encrypted successfully',
      life: 3000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Encryption Failed',
      detail: error instanceof Error ? error.message : 'API request failed',
      life: 5000
    })
  } finally {
    isProcessing.value = false
  }
}

const decryptData = async () => {
  if (!encryptionKey.value || !inputData.value) return

  isProcessing.value = true
  try {
    const payload = {
      key: encryptionKey.value,
      data: inputData.value
    }

    const result = await callEncryptionApi(API_CONFIG.endpoints.decrypt, payload)

    outputData.value = result.decryptedData || result.data
    toast.add({
      severity: 'success',
      summary: 'Decryption Successful',
      detail: 'Data has been decrypted successfully',
      life: 3000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Decryption Failed',
      detail: error instanceof Error ? error.message : 'API request failed',
      life: 5000
    })
  } finally {
    isProcessing.value = false
  }
}

const swapInputOutput = () => {
  if (!outputData.value) return
  const temp = inputData.value
  inputData.value = outputData.value
  outputData.value = temp
}

// JSON Functions
const isValidJson = (str: string): boolean => {
  if (!str.trim()) return false
  try {
    JSON.parse(str)
    return true
  } catch {
    return false
  }
}

const beautifyInputJson = () => {
  if (!isValidJson(inputData.value)) return
  try {
    const parsed = JSON.parse(inputData.value)
    inputData.value = JSON.stringify(parsed, null, 2)
  } catch (error) {
    console.error('Failed to beautify input JSON:', error)
  }
}

const beautifyOutputJson = () => {
  if (!isValidJson(outputData.value)) return
  try {
    const parsed = JSON.parse(outputData.value)
    outputData.value = JSON.stringify(parsed, null, 2)
  } catch (error) {
    console.error('Failed to beautify output JSON:', error)
  }
}

// Utility Functions
const copyToClipboard = async (text: string) => {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    toast.add({
      severity: 'success',
      summary: 'Copied',
      detail: 'Content copied to clipboard',
      life: 2000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Copy Failed',
      detail: 'Failed to copy to clipboard',
      life: 3000
    })
  }
}
</script>

<style scoped>
.encryption-tools-page {
  width: 100%;
  max-width: 100%;
  padding: 0;
}

.page-header {
  margin-bottom: 2rem;
  text-align: center;
  padding: 2rem 0;
  background: var(--surface-ground);
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

.page-header h1 {
  font-size: 2.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  color: var(--text-color);
}

.page-header p {
  font-size: 1.125rem;
  color: var(--text-color-secondary);
  margin: 0;
}

.tool-card {
  background: var(--surface-card);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  border: 1px solid var(--surface-border);
  overflow: hidden;
}

.tool-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 2rem;
  background: var(--surface-ground);
  border-bottom: 1px solid var(--surface-border);
}

.tool-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin: 0;
  color: var(--text-color);
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: right;
}

.encrypt-btn {
  background: var(--green-500);
  border-color: var(--green-500);
}

.encrypt-btn:hover {
  background: var(--green-600);
  border-color: var(--green-600);
}

.decrypt-btn {
  background: var(--blue-500);
  border-color: var(--blue-500);
}

.decrypt-btn:hover {
  background: var(--blue-600);
  border-color: var(--blue-600);
}

.swap-btn {
  color: var(--primary-color);
  border-color: var(--primary-color);
}

.tool-form {
  padding: 2rem;
}

.key-section {
  margin-bottom: 2rem;
  background: var(--surface-ground);
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

.key-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: var(--text-color);
  font-size: 1rem;
}

.key-input-group {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.key-input {
  flex: 1;
  font-family: 'Courier New', monospace;
}

.key-toggle {
  flex-shrink: 0;
}

.data-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.data-column {
  min-width: 0;
}

.data-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--surface-border);
}

.data-header label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  color: var(--text-color);
  font-size: 1rem;
}

.data-actions {
  display: flex;
  gap: 0.25rem;
}

.data-textarea {
  width: 100%;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
  resize: vertical;
  border-radius: var(--border-radius);
  transition: all 0.2s ease;
}

.input-textarea {
  border-color: var(--surface-border);
}

.input-textarea:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 1px var(--primary-color);
}

.output-textarea {
  background-color: var(--surface-ground);
  border-color: var(--surface-border);
}

.json-indicator {
  font-size: 0.8rem;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  margin-top: 0.5rem;
  color: var(--green-600);
  font-weight: 500;
}

.api-status-section {
  padding: 1rem;
  background: var(--surface-ground);
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

.api-status {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.925rem;
  font-weight: 500;
}

.api-status.success {
  color: var(--green-600);
}

.api-status.error {
  color: var(--red-600);
}

.api-duration {
  margin-left: auto;
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
  color: var(--text-color-secondary);
  background: var(--surface-card);
  padding: 0.25rem 0.5rem;
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .page-header {
    padding: 1.5rem 1rem;
  }

  .page-header h1 {
    font-size: 2rem;
  }

  .tool-header {
    flex-direction: column;
    gap: 1.5rem;
    align-items: stretch;
    text-align: center;
  }

  .header-actions {
    justify-content: right;
    flex-wrap: wrap;
  }

  .tool-form {
    padding: 1rem;
  }

  .key-section {
    padding: 1rem;
  }

  .data-section {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }

  .data-header {
    flex-direction: column;
    gap: 0.75rem;
    align-items: stretch;
  }

  .data-actions {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .header-actions {
    flex-direction: column;
  }

  .data-actions {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>