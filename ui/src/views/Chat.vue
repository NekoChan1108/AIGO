<template>
  <div class="flex h-screen bg-white dark:bg-slate-950 overflow-hidden">
    <!-- Sidebar -->
    <aside 
      class="w-64 bg-slate-50 dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col transition-all duration-300"
      :class="{ '-translate-x-full absolute z-50 h-full': isMobile && !sidebarOpen, 'translate-x-0 relative': !isMobile || sidebarOpen }"
    >
      <!-- New Chat Button -->
      <div class="p-4">
        <button 
          @click="startNewChat"
          class="w-full flex items-center justify-between gap-3 px-4 py-3 bg-white dark:bg-slate-800 hover:bg-slate-100 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 rounded-xl transition-all shadow-sm group"
        >
          <span class="flex items-center gap-2 font-medium text-slate-700 dark:text-slate-200">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-5 h-5 text-primary-600">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
            </svg>
            New Chat
          </span>
          <span class="opacity-0 group-hover:opacity-100 transition-opacity">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 text-slate-400">
              <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
            </svg>
          </span>
        </button>
      </div>

      <!-- Session List -->
      <div class="flex-1 overflow-y-auto px-2 space-y-1 custom-scrollbar">
        <div v-if="loadingSessions" class="text-center py-4 text-slate-400 text-sm">Loading...</div>
        <button
          v-for="session in sessions"
          :key="session.id"
          @click="loadSession(session.id)"
          class="w-full text-left px-3 py-3 rounded-lg text-sm transition-colors flex items-center gap-3 group relative overflow-hidden"
          :class="currentSessionId === session.id ? 'bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-white font-medium' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800/50'"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 shrink-0 opacity-70">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 01.865-.501 48.172 48.172 0 003.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z" />
          </svg>
          <span class="truncate">{{ session.title || 'New Chat' }}</span>
          
          <!-- Delete Button (Visible on Hover/Active) -->
          <div 
            v-if="currentSessionId === session.id"
            class="absolute right-2 top-1/2 -translate-y-1/2 flex gap-1"
          >
             <button @click.stop="deleteSession(session.id)" class="p-1 text-slate-400 hover:text-red-500 rounded bg-white dark:bg-slate-800 shadow-sm transition-colors" title="Delete Chat">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                </svg>
             </button>
          </div>
        </button>
      </div>

      <!-- User Profile -->
      <div class="p-4 border-t border-slate-200 dark:border-slate-800">
        <div class="flex items-center gap-3 px-2">
          <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-primary-500 to-indigo-500 flex items-center justify-center text-white font-bold text-xs">
            {{ userInitials }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-slate-900 dark:text-white truncate">{{ username }}</p>
            <p class="text-xs text-slate-500 truncate">Pro Plan</p>
          </div>
          <button @click="logout" class="text-slate-400 hover:text-red-500 transition-colors">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15m3 0l3-3m0 0l-3-3m3 3H9" />
            </svg>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex relative h-full overflow-hidden">
      <!-- Chat Column -->
      <div class="flex-1 flex flex-col relative h-full transition-all duration-300 min-w-0">
        <!-- Chat Header (Desktop & Mobile) -->
        <div class="h-14 min-h-[56px] px-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between bg-white/80 dark:bg-slate-950/80 backdrop-blur z-10">
          <div class="flex items-center gap-3">
             <!-- Mobile Menu Button -->
             <button v-if="isMobile" @click="sidebarOpen = !sidebarOpen" class="text-slate-500">
               <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                 <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
               </svg>
             </button>
             <span class="font-bold text-slate-900 dark:text-white truncate max-w-[200px] md:max-w-md">
               {{ sessions.find(s => s.id === currentSessionId)?.title || 'AIGO Chat' }}
             </span>
          </div>

          <!-- File List Toggle -->
          <div class="relative" v-if="sessionFiles.length > 0">
            <button 
              @click="showFileList = !showFileList"
              class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 transition-colors text-sm font-medium text-slate-700 dark:text-slate-300"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244" />
              </svg>
              <span>{{ sessionFiles.length }}</span>
            </button>

            <!-- File List Popover -->
            <div 
              v-if="showFileList" 
              class="absolute right-0 top-full mt-2 w-64 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-200 dark:border-slate-700 overflow-hidden z-50 animate-fade-in"
            >
              <div class="p-3 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center bg-slate-50 dark:bg-slate-800/50">
                <span class="text-xs font-bold text-slate-500 uppercase">Attached Files</span>
                <button @click="showFileList = false" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-200">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <div class="max-h-64 overflow-y-auto p-1">
                <div 
                  v-for="(file, idx) in sessionFiles" 
                  :key="idx"
                  @click="file.originalFile ? handlePreview(file.originalFile) : null"
                  class="flex items-center gap-2 p-2 rounded-lg transition-colors cursor-pointer group"
                  :class="file.originalFile ? 'hover:bg-slate-100 dark:hover:bg-slate-800' : 'opacity-60 cursor-not-allowed'"
                  :title="!file.originalFile ? 'Preview unavailable for history files' : 'Click to preview'"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-slate-400 group-hover:text-primary-500">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                  </svg>
                  <span class="text-sm text-slate-700 dark:text-slate-300 truncate flex-1">{{ file.name }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Close Popover Overlay -->
          <div v-if="showFileList" @click="showFileList = false" class="fixed inset-0 z-40 bg-transparent cursor-default"></div>
        </div>

        <!-- Messages Area -->
        <div class="flex-1 overflow-y-auto p-4 md:p-8 scroll-smooth" ref="messagesContainer">
          <!-- Empty State -->
          <div v-if="messages.length === 0" class="h-full flex flex-col items-center justify-center text-center space-y-8 max-w-2xl mx-auto animate-fade-in">
            <img src="/logo.svg" alt="AIGO" class="w-24 h-24 mb-4" />
            <h2 class="text-3xl font-bold text-slate-900 dark:text-white">{{ greeting }}</h2>
            
            <!-- Capabilities Cards -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 w-full">
              <button class="p-4 bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-left transition-all group">
                <h3 class="font-medium text-slate-900 dark:text-white mb-1 group-hover:text-primary-600 transition-colors">Analyze Data</h3>
                <p class="text-sm text-slate-500">Upload a file and ask for insights</p>
              </button>
              <button class="p-4 bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-left transition-all group">
                <h3 class="font-medium text-slate-900 dark:text-white mb-1 group-hover:text-primary-600 transition-colors">Code Assistant</h3>
                <p class="text-sm text-slate-500">Help me write a Go function</p>
              </button>
            </div>
          </div>

          <!-- Chat History -->
          <div v-else class="space-y-6 max-w-3xl mx-auto pb-32">
            <div v-for="(msg, index) in messages" :key="index" class="flex gap-4" :class="msg.isUser ? 'flex-row-reverse' : ''">
              <!-- Avatar -->
              <div 
                class="w-8 h-8 rounded-full flex-shrink-0 flex items-center justify-center text-xs font-bold shadow-sm mt-1"
                :class="msg.isUser ? 'bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300' : 'bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700'"
              >
                <img v-if="!msg.isUser" src="/logo.svg" class="w-6 h-6" />
                <span v-else>{{ userInitials }}</span>
              </div>

              <!-- Message Bubble -->
              <div 
                class="px-5 py-3.5 rounded-2xl max-w-[85%] shadow-sm leading-relaxed"
                :class="[
                  msg.isUser 
                    ? 'bg-primary-600 text-white rounded-tr-sm' 
                    : 'bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 text-slate-800 dark:text-slate-200 rounded-tl-sm'
                ]"
              >
                <div v-if="msg.content" class="markdown-body text-sm md:text-base" :class="{'text-white': msg.isUser, 'text-slate-800 dark:text-slate-200': !msg.isUser}" v-html="renderMessage(msg.content)"></div>
                <!-- File Preview Link -->
                <div v-if="msg.files && msg.files.length > 0" class="mt-2 flex flex-wrap gap-2">
                  <div v-for="(file, idx) in msg.files" :key="idx" 
                       class="flex items-center gap-1 bg-slate-50/20 px-2 py-1 rounded text-xs border border-white/20 cursor-pointer hover:bg-slate-50/30 transition-colors"
                       @click.stop="handlePreview(file)">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3 h-3">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                    </svg>
                    <span class="truncate max-w-[150px]">{{ file.name }}</span>
                  </div>
                </div>
                <div v-else-if="!msg.content" class="flex gap-1 items-center h-6">
                  <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce"></span>
                  <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce delay-75"></span>
                  <span class="w-1.5 h-1.5 bg-current rounded-full animate-bounce delay-150"></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Input Area -->
        <div class="absolute bottom-0 left-0 w-full bg-gradient-to-t from-white via-white to-transparent dark:from-slate-950 dark:via-slate-950 p-4 md:p-6 z-20">
          <div class="max-w-3xl mx-auto relative">
            <!-- AIGO Header - Moved inside to avoid blocking -->
            <!-- <div class="absolute -top-8 left-0 flex items-center gap-2 text-xs font-bold text-slate-400 uppercase tracking-wider">
              <span class="w-2 h-2 bg-primary-500 rounded-full animate-pulse"></span>
              AIGO Intelligent Assistant
            </div> -->

            <!-- Input Container -->
            <div class="bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-2xl shadow-xl flex flex-col focus-within:ring-2 focus-within:ring-primary-500/50 focus-within:border-primary-500 transition-all overflow-hidden relative">
              
              <!-- Integrated Header -->
              <div class="absolute top-2 right-4 flex items-center gap-2 text-[10px] font-bold text-slate-300 dark:text-slate-600 uppercase tracking-wider pointer-events-none">
                <span class="w-1.5 h-1.5 bg-primary-500 rounded-full animate-pulse"></span>
                AIGO
              </div>
              
              <!-- File & Model Bar (Top of Input) -->
              <div v-if="files.length > 0" class="px-4 pt-3 flex gap-2 overflow-x-auto">
                <div v-for="(file, idx) in files" :key="idx" class="flex items-center gap-2 bg-slate-100 dark:bg-slate-800 px-3 py-1.5 rounded-lg text-xs font-medium text-slate-700 dark:text-slate-300">
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                  </svg>
                  {{ file.name }}
                  <button @click="removeFile(idx)" class="hover:text-red-500 ml-1">×</button>
                </div>
              </div>

              <!-- Textarea -->
              <textarea 
                v-model="inputMessage"
                @keydown.enter.prevent="sendMessage"
                placeholder="Message AIGO..." 
                rows="1"
                class="w-full px-4 py-4 bg-transparent border-none outline-none resize-none max-h-48 text-slate-900 dark:text-white placeholder-slate-400"
                style="min-height: 56px;"
              ></textarea>

              <!-- Action Bar -->
              <div class="flex items-center justify-between px-3 pb-3">
                <div class="flex items-center gap-2">
                  <!-- File Upload -->
                  <button @click="triggerFileUpload" class="p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors" title="Upload File">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M18.375 12.739l-7.693 7.693a4.5 4.5 0 01-6.364-6.364l10.94-10.94A3 3 0 1119.5 7.372L8.552 18.32m.009-.01l-.01.01m5.699-9.941l-7.81 7.81a1.5 1.5 0 002.112 2.13" />
                    </svg>
                  </button>
                  <input type="file" ref="fileInput" @change="handleFileSelect" class="hidden" multiple />

                  <!-- Voice Input -->
                  <button 
                    @click="toggleVoiceInput"
                    class="p-2 rounded-lg transition-colors relative group"
                    :class="isRecording ? 'text-red-500 bg-red-50 hover:bg-red-100 dark:bg-red-900/20 dark:hover:bg-red-900/30' : 'text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800'"
                    title="Voice Input"
                  >
                    <span v-if="isRecording" class="absolute -top-1 -right-1 w-2.5 h-2.5 bg-red-500 rounded-full animate-ping"></span>
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z" />
                    </svg>
                  </button>
                  
                  <!-- Model Selector (Mock) -->
                   <div class="h-4 w-[1px] bg-slate-200 dark:bg-slate-700 mx-1"></div>
                   <button class="flex items-center gap-1 text-xs font-medium text-slate-500 hover:text-primary-600 px-2 py-1 rounded hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors">
                     <span>GPT-4o</span>
                     <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-3 h-3">
                       <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.25 4.5a.75.75 0 01-1.08 0l-4.25-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
                     </svg>
                   </button>
                </div>

                <!-- Send Button -->
                <button 
                  @click="sendMessage" 
                  :title="!inputMessage.trim() && !isStreaming ? 'look like you haven\'t asked a question yet, try asking something' : 'Send Message'"
                  :class="[!inputMessage.trim() || isStreaming ? 'bg-slate-300 dark:bg-slate-700 opacity-50 cursor-not-allowed' : 'bg-primary-600 text-white hover:bg-primary-700 shadow-md shadow-primary-500/20']"
                  class="p-2 rounded-lg transition-all"
                >
                  <svg v-if="isStreaming" class="animate-spin w-5 h-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
                    <path d="M3.478 2.405a.75.75 0 00-.926.94l2.432 7.905H13.5a.75.75 0 010 1.5H4.984l-2.432 7.905a.75.75 0 00.926.94 60.519 60.519 0 0018.445-8.986.75.75 0 000-1.218A60.517 60.517 0 003.478 2.405z" />
                  </svg>
                </button>
              </div>
            </div>
            <div class="text-center mt-2">
               <p class="text-[10px] text-slate-400 dark:text-slate-600">AI can make mistakes. Please verify important information.</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Preview Column -->
      <div v-if="previewFile" class="w-full md:w-1/2 h-full border-l border-slate-200 dark:border-slate-800 bg-white dark:bg-slate-950 flex flex-col transition-all duration-300 absolute right-0 top-0 z-30 shadow-xl md:static md:shadow-none">
        <!-- Preview Header -->
        <div class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-800">
          <div class="flex items-center gap-2 overflow-hidden">
             <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 text-slate-500">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
             </svg>
             <span class="font-medium text-slate-900 dark:text-white truncate">{{ previewFile.name }}</span>
          </div>
          <button @click="closePreview" class="p-1 hover:bg-slate-100 dark:hover:bg-slate-800 rounded text-slate-500">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Preview Body -->
        <div class="flex-1 overflow-auto bg-slate-50 dark:bg-slate-900 p-4 relative">
          <!-- Image -->
          <div v-if="previewType === 'image'" class="h-full flex items-center justify-center">
            <img :src="previewUrl" class="max-w-full max-h-full object-contain rounded shadow-sm" />
          </div>
          
          <!-- PDF -->
          <iframe v-else-if="previewType === 'pdf'" :src="previewUrl" class="w-full h-full rounded shadow-sm border border-slate-200 dark:border-slate-700"></iframe>
          
          <!-- Text/Code -->
          <div v-else-if="previewType === 'text'" class="h-full bg-white dark:bg-slate-950 p-4 rounded border border-slate-200 dark:border-slate-800 overflow-auto">
            <pre class="text-xs md:text-sm font-mono text-slate-800 dark:text-slate-200 whitespace-pre-wrap break-words">{{ previewContent }}</pre>
          </div>
          
          <!-- Unsupported -->
          <div v-else class="h-full flex flex-col items-center justify-center text-slate-500 gap-4">
             <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1" stroke="currentColor" class="w-16 h-16 opacity-50">
               <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
             </svg>
             <p>Preview not available for this file type</p>
             <a :href="previewUrl" download class="text-primary-600 hover:underline">Download to view</a>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authClient } from '../api/client'
