package org.work.depositservice.service;

import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class SecurityContextService {

    public String getCurrentUserId() {
        return getJwt()
                .map(jwt -> jwt.getSubject())
                .orElseThrow(() -> new RuntimeException("Пользователь не аутентифицирован"));
    }

    public Optional<String> getCurrentUsername() {
        return getJwt()
                .map(jwt -> jwt.getClaimAsString("preferred_username"));
    }

    public boolean hasRole(String role) {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication != null) {
            return authentication.getAuthorities().stream()
                    .anyMatch(authority -> authority.getAuthority().equals("ROLE_" + role));
        }
        return false;
    }

    public boolean isAdmin() {
        return hasRole("ADMIN");
    }

    public boolean canAccessUserData(String userId) {
        return getCurrentUserId().equals(userId) || isAdmin();
    }

    private Optional<Jwt> getJwt() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication instanceof JwtAuthenticationToken) {
            JwtAuthenticationToken jwtAuth = (JwtAuthenticationToken) authentication;
            return Optional.of(jwtAuth.getToken());
        }
        return Optional.empty();
    }

    public <T> Optional<T> getClaim(String claimName, Class<T> claimType) {
        return getJwt().map(jwt -> jwt.getClaim(claimName));
    }
}