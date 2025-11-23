package org.work.authservice.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Service
public class ValidationService {

    @Value("${password.policy.min-length:8}")
    private int minLength;

    @Value("${password.policy.require-uppercase:true}")
    private boolean requireUppercase;

    @Value("${password.policy.require-lowercase:true}")
    private boolean requireLowercase;

    @Value("${password.policy.require-digits:true}")
    private boolean requireDigits;

    @Value("${password.policy.require-special-chars:true}")
    private boolean requireSpecialChars;

    public void validatePassword(String password) {
        if (password == null || password.length() < minLength) {
            throw new RuntimeException("Password must be at least " + minLength + " characters long");
        }

        if (requireUppercase && !password.matches(".*[A-Z].*")) {
            throw new RuntimeException("Password must contain at least one uppercase letter");
        }

        if (requireLowercase && !password.matches(".*[a-z].*")) {
            throw new RuntimeException("Password must contain at least one lowercase letter");
        }

        if (requireDigits && !password.matches(".*\\d.*")) {
            throw new RuntimeException("Password must contain at least one digit");
        }

        if (requireSpecialChars && !password.matches(".*[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>\\/?].*")) {
            throw new RuntimeException("Password must contain at least one special character");
        }
    }

    public void validateUsername(String username) {
        if (username == null || username.length() < 3) {
            throw new RuntimeException("Username must be at least 3 characters long");
        }

        if (!username.matches("^[a-zA-Z0-9._-]+$")) {
            throw new RuntimeException("Username can only contain letters, numbers, dots, dashes and underscores");
        }
    }
}