import { useWindowSize } from '@vueuse/core'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css'

const md:any = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }

    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})

const renderMessage = (content: string) => {
  return md.render(content)
}

const router = useRouter()
const { width } = useWindowSize()
const isMobile = computed(() => width.value < 768)
const sidebarOpen = ref(false)

// State
const sessions = ref<any[]>([])
const currentSessionId = ref('')
const messages = ref<any[]>([])
const inputMessage = ref('')
const isStreaming = ref(false)
const loadingSessions = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)
const files = ref<File[]>([])
const fileInput = ref<HTMLInputElement | null>(null)
const showFileList = ref(false)

// Computed Session Files
const sessionFiles = computed(() => {
  const allFiles: any[] = []
  messages.value.forEach(msg => {
    if (msg.files && msg.files.length > 0) {
      msg.files.forEach((f: any) => {
        allFiles.push({
          name: f.name,
          originalFile: f instanceof File ? f : null
        })
      })
    }
  })
  return allFiles
})

// User Info
const username = ref(localStorage.getItem('username') || 'User')
const userInitials = computed(() => username.value.charAt(0).toUpperCase())

// Greeting Logic
const greeting = computed(() => {
  const hour = new Date().getHours()
  
  let timeGreeting = 'Good morning'
  if (hour >= 12 && hour < 18) {
    timeGreeting = 'Good afternoon'
  } else if (hour >= 18) {
    timeGreeting = 'Good evening'
  }
  
  return `${timeGreeting} ${username.value} !`
})

