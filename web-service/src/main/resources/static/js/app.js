// Bank Online Application JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('Bank Online application loaded');

    // Initialize all components
    initApplication();
});

function initApplication() {
    initAnimations();
    initFormValidation();
    initNotifications();
    initNavigation();
    initUserData();
    initServiceStatus();

    // Add global error handler
    window.addEventListener('error', handleGlobalError);
}

// Animation system
function initAnimations() {
    // Fade in elements
    const fadeElements = document.querySelectorAll('.fade-in');
    fadeElements.forEach((el, index) => {
        el.style.opacity = '0';
        el.style.transform = 'translateY(20px)';

        setTimeout(() => {
            el.style.transition = 'opacity 0.6s ease-out, transform 0.6s ease-out';
            el.style.opacity = '1';
            el.style.transform = 'translateY(0)';
        }, 100 + index * 100);
    });

    // Card hover effects
    const cards = document.querySelectorAll('.card, .dashboard-card, .service-item');
    cards.forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-4px) scale(1.02)';
            this.style.boxShadow = '0 12px 24px rgba(0, 0, 0, 0.15)';
        });

        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0) scale(1)';
            this.style.boxShadow = '';
        });
    });
}

// Form validation system
function initFormValidation() {
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        // Real-time validation
        const inputs = form.querySelectorAll('input[required]');
        inputs.forEach(input => {
            input.addEventListener('blur', function() {
                validateField(this);
            });

            input.addEventListener('input', function() {
                clearFieldError(this);
            });
        });

        // Form submission validation
        form.addEventListener('submit', function(e) {
            if (!validateForm(this)) {
                e.preventDefault();
                showNotification('Пожалуйста, исправьте ошибки в форме', 'error');
            }
        });
    });
}

function validateForm(form) {
    const inputs = form.querySelectorAll('input[required]');
    let isValid = true;

    inputs.forEach(input => {
        if (!validateField(input)) {
            isValid = false;
        }
    });

    return isValid;
}

function validateField(input) {
    const value = input.value.trim();
    let isValid = true;
    let message = '';

    // Required field validation
    if (!value) {
        message = 'Это поле обязательно для заполнения';
        isValid = false;
    } else {
        // Field-specific validation
        switch(input.name) {
            case 'username':
                if (value.length < 3) {
                    message = 'Имя пользователя должно содержать минимум 3 символа';
                    isValid = false;
                } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
                    message = 'Имя пользователя может содержать только буквы, цифры и подчёркивания';
                    isValid = false;
                }
                break;

            case 'password':
                if (value.length < 8) {
                    message = 'Пароль должен содержать минимум 8 символов';
                    isValid = false;
                } else if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(value)) {
                    message = 'Пароль должен содержать буквы в верхнем и нижнем регистре и цифры';
                    isValid = false;
                }
                break;

            case 'email':
                if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
                    message = 'Введите корректный email адрес';
                    isValid = false;
                }
                break;
        }
    }

    if (!isValid) {
        showFieldError(input, message);
    } else {
        clearFieldError(input);
        showFieldSuccess(input);
    }

    return isValid;
}

function showFieldError(input, message) {
    clearFieldError(input);

    const errorDiv = document.createElement('div');
    errorDiv.className = 'field-error';
    errorDiv.textContent = message;

    input.classList.add('error');
    input.parentNode.appendChild(errorDiv);
}

function showFieldSuccess(input) {
    input.classList.remove('error');
    input.classList.add('success');

    setTimeout(() => {
        input.classList.remove('success');
    }, 2000);
}

function clearFieldError(input) {
    const existingError = input.parentNode.querySelector('.field-error');
    if (existingError) {
        existingError.remove();
    }
    input.classList.remove('error');
}

// Notification system
function initNotifications() {
    // Auto-hide success alerts after 5 seconds
    const alerts = document.querySelectorAll('.alert:not(.alert-error)');
    alerts.forEach(alert => {
        setTimeout(() => {
            fadeOutElement(alert);
        }, 5000);
    });

    // Close button for alerts
    const alertCloseButtons = document.querySelectorAll('.alert-close');
    alertCloseButtons.forEach(btn => {
        btn.addEventListener('click', function() {
            fadeOutElement(this.parentElement);
        });
    });
}

function showNotification(message, type = 'info', duration = 5000) {
    const alert = document.createElement('div');
    alert.className = `alert alert-${type} fade-in`;
    alert.innerHTML = `
        <span>${getNotificationIcon(type)}</span>
        <span>${message}</span>
        <button class="alert-close">&times;</button>
    `;

    // Add styles if not present
    if (!document.querySelector('#alert-styles')) {
        const styles = document.createElement('style');
        styles.id = 'alert-styles';
        styles.textContent = `
            .alert-close {
                background: none;
                border: none;
                font-size: 1.2rem;
                cursor: pointer;
                margin-left: auto;
                opacity: 0.7;
            }
            .alert-close:hover {
                opacity: 1;
            }
            .field-error {
                color: var(--error);
                font-size: 0.875rem;
                margin-top: 0.25rem;
            }
            .form-input.error {
                border-color: var(--error);
            }
            .form-input.success {
                border-color: var(--success);
            }
        `;
        document.head.appendChild(styles);
    }

    document.body.appendChild(alert);

    // Add close functionality
    const closeBtn = alert.querySelector('.alert-close');
    closeBtn.addEventListener('click', () => fadeOutElement(alert));

    if (duration > 0 && type !== 'error') {
        setTimeout(() => fadeOutElement(alert), duration);
    }

    return alert;
}

