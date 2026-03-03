<template>
  <div class="gopher-scene">
    <div class="glow-bg"></div>
    <div class="gopher-wrapper" ref="gopherRef">
      <svg 
        xmlns="http://www.w3.org/2000/svg" 
        xmlns:xlink="http://www.w3.org/1999/xlink"
        viewBox="0 0 280 510" 
        preserveAspectRatio="xMidYMid meet"
        class="gopher-svg"
      >
        <defs>
          <clipPath id="clip-main"><rect width="280" height="510" x="0" y="0"/></clipPath>
          <clipPath id="clip-content"><path d="M0,0 L242,0 L242,400 L0,400z"/></clipPath>
        </defs>
        
        <g clip-path="url(#clip-main)">
          <g clip-path="url(#clip-content)" transform="translate(19, -49.072)">
            
            <!-- LEFT EAR (ID 74) -->
            <g class="animate-ear-left" transform="translate(0, 62.5)">
              <g transform="translate(70.195, 16.534)">
                <path fill="rgb(125,64,255)" d="M12.101,14.604 C13.247,0.965 10.573,-16.691 -1.337,-15.362 C-13.247,-14.033 -12.616,-0.145 2.74,16.534 C2.74,16.534 12.101,14.604 12.101,14.604z"/>
              </g>
            </g>
            
            <!-- RIGHT EAR (ID 75) -->
            <g class="animate-ear-right" transform="translate(0, 62.5)">
              <g transform="translate(173.126, 15.948)">
                <path fill="rgb(125,64,255)" d="M-12.611,13.783 C-13.423,0.121 -10.319,-17.465 1.551,-15.839 C13.423,-14.213 12.456,-0.348 -3.304,15.948 C-3.304,15.948 -12.611,13.783 -12.611,13.783z"/>
              </g>
            </g>

            <!-- LEFT ARM (ID 76) -->
            <g transform="translate(0, 62.5)">
              <g transform="translate(14.591, 157.814)">
                <path fill="rgb(125,64,255)" d="M13.165,-8.016 C4.622,-11.883 3.507,-11.93 -4.677,-15.73 C-15.083,-20.551 -18.13,-12.768 -9.767,-0.652 C-4.965,6.311 0.74,12.994 14.591,20.551 C14.591,20.551 13.165,-8.016 13.165,-8.016z"/>
              </g>
            </g>

            <!-- RIGHT ARM (ID 77) -->
            <g transform="translate(0, 62.5)">
              <g transform="translate(210.147, 150.17)">
                <path fill="rgb(125,64,255)" d="M-31.853,3.067 C-31.853,3.067 -4.431,-4.28 20.06,-19.771 C29.335,-25.64 36.023,-19.161 28.79,-7.07 C20.535,6.717 6.193,20.872 -7.105,25.64 C-7.105,25.64 -31.853,3.067 -31.853,3.067z"/>
              </g>
            </g>

            <!-- BODY (ID 78) -->
            <g transform="translate(0, 62.5)">
              <g transform="translate(118.665, 154.504)">
                <path fill="rgb(143,112,255)" d="M-93.949,-11.044 C-96.521,-48.904 -98.352,-75.415 -81.453,-100.246 C-58.178,-134.449 -17.865,-138.379 -9.331,-139.207 C0.328,-140.146 37.352,-143.757 65.07,-117.816 C86.755,-97.518 89.709,-65.09 86.303,-22.963 C81.695,34.078 98.352,55.219 88.048,93.164 C79.372,125.131 41.928,138.926 17.327,140.872 C-19.089,143.757 -41.701,137.473 -61.943,115.346 C-86.021,89.022 -91.143,30.434 -93.956,-11.036 C-93.956,-11.036 -93.949,-11.044 -93.949,-11.044z"/>
              </g>
            </g>

            <!-- EYES GROUP -->
            <g class="animate-eyes-group">
              <!-- SCLERA (ID 80) -->
              <g transform="translate(0, 62.5)">
                <g transform="translate(155.898, 88.698)">
                  <path fill="rgb(255,255,255)" d="M4.93,-38.419 C25.687,-36.012 41.144,-17.16 39.522,3.818 C37.893,24.929 19.537,40.92 -0.893,40.216 C-23.326,39.451 -41.144,18.801 -38.564,-3.74 C-36.148,-24.866 -16.583,-40.92 4.93,-38.427 C4.93,-38.427 4.93,-38.419 4.93,-38.419z"/>
                </g>
                <g transform="translate(52.033, 90.183)">
                  <path fill="rgb(255,255,255)" d="M4.93,-38.419 C25.687,-36.012 41.144,-17.16 39.523,3.818 C37.894,24.929 19.537,40.92 -0.892,40.216 C-23.325,39.451 -41.144,18.801 -38.564,-3.74 C-36.148,-24.866 -16.583,-40.92 4.93,-38.427 C4.93,-38.427 4.93,-38.419 4.93,-38.419z"/>
                </g>
              </g>
              
              <!-- PUPILS (ID 81) -->
              <g transform="translate(-1.303, 67.697)">
                <g transform="translate(139.96, 76.111)">
                  <path fill="rgb(18,21,23)" d="M1.177,-9.184 C1.177,-9.184 1.185,-9.184 1.185,-9.184 C-3.959,-9.778 -8.636,-5.94 -9.213,-0.891 C-9.829,4.494 -5.573,9.434 -0.21,9.613 C4.669,9.778 9.058,5.956 9.447,0.914 C9.829,-4.103 6.135,-8.605 1.177,-9.184z M0,3.181 C-1.707,3.126 -3.063,1.555 -2.868,-0.164 C-2.681,-1.774 -1.192,-2.993 0.445,-2.806 C0.445,-2.806 0.445,-2.798 0.445,-2.798 C2.027,-2.618 3.204,-1.18 3.079,0.414 C2.955,2.024 1.559,3.236 0,3.181z"/>
                </g>
                <g transform="translate(37.142, 77.267)">
                  <path fill="rgb(18,21,23)" d="M1.177,-9.184 C1.177,-9.184 1.185,-9.184 1.185,-9.184 C-3.96,-9.777 -8.637,-5.939 -9.214,-0.891 C-9.829,4.495 -5.573,9.435 -0.211,9.614 C4.669,9.777 9.057,5.956 9.447,0.915 C9.829,-4.102 6.134,-8.604 1.177,-9.184z M-0.008,3.182 C-1.715,3.127 -3.071,1.555 -2.876,-0.163 C-2.689,-1.773 -1.201,-2.993 0.436,-2.805 C0.436,-2.805 0.436,-2.798 0.436,-2.798 C2.018,-2.617 3.195,-1.18 3.071,0.415 C2.946,2.025 1.551,3.236 -0.008,3.182z"/>
                </g>
              </g>
            </g>

            <!-- TEETH (ID 79) - 上移15px，位于双眼正下方中间 -->
            <g class="animate-teeth">
              <!-- 白色牙齿部分 -->
              <g transform="translate(-1, 62.5)">
                <g transform="translate(105.241, 68.826)"> <!-- 原123.826-15=108.826 -->
                  <path fill="rgb(255,255,255)" d="M-2.825,-6.399 C-2.825,-6.399 2.309,-6.604 2.309,-6.604 C4.471,-6.69 6.296,-5.002 6.382,-2.832 C6.382,-2.832 6.587,2.315 6.587,2.315 C6.673,4.484 4.989,6.313 2.825,6.399 C2.825,6.399 -2.309,6.604 -2.309,6.604 C-4.472,6.69 -6.296,5.002 -6.382,2.832 C-6.382,2.832 -6.586,-2.314 -6.586,-2.314 C-6.673,-4.484 -4.989,-6.313 -2.825,-6.399z"/>
                </g>
                <!-- 紫色嘴巴上沿 -->
                <g transform="translate(105.294, 64.327)"> <!-- 原117.327-15=102.327 -->
                  <path fill="rgb(126,79,255)" d="M-10.581,8.766 C-10.581,8.766 -6.349,8.016 -6.349,8.016 C-2.109,7.258 2.233,7.217 6.489,7.883 C6.489,7.883 10.737,8.547 10.737,8.547 C14.556,9.148 17.994,6.163 17.955,2.286 C17.884,-4.216 12.576,-9.437 6.091,-9.367 C6.091,-9.367 -6.302,-9.242 -6.302,-9.242 C-12.787,-9.171 -17.994,-3.85 -17.924,2.654 C-17.885,6.531 -14.385,9.438 -10.581,8.758 C-10.581,8.758 -10.581,8.766 -10.581,8.766z"/>
                </g>
                <!-- 深色嘴巴下沿 -->
                <g transform="translate(104.624, 53.011)"> <!-- 原108.011-15=93.011 -->
                  <path fill="rgb(97,38,255)" d="M-7.323,0.84 C-7.261,2.341 -2.701,2.974 0.105,2.927 C2.904,2.755 7.402,1.755 7.338,0.262 C7.198,-2.466 3.348,-2.974 -0.128,-2.88 C-3.598,-2.7 -7.401,-1.887 -7.323,0.832 C-7.323,0.832 -7.323,0.84 -7.323,0.84z"/>
                </g>
              </g>
            </g>

          </g>
        </g>
      </svg>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const gopherRef = ref(null)