// Load Sessions
const loadSessions = async () => {
  loadingSessions.value = true
  try {
    const res: any = await authClient.get('/llm/session/list')
    if (res.code === 200 && res.data) {
      console.log('Sessions:', res.data)
      sessions.value = res.data
    }
  } catch (e) {
    console.error('Failed to load sessions', e)
  } finally {
    loadingSessions.value = false
  }
}

// Load Chat History
const loadSession = async (id: string) => {
  currentSessionId.value = id
  sidebarOpen.value = false // Close sidebar on mobile
  try {
    const res: any = await authClient.post('/llm/chat/history', { session_id: id })
    if (res.code === 200 && res.data) {
      // Transform backend message format to frontend format
      messages.value = res.data.map((m: any) => ({
        content: m.content,
        isUser: m.is_user,
        files: m.files ? m.files.map((name: string) => ({ name })) : []
      }))
      scrollToBottom()
    }
  } catch (e) {
    console.error('Failed to load history', e)
  }
}

const startNewChat = () => {
  currentSessionId.value = ''
  messages.value = []
  sidebarOpen.value = false
}

const deleteSession = async (id: string) => {
  if (!confirm('Are you sure you want to delete this chat?')) return
  
  try {
    const res: any = await authClient.delete('/llm/session/delete', { session_id: id })
    if (res.code === 200) {
      sessions.value = sessions.value.filter(s => s.id !== id)
      if (currentSessionId.value === id) {
        startNewChat()
      }
    }
  } catch (e) {
    console.error('Failed to delete session', e)
  }
}

