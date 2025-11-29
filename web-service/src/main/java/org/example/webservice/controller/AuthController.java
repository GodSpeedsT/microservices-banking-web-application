package org.example.webservice.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.dto.UserResponse;
import org.example.webservice.service.ApiService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.servlet.mvc.support.RedirectAttributes;

@Controller
@RequestMapping("/auth")
@RequiredArgsConstructor
@Slf4j
public class AuthController {

    private final ApiService apiService;

    // УДАЛИТЕ эти методы - они будут обрабатываться статическими файлами
    // @GetMapping("/login")
    // public String loginForm(Model model) {
    //     return "login";
    // }

    // @GetMapping("/register")
    // public String registerForm(Model model) {
    //     return "register";
    // }

    // Оставьте только POST методы для обработки форм
    @PostMapping("/register")
    public String register(@RequestParam @Valid String username,
                           @RequestParam @Valid String password,
                           RedirectAttributes redirectAttributes,
                           Model model) {
        try {
            AuthRequest authRequest = new AuthRequest(username, password);
            apiService.registerUser(authRequest);

            redirectAttributes.addFlashAttribute("success",
                    "Регистрация успешна! Теперь вы можете войти в систему.");
            return "redirect:/auth/login";
        } catch (Exception e) {
            log.error("Registration failed: {}", e.getMessage());
            model.addAttribute("error", "Ошибка при регистрации: " + e.getMessage());
            return "forward:/register.html";
        }
    }
}