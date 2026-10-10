import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth.store'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/kasir',
    name: 'Kasir',
    component: () => import('@/views/kasir/TransactionView.vue'),
    meta: { requiresAuth: true, roles: ['kasir', 'owner'] },
  },
  {
    path: '/kasir/history',
    name: 'TransactionHistory',
    component: () => import('@/views/kasir/TransactionHistoryView.vue'),
    meta: { requiresAuth: true, roles: ['kasir', 'owner'] },
  },
  {
    path: '/owner',
    name: 'Dashboard',
    component: () => import('@/views/owner/DashboardView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/products',
    name: 'Products',
    component: () => import('@/views/owner/ProductsView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/inventory',
    name: 'Inventory',
    component: () => import('@/views/owner/InventoryView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/users',
    name: 'Users',
    component: () => import('@/views/owner/UsersView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/transactions',
    name: 'Transactions',
    component: () => import('@/views/owner/TransactionsView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/reports',
    name: 'Reports',
    component: () => import('@/views/owner/ReportsView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/owner/settings',
    name: 'Settings',
    component: () => import('@/views/owner/SettingsView.vue'),
    meta: { requiresAuth: true, roles: ['owner'] },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.roles && !to.meta.roles.includes(authStore.user?.role || '')) {
    next('/login')
  } else {
    next()
  }
})

export default router