onMounted(() => {
  setTimeout(() => gopherRef.value?.classList.add('crash-in'), 100)
})
</script>

<style scoped>
.gopher-scene {
  position: relative;
  width: 320px;
  height: 420px;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  overflow: visible;
}

.glow-bg {
  position: absolute;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(125,64,255,0.12) 0%, rgba(125,64,255,0.05) 40%, transparent 70%);
  border-radius: 50%;
  animation: pulse 4s ease-in-out infinite;
  z-index: 0;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 0.6; }
  50% { transform: scale(1.2); opacity: 0.3; }
}

.gopher-wrapper {
  position: relative;
  z-index: 10;
  width: 280px;
  height: 510px;
  opacity: 0;
  transform: translate(200px, -200px) rotate(45deg) scale(0.3);
  transform-origin: center center;
}

.gopher-wrapper.crash-in {
  animation: crash-in 1.2s cubic-bezier(0.34, 1.56, 0.64, 1) forwards,
             float 6s ease-in-out infinite 1.2s;
}

@keyframes crash-in {
  0% {
    opacity: 0;
    transform: translate(200px, -200px) rotate(45deg) scale(0.3);
  }
  60% {
    opacity: 1;
    transform: translate(-20px, 20px) rotate(-5deg) scale(1.1);
  }
  80% {
    transform: translate(10px, -10px) rotate(2deg) scale(0.95);
  }
  100% {
    opacity: 1;
    transform: translate(0, 0) rotate(0deg) scale(1);
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.gopher-svg {
  width: 280px;
  height: 510px;
  overflow: visible;
  filter: drop-shadow(0 20px 30px rgba(0,0,0,0.2));
}

/* Ear Animations */
.animate-ear-left {
  animation: ear-wiggle-l 2.5s ease-in-out infinite;
  transform-box: fill-box;
  transform-origin: center;
}

.animate-ear-right {
  animation: ear-wiggle-r 2.5s ease-in-out infinite 0.3s;
  transform-box: fill-box;
  transform-origin: center;
}

@keyframes ear-wiggle-l {
  0%, 85%, 100% { transform: translate(0px, 62.5px) rotate(0deg); }
  90% { transform: translate(0px, 62.5px) rotate(-10deg); }
  95% { transform: translate(0px, 62.5px) rotate(5deg); }
}

@keyframes ear-wiggle-r {
  0%, 85%, 100% { transform: translate(0px, 62.5px) rotate(0deg); }
  90% { transform: translate(0px, 62.5px) rotate(10deg); }
  95% { transform: translate(0px, 62.5px) rotate(-5deg); }
}

/* Eyes Blink */
.animate-eyes-group {
  transform-box: fill-box;
  transform-origin: 100px 85px;
  animation: blink 4s ease-in-out infinite;
}

@keyframes blink {
  0%, 45%, 55%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(0.1); }
}

/* Teeth Animation - 调整后的位置 */
.animate-teeth {
  transform-box: fill-box;
  transform-origin: 105px 100px; /* 调整origin到新位置 */
  animation: teeth-chatter 3s ease-in-out infinite;
}

@keyframes teeth-chatter {
  0%, 90%, 100% { transform: translate(-1px, 62.5px) translateY(0); }
  94% { transform: translate(-1px, 62.5px) translateY(-2px); }
}

@media (max-width: 640px) {
  .gopher-scene { transform: scale(0.75); }
}
</style>