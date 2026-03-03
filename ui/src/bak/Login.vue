<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 dark:bg-slate-950 px-4 py-20">
    <div v-motion-pop class="relative max-w-md w-full">
      <!-- Gopher positioned above card -->
      <div class="absolute -top-32 left-1/2 transform -translate-x-1/2 z-50 pointer-events-none">
        <InteractiveGopher :mode="gopherMode" />
      </div>
      
      <!-- Card -->
      <div class="relative bg-white dark:bg-slate-900 rounded-2xl shadow-2xl p-8 border border-slate-100 dark:border-slate-800 pt-16 mt-16">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-slate-900 dark:text-white">Welcome Back</h1>
          <p class="text-slate-500 mt-2">Sign in to continue to AIGO</p>
        </div>

        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Username</label>
            <input 
              v-model="form.username" 
              @focus="gopherMode = 'reading'"
              @blur="gopherMode = 'idle'"
              type="text" 
              class="w-full px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all" 
              required 
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">Password</label>
            <input 
              v-model="form.password" 
              @focus="gopherMode = 'shy'"
              @blur="gopherMode = 'idle'"
              type="password" 
              class="w-full px-4 py-2 rounded-lg border border-slate-300 dark:border-slate-700 bg-transparent focus:ring-2 focus:ring-primary-500 focus:border-transparent outline-none transition-all" 
              required 
            />
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
            {{ loading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>

        <div class="mt-6 text-center text-sm text-slate-500">
          Don't have an account? 
          <router-link to="/register" class="text-primary-600 font-semibold hover:underline hover:text-primary-700 transition-colors">Sign up</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { publicClient } from '../api/client'
import { useRouter } from 'vue-router'
import { ref, reactive } from 'vue'
import InteractiveGopher from '../components/InteractiveGopher.vue'

const router = useRouter()
const loading = ref(false)
const gopherMode = ref('idle')
const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  loading.value = true
  try {
    const data: any = await publicClient.post('/auth/login', form)
    
    if (data.code === 200) {
      localStorage.setItem('access_token', data.data.access_token)
      localStorage.setItem('username', data.data.username)
      router.push('/chat') 
    } else {
      alert(data.msg || 'Login failed')
    }
  } catch (e: any) {
    console.error(e)
    alert(e.message || 'Network error')
  } finally {
    loading.value = false
  }
}
</script>