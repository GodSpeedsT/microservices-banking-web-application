package org.work.authservice.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.work.authservice.dto.AuthResponse;
import org.work.authservice.entity.User;
import org.work.authservice.security.JwtUtil;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {
    private final UserService userService;
    private final JwtUtil jwtUtil;
    private final PasswordEncoder passwordEncoder; // Добавляем напрямую

    public AuthResponse login(String username, String password) {
        User user = userService.findByUsername(username)
                .orElseThrow(() -> {
                    log.warn("Login failed: user not found - {}", username);
                    return new RuntimeException("Invalid username or password");
                });

        // Используем passwordEncoder напрямую
        if (!passwordEncoder.matches(password, user.getPassword())) {
            log.warn("Login failed: invalid password for user - {}", username);
            throw new RuntimeException("Invalid username or password");
        }

        String accessToken = jwtUtil.generateToken(user);
        String refreshToken = jwtUtil.generateRefreshToken(user);

        log.info("User logged in successfully: {}", username);
        return new AuthResponse(accessToken, refreshToken);
    }

    public AuthResponse refreshToken(String refreshToken) {
        if (!jwtUtil.validateToken(refreshToken) || !"refresh".equals(jwtUtil.getTokenType(refreshToken))) {
            throw new RuntimeException("Invalid refresh token");
        }

        String username = jwtUtil.getUsernameFromToken(refreshToken);
        User user = userService.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("User not found"));

        String newAccessToken = jwtUtil.generateToken(user);
        log.info("Token refreshed for user: {}", username);
        return new AuthResponse(newAccessToken, refreshToken);
    }
}