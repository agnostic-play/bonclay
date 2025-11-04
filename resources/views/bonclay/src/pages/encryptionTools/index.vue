<template>
  <div class="min-h-screen p-5 pt-0">
    <div class="x-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent>
          <div class="mx-auto px-6 py-4">
            <!-- Header -->
            <h1 class="text-2xl font-semibold mb-8 flex items-center gap-2 text-gray-800">
              <Shield class="w-5 h-5 text-primary"/>
              Encryption Tools
            </h1>

            <!-- Encryption Type -->
            <div class="mb-6 space-y-3">
              <Label for="type" class="flex items-center gap-2">
                <LockKeyhole class="w-4 h-4"/>
                Encryption Type
              </Label>
              <select
                  id="type"
                  v-model="form.type"
                  class="w-full border rounded-md h-10 px-3 focus:ring-2 focus:ring-blue-200 outline-none"
              >
                <option v-for="opt in encryptionTypes" :key="opt" :value="opt">
                  {{ opt }}
                </option>
              </select>
            </div>

            <!-- Output Format -->
            <div class="mb-6 space-y-3">
              <Label for="outputFormat" class="flex items-center gap-2">
                <FileText class="w-4 h-4"/>
                Expected Output
              </Label>
              <select
                  id="outputFormat"
                  v-model="form.outputFormat"
                  class="w-full border rounded-md h-10 px-3 focus:ring-2 focus:ring-blue-200 outline-none"
              >
                <option value="base64">Base64</option>
                <option value="hex">Hex</option>
              </select>
            </div>

            <!-- Encryption Key -->
            <div class="mb-8 space-y-3">
              <Label for="key" class="flex items-center gap-2">
                <KeyRound class="w-4 h-4"/>
                Encryption Key
              </Label>
              <div class="relative">
                <Input
                    id="key"
                    v-model="form.key"
                    placeholder="Enter your encryption key"
                    :type="showKey ? 'text' : 'password'"
                    class="pr-10 focus:ring-2 focus:ring-blue-200"
                />
                <button
                    type="button"
                    class="absolute inset-y-0 right-0 px-3 text-gray-400 hover:text-gray-600 transition"
                    @click="toggleKeyVisibility"
                >
                  <Eye v-if="!showKey" class="w-4 h-4"/>
                  <EyeOff v-else class="w-4 h-4"/>
                </button>
              </div>
            </div>

            <!-- Key Derive Method (hidden for AES Imeg) -->
            <div v-if="form.type !== 'AES Imeg'" class="mb-8 space-y-3">
              <Label for="kdm" class="flex items-center gap-2">
                <KeyRound class="w-4 h-4"/>
                Key Derive Method
              </Label>
              <select
                  id="kdm"
                  v-model="form.keyDeriveMethod"
                  class="w-full border rounded-md h-10 px-3 focus:ring-2 focus:ring-blue-200 outline-none"
              >
                <option v-for="m in keyDeriveMethods" :key="m" :value="m">
                  {{ m }}
                </option>
              </select>
            </div>

            <!-- IV (only when needed) -->
            <div v-if="needsIV" class="mb-8 space-y-3">
              <Label for="iv" class="flex items-center gap-2">
                <PlusCircle class="w-4 h-4"/>
                Initialization Vector (IV)
              </Label>
              <Input
                  id="iv"
                  v-model="form.iv"
                  placeholder="Required for AES CBC IV / AES Imeg"
                  class="focus:ring-2 focus:ring-blue-200"
              />
            </div>

            <!-- Input / Output -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
              <!-- Input -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <Label for="input" class="flex items-center gap-2 text-lg font-semibold">
                    <FilePenLine class="w-4 h-4 text-gray-600"/>
                    Input Data
                  </Label>
                  <div class="flex gap-2">
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!form.input"
                        @click="copyText('input')"
                        class="disabled:opacity-50"
                        title="Copy Input"
                    >
                      <Copy class="w-4 h-4"/>
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!form.input"
                        @click="clearText('input')"
                        class="disabled:opacity-50"
                        title="Clear Input"
                    >
                      <Trash2 class="w-4 h-4"/>
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!isValidJson(form.input)"
                        @click="formatJson('input')"
                        class="disabled:opacity-50"
                        title="Format JSON"
                    >
                      <Braces class="w-4 h-4"/>
                    </Button>
                  </div>
                </div>

                <Textarea
                    id="input"
                    v-model="form.input"
                    placeholder="Enter the data you want to encrypt or decrypt..."
                    class="font-mono h-[400px] resize-none focus:ring-2 focus:ring-blue-200"
                />
              </div>

              <!-- Output -->
              <div class="space-y-3">
                <div class="flex items-center justify-between">
                  <Label for="output" class="flex items-center gap-2 text-lg font-semibold">
                    <FileOutput class="w-4 h-4 text-gray-600"/>
                    Output Data
                  </Label>
                  <div class="flex gap-2">
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!form.output"
                        @click="copyText('output')"
                        class="disabled:opacity-50"
                        title="Copy Output"
                    >
                      <Copy class="w-4 h-4"/>
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!form.output"
                        @click="clearText('output')"
                        class="disabled:opacity-50"
                        title="Clear Output"
                    >
                      <Trash2 class="w-4 h-4"/>
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        :disabled="!isValidJson(form.output)"
                        @click="formatJson('output')"
                        class="disabled:opacity-50"
                        title="Format JSON"
                    >
                      <Braces class="w-4 h-4"/>
                    </Button>
                  </div>
                </div>

                <Textarea
                    id="output"
                    v-model="form.output"
                    placeholder="Encrypted or decrypted data will appear here..."
                    class="font-mono h-[400px] bg-gray-50 resize-none text-gray-700"
                    readonly
                />
              </div>
            </div>

            <!-- Buttons -->
            <div class="flex flex-wrap justify-end gap-3 mt-10">
              <Button
                  :disabled="!form.key"
                  class="flex items-center gap-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="form.key ? 'bg-green-500 hover:bg-green-600 text-white' : 'bg-green-200 text-gray-500'"
                  @click="handleEncrypt"
              >
                <LockKeyhole/>
                Encrypt
              </Button>

              <Button
                  :disabled="!form.key"
                  class="flex items-center gap-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="form.key ? 'bg-blue-500 hover:bg-blue-600 text-white' : 'bg-blue-200 text-gray-500'"
                  @click="handleDecrypt"
              >
                <LockKeyholeOpen/>
                Decrypt
              </Button>

              <Button
                  variant="outline"
                  class="flex items-center gap-2 border-gray-300 text-gray-700 hover:bg-gray-100"
                  @click="swapInputOutput"
              >
                <ArrowRightLeft/>
                Swap
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { toast } from 'vue-sonner'
import {
  Eye,
  EyeOff,
  Shield,
  LockKeyhole,
  LockKeyholeOpen,
  ArrowRightLeft,
  FileText,
  KeyRound,
  PlusCircle,
  FilePenLine,
  FileOutput,
  Copy,
  Trash2,
  Braces
} from 'lucide-vue-next'
import { Card, CardContent } from '@/components/ui/card'
import encryptionToolServices from '@/api/EncryptionToolServices.ts'
import type { EncryptionToolReq } from '@/types/api.ts'
import { useApiFeedback } from '@/composables/useApiFeedback.ts'

