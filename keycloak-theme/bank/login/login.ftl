<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" class="h-full">
<head>
    <meta charset="utf-8">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="robots" content="noindex, nofollow">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <title>${msg("loginTitle", (realmDisplayName!'Global Bank'))}</title>

    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        gray: {
                            50: '#f9fafb', 100: '#f3f4f6', 200: '#e5e7eb', 300: '#d1d5db',
                            400: '#9ca3af', 500: '#6b7280', 600: '#4b5563', 700: '#374151',
                            800: '#1f2937', 900: '#111827'
                        },
                        blue: {
                            50: '#eff6ff', 100: '#dbeafe', 200: '#bfdbfe', 300: '#93c5fd',
                            400: '#60a5fa', 500: '#3b82f6', 600: '#2563eb', 700: '#1d4ed8',
                            800: '#1e40af', 900: '#1e3a8a'
                        },
                        green: {
                            600: '#059669', 700: '#047857'
                        }
                    },
                    fontFamily: {
                        sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Arial', 'sans-serif']
                    },
                    boxShadow: {
                        'soft': '0 4px 20px rgba(0, 0, 0, 0.08)',
                        'elevated': '0 10px 40px rgba(0, 0, 0, 0.12)'
                    }
                }
            }
        }
    </script>

    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

        body {
            font-family: 'Inter', ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
            background: linear-gradient(135deg, #f8fafc 0%, #eef2ff 100%);
        }

        .gradient-bg {
            background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 100%);
        }

        .card-hover {
            transition: all 0.3s ease;
        }

        .card-hover:hover {
            transform: translateY(-2px);
            box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
        }

        .input-focus:focus {
            box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
        }
    </style>
</head>

<body class="h-full">
<div class="min-h-screen flex items-center justify-center p-4">
    <div class="max-w-md w-full">
        <div class="text-center mb-8">
            <div class="flex items-center justify-center mb-6">
                <div class="ml-4 text-left">
                    <h1 class="text-2xl font-bold text-gray-900">Global Bank</h1>
                </div>
            </div>
            <div class="h-1 w-24 bg-blue-600 rounded-full mx-auto"></div>
        </div>

        <div class="bg-white rounded-2xl shadow-elevated border border-gray-200 overflow-hidden card-hover">
            <div class="px-8 py-10">
                <div class="text-center mb-8">
                    <h2 class="text-2xl font-bold text-gray-900 mb-2">${msg("loginTitle", (realmDisplayName!'Клиентский портал'))}</h2>
                    <p class="text-gray-600">Войдите в аккаунт</p>
                </div>

                <form id="kc-form-login" class="${kcFormClass!}" action="${url.loginAction}" method="post">
                    <#if message?has_content>
                        <div class="mb-6 p-4 rounded-lg bg-red-50 border border-red-200" role="alert">
                            <div class="flex items-center">
                                <svg class="h-5 w-5 text-red-500 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                                <span class="font-medium text-red-700">${kcSanitize(message.summary)}</span>
                            </div>
                        </div>
                    </#if>

                    <div class="mb-6">
                        <label for="username" class="block text-sm font-medium text-gray-700 mb-3">
                            ${msg("username")}
                            <span class="text-red-500">*</span>
                        </label>
                        <div class="relative rounded-lg shadow-sm">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                                </svg>
                            </div>
                            <input id="username" name="username" value="${(login.username!'')}" type="text"
                                   class="block w-full rounded-lg border-gray-300 pl-10 pr-4 py-3 text-gray-900 placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 input-focus transition-colors"
                                   autocomplete="username"
                                   aria-invalid="<#if message?has_content && message.type='error'>true</#if>"
                                   placeholder="${msg("usernamePlaceholder")}" />
                        </div>
                    </div>

                    <div class="mb-8">
                        <label for="password" class="block text-sm font-medium text-gray-700 mb-3">
                            ${msg("password")}
                            <span class="text-red-500">*</span>
                        </label>
                        <div class="relative rounded-lg shadow-sm">
                            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                                </svg>
                            </div>
                            <input id="password" name="password" type="password"
                                   class="block w-full rounded-lg border-gray-300 pl-10 pr-4 py-3 text-gray-900 placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 input-focus transition-colors"
                                   autocomplete="current-password"
                                   aria-invalid="<#if message?has_content && message.type='error'>true</#if>"
                                   placeholder="${msg("passwordPlaceholder")}" />
                        </div>
                    </div>

                    <button type="submit" name="login" id="kc-login"
                            class="w-full py-3.5 px-4 border border-transparent rounded-lg text-base font-semibold text-white gradient-bg hover:opacity-95 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-600 transition-all shadow-md hover:shadow-lg">
                        <div class="flex items-center justify-center">
                            <svg class="h-5 w-5 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                            </svg>
                            ${msg("doLogIn")}
                        </div>
                    </button>
                </form>

                <#if realm.password && realm.identityProviders?? && realm.identityProviders?has_content>
                    <div class="relative my-8">
                        <div class="absolute inset-0 flex items-center">
                            <div class="w-full border-t border-gray-200"></div>
                        </div>
                        <div class="relative flex justify-center text-sm">
                            <span class="px-4 bg-white text-gray-500">${msg("socialLoginTitle")}</span>
                        </div>
                    </div>

                    <div id="kc-social-providers" class="grid grid-cols-1 gap-3">
                        <#list realm.identityProviders as idp>
                            <a id="sso-${idp.alias}" href="${url.loginSocialUrl?replace("PROVIDER_ID", idp.alias)}"
                               class="w-full flex items-center justify-center py-3 px-4 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors card-hover">
                                <svg class="h-5 w-5 mr-3 text-gray-500" fill="currentColor" viewBox="0 0 20 20">
                                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clip-rule="evenodd" />
                                </svg>
                                ${idp.displayName!}
                            </a>
                        </#list>
                    </div>
                </#if>

                <div class="mt-8 pt-6 border-t border-gray-200">
                    <div class="flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0">
                        <#if realm.resetPasswordAllowed>
                            <a href="${url.loginResetCredentialsUrl}"
                               class="text-sm font-medium text-blue-600 hover:text-blue-800 transition-colors flex items-center">
                                <svg class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                                </svg>
                                ${msg("doForgotPassword")}
                            </a>
                        </#if>

                        <#if realm.registrationAllowed>
                            <a href="${url.registrationUrl}"
                               class="text-sm font-medium text-blue-600 hover:text-blue-800 transition-colors flex items-center">
                                <svg class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                                </svg>
                                ${msg("doRegister")}
                            </a>
                        </#if>
                    </div>
                </div>
            </div>

            <div class="bg-gray-50 px-8 py-5 border-t border-gray-200">
            </div>
        </div>

        <div class="mt-6 p-4 bg-blue-50 rounded-xl border border-blue-100">
            <div class="flex items-start">
                <svg class="h-5 w-5 text-blue-600 mr-3 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <div>
                    <p class="text-sm text-gray-700">
                        Впервые используете интернет-банк?
                        <a href="#" class="font-medium text-blue-600 hover:text-blue-800">Запросите доступ</a>
                    </p>
                </div>
            </div>
        </div>
    </div>
</div>

<div class="hidden lg:block fixed bottom-0 left-0 right-0 bg-gray-900 text-white py-4">
    <div class="max-w-6xl mx-auto px-8">
        <div class="flex items-center justify-between text-sm">
            <div class="flex items-center space-x-8">
                <div class="flex items-center">
                    <svg class="h-5 w-5 text-green-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>

                </div>
                <div class="flex items-center">
                    <svg class="h-5 w-5 text-green-400 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>