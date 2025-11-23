// Основной JavaScript для банковского приложения
document.addEventListener('DOMContentLoaded', function() {
    // Инициализация приложения
    initApp();
});

function initApp() {
    // Добавление обработчиков событий
    initForms();
    initNotifications();
}

function initForms() {
    // Валидация форм
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            if (!validateForm(this)) {
                e.preventDefault();
            }
        });
    });
}

function validateForm(form) {
    const inputs = form.querySelectorAll('input[required]');
    let isValid = true;

    inputs.forEach(input => {
        if (!input.value.trim()) {
            showFieldError(input, 'Это поле обязательно для заполнения');
            isValid = false;
        } else {
            clearFieldError(input);
        }

        // Специфичная валидация для username
        if (input.name === 'username' && input.value.trim()) {
            if (input.value.length < 3) {
                showFieldError(input, 'Имя пользователя должно содержать минимум 3 символа');
                isValid = false;
            }
        }

        // Специфичная валидация для password
        if (input.name === 'password' && input.value.trim()) {
            if (input.value.length < 8) {
                showFieldError(input, 'Пароль должен содержать минимум 8 символов');
                isValid = false;
            }
        }
    });

    return isValid;
}

function showFieldError(input, message) {
    clearFieldError(input);

    const errorDiv = document.createElement('div');
    errorDiv.className = 'field-error';
    errorDiv.style.color = '#dc3545';
    errorDiv.style.fontSize = '0.875rem';
    errorDiv.style.marginTop = '0.25rem';
    errorDiv.textContent = message;

    input.style.borderColor = '#dc3545';
    input.parentNode.appendChild(errorDiv);
}

function clearFieldError(input) {
    const existingError = input.parentNode.querySelector('.field-error');
    if (existingError) {
        existingError.remove();
    }
    input.style.borderColor = '';
}

function initNotifications() {
    // Автоматическое скрытие уведомлений через 5 секунд
    const alerts = document.querySelectorAll('.alert');
    alerts.forEach(alert => {
        setTimeout(() => {
            alert.style.opacity = '0';
            alert.style.transition = 'opacity 0.5s';
            setTimeout(() => alert.remove(), 500);
        }, 5000);
    });
}

// Утилиты для работы с API
const BankApp = {
    // Сохранение токена
    setToken(token) {
        localStorage.setItem('bank_token', token);
    },

    // Получение токена
    getToken() {
        return localStorage.getItem('bank_token');
    },

    // Удаление токена
    removeToken() {
        localStorage.removeItem('bank_token');
    },

    // Проверка авторизации
    isAuthenticated() {
        return !!this.getToken();
    }
};