const encryptionTypes = ['AES GCM', 'AES CBC', 'AES CBC IV', 'AES Imeg']
const keyDeriveMethods = ['No SUM', 'sha256', 'sha512', 'md5', 'sha1', 'hmac-sha256']

const form = ref({
  type: 'AES GCM',            // matches Go normalizeEncType
  key: '',
  iv: '',
  // AAD removed
  input: '',
  output: '',
  expected: '',
  outputFormat: 'base64',     // affects ENCRYPT only (decrypt returns base64 plaintext)
  keyDeriveMethod: 'No SUM'
})

const showKey = ref(false)

function toggleKeyVisibility() {
  showKey.value = !showKey.value
}

function swapInputOutput() {
  const tmp = form.value.input
  form.value.input = form.value.output
  form.value.output = tmp
}

function clearText(field: 'input' | 'output') {
  form.value[field] = ''
}

function copyText(field: 'input' | 'output') {
  if (!form.value[field]) return
  navigator.clipboard.writeText(form.value[field])
  toast.success(`${field === 'input' ? 'Input' : 'Output'} copied!`)
}

function isValidJson(text: string) {
  if (!text) return false
  try {
    JSON.parse(text)
    return true
  } catch {
    return false
  }
}

function formatJson(field: 'input' | 'output') {
  try {
    const obj = JSON.parse(form.value[field])
    form.value[field] = JSON.stringify(obj, null, 2)
    toast.success('JSON formatted!')
  } catch {
    toast.error('Invalid JSON')
  }
}

const normalize = (str: string) => str?.replace(/[\s:]/g, '').trim().toLowerCase()
const isExpectedMatch = computed(() => {
  if (!form.value.expected || !form.value.output) return false
  return normalize(form.value.expected) === normalize(form.value.output)
})

// Modes that require IV per backend (AES CBC IV, AES Imeg)
const needsIV = computed(() => ['AES CBC IV', 'AES Imeg'].includes(form.value.type))

const { withApiFeedback } = useApiFeedback()

async function callAPI(action: 'encrypt' | 'decrypt') {
  try {
    // Build payload exactly as the Go service expects
    const res = await withApiFeedback(
        encryptionToolServices.actEncryptDecrypt(action, {
          encryptionType: form.value.type,
          type: form.value.type,
          key: form.value.key,
          keyDeriveMethod: form.value.keyDeriveMethod,
          rawValue: form.value.input,
          outputFormat: form.value.outputFormat,
          iv: needsIV.value ? form.value.iv : undefined
        }),
        { errorMessage: `Failed to ${action}` }
    )
    if (action === 'decrypt') {
      form.value.output = atob(res.data?.output) || ''
      return
    }
    form.value.output = res.data?.output || ''
  } catch (err: any) {
    console.error(err)
    toast.error(err.message || 'API call failed.')
  }
}

function handleEncrypt() {
  if (!form.value.key || !form.value.input) {
    toast.error('Please provide both key and input.')
    return
  }

  if (needsIV.value && !form.value.iv) {
    toast.error('IV is required for this encryption type.')
    return
  }

  callAPI('encrypt')
}

function handleDecrypt() {
  if (!form.value.key || !form.value.input) {
    toast.error('Please provide both key and input.')
    return
  }
  if (needsIV.value && !form.value.iv) {
    toast.error('IV is required for this encryption type.')
    return
  }
  callAPI('decrypt')
}
</script>
