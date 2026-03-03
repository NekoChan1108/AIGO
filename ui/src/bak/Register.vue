<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-950 px-4 py-24">
    <div v-motion-pop class="relative max-w-md w-full">
      <!-- Gopher positioned above card -->
      <div class="absolute -top-32 left-1/2 transform -translate-x-1/2 z-50 pointer-events-none">
        <InteractiveGopher />
      </div>
      
      <!-- Card -->
      <div class="relative bg-white dark:bg-slate-900 rounded-2xl shadow-2xl p-8 border border-slate-100 dark:border-slate-800 pt-16 mt-16">
        <div class="text-center mb-6">
          <h1 class="text-3xl font-bold text-slate-900 dark:text-white">Create Account</h1>
          <p class="text-slate-500 mt-2">Join AIGO Platform today</p>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Username</label>
            <input 
              v-model="form.username" 
              type="text" 
              class="w-full px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all" 
              required 
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Email</label>
            <div class="flex gap-2">
              <input 
                v-model="form.email" 
                :disabled="isEmailLocked" 
                type="email" 
                class="w-full px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all disabled:opacity-50 disabled:cursor-not-allowed" 
                required 
              />
              <button 
                v-if="isEmailLocked" 
                type="button" 
                @click="resetEmail" 
                class="px-3 py-2 text-sm text-primary-600 hover:text-primary-700 dark:text-primary-400 font-medium whitespace-nowrap hover:bg-primary-50 dark:hover:bg-primary-900/30 rounded-lg transition-colors"
              >
                Change
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Verification Code</label>
            <div class="flex gap-2">
              <input 
                v-model="form.verification_code" 
                type="text" 
                class="flex-grow px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all" 
                required 
              />
              <button 
                type="button" 
                @click="sendCode" 
                :disabled="codeSending || codeTimer > 0" 
                class="px-4 py-2 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-700 disabled:opacity-50 transition-colors whitespace-nowrap font-medium"
              >
                {{ codeTimer > 0 ? `${codeTimer}s` : (codeSending ? 'Sending...' : 'Get Code') }}
              </button>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Password</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="w-full px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all" 
              required 
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Confirm Password</label>
            <input 
              v-model="form.confirm_password" 
              type="password" 
              class="w-full px-4 py-2 rounded-lg border bg-transparent outline-none transition-all" 
              :class="[
                passwordMismatch 
                  ? 'border-red-500 focus:ring-2 focus:ring-red-500 text-red-600 focus:border-red-500' 
                  : 'border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-primary-500 focus:border-transparent'
              ]"
              required 
            />
            <p v-if="passwordMismatch" class="mt-1 text-xs text-red-500 font-medium animate-pulse">
              Passwords do not match
            </p>
          </div>

          <button 
            type="submit" 
            :disabled="loading" 
            class="w-full bg-primary-600 hover:bg-primary-700 text-white font-semibold py-2.5 rounded-lg transition-all transform active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center shadow-lg shadow-primary-500/30 hover:shadow-primary-500/50"
          >
            <svg v-if="loading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Creating Account...' : 'Sign Up' }}
          </button>
        </form>

        <div class="mt-6 text-center text-sm text-slate-500">
          Already have an account? 
          <router-link to="/login" class="text-primary-600 font-semibold hover:underline hover:text-primary-700 transition-colors">Log in</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { publicClient } from '../api/client'
import InteractiveGopher from '../components/InteractiveGopher.vue'

const router = useRouter()
const loading = ref(false)
const codeSending = ref(false)
const codeTimer = ref(0)
const isEmailLocked = ref(false)

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  verification_code: ''
})

const passwordMismatch = computed(() => {
  return form.confirm_password.length > 0 && form.password !== form.confirm_password
})

const sendCode = async () => {
  if (!form.email) {
    alert('Please enter your email first')
    return
  }
  
  codeSending.value = true
  try {
    await publicClient.post('/auth/verification', { email: form.email })
    
    alert('Verification code sent!')
    isEmailLocked.value = true
    codeTimer.value = 60
    const timer = setInterval(() => {
      codeTimer.value--
      if (codeTimer.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (e: any) {
    alert(e.message || 'Failed to send code')
  } finally {
    codeSending.value = false
  }
}

const resetEmail = () => {
  isEmailLocked.value = false
  codeTimer.value = 0
  form.verification_code = ''
}

const handleRegister = async () => {
  if (passwordMismatch.value) {
    alert('Passwords do not match')
    return
  }
  loading.value = true
  try {
    await publicClient.post('/auth/register', form)
    
    alert('Registration successful! Please login.')
    router.push('/login')
  } catch (e: any) {
    alert(e.message || 'Registration failed')
  } finally {
    loading.value = false
  }
}
</script>