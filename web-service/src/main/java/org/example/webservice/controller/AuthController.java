package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.service.ApiService;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/auth")
public class AuthController {

    private final ApiService apiService;
    private final AuthService authService;

    public AuthController(ApiService apiService, AuthService authService) {
        this.apiService = apiService;
        this.authService = authService;
    }

    @GetMapping("/login")
    public String loginForm(Model model, HttpSession session) {
        if (authService.isAuthenticated(session)) {
            return "redirect:/dashboard";
        }
        model.addAttribute("title", "Вход в систему");
        return "login";
    }

    @PostMapping("/login")
    public String login(@RequestParam String username,
                        @RequestParam String password,
                        HttpSession session,
                        Model model) {
        try {
            // Пытаемся использовать реальный сервис
            AuthRequest authRequest = new AuthRequest(username, password);
            Map<String, Object> response = (Map<String, Object>) apiService.loginUser(authRequest);

            if (response != null && response.containsKey("accessToken")) {
                String token = (String) response.get("accessToken");
                authService.loginUser(session, token, username);
                return "redirect:/dashboard";
            } else {
                model.addAttribute("error", "Неверные учетные данные");
                return "login";
            }
        } catch (Exception e) {
            // ЗАГЛУШКА: если auth-service недоступен, используем локальную аутентификацию
            if (isAuthServiceUnavailable(e)) {
                return handleStubLogin(username, password, session, model);
            }
            model.addAttribute("error", "Ошибка при входе: " + e.getMessage());
            return "login";
        }
    }

    @GetMapping("/register")
    public String registerForm(Model model) {
        model.addAttribute("title", "Регистрация");
        return "register";
    }

    @PostMapping("/register")
    public String register(@RequestParam String username,
                           @RequestParam String password,
                           Model model) {
        try {
            // Пытаемся использовать реальный сервис
            AuthRequest authRequest = new AuthRequest(username, password);
            Object response = apiService.registerUser(authRequest, null);
            model.addAttribute("success", "Регистрация успешна! Теперь вы можете войти.");
            return "login";
        } catch (Exception e) {
            // ЗАГЛУШКА: если auth-service недоступен, имитируем успешную регистрацию
            if (isAuthServiceUnavailable(e)) {
                model.addAttribute("success",
                        "Регистрация успешна! (режим заглушки). Логин: " + username);
                return "login";
            }
            model.addAttribute("error", "Ошибка при регистрации: " + e.getMessage());
            return "register";
        }
    }

    @GetMapping("/logout")
    public String logout(HttpSession session) {
        authService.logoutUser(session);
        return "redirect:/";
    }

    /**
     * ЗАГЛУШКА: Локальная аутентификация когда auth-service недоступен
     */
    private String handleStubLogin(String username, String password,
                                   HttpSession session, Model model) {
        // Простая проверка для демонстрации
        if (isValidStubCredentials(username, password)) {
            // Создаем заглушечный токен
            String stubToken = "stub-token-" + System.currentTimeMillis();
            authService.loginUser(session, stubToken, username);

            model.addAttribute("warning",
                    "⚠️ Авторизация в режиме заглушки. Auth-service временно недоступен.");
            return "redirect:/dashboard";
        } else {
            model.addAttribute("error",
                    "Неверные учетные данные для демо-режима. Попробуйте: demo/demo");
            return "login";
        }
    }

    /**
     * Простые демо-учетные данные
     */
    private boolean isValidStubCredentials(String username, String password) {
        return "demo".equals(username) && "demo".equals(password) ||
                "admin".equals(username) && "admin".equals(password) ||
                "user".equals(username) && "user".equals(password);
    }

    /**
     * Проверяем, что ошибка связана с недоступностью auth-service
     */
    private boolean isAuthServiceUnavailable(Exception e) {
        String message = e.getMessage().toLowerCase();
        return message.contains("connection") ||
                message.contains("unavailable") ||
                message.contains("timeout") ||
                message.contains("refused") ||
                message.contains("503") ||
                message.contains("500");
    }
}