import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UploadView from '@/views/UploadView.vue'
import SharedView from '@/views/SharedView.vue'
import FileShareView from '@/views/FileShareView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'files',
      component: HomeView,
      meta: { title: 'My Files' }
    },
    {
      path: '/upload',
      name: 'upload',
      component: UploadView,
      meta: { title: 'Upload Files' }
    },
    {
      path: '/shared',
      name: 'shared',
      component:SharedView,
      meta: { title: 'Shared Files' }
    },
    {
      path: '/share/:fileId',
      name: 'file-share',
      component: FileShareView,
      meta: { title: 'Shared File' },
      props: true
    }
  ],
})

// Update page title based on route
router.beforeEach((to, from, next) => {
  document.title = `FileShare - ${to.meta.title || 'File Sharing App'}`
  next()
})

export default router
