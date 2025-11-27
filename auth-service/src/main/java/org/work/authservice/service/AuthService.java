package org.work.authservice.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.authservice.entity.User;

import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {

    private final UserService userService;
    // PasswordEncoder оставлен для возможности будущих внутренних проверок, но в основном используется в UserService
    private final PasswordEncoder passwordEncoder;

    /**
     * Методы login() и refreshToken() удалены.
     * * Теперь:
     * 1. Регистрация: /auth/register (остается здесь)
     * 2. Логин (получение токена):
     * - Для браузера: через стандартную форму /login (обработку берет на себя Spring Security)
     * - Для API: через стандартный эндпоинт OAuth2 /oauth2/token
     */

    @Transactional
    public User register(String username, String rawPassword) {
        log.info("Attempting registration for user: {}", username);

        // Вся логика валидации и сохранения перенесена в UserService,
        // чтобы AuthService был максимально тонким.
        User newUser = userService.registerUser(username, rawPassword);

        log.info("User registered successfully: {}", username);
        return newUser;
    }

    // Этот метод теперь нужен для UserService, чтобы он мог находить пользователей
    public Optional<User> findByUsername(String username) {
        return userService.findByUsername(username);
    }
}