// File Upload
const triggerFileUpload = () => {
  fileInput.value?.click()
}

const handleFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    Array.from(target.files).forEach(file => {
      files.value.push(file)
    })
    // Reset input so same file can be selected again
    target.value = ''
  }
}

const removeFile = (index: number) => {
  files.value.splice(index, 1)
}

// File Preview Logic
const previewFile = ref<File | null>(null)
const previewUrl = ref('')
const previewType = ref('unsupported')
const previewContent = ref('')

const handlePreview = async (file: File) => {
  previewFile.value = file
  previewUrl.value = ''
  previewContent.value = ''
  
  if (file.type.startsWith('image/')) {
    previewType.value = 'image'
    previewUrl.value = URL.createObjectURL(file)
  } else if (file.type === 'application/pdf') {
    previewType.value = 'pdf'
    previewUrl.value = URL.createObjectURL(file)
  } else {
    // Try to read as text for preview
    try {
       // Only try for small files (< 1MB) to avoid freezing
       if (file.size < 1024 * 1024) {
           const text = await file.text()
           previewType.value = 'text'
           previewContent.value = text
       } else {
           previewType.value = 'unsupported'
           previewUrl.value = URL.createObjectURL(file)
       }
    } catch (e) {
      previewType.value = 'unsupported'
      previewUrl.value = URL.createObjectURL(file)
    }
  }
}

