package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/")
@RequiredArgsConstructor
public class BankController {

    private final AuthService authService;

    @GetMapping
    public String index(HttpSession session, Model model) {
        boolean isAuthenticated = authService.isAuthenticated(session);
        model.addAttribute("isAuthenticated", isAuthenticated);

        if (isAuthenticated) {
            model.addAttribute("username", authService.getUsername(session));
        }

        return "index";
    }
}