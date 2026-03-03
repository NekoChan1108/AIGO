<template>
  <div class="terminal-container">
    <div class="terminal-header">
      <div class="terminal-controls">
        <span class="control close"></span>
        <span class="control minimize"></span>
        <span class="control maximize"></span>
      </div>
      <div class="terminal-title">AIGO System Core</div>
    </div>
    <div class="terminal-body">
      <div class="code-line">
        <span class="prompt">></span>
        <p v-html="formattedText" class="typing-content"></p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'

const props = defineProps<{
  fullText: string
}>()

const displayedText = ref('')
const isTyping = ref(true)

const formattedText = computed(() => {
  let text = displayedText.value
  // Highlight keywords
  const keywords = ['Go', 'Gin', 'Milvus', 'Redis', 'Eino', 'RAG', 'Agent', 'Kafka', 'MongoDB', 'MySQL', 'VolcEngine Ark', 'Vue 3', 'Vite', 'Tailwind CSS', 'Pinia', 'VueUse']
  keywords.forEach(kw => {
    const reg = new RegExp(`\\b${kw}\\b`, 'g')
    text = text.replace(reg, `<span class="highlight">${kw}</span>`)
  })
  
  if (isTyping.value) {
    text += '<span class="caret"></span>'
  }
  return text
})

const typeText = () => {
  displayedText.value = ''
  isTyping.value = true
  let i = 0
  const interval = setInterval(() => {
    if (i < props.fullText.length) {
      displayedText.value += props.fullText.charAt(i)
      i++
    } else {
      clearInterval(interval)
      isTyping.value = false
    }
  }, 30)
}

onMounted(() => {
  setTimeout(typeText, 300)
})

watch(() => props.fullText, () => {
  typeText()
})
</script>

<style>
.highlight {
  color: #818cf8; /* indigo-400 */
  font-weight: 600;
}

.dark .highlight {
  color: #a5b4fc; /* indigo-300 */
}

.caret {
  display: inline-block;
  width: 8px;
  height: 1.1em;
  background-color: #6366f1;
  animation: blink 1s infinite;
  vertical-align: text-bottom;
  margin-left: 4px;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>

<style scoped>
.terminal-container {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(226, 232, 240, 0.8);
  border-radius: 12px;
  width: 100%;
  max-width: 700px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.dark .terminal-container {
  background: rgba(15, 23, 42, 0.8);
  border-color: rgba(51, 65, 85, 0.8);
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
}

.terminal-header {
  background: rgba(241, 245, 249, 0.5);
  padding: 10px 16px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
}

.dark .terminal-header {
  background: rgba(30, 41, 59, 0.5);
  border-bottom-color: rgba(51, 65, 85, 0.8);
}

.terminal-controls {
  display: flex;
  gap: 8px;
}

.control {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.close { background: #ff5f56; }
.minimize { background: #ffbd2e; }
.maximize { background: #27c93f; }

.terminal-title {
  margin-left: 20px;
  font-size: 12px;
  color: #64748b;
  font-family: ui-sans-serif, system-ui, -apple-system, sans-serif;
  letter-spacing: 0.05em;
  text-transform: uppercase;
}

.terminal-body {
  padding: 12px 16px;
  font-family: 'Fira Code', 'JetBrains Mono', 'SF Mono', monospace;
  font-size: 14px;
  line-height: 1.6;
  min-height: 80px;
}

.code-line {
  display: flex;
  gap: 12px;
}

.prompt {
  color: #6366f1;
  font-weight: bold;
  user-select: none;
}

.typing-content {
  margin: 0;
  color: #334155;
  white-space: pre-wrap;
  word-break: break-word;
}

.dark .typing-content {
  color: #e2e8f0;
}
</style>