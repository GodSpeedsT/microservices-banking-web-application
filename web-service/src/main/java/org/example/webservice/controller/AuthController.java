package org.example.webservice.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.dto.UserResponse;
import org.example.webservice.service.ApiService;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
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

    // GET /auth/login
    // Эта страница отображает форму. Spring Security сам инициирует редирект на Auth Server
    @GetMapping("/login")
    public String loginForm(Model model) {
        // Если пользователь уже аутентифицирован, отправляем его в кабинет
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth != null && auth.isAuthenticated() && !"anonymousUser".equals(auth.getPrincipal())) {
            return "redirect:/dashboard";
        }

        model.addAttribute("title", "Вход в систему");
        // Статус isAuthenticated будет определяться в layout.html через th:authorize
        return "login";
    }

    // GET /auth/register
    @GetMapping("/register")
    public String registerForm(Model model) {
        model.addAttribute("title", "Регистрация");
        return "register";
    }

    // POST /auth/register
    @PostMapping("/register")
    public String register(@RequestParam @Valid String username,
                           @RequestParam @Valid String password,
                           RedirectAttributes redirectAttributes,
                           Model model) {
        try {
            AuthRequest authRequest = new AuthRequest(username, password);
            apiService.registerUser(authRequest); // Регистрируем через Auth Service

            redirectAttributes.addFlashAttribute("success",
                    "Регистрация успешна! Теперь вы можете войти в систему.");
            return "redirect:/auth/login";
        } catch (Exception e) {
            log.error("Registration failed: {}", e.getMessage());
            model.addAttribute("error", "Ошибка при регистрации: " + e.getMessage());
            return "register";
        }
    }

    // GET /auth/logout
    // Просто редиректим на стандартный эндпоинт Spring Security, который очистит сессию.
    @GetMapping("/logout")
    public String logout(RedirectAttributes redirectAttributes) {
        redirectAttributes.addFlashAttribute("success", "Вы успешно вышли из системы.");
        return "redirect:/logout";
    }
}