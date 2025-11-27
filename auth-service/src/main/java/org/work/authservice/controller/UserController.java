package org.work.authservice.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.work.authservice.dto.UserResponse;
import org.work.authservice.entity.User;
import org.work.authservice.service.UserService;

import java.util.List;
import java.util.Map;

@RestController
@RequiredArgsConstructor
@Slf4j
public class UserController {

    private final UserService userService;

    @GetMapping("/userinfo")
    public ResponseEntity<Map<String, Object>> getUserInfo(Authentication authentication) {
        String username = authentication.getName();

        log.info("Fetching user info for: {}", username);

        User user = userService.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("User not found despite valid token"));

        List<String> roles = authentication.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .toList();

        // Стандартные OIDC claims
        Map<String, Object> userInfo = Map.of(
                "sub", user.getUsername(),
                "preferred_username", user.getUsername(),
                "roles", roles
        );

        return ResponseEntity.ok(userInfo);
    }

    // Ваш кастомный endpoint (оставьте для обратной совместимости)
    @GetMapping("/auth/user/info")
    public ResponseEntity<UserResponse> getCustomUserInfo(Authentication authentication) {
        String username = authentication.getName();

        log.info("Fetching custom user info for: {}", username);

        User user = userService.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("User not found despite valid token"));

        List<String> roles = authentication.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .toList();

        UserResponse response = new UserResponse(
                user.getId(),
                user.getUsername(),
                roles
        );

        return ResponseEntity.ok(response);
    }
}