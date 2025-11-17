package org.work.authservice.controller;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.authservice.dto.AuthRequest;
import org.work.authservice.dto.AuthResponse;
import org.work.authservice.dto.RefreshRequest;
import org.work.authservice.entity.Role;
import org.work.authservice.entity.User;
import org.work.authservice.service.AuthService;
import org.work.authservice.service.UserService;

import java.util.Set;

@RestController
@RequestMapping("/auth")
public class AuthController {

    private final AuthService authService;
    private final UserService userService;
    private static final Logger log = LoggerFactory.getLogger(AuthController.class);

    public AuthController(AuthService authService, UserService userService) {
        this.authService = authService;
        this.userService = userService;
    }

    @PostMapping("/register")
    public ResponseEntity<User> register(@RequestBody AuthRequest request) {
        Set<Role> defaultRoles = Set.of(new Role("ROLE_USER"));
        User user = userService.registerUser(request.getUsername(), request.getPassword(), defaultRoles);
        return ResponseEntity.ok(user);
    }

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@RequestBody AuthRequest request) {
        return ResponseEntity.ok(authService.login(request.getUsername(), request.getPassword()));
    }

    @PostMapping("/refresh")
    public ResponseEntity<AuthResponse> refresh(@RequestBody RefreshRequest request) {
        return ResponseEntity.ok(authService.refreshToken(request.getRefreshToken()));
    }

    @PostMapping("/logout")
    public ResponseEntity<String> logout() {
        // В stateless приложении logout обычно обрабатывается на клиенте
        return ResponseEntity.ok("Logged out successfully");
    }
}