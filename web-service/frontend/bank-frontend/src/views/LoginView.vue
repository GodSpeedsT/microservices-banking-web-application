<script setup lang="ts">
const handleLogin = () => {
  window.location.href = '/oauth2/authorization/gateway';
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-bank-100">
    <div class="max-w-md w-full bg-white rounded-xl shadow-2xl border border-bank-200 overflow-hidden">
      <div class="px-8 py-12">
        <div class="text-center mb-10">
          <div class="mx-auto h-16 w-16 bg-bank-900 rounded-full flex items-center justify-center mb-6 shadow-lg">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-8 h-8 text-white">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 21v-8.25M15.75 21v-8.25M8.25 21v-8.25M3 9l9-6 9 6m-1.5 12V10.332A48.36 48.36 0 0012 9.75c-2.551 0-5.056.2-7.5.582V21M3 21h18M12 6.75h.008v.008H12V6.75z" />
            </svg>
          </div>
          <h2 class="text-3xl font-extrabold text-bank-900 tracking-tight">Corporate Bank</h2>
          <p class="text-sm text-bank-500 mt-2 font-medium">Secure Client Portal Access</p>
        </div>

        <button
            @click="handleLogin"
            class="w-full flex justify-center py-3.5 px-4 border border-transparent rounded-lg shadow-md text-sm font-semibold text-white bg-bank-900 hover:bg-bank-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-bank-600 transition-all duration-200 ease-in-out transform hover:-translate-y-0.5"
        >
          Войти через Keycloak SSO
        </button>
      </div>
      <div class="bg-bank-50 px-8 py-4 border-t border-bank-200 text-center">
        <p class="text-xs text-bank-400 font-medium">© 2026 Corporate Bank System. Protected Area.</p>
      </div>
    </div>
  </div>
</template>