const closePreview = () => {
  if (previewUrl.value && previewType.value !== 'text') {
    URL.revokeObjectURL(previewUrl.value)
  }
  previewFile.value = null
}

const createObjectUrl = (file: File) => {
  if (file.type.startsWith('image/') || file.type === 'application/pdf') {
    return URL.createObjectURL(file)
  }
  // Force text/plain for other files to allow browser preview instead of download
  const blob = new Blob([file], { type: 'text/plain' })
  return URL.createObjectURL(blob)
}

// Send Message & Stream Response
const sendMessage = async () => {
  if (!inputMessage.value.trim()) return
  if (isStreaming.value) return

  const question = inputMessage.value
  inputMessage.value = ''
  
  // Add user message immediately
  const currentFiles = [...files.value]
  messages.value.push({ content: question, isUser: true, files: currentFiles })
  scrollToBottom()
  
  // Add placeholder for AI response
  const aiMessageIndex = messages.value.length
  messages.value.push({ content: '', isUser: false })
  
  isStreaming.value = true
  
  try {
    const token = localStorage.getItem('access_token')
    
    let url = ''
    let body: any = {}
    let isNewSession = false

    if (currentSessionId.value) {
      // Existing Session: Upload files first, then send message
      if (files.value.length > 0) {
        const formData = new FormData()
        formData.append('session_id', currentSessionId.value)
        files.value.forEach(file => formData.append('file', file))
        try {
          await authClient.post('/file/upload', formData)
          files.value = []
        } catch (e) {
          console.error('File upload failed', e)
        }
      }
      
      url = '/api/v1/llm/chat/send/stream'
      body = {
        question,
        session_id: currentSessionId.value,
        llm_type: 'ark',
        files: files.value.map(f => f.name)
      }
    } else {
      // New Session: Create session first, then upload files once we get the ID
      isNewSession = true
      url = '/api/v1/llm/chat/create/stream'
      body = {
        question,
        llm_type: 'ark',
        files: files.value.map(f => f.name)
      }
    }

    // Use native fetch for streaming
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(body)
    })

    if (!response.ok) throw new Error(response.statusText)
    
    const reader = response.body?.getReader()
    const decoder = new TextDecoder()
    
    if (!reader) throw new Error('No reader')

    let filesUploaded = false
    const pendingFiles = [...files.value] // Copy files to upload
    if (isNewSession) files.value = [] // Clear UI immediately

    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      
      const chunk = decoder.decode(value, { stream: true })
      const lines = chunk.split('\n')
      
      for (const line of lines) {
        if (line.startsWith('data:')) {
            try {
                const dataStr = line.slice(5).trim()
                if (!dataStr) continue
                if (dataStr === '[DONE]') break
                
                let handled = false
                // Check if it's a JSON control message (like session_id)
                if (dataStr.startsWith('{')) {
                    try {
                        const data = JSON.parse(dataStr)
                        console.log('Stream JSON:', data)

                        // Handle Session ID from Backend (New Session Flow)
                        if (data.session_id && isNewSession && !currentSessionId.value) {
                            currentSessionId.value = data.session_id
                            
                            // Optimistically add to session list with correct title
                            const newSession = {
                              id: data.session_id,
                              title: question,
                              created_at: new Date().toISOString()
                            }
                            sessions.value.unshift(newSession)

                            // Trigger async file upload now that we have the ID
                            if (pendingFiles.length > 0 && !filesUploaded) {
                                filesUploaded = true
                                const formData = new FormData()
                                formData.append('session_id', data.session_id)
                                pendingFiles.forEach(file => formData.append('file', file))
                                // Non-blocking upload
                                authClient.post('/file/upload', formData).catch(e => console.error('Async upload failed', e))
                            }
                            handled = true
                        }
                    } catch (e) {
                        // Not valid JSON, treat as content
                    }
                }

                // If not handled as a control message, treat as content
                if (!handled) {
                    messages.value[aiMessageIndex].content += dataStr
                }
                
                scrollToBottom()
            } catch (e) {
                console.error('Parse error', e)
            }
        }
      }
    }

  } catch (e) {
    console.error('Stream error', e)
    messages.value[aiMessageIndex].content += '\n[Error: Failed to get response]'
  } finally {
    isStreaming.value = false
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const logout = () => {
  localStorage.removeItem('access_token')
  router.push('/')
}

onMounted(() => {
  loadSessions()
})

// Voice Input Logic
const isRecording = ref(false)
let recognition: any = null

const toggleVoiceInput = () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    startRecording()
  }
}

