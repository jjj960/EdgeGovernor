import Vue from 'vue'
import Router from 'vue-router'
import sourceSystem from '@/views/sourceSystem'
import Login from '@/views/Login'
Vue.use(Router)

const router =  new Router({
    routes: [
        {
            path: '/',
            name: 'Login',
            component: Login
        },
        {
            path: '/sourceSystem',
            name: 'sourceSystem',
            component: sourceSystem
        }
    ]
})

export default router

router.beforeEach((to, from, next) => {
    // 1.如果访问的是登录页面（无需权限），直接放行
    if (to.path === '/') return next()
    // 2.如果访问的是有登录权限的页面，先要获取token
    const tokenStr = window.sessionStorage.getItem('token')
    // 2.1如果token为空，强制跳转到登录页面；否则，直接放行
    if (!tokenStr) {
      alert("Please Login！！")
      return next('/')
    }
    next()
})


