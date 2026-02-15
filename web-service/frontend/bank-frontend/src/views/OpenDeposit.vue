<template>
  <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Заголовок страницы -->
    <div class="mb-10">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Открытие нового вклада</h1>
          <p class="text-gray-600 mt-2">Выберите подходящий сберегательный продукт и укажите параметры</p>
        </div>
        <button
            @click="$router.push('/')"
            class="text-gray-700 hover:text-gray-900 font-medium px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          ← Назад
        </button>
      </div>
      <div class="h-1 w-24 bg-blue-700 rounded-full"></div>
    </div>

    <!-- Состояние загрузки -->
    <div v-if="loading" class="space-y-6 animate-pulse">
      <div class="h-48 bg-gray-200 rounded-xl"></div>
      <div class="h-64 bg-gray-200 rounded-xl"></div>
    </div>

    <!-- Ошибка загрузки -->
    <div v-else-if="error && depositTypes.length === 0"
         class="p-6 rounded-xl bg-red-50 border border-red-200 text-red-700">
      <div class="flex items-start">
        <svg class="h-5 w-5 mr-3 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
<!--        <div>-->
<!--          <h3 class="font-medium mb-1">Ошибка загрузки</h3>-->
<!--          <p>{{ error }}</p>-->
<!--          <button @click="location.reload()" class="text-sm underline mt-2">Попробовать снова</button>-->
<!--        </div>-->
      </div>
    </div>

    <!-- Основной контент -->
    <div v-else class="grid lg:grid-cols-3 gap-8">
      <!-- Левая колонка - выбор типа вклада -->
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div class="mb-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-2">1. Выберите тип вклада</h2>
            <p class="text-gray-600 text-sm">Все вклады защищены системой страхования до 1.4 млн ₽</p>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div
                v-for="type in depositTypes"
                :key="type.id"
                @click="selectedTypeId = type.id"
                class="relative cursor-pointer group rounded-xl border-2 p-6 transition-all duration-200 hover:shadow-lg"
                :class="[
                selectedTypeId === type.id
                  ? 'border-blue-600 bg-blue-50 ring-2 ring-blue-100'
                  : 'border-gray-200 hover:border-gray-300 bg-white'
              ]"
            >
              <!-- Бейдж "Популярный" -->
              <div v-if="type.id === 1" class="absolute -top-2 left-4">
                <span class="bg-green-600 text-white text-xs font-bold px-3 py-1 rounded-full">
                  Популярный
                </span>
              </div>

              <!-- Заголовок и ставка -->
              <div class="flex justify-between items-start mb-6">
                <div>
                  <h3 class="font-bold text-lg text-gray-900 mb-1">{{ type.name }}</h3>
                  <p class="text-sm text-gray-500">Срок: {{ type.months }} месяцев</p>
                </div>
                <div class="text-right">
                  <div class="text-2xl font-bold text-green-600">
                    {{ type.interestRate }}%
                  </div>
                  <div class="text-xs text-gray-500">годовых</div>
                </div>
              </div>

              <!-- Условия -->
              <div class="space-y-3 mb-6">
                <div class="flex items-center text-sm text-gray-700">
                  <svg class="h-4 w-4 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span>Ежемесячная капитализация</span>
                </div>
                <div class="flex items-center text-sm text-gray-700">
                  <svg class="h-4 w-4 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <span>Пополнение без ограничений</span>
                </div>
              </div>

              <!-- Минимальная сумма -->
              <div class="pt-4 border-t border-gray-200">
                <div class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-1">
                  Минимальная сумма
                </div>
                <div class="text-lg font-semibold text-gray-900">
                  {{ type.minAmount.toLocaleString() }} {{ type.currency }}
                </div>
              </div>

              <!-- Иконка выбора -->
              <div v-if="selectedTypeId === type.id" class="absolute top-6 right-6">
                <div class="h-6 w-6 bg-blue-600 rounded-full flex items-center justify-center">
                  <svg class="h-4 w-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Дополнительная информация -->
        <div class="bg-gray-50 rounded-xl border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">Важная информация</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Страхование вкладов</h4>
              <p class="text-sm text-gray-600">
                Все вклады застрахованы Агентством по страхованию вкладов на сумму до 1 400 000 ₽.
              </p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-700 mb-2">Налогообложение</h4>
              <p class="text-sm text-gray-600">
                Налог на доход от вклада уплачивается, если ставка превышает ключевую ставку ЦБ РФ на 5%.
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Правая колонка - форма -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-xl shadow-lg border border-gray-200 sticky top-6">
          <!-- Шапка формы -->
          <div class="p-6 border-b border-gray-200 bg-gray-50 rounded-t-xl">
            <div class="flex items-center">
              <div class="h-10 w-10 bg-blue-700 rounded-lg flex items-center justify-center mr-3">
                <svg class="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <h2 class="text-lg font-semibold text-gray-900">2. Параметры вклада</h2>
                <p class="text-sm text-gray-600">Заполните данные для открытия</p>
              </div>
            </div>
          </div>

          <!-- Тело формы -->
          <div class="p-6">
            <!-- Состояние "выберите тип" -->
            <div v-if="!selectedType" class="text-center py-8 text-gray-500">
              <div class="h-16 w-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg class="h-8 w-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              </div>
              <p class="font-medium text-gray-700 mb-1">Сначала выберите тип вклада</p>
              <p class="text-sm">Кликните на один из продуктов слева</p>
            </div>

            <!-- Форма с выбранным типом -->
            <div v-else>
              <!-- Информация о выбранном вкладе -->
              <div class="mb-6 p-4 bg-blue-50 rounded-lg border border-blue-100">
                <div class="flex items-center justify-between mb-2">
                  <span class="font-semibold text-gray-900">{{ selectedType.name }}</span>
                  <span class="text-lg font-bold text-green-600">{{ selectedType.interestRate }}%</span>
                </div>
                <div class="flex items-center text-sm text-gray-600">
                  <svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Срок: {{ selectedType.months }} месяцев
                </div>
              </div>

              <!-- Поле суммы -->
              <div class="mb-8">
                <label class="block text-sm font-medium text-gray-700 mb-3">
                  Сумма вклада
                  <span class="text-red-500">*</span>
                </label>
                <div class="relative rounded-lg shadow-sm">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <span class="text-gray-500 sm:text-sm">₽</span>
                  </div>
                  <input
                      type="number"
                      v-model="amount"
                      :min="selectedType.minAmount"
                      class="block w-full rounded-lg border-gray-300 pl-10 pr-4 py-3 text-gray-900 placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
                      placeholder="Введите сумму"
                      :class="{
                      'border-red-300': amount && !isValid,
                      'border-green-300': isValid
                    }"
                  />
                  <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
                    <svg v-if="amount && isValid" class="h-5 w-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                    </svg>
                    <svg v-else-if="amount && !isValid" class="h-5 w-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                <div class="mt-2 flex justify-between text-xs">
                  <span class="text-gray-500">
                    Мин. сумма: {{ selectedType.minAmount.toLocaleString() }} {{ selectedType.currency }}
                  </span>
                  <span v-if="isValid" class="text-green-600 font-medium">
                    ✓ Готово к открытию
                  </span>
                  <span v-else-if="amount && !isValid" class="text-red-600 font-medium">
                    Сумма меньше минимальной
                  </span>
                </div>
              </div>

              <!-- Расчет доходности -->
              <div class="mb-8 bg-gray-50 rounded-lg p-5 border border-gray-200">
                <h4 class="font-medium text-gray-900 mb-4">Расчет доходности</h4>
                <div class="space-y-3">
                  <div class="flex justify-between text-sm">
                    <span class="text-gray-600">Сумма вклада</span>
                    <span class="font-medium text-gray-900">
                      {{ amount ? Number(amount).toLocaleString() : '0' }} ₽
                    </span>
                  </div>
                  <div class="flex justify-between text-sm">
                    <span class="text-gray-600">Годовая ставка</span>
                    <span class="font-medium text-gray-900">{{ selectedType.interestRate }}%</span>
                  </div>
                  <div class="pt-3 border-t border-gray-200">
                    <div class="flex justify-between items-center">
                      <span class="text-gray-900 font-medium">Доход за {{ selectedType.months }} мес.</span>
                      <span class="text-lg font-bold text-green-600">
                        {{ calculateProfit() }} ₽
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Кнопка отправки -->
              <button
                  @click="handleSubmit"
                  :disabled="!isValid || submitting"
                  class="w-full py-4 px-4 border border-transparent rounded-lg text-base font-semibold text-white bg-blue-700 hover:bg-blue-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-md hover:shadow-lg"
              >
                <div class="flex items-center justify-center">
                  <svg v-if="submitting" class="animate-spin h-5 w-5 mr-3 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                  </svg>
                  <span>{{ submitting ? 'Оформление вклада...' : 'Открыть вклад' }}</span>
                </div>
              </button>

              <!-- Сообщение об ошибке -->
              <div v-if="error" class="mt-4 p-3 rounded-lg bg-red-50 border border-red-200">
                <div class="flex">
                  <svg class="h-5 w-5 text-red-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p class="text-sm text-red-700">{{ error }}</p>
                </div>
              </div>

              <!-- Информационное сообщение -->
              <div class="mt-6 p-4 bg-gray-50 rounded-lg border border-gray-200">
                <div class="flex">
                  <svg class="h-5 w-5 text-blue-500 mr-2 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p class="text-xs text-gray-600">
                    Вклад будет открыт моментально. Вы сможете управлять им в разделе "Мои вклады".
                    Средства застрахованы АСВ в установленном порядке.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import api from '@/services/api';
