package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

@Controller
@RequestMapping("/deposits")
public class DepositController {

    private final AuthService authService;

    public DepositController(AuthService authService) {
        this.authService = authService;
    }

    @GetMapping
    public String deposits(Model model, HttpSession session) {
        if (!authService.isAuthenticated(session)) {
            return "redirect:/auth/login";
        }

        // ЗАГЛУШКА: вместо реальных данных из deposit-service
        List<Map<String, Object>> stubDeposits = createStubDeposits();

        model.addAttribute("title", "Мои депозиты");
        model.addAttribute("username", authService.getUsername(session));
        model.addAttribute("deposits", stubDeposits);
        model.addAttribute("message", "Депозитный сервис временно недоступен. Отображаются тестовые данные.");
        model.addAttribute("serviceStatus", "stub");

        return "deposits";
    }

    private List<Map<String, Object>> createStubDeposits() {
        List<Map<String, Object>> deposits = new ArrayList<>();

        // Тестовые данные депозитов
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