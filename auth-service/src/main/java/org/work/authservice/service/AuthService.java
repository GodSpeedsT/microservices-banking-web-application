package org.work.authservice.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.authservice.entity.User;

import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {

    private final UserService userService;
    private final PasswordEncoder passwordEncoder;

    @Transactional
    public User register(String username, String rawPassword) {
        log.info("Attempting registration for user: {}", username);

        User newUser = userService.registerUser(username, rawPassword);

        log.info("User registered successfully: {}", username);
        return newUser;
    }

    public Optional<User> findByUsername(String username) {
        return userService.findByUsername(username);
    }
}