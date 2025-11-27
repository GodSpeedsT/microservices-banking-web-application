package org.example.webservice.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/")
@RequiredArgsConstructor
public class BankController {

    @GetMapping
    public String index(Model model) {
        // Проверка аутентификации в контексте Spring Security
        boolean isAuthenticated = SecurityContextHolder.getContext().getAuthentication().isAuthenticated();

        // В Thymeleaf эта переменная больше не нужна, но оставим для совместимости с шаблоном index.html
        model.addAttribute("isAuthenticated", isAuthenticated);

        if (isAuthenticated) {
            // Имя пользователя будет доступно через #authentication.name в шаблоне
            model.addAttribute("username", SecurityContextHolder.getContext().getAuthentication().getName());
        }

        return "index";
    }
}