const startRecording = () => {
  if (!('webkitSpeechRecognition' in window) && !('SpeechRecognition' in window)) {
    alert('Speech recognition is not supported in this browser. Please use Chrome or Edge.')
    return
  }

  const SpeechRecognition = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition
  recognition = new SpeechRecognition()
  recognition.continuous = true
  recognition.interimResults = true
  recognition.lang = 'zh-CN'

  recognition.onstart = () => {
    isRecording.value = true
  }

  recognition.onend = () => {
    isRecording.value = false
  }

  recognition.onresult = (event: any) => {
    for (let i = event.resultIndex; i < event.results.length; ++i) {
      if (event.results[i].isFinal) {
        inputMessage.value = (inputMessage.value + event.results[i][0].transcript).trim() + ' '
      }
    }
  }

  recognition.onerror = (event: any) => {
    console.error('Speech recognition error', event.error)
    isRecording.value = false
  }

  recognition.start()
}

const stopRecording = () => {
  if (recognition) {
    recognition.stop()
  }
  isRecording.value = false
}
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 20px;
}

/* Markdown Styles */
:deep(.markdown-body) {
  line-height: 1.6;
}
:deep(.markdown-body p) {
  margin-bottom: 0.75em;
}
:deep(.markdown-body p:last-child) {
  margin-bottom: 0;
}
:deep(.markdown-body h1), :deep(.markdown-body h2), :deep(.markdown-body h3) {
  font-weight: 600;
  margin-top: 1em;
  margin-bottom: 0.5em;
  line-height: 1.3;
}
:deep(.markdown-body h1) { font-size: 1.5em; }
:deep(.markdown-body h2) { font-size: 1.25em; }
:deep(.markdown-body h3) { font-size: 1.1em; }
:deep(.markdown-body ul), :deep(.markdown-body ol) {
  padding-left: 1.5em;
  margin-bottom: 0.75em;
}
:deep(.markdown-body ul) { list-style-type: disc; }
:deep(.markdown-body ol) { list-style-type: decimal; }
:deep(.markdown-body code) {
  background-color: rgba(100, 116, 139, 0.1);
  padding: 0.2em 0.4em;
  border-radius: 0.25em;
  font-family: monospace;
  font-size: 0.9em;
}
:deep(.markdown-body pre) {
  background-color: #1e293b;
  color: #e2e8f0;
  padding: 1em;
  border-radius: 0.5em;
  overflow-x: auto;
  margin-bottom: 0.75em;
}
:deep(.markdown-body pre code) {
  background-color: transparent;
  padding: 0;
  color: inherit;
}
:deep(.markdown-body a) {
  color: #3b82f6;
  text-decoration: underline;
}
:deep(.markdown-body blockquote) {
  border-left: 3px solid #cbd5e1;
  padding-left: 1em;
  color: #64748b;
  margin-bottom: 0.75em;
}
</style>