import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', component: () => import('./views/Dashboard.vue') },
  { path: '/agent-management', component: () => import('./views/AgentManagement.vue') },
  { path: '/multi-person-chat', component: () => import('./views/MultiPersonChat.vue') },
  { path: '/execution-monitor', component: () => import('./views/ExecutionMonitor.vue') },
  { path: '/model-management', component: () => import('./views/ModelManagement.vue') },
  { path: '/orchestration-list', component: () => import('./views/OrchestrationList.vue') },
  { path: '/orchestration-editor/:id', component: () => import('./views/WorkflowEditorView.vue') },
  { path: '/schedule-tasks', component: () => import('./views/ScheduleTasks.vue') },
  { path: '/skill-management', component: () => import('./views/SkillManagement.vue') },
  { path: '/single-person-chat', component: () => import('./views/SinglePersonChat.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')
