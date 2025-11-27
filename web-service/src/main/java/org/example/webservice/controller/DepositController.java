package org.example.webservice.controller;

import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Controller
@RequestMapping("/deposits")
// Убрана инъекция AuthService
public class DepositController {

    @GetMapping
    public String deposits(Model model, @AuthenticationPrincipal OidcUser principal) {

        String username = principal.getClaimAsString("preferred_username");
        if (username == null) {
            username = principal.getSubject();
        }

        // ЗАГЛУШКА: вместо реальных данных из deposit-service
        List<Map<String, Object>> stubDeposits = createStubDeposits();

        model.addAttribute("title", "Мои депозиты");
        model.addAttribute("username", username);
        model.addAttribute("deposits", stubDeposits);
        model.addAttribute("message", "Депозитный сервис временно недоступен. Отображаются тестовые данные.");
        model.addAttribute("serviceStatus", "stub");

        return "deposits";
    }

    private List<Map<String, Object>> createStubDeposits() {
        // ... (метод без изменений)
        List<Map<String, Object>> deposits = new ArrayList<>();

        deposits.add(Map.of(
                "id", "DEP001",
                "type", "Накопительный",
                "amount", 50000.0,
                "interestRate", 5.5,
                "startDate", "2024-01-15",
                "endDate", "2025-01-15",
                "status", "ACTIVE",
                "currency", "RUB"
        ));

        deposits.add(Map.of(
                "id", "DEP002",
                "type", "Срочный",
                "amount", 100000.0,
                "interestRate", 7.0,
                "startDate", "2024-02-01",
                "endDate", "2026-02-01",
                "status", "ACTIVE",
                "currency", "RUB"
        ));

        return deposits;
    }
}