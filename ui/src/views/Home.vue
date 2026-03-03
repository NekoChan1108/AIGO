<template>
  <div class="min-h-screen flex flex-col relative overflow-hidden">
    <!-- Navbar -->
    <nav class="fixed w-full z-50 backdrop-blur-md bg-white/70 dark:bg-slate-950/70 border-b border-slate-200 dark:border-slate-800">
      <div class="container mx-auto px-6 h-16 flex items-center justify-between">
        <router-link to="/" class="flex items-center gap-3 group">
          <img 
            :src="logo" 
            alt="AIGO Logo" 
            class="h-10 w-auto group-hover:rotate-12 transition-transform duration-300 drop-shadow-md mix-blend-multiply dark:mix-blend-normal rounded-lg"
          />
          <span class="text-2xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary-600 to-indigo-600">
            AIGO
          </span>
        </router-link>
        <div class="space-x-4 flex items-center">
          <!-- Theme Toggle -->
          <button 
            @click="toggleTheme" 
            class="p-2 rounded-full text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
            :title="isDark ? 'Switch to Light Mode' : 'Switch to Dark Mode'"
          >
            <!-- Sun Icon (for Dark Mode) -->
            <svg v-if="isDark" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z" />
            </svg>
            <!-- Moon Icon (for Light Mode) -->
            <svg v-else xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z" />
            </svg>
          </button>

          <router-link to="/login" class="text-slate-600 hover:text-primary-600 font-medium transition-colors">Log In</router-link>
          <router-link to="/register" class="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg font-medium transition-all shadow-lg shadow-primary-500/30">
            Sign Up
          </router-link>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <main class="flex-grow pt-32 pb-0 px-6">
      <div class="container mx-auto grid lg:grid-cols-2 gap-12 items-center">
        <!-- Left: Text -->
        <div v-motion-slide-visible-once-left class="space-y-8 relative z-10">
          <div class="inline-flex items-center space-x-2 bg-primary-50 dark:bg-primary-900/30 text-primary-600 px-3 py-1 rounded-full text-sm font-medium">
            <span class="relative flex h-2 w-2">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2 w-2 bg-primary-500"></span>
            </span>
            <span>Powered by VolcEngine Ark & Eino</span>
          </div>
          
          <h1 class="text-5xl lg:text-7xl font-bold leading-tight">
            Enterprise-Grade <br/>
            <span class="text-primary-600">Agent Platform</span>
          </h1>
          <p class="text-xl text-slate-500 dark:text-slate-400 max-w-lg">
            Built with Go & Gin. Features vector retrieval, multi-model agent system, and real-time streaming responses.
          </p>
          <div class="flex items-center space-x-4">
            <router-link to="/register" class="bg-slate-900 dark:bg-white text-white dark:text-slate-900 px-8 py-3 rounded-full font-semibold hover:scale-105 transition-transform shadow-xl">
              Start Building ->
            </router-link>
            <a href="https://github.com/NekoChan1108/AIGO" target="_blank" class="text-slate-600 dark:text-slate-300 font-medium hover:underline flex items-center">
              <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
              View on GitHub
            </a>
          </div>
        </div>

        <!-- Right: Animated Gopher -->
        <div v-motion-slide-visible-once-right class="relative z-10 flex justify-center items-center h-[500px]">
          <div class="gopher-container relative w-[400px] h-[400px]">
            <!-- Animated Background -->
            <div class="gopher-bg absolute inset-0 bg-radial-gradient opacity-0"></div>
            
            <!-- Gopher Image -->
            <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-20 scale-90">
              <InteractiveGopher />
            </div>

            <!-- Surrounding AIGO Elements -->
            <div class="absolute inset-0 animate-spin-slow z-10">
               <div class="absolute top-1/2 -right-8 -translate-y-1/2 bg-white dark:bg-slate-800 px-4 py-2 rounded-xl shadow-lg border border-primary-100 dark:border-primary-900">
                <span class="font-bold text-indigo-600">Gin</span>
              </div>
              <div class="absolute -top-8 left-1/2 -translate-x-1/2 bg-white dark:bg-slate-800 px-4 py-2 rounded-xl shadow-lg border border-primary-100 dark:border-primary-900">
                <span class="font-bold text-indigo-600">Agent</span>
              </div>
              <div class="absolute -bottom-8 left-1/2 -translate-x-1/2 bg-white dark:bg-slate-800 px-4 py-2 rounded-xl shadow-lg border border-primary-100 dark:border-primary-900">
                <span class="font-bold text-purple-600">RAG</span>
              </div>
               <div class="absolute top-[15%] right-[15%] -translate-y-1/2 translate-x-1/2 bg-white dark:bg-slate-800 px-4 py-2 rounded-xl shadow-lg border border-primary-100 dark:border-primary-900">
                <span class="font-bold text-indigo-600">Eino</span>
                </div>
              <div class="absolute top-1/2 -left-8 -translate-y-1/2 bg-white dark:bg-slate-800 px-4 py-2 rounded-xl shadow-lg border border-primary-100 dark:border-primary-900">
                <span class="font-bold text-blue-600">Go</span>
              </div>
            </div>
            
            <!-- Orbit Rings -->
            <div class="absolute inset-0 rounded-full border border-dashed border-primary-200/50 dark:border-primary-700/30 animate-spin-reverse-slower pointer-events-none"></div>
            <div class="absolute inset-12 rounded-full border border-dotted border-indigo-200/50 dark:border-indigo-700/30 animate-spin-slow pointer-events-none"></div>
          </div>
        </div>
      </div>
    </main>

    <!-- System Insight Section (Typing Bubble) -->
    <section class="container mx-auto px-6 pt-0 pb-10 flex justify-center -mt-28 relative z-20">
      <div class="w-full max-w-4xl">
        <div class="mb-3 flex items-center justify-between px-2">
          <h3 class="text-sm font-semibold text-slate-400 uppercase tracking-widest">System Intelligence Overview</h3>
          <div class="flex gap-2">
            <span class="w-2 h-2 rounded-full bg-primary-500 animate-pulse"></span>
            <span class="text-[10px] text-primary-500 font-mono">CORE_SYSTEM_READY</span>
          </div>
        </div>
        <TypingBubble :full-text="gopherWisdom" />
      </div>
    </section>

    <!-- Tech Stack Carousel -->
    <section class="bg-slate-50 dark:bg-slate-900 py-12 border-t border-slate-200 dark:border-slate-800 overflow-hidden">
      <div class="container mx-auto px-6 mb-8 text-center">
        <h2 class="text-3xl font-bold mb-4">Powered by Modern Tech Stack</h2>
        <p class="text-slate-500 dark:text-slate-400">Built with industry-leading performance and scalability tools</p>
      </div>
      
      <div class="relative w-full overflow-hidden group/track">
        <div class="flex animate-scroll w-max hover:[animation-play-state:paused]">
          <!-- Double the list for infinite scroll effect -->
          <div class="flex space-x-16 px-8" v-for="i in 6" :key="i">
            <a 
              v-for="(tech, index) in techStack" 
              :key="index" 
              :href="tech.url"
              target="_blank"
              class="flex flex-col items-center justify-center space-y-4 group transition-all duration-300 min-w-[120px] cursor-pointer"
            >
              <div class="w-20 h-20 flex items-center justify-center bg-white dark:bg-slate-800 rounded-2xl shadow-sm group-hover:shadow-xl group-hover:-translate-y-2 transition-all p-4 border border-slate-100 dark:border-slate-700">
                <img v-if="tech.logo" :src="tech.logo" :alt="tech.name" class="w-full h-full object-contain" />
                <span v-else class="text-2xl">{{ tech.logo }}</span>
              </div>
              <span class="font-semibold text-slate-700 dark:text-slate-300 group-hover:text-primary-600 transition-colors">{{ tech.name }}</span>
            </a>
          </div>
        </div>
        
        <!-- Gradient masks -->
        <div class="absolute top-0 left-0 h-full w-24 bg-gradient-to-r from-slate-50 dark:from-slate-900 to-transparent z-10 pointer-events-none"></div>
        <div class="absolute top-0 right-0 h-full w-24 bg-gradient-to-l from-slate-50 dark:from-slate-900 to-transparent z-10 pointer-events-none"></div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import InteractiveGopher from '../components/InteractiveGopher.vue'