import type { DepositType, DepositRequest } from '@/types';

const router = useRouter();

// Состояние
const depositTypes = ref<DepositType[]>([]);
const loading = ref(true);
const submitting = ref(false);
const error = ref<string | null>(null);

// Форма
const selectedTypeId = ref<number | null>(null);
const amount = ref<string>('');

// Вычисляем выбранный объект типа для отображения условий
const selectedType = computed(() =>
    depositTypes.value.find(t => t.id === selectedTypeId.value)
);

// Валидация
const isValid = computed(() => {
  if (!selectedType.value) return false;
  const val = parseFloat(amount.value);
  return !isNaN(val) && val >= selectedType.value.minAmount;
});

// Расчет дохода
const calculateProfit = () => {
  if (!selectedType.value || !amount.value) return '0';
  const principal = parseFloat(amount.value);
  const rate = selectedType.value.interestRate / 100;
  const months = selectedType.value.months;

  // Простой расчет с ежемесячной капитализацией
  const monthlyRate = rate / 12;
  const profit = principal * Math.pow(1 + monthlyRate, months) - principal;

  return profit.toFixed(0);
};

// Загрузка типов вкладов
onMounted(async () => {
  try {
    const response = await api.get<DepositType[]>('/deposit-types/active');
    depositTypes.value = response.data;
    // Автовыбор первого вклада
    if (depositTypes.value.length > 0) {
      selectedTypeId.value = depositTypes.value[0].id;
    }
  } catch (e) {
    error.value = 'Не удалось загрузить доступные вклады. Попробуйте позже.';
  } finally {
    loading.value = false;
  }
});

// При смене типа вклада сбрасываем или корректируем сумму
watch(selectedTypeId, () => {
  error.value = null;
});

// Отправка формы
const handleSubmit = async () => {
  if (!isValid.value || !selectedType.value) return;

  submitting.value = true;
  error.value = null;

  const payload: DepositRequest = {
    depositTypeId: selectedType.value.id,
    amount: parseFloat(amount.value),
    currency: selectedType.value.currency
  };

  try {
    await api.post('/deposits', payload);
    // Успех -> редирект на главную
    router.push('/');
  } catch (e: any) {
    console.error(e);
    error.value = e.response?.data?.message || 'Ошибка при открытии вклада. Проверьте данные.';
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
/* Плавные анимации */
.transition-all {
  transition: all 0.3s ease;
}

.transition-colors {
  transition: background-color 0.3s ease, border-color 0.3s ease, color 0.3s ease;
}

/* Кастомная анимация для спиннера */
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Эффект при наведении на карточки вкладов */
.group:hover {
  transform: translateY(-2px);
}
</style>