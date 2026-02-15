<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
    <nav class="bg-white shadow-md border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <div class="h-10 w-10 bg-blue-700 rounded-lg flex items-center justify-center font-bold text-white mr-3">
              GB
            </div>
            <span class="text-xl font-bold text-gray-900">Global Bank</span>
          </div>

          <div class="hidden md:flex items-center space-x-8">
            <a href="#" class="text-gray-700 hover:text-blue-600 font-medium transition-colors">Частным лицам</a>
            <a href="#" class="text-gray-700 hover:text-blue-600 font-medium transition-colors">Бизнесу</a>
            <a href="#" class="text-gray-700 hover:text-blue-600 font-medium transition-colors">Инвестиции</a>
            <a href="#" class="text-gray-700 hover:text-blue-600 font-medium transition-colors">О банке</a>
          </div>

          <div class="flex items-center space-x-4">
            <button class="text-gray-700 hover:text-blue-600 font-medium px-4 py-2 transition-colors">
              Отделения и банкоматы
            </button>
            <button @click="login"
                    class="bg-blue-700 hover:bg-blue-800 text-white font-medium px-6 py-2.5 rounded-lg transition-all shadow-sm hover:shadow">
              Войти в интернет-банк
            </button>
          </div>
        </div>
      </div>
    </nav>

    <section class="bg-gradient-to-r from-blue-900 to-blue-800 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h1 class="text-4xl md:text-5xl font-bold leading-tight mb-6">
              Надежный банк для ваших финансовых целей
            </h1>
            <p class="text-xl text-blue-100 mb-8">
              Свежий банк с современными решениями, но космическими амбициями
            </p>
            <div class="flex flex-col sm:flex-row gap-4">
              <button @click="$router.push('/register')"
                      class="bg-white text-blue-900 hover:bg-blue-50 font-bold px-8 py-4 rounded-lg text-lg transition-all shadow-lg hover:shadow-xl">
                Открыть счет
              </button>
              <button class="bg-transparent border-2 border-white hover:bg-white/10 text-white font-bold px-8 py-4 rounded-lg text-lg transition-all">
                Подробнее о банке
              </button>
            </div>
          </div>

          <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
            <h3 class="text-2xl font-bold mb-6">Global Bank в цифрах</h3>
            <div class="grid grid-cols-2 gap-6">
              <div>
                <div class="text-4xl font-bold mb-2">4.8 млн</div>
                <div class="text-blue-200">Клиентов по всему миру</div>
              </div>
              <div>
                <div class="text-4xl font-bold mb-2">€450 млрд</div>
                <div class="text-blue-200">Активы под управлением</div>
              </div>
              <div>
                <div class="text-4xl font-bold mb-2">156 лет</div>
                <div class="text-blue-200">На финансовом рынке</div>
              </div>
              <div>
                <div class="text-4xl font-bold mb-2">ААА</div>
                <div class="text-blue-200">Кредитный рейтинг</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">Банковские продукты</h2>
          <p class="text-xl text-gray-600">Выберите подходящий финансовый инструмент</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-xl transition-shadow">
            <div class="p-8">
              <div class="h-12 w-12 bg-blue-100 rounded-lg flex items-center justify-center mb-6">
                <svg class="h-6 w-6 text-blue-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">Сберегательные счета</h3>
              <p class="text-gray-600 mb-6">Высокие процентные ставки с гарантированной доходностью</p>
              <ul class="space-y-3 mb-8">
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  До 7.5% годовых
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Ежемесячная капитализация
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Страхование вкладов
                </li>
              </ul>
              <button class="w-full bg-blue-700 hover:bg-blue-800 text-white font-medium py-3.5 rounded-lg transition-colors">
                Открыть вклад
              </button>
            </div>
          </div>

          <div class="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-xl transition-shadow">
            <div class="p-8">
              <div class="h-12 w-12 bg-green-100 rounded-lg flex items-center justify-center mb-6">
                <svg class="h-6 w-6 text-green-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">Кредитные карты</h3>
              <p class="text-gray-600 mb-6">Премиальные карты с льготным периодом до 120 дней</p>
              <ul class="space-y-3 mb-8">
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  0% на все покупки 4 месяца
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Кешбэк до 10%
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Бесплатное обслуживание
                </li>
              </ul>
              <button class="w-full bg-green-700 hover:bg-green-800 text-white font-medium py-3.5 rounded-lg transition-colors">
                Оформить карту
              </button>
            </div>
          </div>

          <div class="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-xl transition-shadow">
            <div class="p-8">
              <div class="h-12 w-12 bg-purple-100 rounded-lg flex items-center justify-center mb-6">
                <svg class="h-6 w-6 text-purple-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
              <h3 class="text-xl font-bold text-gray-900 mb-4">Инвестиции</h3>
              <p class="text-gray-600 mb-6">Доверительное управление и индивидуальные портфели</p>
              <ul class="space-y-3 mb-8">
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Средняя доходность 12-15%
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Личный финансовый советник
                </li>
                <li class="flex items-center text-gray-700">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  Минимальная сумма от 100 000 ₽
                </li>
              </ul>
              <button class="w-full bg-purple-700 hover:bg-purple-800 text-white font-medium py-3.5 rounded-lg transition-colors">
                Начать инвестировать
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-gray-900 text-white py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl font-bold mb-4">Безопасность — наш приоритет</h2>
          <p class="text-xl text-gray-300">Мы используем передовые технологии для защиты ваших средств</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div class="text-center">
            <div class="h-16 w-16 bg-blue-800 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Шифрование 256-bit</h3>
            <p class="text-gray-300">Военная степень защиты данных</p>
          </div>

          <div class="text-center">
            <div class="h-16 w-16 bg-blue-800 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Страхование вкладов</h3>
            <p class="text-gray-300">До 1.4 млн рублей по госпрограмме</p>
          </div>

          <div class="text-center">
            <div class="h-16 w-16 bg-blue-800 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Мгновенные переводы</h3>
            <p class="text-gray-300">Переводы 24/7 в любую точку мира</p>
          </div>

          <div class="text-center">
            <div class="h-16 w-16 bg-blue-800 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636l-3.536 3.536m0 5.656l3.536 3.536M9.172 9.172L5.636 5.636m3.536 9.192l-3.536 3.536M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-5 0a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold mb-2">Поддержка 24/7</h3>
            <p class="text-gray-300">Круглосуточная помощь клиентам</p>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 class="text-3xl font-bold text-gray-900 mb-6">Мобильный банк</h2>
            <p class="text-xl text-gray-600 mb-8">
              Управляйте финансами с телефона. Быстрые переводы, оплата счетов и контроль расходов в одном приложении.
            </p>
            <ul class="space-y-4 mb-8">
              <li class="flex items-center">
                <svg class="h-6 w-6 text-green-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span class="text-gray-700">Биометрия для входа</span>
              </li>
              <li class="flex items-center">
                <svg class="h-6 w-6 text-green-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span class="text-gray-700">QR-платежи и NFC</span>
              </li>
              <li class="flex items-center">
                <svg class="h-6 w-6 text-green-500 mr-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                <span class="text-gray-700">Умные уведомления</span>
              </li>
            </ul>
            <div class="flex gap-4">
              <button class="flex items-center bg-gray-900 hover:bg-black text-white px-6 py-3 rounded-lg transition-colors">
                <svg class="h-6 w-6 mr-3" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M17.52 17.97c-.19 0-1.85-.73-2.13-1.02-.27-.29-.73-.88-1.02-1.02-.28-.14-1.02-.46-1.02-.73s.15-.59.44-.88c.29-.29.58-.58.73-.73.29-.29.88-1.02 1.02-1.31.14-.29.29-.73.15-1.02-.15-.29-.73-.58-1.17-.73-.44-.15-1.02-.15-1.46 0-.44.15-1.46.29-2.34 1.02-.88.73-1.75 2.05-1.9 3.36-.15 1.31.73 2.63 1.46 3.36.73.73 3.36 1.31 4.96 1.31 1.6 0 2.19-.15 2.63-.44.44-.29.58-.73.58-1.17 0-.44-.29-.88-.73-1.02-.44-.15-1.02-.29-1.46-.44z"/>
                </svg>
                App Store
              </button>
              <button class="flex items-center bg-gray-900 hover:bg-black text-white px-6 py-3 rounded-lg transition-colors">
                <svg class="h-6 w-6 mr-3" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M3.609 1.814L13.792 12 3.61 22.186a.996.996 0 00.29 1.624l7.888 4.866c.34.21.764.21 1.104 0l7.888-4.866a.996.996 0 00.29-1.624L14.208 12l10.182-10.186a.996.996 0 00-.29-1.624l-7.888-4.866a1.001 1.001 0 00-1.104 0L3.899.19a.996.996 0 00-.29 1.624z"/>
                </svg>
                Google Play
              </button>
            </div>
          </div>
          <div class="relative">
            <div class="relative z-10">
              <img src="https://images.unsplash.com/photo-1563013544-824ae1b704d3?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=800&q=80"
                   alt="Мобильное приложение банка"
                   class="rounded-2xl shadow-2xl">
            </div>
            <div class="absolute -bottom-6 -right-6 w-64 h-64 bg-blue-600 rounded-full opacity-10"></div>
          </div>
        </div>
      </div>
    </section>

    <!-- Футер -->
    <footer class="bg-gray-900 text-white pt-12 pb-8">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div>
            <div class="flex items-center mb-6">
              <div class="h-10 w-10 bg-blue-700 rounded-lg flex items-center justify-center font-bold mr-3">
                GB
              </div>
              <span class="text-xl font-bold">Global Bank</span>
            </div>
            <p class="text-gray-400">
              Международный финансовый институт с 2025 года.
            </p>
          </div>

          <div>
            <h3 class="text-lg font-bold mb-6">Продукты</h3>
            <ul class="space-y-3 text-gray-400">
              <li><a href="#" class="hover:text-white transition-colors">Вклады и счета</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Кредитные карты</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Ипотека</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Инвестиции</a></li>
            </ul>
          </div>

          <div>
            <h3 class="text-lg font-bold mb-6">Клиентам</h3>
            <ul class="space-y-3 text-gray-400">
              <li><a href="#" class="hover:text-white transition-colors">Отделения и банкоматы</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Курсы валют</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Тарифы и документы</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Частые вопросы</a></li>
            </ul>
          </div>

          <div>
            <h3 class="text-lg font-bold mb-6">Контакты</h3>
            <ul class="space-y-3 text-gray-400">
              <li>Круглосуточная поддержка: +375 (33)-606-36-07</li>
              <li>Для бизнеса: +375 (33)-606-36-07</li>
              <li>Эл. почта: masterov_k@bk.ru</li>
              <li>Минск, Острожских 12</li>
            </ul>
          </div>
        </div>

        <div class="border-t border-gray-800 pt-8">
          <div class="flex flex-col md:flex-row justify-between items-center">
            <p class="text-gray-500 text-sm mb-4 md:mb-0">
              © 2026 Global Bank PLC. Все права защищены.
            </p>
            <div class="flex gap-6 text-sm text-gray-500">
              <a href="#" class="hover:text-white transition-colors">Конфиденциальность</a>
              <a href="#" class="hover:text-white transition-colors">Условия использования</a>
              <a href="#" class="hover:text-white transition-colors">Карта сайта</a>
            </div>
          </div>
          <p class="text-center text-gray-600 text-xs mt-6">
            Global Bank PLC зарегистрирован в Англии и Уэльсе под номером 12345678.
            Регулируется Управлением по финансовому регулированию и надзору Великобритании и Управлением по финансовому поведению.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import {RouterView} from "vue-router";

const router = useRouter();

const login = () => {
  window.location.href = 'http://localhost:8080/oauth2/authorization/keycloak';
}

</script>

<style scoped>
.transition-all {
  transition: all 0.3s ease;
}

.transition-colors {
  transition: background-color 0.3s ease, color 0.3s ease;
}

.transition-shadow {
  transition: box-shadow 0.3s ease;
}

html {
  scroll-behavior: smooth;
}

.hover\:shadow-xl:hover {
  transform: translateY(-4px);
}
</style>