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

            <!-- Mode Select (Encrypt / Decrypt) -->
            <div class="mb-8 space-y-3">
              <Label class="flex items-center gap-2">
                <LockKeyhole class="w-4 h-4"/>
                Mode
              </Label>
              <div class="flex gap-3">
                <Button
                    :variant="mode === 'encrypt' ? 'default' : 'outline'"
                    class="flex-1"
                    @click="mode = 'encrypt'"
                >
                  <LockKeyhole class="w-4 h-4 mr-2"/>Encrypt
                </Button>
                <Button
                    :variant="mode === 'decrypt' ? 'default' : 'outline'"
                    class="flex-1"
                    @click="mode = 'decrypt'"
                >
                  <LockKeyholeOpen class="w-4 h-4 mr-2"/>Decrypt
                </Button>
              </div>
            </div>

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
                  <div class="flex items-center gap-2">
                    <select
                        v-if="mode === 'decrypt'"
                        v-model="form.inputType"
                        aria-label="Input Format"
                        class="border rounded-md h-9 px-2 text-sm focus:ring-2 focus:ring-blue-200 outline-none"
                    >
                      <option v-for="opt in inputTypes" :key="opt" :value="opt">{{ opt }}</option>
                    </select>
                    <Button variant="outline" size="icon" :disabled="!form.input" @click="clearText('input')" title="Clear Input">
                      <Trash2 class="w-4 h-4"/>
                    </Button>
                    <Button variant="outline" size="icon" :disabled="!isValidJson(form.input)" @click="formatJson('input')" title="Beautify JSON">
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
                  <div class="flex items-center gap-2">
                    <select
                        v-if="mode === 'encrypt'"
                        v-model="form.outputFormat"
                        aria-label="Output Format"
                        class="border rounded-md h-9 px-2 text-sm focus:ring-2 focus:ring-blue-200 outline-none"
                    >
                      <option value="base64">Base64</option>
                      <option value="hex">Hex</option>
                    </select>
                    <Button variant="outline" size="icon" :disabled="!form.output" @click="clearText('output')" title="Clear Output">
                      <Trash2 class="w-4 h-4"/>
                    </Button>
                    <Button variant="outline" size="icon" :disabled="!isValidJson(form.output)" @click="formatJson('output')" title="Beautify JSON">
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

            <!-- Action Buttons -->
            <div class="flex justify-end gap-3">
              <Button
                  :disabled="!form.key"
                  class="flex items-center gap-2"
                  :class="mode === 'encrypt' ? 'bg-green-500 hover:bg-green-600 text-white' : 'bg-blue-500 hover:bg-blue-600 text-white'"
                  @click="mode === 'encrypt' ? handleEncrypt() : handleDecrypt()"
              >
                <component :is="mode === 'encrypt' ? LockKeyhole : LockKeyholeOpen" class="w-4 h-4"/>
                {{ mode === 'encrypt' ? 'Encrypt' : 'Decrypt' }}
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
  KeyRound,
  PlusCircle,
  FilePenLine,
  FileOutput,
  Trash2,
  Braces
} from 'lucide-vue-next'
import { Card, CardContent } from '@/components/ui/card'
import encryptionToolServices from '@/api/EncryptionToolServices.ts'
import { useApiFeedback } from '@/composables/useApiFeedback.ts'

const encryptionTypes = ['AES GCM', 'AES CBC', 'AES CBC IV', 'AES Imeg']
const keyDeriveMethods = ['No SUM', 'sha256', 'sha512', 'md5', 'sha1', 'hmac-sha256']
const inputTypes = ['Plain Text', 'Base64', 'Hex'] as const

type InputType = typeof inputTypes[number]
const mode = ref<'encrypt' | 'decrypt'>('encrypt')

const form = ref({
  type: 'AES GCM',
  key: '',
  iv: '',
  input: '',
  output: '',
  outputFormat: 'base64',
  keyDeriveMethod: 'No SUM',
  inputType: 'Plain Text' as InputType
})

const showKey = ref(false)
function toggleKeyVisibility() { showKey.value = !showKey.value }
function swapInputOutput() { const tmp = form.value.input; form.value.input = form.value.output; form.value.output = tmp }
function clearText(field: 'input' | 'output') { form.value[field] = '' }
function isValidJson(text: string) { if (!text) return false; try { JSON.parse(text); return true } catch { return false } }
function formatJson(field: 'input' | 'output') { try { const obj = JSON.parse(form.value[field]); form.value[field] = JSON.stringify(obj, null, 2); toast.success('JSON beautified!') } catch { toast.error('Invalid JSON') } }
const needsIV = computed(() => ['AES CBC IV', 'AES Imeg'].includes(form.value.type))
const { withApiFeedback } = useApiFeedback()

async function callAPI(action: 'encrypt' | 'decrypt') {
  try {
    const res = await withApiFeedback(
        encryptionToolServices.actEncryptDecrypt(action, {
          encryptionType: form.value.type.replace(/\s/g, ''),
          type: form.value.type.replace(/\s/g, ''),
          key: form.value.key,
          keyDeriveMethod: form.value.keyDeriveMethod,
          rawValue: form.value.input,
          outputFormat: form.value.outputFormat,
          inputFormat: form.value.inputType,
          iv: needsIV.value ? form.value.iv : undefined
        }),
        { errorMessage: `Failed to ${action}` }
    )

    if (action === 'decrypt') {
      try {
        form.value.output = res.data?.output
            ? atob(res.data.output)
            : ''
      } catch (e) {
        form.value.output = 'Invalid Base64'
      }
    }else{
      form.value.output = res.data?.output || ''
    }
  } catch (err: any) { toast.error(err.message || 'API call failed.') }
}

function handleEncrypt() { if (!guardCommon()) return; callAPI('encrypt') }
function handleDecrypt() { if (!guardCommon()) return; callAPI('decrypt') }
function guardCommon() {
  if (!form.value.key || !form.value.input) { toast.error('Please provide both key and input.'); return false }
  if (needsIV.value && !form.value.iv) { toast.error('IV is required for this encryption type.'); return false }
  return true
}
</script>
