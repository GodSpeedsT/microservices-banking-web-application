package org.example.webservice.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.HashMap;
import java.util.Map;

@Controller
@RequestMapping("/dashboard")
@RequiredArgsConstructor
@Slf4j
public class DashboardController {

    // Класс AuthService был удален, инъекция больше не нужна.

    @GetMapping
    // Spring Security гарантирует, что здесь будет аутентифицированный пользователь
    public String dashboard(Model model, @AuthenticationPrincipal OidcUser principal) {

        // Получаем имя пользователя из токена (claim 'preferred_username' или 'sub')
        // Используем getPreferredUsername, если Auth Server его предоставляет, иначе getSubject
        String username = principal.getClaimAsString("preferred_username");
        if (username == null) {
            username = principal.getSubject();
        }

        // Добавляем информацию для dashboard
        model.addAttribute("title", "Личный кабинет");
        model.addAttribute("username", username);
        // model.addAttribute("isAuthenticated", true); - больше не нужно

        // Статус сервисов (можно расширить)
        Map<String, String> servicesStatus = new HashMap<>();
        servicesStatus.put("Auth Service", "✅ Online");
        servicesStatus.put("Deposit Service", "⏳ Starting");
        servicesStatus.put("Transaction Service", "⏳ Starting");

        model.addAttribute("servicesStatus", servicesStatus);
        model.addAttribute("isStubMode", true);

        return "dashboard";
    }
}