import TypingBubble from '../components/TypingBubble.vue'
const logo = '/logo.svg'

// Theme Logic
const isDark = useDark()
const toggleDark = useToggle(isDark)

const toggleTheme = () => {
  // Check for browser support
  if (!document.startViewTransition) {
    toggleDark()
    return
  }
  document.startViewTransition(() => {
    toggleDark()
  })
}

const gopherWisdom = ref("AIGO is built with a modern, practical stack: The backend uses Go (Gin) for core APIs, Eino for LLM orchestration, and VolcEngine Ark for model access. We rely on Milvus for vector retrieval (RAG), Redis for session caching, and MySQL/MongoDB for data persistence. The frontend is powered by Vue 3, Vite, and Tailwind CSS, using Pinia for state management and VueUse for smooth UI interactions. It's a clean, efficient architecture designed for intelligent, document-aware conversations.")

const techStack = [
  { name: 'Go', logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original.svg', url: 'https://go.dev/' },
  { name: 'Gin Web', logo: 'https://gin-gonic.com/_astro/gin.D6H2T_2v_ZD2G7l.webp', url: 'https://gin-gonic.com/' },
  { name: 'Eino', logo: 'https://avatars.githubusercontent.com/u/79236453?s=48&v=4', url: 'https://www.cloudwego.io/zh/docs/eino/' },
  { name: 'Milvus', logo: 'https://milvus.io/images/layout/milvus-logo.svg', url: 'https://milvus.io/' },
  { name: 'Redis', logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/redis/redis-original.svg', url: 'https://redis.io/' },
  { name: 'MySQL', logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/mysql/mysql-original.svg', url: 'https://www.mysql.com/' },
  { name: 'Kafka', logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/apachekafka/apachekafka-original.svg', url: 'https://kafka.apache.org/' },
  { name: 'VolcEngine', logo: 'https://portal.volccdn.com/obj/volcfe/misc/favicon.png', url: 'https://www.volcengine.com/' },
  { name: 'MongoDB', logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/mongodb/mongodb-original.svg', url: 'https://www.mongodb.com/' },
]
</script>

<style scoped>
.bg-radial-gradient {
  background-image: radial-gradient(#7d40ff, transparent 70%);
  animation: appear 1s ease-out forwards;
}

@keyframes appear {
  from { opacity: 0; transform: scale(0.8); }
  to { opacity: 0.2; transform: scale(1); }
}

@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes spin-reverse-slower {
  from { transform: rotate(360deg); }
  to { transform: rotate(0deg); }
}

.animate-spin-slow {
  animation: spin-slow 20s linear infinite;
}

.animate-spin-reverse-slower {
  animation: spin-reverse-slower 30s linear infinite;
}
</style>