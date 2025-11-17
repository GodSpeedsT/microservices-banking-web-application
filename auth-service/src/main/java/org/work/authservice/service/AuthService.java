package org.work.authservice.service;

import org.springframework.stereotype.Service;
import org.work.authservice.dto.AuthResponse;
import org.work.authservice.entity.User;
import org.work.authservice.security.JwtUtil;

@Service
public class AuthService {
    private final UserService userService;
    private final JwtUtil jwtUtil;

    public AuthService(UserService userService, JwtUtil jwtUtil) {
        this.userService = userService;
        this.jwtUtil = jwtUtil;
    }

    public AuthResponse login(String username, String password) {
        User user = userService.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("Пользователь не найден"));

        if (!userService.getPasswordEncoder().matches(password, user.getPassword())) {
            throw new RuntimeException("Пароль неверный");
        }

        String accessToken = jwtUtil.generateToken(user);
        String refreshToken = jwtUtil.generateRefreshToken(user);

        return new AuthResponse(accessToken, refreshToken);
    }

    public AuthResponse refreshToken(String refreshToken) {
        if (!jwtUtil.validateToken(refreshToken)) {
            throw new RuntimeException("Недействительный refresh token");
        }

        String username = jwtUtil.getUsernameFromToken(refreshToken);
        User user = userService.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("Пользователь не найден"));

        String newAccessToken = jwtUtil.generateToken(user);
        return new AuthResponse(newAccessToken, refreshToken);
    }
}