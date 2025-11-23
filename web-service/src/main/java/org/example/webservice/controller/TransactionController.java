package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/transactions")
public class TransactionController {

    private final AuthService authService;

    public TransactionController(AuthService authService) {
        this.authService = authService;
    }

    @GetMapping
    public String transactions(Model model, HttpSession session) {
        if (!authService.isAuthenticated(session)) {
            return "redirect:/auth/login";
        }
        model.addAttribute("title", "История операций");
        model.addAttribute("username", authService.getUsername(session));
        return "transactions";
    }
}
