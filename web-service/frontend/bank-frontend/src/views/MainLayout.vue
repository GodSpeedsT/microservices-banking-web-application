<script setup lang="ts">
import { onMounted } from 'vue';
import { useUserStore } from '@/stores/user';
import {RouterView} from "vue-router";

const userStore = useUserStore();

onMounted(() => {
  userStore.fetchUser();
});

const logout = () => {
  const form = document.createElement('form');
  form.method = 'POST';
  form.action = '/logout';
  document.body.appendChild(form);
  form.submit();
};
</script>

<template>
  <div class="min-h-screen bg-gray-50 flex font-sans">

    <aside class="w-80 bg-gray-900 text-white flex-shrink-0 hidden lg:flex flex-col shadow-lg z-20 border-r border-gray-800">
      <div class="h-24 flex items-center px-6 border-b border-gray-800 bg-gray-950">
        <div class="flex items-center gap-4">
          <div class="h-12 w-12 bg-blue-700 rounded-xl flex items-center justify-center font-bold border-2 border-blue-600">
            <span class="text-xl">GB</span>
          </div>
          <div class="flex flex-col">
            <span class="font-bold text-xl tracking-tight text-white">GLOBAL BANK</span>
            <span class="text-xs text-gray-400 mt-1">Private Banking</span>
          </div>
        </div>
      </div>

      <nav class="flex-1 py-8 px-4 space-y-1">
        <div class="px-4 py-3 text-xs uppercase tracking-wider text-gray-500 font-semibold mb-2">
          Основное
        </div>

        <router-link to="#"
                     class="group flex items-center px-4 py-4 text-base font-medium rounded-lg transition-all duration-200 border-l-4 border-transparent"
                     active-class="bg-blue-900/30 text-white border-blue-500"
                     :class="{'text-gray-300 hover:bg-gray-800 hover:text-white': $route.path !== '/'}"
        >
          <div class="mr-4 h-6 w-6 flex items-center justify-center">
            <svg class="h-5 w-5 opacity-90" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          Мои счета и вклады
        </router-link>

        <router-link to="/open"
                     class="group flex items-center px-4 py-4 text-base font-medium rounded-lg transition-all duration-200 border-l-4 border-transparent"
                     active-class="bg-blue-900/30 text-white border-blue-500"
                     :class="{'text-gray-300 hover:bg-gray-800 hover:text-white': $route.path !== '/open'}"
        >
          <div class="mr-4 h-6 w-6 flex items-center justify-center">
            <svg class="h-5 w-5 opacity-90" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          Открыть новый вклад
        </router-link>

        <div class="px-4 py-3 text-xs uppercase tracking-wider text-gray-500 font-semibold mt-8 mb-2">
          Операции
        </div>

        <a href="#" class="group flex items-center px-4 py-4 text-base font-medium rounded-lg text-gray-300 hover:bg-gray-800 hover:text-white transition-all duration-200 border-l-4 border-transparent">
          <div class="mr-4 h-6 w-6 flex items-center justify-center">
            <svg class="h-5 w-5 opacity-90" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
          </div>
          Переводы
        </a>

        <a href="#" class="group flex items-center px-4 py-4 text-base font-medium rounded-lg text-gray-300 hover:bg-gray-800 hover:text-white transition-all duration-200 border-l-4 border-transparent">
          <div class="mr-4 h-6 w-6 flex items-center justify-center">
            <svg class="h-5 w-5 opacity-90" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          Выписки
        </a>
      </nav>

      <div class="p-6 bg-gray-950 border-t border-gray-800">
        <div class="flex items-center mb-4">
          <div class="h-12 w-12 rounded-full bg-blue-800 flex items-center justify-center text-base font-bold border-2 border-blue-700">
            {{ userStore.username.charAt(0).toUpperCase() }}
          </div>
          <div class="ml-4">
            <p class="text-base font-semibold text-white">{{ userStore.username }}</p>
            <p class="text-sm text-gray-400">Премиум клиент</p>
          </div>
        </div>
        <button @click="logout" class="flex items-center justify-center w-full py-3 text-sm font-medium text-gray-300 hover:text-white hover:bg-gray-800 rounded-lg border border-gray-700 hover:border-gray-600 transition-all duration-200">
          <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
          Выйти из системы
        </button>
      </div>
    </aside>

    <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-gray-50">
      <header class="lg:hidden bg-gray-900 text-white px-6 py-4 flex justify-between items-center shadow-md border-b border-gray-800">
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 bg-blue-700 rounded-lg flex items-center justify-center font-bold">GB</div>
          <span class="font-bold">GLOBAL BANK</span>
        </div>
        <button @click="logout" class="text-sm text-gray-300 hover:text-white">Выйти</button>
      </header>

      <header class="hidden lg:flex bg-white h-20 items-center justify-between px-8 border-b border-gray-200 sticky top-0 z-10 shadow-sm">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 tracking-tight">
            Личный кабинет
          </h1>
          <p class="text-sm text-gray-600 mt-1">Управление вашими финансовыми продуктами</p>
        </div>
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-full bg-blue-100 flex items-center justify-center text-blue-700 font-bold border border-blue-200">
              {{ userStore.username.charAt(0).toUpperCase() }}
            </div>
            <div class="flex flex-col">
              <div v-if="userStore.loading" class="h-4 w-32 bg-gray-200 rounded animate-pulse"></div>
              <div v-else class="text-base font-semibold text-gray-900">
                {{ userStore.username }}
              </div>
              <div class="text-sm text-gray-500">Сегодня: {{ new Date().toLocaleDateString('ru-RU') }}</div>
            </div>
          </div>
          <div class="h-8 w-px bg-gray-200"></div>
          <button class="text-gray-600 hover:text-gray-900 p-2 rounded-lg hover:bg-gray-100 transition-colors">
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
          </button>
        </div>
      </header>

      <div class="flex-1 overflow-auto p-6 lg:p-8 bg-gray-50">
        <div class="max-w-7xl mx-auto">
          <div class="mb-8 grid grid-cols-1 md:grid-cols-3 gap-6">
            <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-semibold text-gray-900">Общий баланс</h3>
                <span class="text-2xl font-bold text-blue-700">— ₽</span>
              </div>
              <p class="text-sm text-gray-600">Загружается...</p>
            </div>
            <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-semibold text-gray-900">Активные вклады</h3>
                <span class="text-2xl font-bold text-green-600">0</span>
              </div>
              <p class="text-sm text-gray-600">Ставки от 5,5% годовых</p>
            </div>
            <div class="bg-white rounded-xl shadow-sm p-6 border border-gray-200">
              <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-semibold text-gray-900">Срок обслуживания</h3>
                <span class="text-2xl font-bold text-gray-900">—</span>
              </div>
              <p class="text-sm text-gray-600">Загружается...</p>
            </div>
          </div>

          <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <component :is="Component" />
              </transition>
            </router-view>
          </div>

          <!-- Футер -->
          <div class="mt-8 text-center text-sm text-gray-500">
            <p>Global Bank PLC. Лицензия ЦБ РФ № 1234 от 01.01.2023</p>
            <p class="mt-2">Ваши средства защищены системой страхования вкладов</p>
            <div class="mt-4 flex justify-center gap-6">
              <a href="#" class="text-gray-500 hover:text-gray-700">Безопасность</a>
              <a href="#" class="text-gray-500 hover:text-gray-700">Тарифы</a>
              <a href="#" class="text-gray-500 hover:text-gray-700">Контакты</a>
              <a href="#" class="text-gray-500 hover:text-gray-700">О банке</a>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}
.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>