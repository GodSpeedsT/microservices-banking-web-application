//package org.example.webservice.controller;
//
//import jakarta.servlet.http.HttpSession;
//import lombok.RequiredArgsConstructor;
//import lombok.extern.slf4j.Slf4j;
//import org.example.webservice.service.AuthService;
//import org.springframework.stereotype.Controller;
//import org.springframework.ui.Model;
//import org.springframework.web.bind.annotation.GetMapping;
//import org.springframework.web.bind.annotation.RequestMapping;
//
//import java.util.HashMap;
//import java.util.Map;
//
//@Controller
//@RequestMapping("/dashboard")
//@RequiredArgsConstructor
//@Slf4j
//public class DashboardController {
//
//    private final AuthService authService;
//
//    @GetMapping
//    public String dashboard(Model model, HttpSession session) {
//        if (!authService.isAuthenticated(session)) {
//            return "redirect:/auth/login";
//        }
//
//        String username = authService.getUsername(session);
//
//        // Добавляем информацию для dashboard
//        model.addAttribute("title", "Личный кабинет");
//        model.addAttribute("username", username);
//        model.addAttribute("isAuthenticated", true);
//
//        // Статус сервисов (можно расширить)
//        Map<String, String> servicesStatus = new HashMap<>();
//        servicesStatus.put("Auth Service", "✅ Online");
//        servicesStatus.put("Deposit Service", "⏳ Starting");
//        servicesStatus.put("Transaction Service", "⏳ Starting");
//
//        model.addAttribute("servicesStatus", servicesStatus);
//        model.addAttribute("isStubMode", true); // Пока работаем в режиме заглушки
//
//        return "dashboard";
//    }
//}