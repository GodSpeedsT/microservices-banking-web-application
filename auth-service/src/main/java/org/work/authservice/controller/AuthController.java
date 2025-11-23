package org.work.authservice.controller;

import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.authservice.dto.AuthRequest;
import org.work.authservice.dto.AuthResponse;
import org.work.authservice.dto.RefreshRequest;
import org.work.authservice.dto.UserResponse;
import org.work.authservice.entity.User;
import org.work.authservice.service.AuthService;
import org.work.authservice.service.UserService;

@RestController
@RequestMapping("/auth")
@RequiredArgsConstructor
@Slf4j
public class AuthController {
    private final AuthService authService;
    private final UserService userService;

    @PostMapping("/register")
    public ResponseEntity<UserResponse> register(@Valid @RequestBody AuthRequest request) {
        log.info("Registration attempt for user: {}", request.getUsername());
        User user = userService.registerUser(request.getUsername(), request.getPassword());
        UserResponse response = new UserResponse(user.getId(), user.getUsername(),
                user.getRoles().stream().map(role -> role.getName()).toList());
        return ResponseEntity.ok(response);
    }

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@Valid @RequestBody AuthRequest request) {
        log.info("Login attempt for user: {}", request.getUsername());
        return ResponseEntity.ok(authService.login(request.getUsername(), request.getPassword()));
    }

    @PostMapping("/refresh")
    public ResponseEntity<AuthResponse> refresh(@Valid @RequestBody RefreshRequest request) {
        log.info("Token refresh attempt");
        return ResponseEntity.ok(authService.refreshToken(request.getRefreshToken()));
    }

    @PostMapping("/logout")
    public ResponseEntity<String> logout() {
        return ResponseEntity.ok("Logged out successfully");
    }
}