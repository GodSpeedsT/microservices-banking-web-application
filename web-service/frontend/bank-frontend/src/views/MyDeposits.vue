<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '@/services/api';
import type { DepositResponse } from '@/types';

const deposits = ref<DepositResponse[]>([]);
const loading = ref(true);

const fetchDeposits = async () => {
  try {
    // ВАЖНО: Убедитесь, что ваш Java Controller возвращает JSON с полями, совпадающими с интерфейсом
    const response = await api.get<DepositResponse[]>('/deposits/my');
    deposits.value = response.data;
  } catch (e) {
    console.error("Ошибка загрузки", e);
  } finally {
    loading.value = false;
  }
};

const closeDeposit = async (id: number) => {
  if(!confirm('Вы подтверждаете досрочное закрытие вклада? Проценты могут быть потеряны.')) return;

  try {
    await api.post(`/deposits/${id}/close`);
    await fetchDeposits();
  } catch (e) {
    alert('Не удалось закрыть вклад. Попробуйте позже.');
  }
};

onMounted(fetchDeposits);
</script>

<template>
  <div>
    <div v-if="loading" class="flex justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>

    <div v-else-if="deposits.length === 0" class="bg-white p-12 rounded-xl shadow-sm border border-bank-200 text-center">
      <div class="mx-auto h-12 w-12 text-bank-300 mb-4">
        <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </div>
      <p class="text-bank-900 font-medium text-lg mb-2">Активных вкладов не найдено</p>
      <p class="text-bank-500 mb-6">Откройте свой первый вклад на выгодных условиях.</p>
      <router-link to="/open" class="inline-flex items-center px-6 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 transition-colors">
        Открыть вклад
      </router-link>
    </div>

    <div v-else class="bg-white rounded-xl shadow-sm border border-bank-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-bank-200">
          <thead class="bg-bank-50">
          <tr>
            <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-bank-500 uppercase tracking-wider">Номер / Тип</th>
            <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-bank-500 uppercase tracking-wider">Баланс</th>
            <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-bank-500 uppercase tracking-wider">Ставка</th>
            <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-bank-500 uppercase tracking-wider">Окончание</th>
            <th scope="col" class="relative px-6 py-4">
              <span class="sr-only">Действия</span>
            </th>
          </tr>
          </thead>
          <tbody class="bg-white divide-y divide-bank-200">
          <tr v-for="deposit in deposits" :key="deposit.id" class="hover:bg-bank-50/50 transition-colors group">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-bank-900">Вклад #{{ deposit.id }}</div>
              <div class="text-xs text-bank-500 mt-0.5">{{ deposit.depositType || 'Стандартный' }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-bold text-bank-900">{{ deposit.balance?.toLocaleString() }} <span class="text-bank-500 font-normal">{{ deposit.currency }}</span></div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  {{ deposit.interestRate }}%
                </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-bank-600">
              {{ new Date(deposit.maturityDate).toLocaleDateString() }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button @click="closeDeposit(deposit.id)" class="text-bank-400 hover:text-red-600 transition-colors font-medium">
                Закрыть
              </button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>