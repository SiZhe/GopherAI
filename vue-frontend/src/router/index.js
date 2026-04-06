import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Menu from '../views/Menu.vue'
import AIChat from '../views/AIChat.vue'
import api from '../utils/api'
//import ImageRecognition from '../views/ImageRecognition.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/menu',
    name: 'Menu',
    component: Menu,
    meta: { requiresAuth: true }
  },
  {
    path: '/ai-chat',
    name: 'AIChat',
    component: AIChat,
    meta: { requiresAuth: true }
  },
  // {
  //   path: '/image-recognition',
  //   name: 'ImageRecognition',
  //   component: ImageRecognition,
  //   meta: { requiresAuth: true }
  // }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

router.beforeEach(async (to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    try {
      // 前端无法读取 httpOnly cookie，所以这里用一个受保护接口探测登录态
      const res = await api.get('/AI/chat/sessions')
      const statusCode = res?.data?.status_code
      if (statusCode === 1000) {
        next()
        return
      }
      next('/login')
      return
    } catch (e) {
      next('/login')
      return
    }
  }
  next()
})

export default router

