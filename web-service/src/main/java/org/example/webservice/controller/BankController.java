package org.example.webservice.controller;

import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/")
public class BankController {

    @GetMapping
    public String index(Model model) {
        model.addAttribute("appName", "Банк Онлайн");
        return "index";
    }

    @GetMapping("/login")
    public String login(Model model) {
        model.addAttribute("title", "Вход в систему");
        return "login";
    }

    @GetMapping("/register")
    public String register(Model model) {
        model.addAttribute("title", "Регистрация");
        return "register";
    }

    @GetMapping("/dashboard")
    public String dashboard(Model model) {
        model.addAttribute("title", "Личный кабинет");
        return "dashboard";
    }

    @GetMapping("/deposits")
    public String deposits(Model model) {
        model.addAttribute("title", "Мои депозиты");
        return "deposits";
    }

    @GetMapping("/transactions")
    public String transactions(Model model) {
        model.addAttribute("title", "История операций");
        return "transactions";
    }
}