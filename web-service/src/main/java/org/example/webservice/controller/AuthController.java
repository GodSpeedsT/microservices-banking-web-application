package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.dto.AuthResponse;
import org.example.webservice.service.ApiService;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

@Controller
@RequestMapping("/auth")
@RequiredArgsConstructor
@Slf4j
public class AuthController {

    private final ApiService apiService;
    private final AuthService authService;

    @GetMapping("/login")
    public String loginForm(Model model, HttpSession session) {
        if (authService.isAuthenticated(session)) {
            return "redirect:/dashboard";
        }
        model.addAttribute("title", "Вход в систему");
        model.addAttribute("isAuthenticated", false);
        return "login";
    }

    @PostMapping("/login")
    public String login(@RequestParam String username,
                        @RequestParam String password,
                        HttpSession session,
                        Model model,
                        RedirectAttributes redirectAttributes) {
        try {
            log.info("Login attempt for user: {}", username);

            AuthRequest authRequest = new AuthRequest(username, password);
            AuthResponse response = apiService.loginUser(authRequest);

            if (response != null && response.getAccessToken() != null) {
                authService.loginUser(session, response.getAccessToken(), username);
                log.info("User {} successfully logged in", username);

                redirectAttributes.addFlashAttribute("success",
                        "Добро пожаловать, " + username + "!");
                return "redirect:/dashboard";
            } else {
                model.addAttribute("error", "Неверные учетные данные");
                return "login";
            }
        } catch (Exception e) {
            log.warn("Auth service unavailable, using stub mode: {}", e.getMessage());
            return handleStubLogin(username, password, session, model, redirectAttributes);
        }
    }

    @GetMapping("/register")
    public String registerForm(Model model) {
        model.addAttribute("title", "Регистрация");
        model.addAttribute("isAuthenticated", false);
        return "register";
    }

    @PostMapping("/register")
    public String register(@RequestParam String username,
                           @RequestParam String password,
                           Model model,
                           RedirectAttributes redirectAttributes) {
        try {
            log.info("Registration attempt for user: {}", username);

            AuthRequest authRequest = new AuthRequest(username, password);
            Object response = apiService.registerUser(authRequest);

            redirectAttributes.addFlashAttribute("success",
                    "Регистрация успешна! Теперь вы можете войти в систему.");
            return "redirect:/auth/login";

        } catch (Exception e) {
            log.warn("Auth service unavailable during registration: {}", e.getMessage());

            if (isAuthServiceUnavailable(e)) {
                redirectAttributes.addFlashAttribute("warning",
                        "Регистрация выполнена в демо-режиме. Auth-service временно недоступен.");
                return "redirect:/auth/login";
            }

            model.addAttribute("error", "Ошибка при регистрации: " + e.getMessage());
            return "register";
        }
    }

    @GetMapping("/logout")
    public String logout(HttpSession session, RedirectAttributes redirectAttributes) {
        String username = authService.getUsername(session);
        authService.logoutUser(session);

        redirectAttributes.addFlashAttribute("success",
                "Вы успешно вышли из системы.");
        log.info("User {} logged out", username);

        return "redirect:/";
    }

    private String handleStubLogin(String username, String password,
                                   HttpSession session, Model model,
                                   RedirectAttributes redirectAttributes) {
        if (isValidStubCredentials(username, password)) {
            String stubToken = "stub-token-" + System.currentTimeMillis();
            authService.loginUser(session, stubToken, username);

            redirectAttributes.addFlashAttribute("warning",
                    "⚠️ Авторизация в режиме заглушки. Auth-service временно недоступен.");
            return "redirect:/dashboard";
        } else {
            model.addAttribute("error",
                    "Неверные учетные данные. Для демо-режима используйте: demo/demo");
            return "login";
        }
    }

    private boolean isValidStubCredentials(String username, String password) {
        return "demo".equals(username) && "demo".equals(password) ||
                "admin".equals(username) && "admin".equals(password) ||
                "user".equals(username) && "user".equals(password);
    }

    private boolean isAuthServiceUnavailable(Exception e) {
        if (e == null || e.getMessage() == null) return true;

        String message = e.getMessage().toLowerCase();
        return message.contains("connection") ||
                message.contains("unavailable") ||
                message.contains("timeout") ||
                message.contains("refused") ||
                message.contains("503") ||
                message.contains("500") ||
                message.contains("i/o error");
    }
}