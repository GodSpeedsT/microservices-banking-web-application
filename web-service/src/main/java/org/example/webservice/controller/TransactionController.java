package org.example.webservice.controller;

import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

@Controller
@RequestMapping("/transactions")
// Убрана инъекция AuthService
public class TransactionController {

    @GetMapping
    public String transactions(Model model, @AuthenticationPrincipal OidcUser principal) {

        String username = principal.getClaimAsString("preferred_username");
        if (username == null) {
            username = principal.getSubject();
        }

        model.addAttribute("title", "История операций");
        model.addAttribute("username", username);
        return "transactions";
    }
}