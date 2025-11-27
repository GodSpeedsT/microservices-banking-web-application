package org.work.authservice.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.authservice.dto.AuthRequest;
import org.work.authservice.dto.UserResponse;
import org.work.authservice.entity.User;
import org.work.authservice.service.UserService; // Теперь используем UserService напрямую для регистрации

@RestController
@RequestMapping("/auth")
@RequiredArgsConstructor
@Slf4j
public class AuthController {

    // Используем UserService напрямую для регистрации
    private final UserService userService;

    @PostMapping("/register")
    public ResponseEntity<UserResponse> register(@Valid @RequestBody AuthRequest request) {
        log.info("Registration attempt for user: {}", request.getUsername());

        User user = userService.registerUser(request.getUsername(), request.getPassword());

        UserResponse response = new UserResponse(user.getId(), user.getUsername(),
                user.getRoles().stream().map(role -> role.getName()).toList());

        return ResponseEntity.ok(response);
    }
}