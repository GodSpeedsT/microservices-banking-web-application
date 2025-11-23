package org.example.webservice.controller;

import jakarta.servlet.http.HttpSession;
import lombok.RequiredArgsConstructor;
import org.example.webservice.service.AuthService;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/dashboard")
@RequiredArgsConstructor
public class UserController {

    private final AuthService authService;

    @GetMapping
    public String dashboard(Model model, HttpSession session) {
        if (!authService.isAuthenticated(session)) {
            return "redirect:/auth/login";
        }

        String username = authService.getUsername(session);

        model.addAttribute("title", "Личный кабинет");
        model.addAttribute("username", username);
        model.addAttribute("isAuthenticated", true);
        model.addAttribute("isStubMode", true);

        return "dashboard";
    }
}