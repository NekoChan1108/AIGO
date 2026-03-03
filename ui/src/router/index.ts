import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Chat from '../views/Chat.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Home },
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    { 
      path: '/chat', 
      component: Chat,
      meta: { requiresAuth: true }
    },
    // Catch-all for 404
    { path: '/:pathMatch(.*)*', redirect: '/' }
  ]
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('access_token')
  
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (token && (to.path === '/login' || to.path === '/register' || to.path === '/')) {
    // If logged in and trying to access login/register/root, redirect to chat
    next('/chat')
  } else {
    next()
  }
})

export default router