function getNotificationIcon(type) {
    const icons = {
        success: '✅',
        error: '❌',
        warning: '⚠️',
        info: 'ℹ️'
    };
    return icons[type] || 'ℹ️';
}

// Navigation system
function initNavigation() {
    // Mobile menu toggle (if needed in future)
    const mobileMenuToggle = document.querySelector('.mobile-menu-toggle');
    if (mobileMenuToggle) {
        mobileMenuToggle.addEventListener('click', function() {
            const navMenu = document.querySelector('.nav-menu');
            navMenu.classList.toggle('active');
        });
    }

    // Active link highlighting
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-link');

    navLinks.forEach(link => {
        const linkPath = link.getAttribute('href');
        if (currentPath === linkPath ||
            (currentPath.startsWith(linkPath) && linkPath !== '/')) {
            link.classList.add('active');
        } else {
            link.classList.remove('active');
        }
    });
}

// User data management
function initUserData() {
    // Try to get username from backend or use fallback
    const usernameElement = document.getElementById('username');
    const welcomeUsernameElement = document.getElementById('welcomeUsername');

    if (usernameElement || welcomeUsernameElement) {
        loadUserData().then(userData => {
            const username = userData?.username || 'Пользователь';
            if (usernameElement) usernameElement.textContent = username;
            if (welcomeUsernameElement) welcomeUsernameElement.textContent = username;
        }).catch(error => {
            console.log('User data not available, using fallback');
        });
    }
}

async function loadUserData() {
    try {
        // This would be an API call in a real app
        // For now, we'll simulate it
        return await simulateUserDataLoad();
    } catch (error) {
        console.error('Failed to load user data:', error);
        throw error;
    }
}

function simulateUserDataLoad() {
    return new Promise((resolve) => {
        setTimeout(() => {
            // Simulate user data
            resolve({
                username: 'Тестовый Пользователь',
                id: 'user-123',
                roles: ['USER']
            });
        }, 500);
    });
}

// Service status monitoring
function initServiceStatus() {
    const serviceStatusElements = document.querySelectorAll('.service-status');

    if (serviceStatusElements.length > 0) {
        checkServicesStatus();

        // Update status every 30 seconds
        setInterval(checkServicesStatus, 30000);
    }
}

async function checkServicesStatus() {
    const services = [
        { name: 'Auth Service', url: '/actuator/health' },
        { name: 'Deposit Service', url: '/deposits/health' },
        { name: 'Transaction Service', url: '/transactions/health' }
    ];

    for (const service of services) {
        await updateServiceStatus(service.name, service.url);
    }
}

async function updateServiceStatus(serviceName, healthUrl) {
    try {
        const response = await fetch(healthUrl, {
            method: 'GET',
            headers: { 'Accept': 'application/json' }
        });

        const status = response.ok ? 'online' : 'error';
        updateServiceUI(serviceName, status);
    } catch (error) {
        updateServiceUI(serviceName, 'offline');
    }
}

function updateServiceUI(serviceName, status) {
    const serviceItems = document.querySelectorAll('.service-item');

    serviceItems.forEach(item => {
        const nameElement = item.querySelector('.service-name');
        const statusElement = item.querySelector('.service-status');

        if (nameElement && nameElement.textContent.includes(serviceName)) {
            statusElement.className = 'service-status';
            statusElement.classList.add(`status-${status}`);

            const statusIcons = {
                online: '✅ Online',
                starting: '⏳ Starting',
                offline: '❌ Offline',
                error: '⚠️ Error'
            };

            statusElement.textContent = statusIcons[status] || '❓ Unknown';
        }
    });
}

// Utility functions
function fadeOutElement(element, duration = 500) {
    element.style.transition = `opacity ${duration}ms ease-out`;
    element.style.opacity = '0';

    setTimeout(() => {
        if (element.parentNode) {
            element.parentNode.removeChild(element);
        }
    }, duration);
}

function formatCurrency(amount, currency = 'RUB', locale = 'ru-RU') {
    return new Intl.NumberFormat(locale, {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 2
    }).format(amount);
}

function formatDate(date, locale = 'ru-RU') {
    return new Intl.DateTimeFormat(locale, {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    }).format(new Date(date));
}

// Global error handling
function handleGlobalError(event) {
    console.error('Global error:', event.error);

    // Don't show error notifications for minor issues
    if (event.error && event.error.message &&
        !event.error.message.includes('ResizeObserver') &&
        !event.error.message.includes('Script error')) {

        showNotification('Произошла непредвиденная ошибка', 'error', 0);
    }
}

// API service functions (for future use)
const BankAPI = {
    async get(url, options = {}) {
        try {
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    },

    async post(url, data, options = {}) {
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                body: JSON.stringify(data),
                ...options
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }
};

// Make utility functions globally available
window.BankApp = {
    showNotification,
    formatCurrency,
    formatDate,
    BankAPI
};

// Export for module usage (if needed)
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { BankApp };
}