package org.work.authservice.security;

import org.springframework.security.oauth2.jwt.*;
import org.springframework.stereotype.Component;
import org.work.authservice.entity.User;

import java.time.Instant;

@Component
public class JwtUtil {
    private final JwtEncoder jwtEncoder;
    private final JwtDecoder jwtDecoder;

    public JwtUtil(JwtEncoder jwtEncoder, JwtDecoder jwtDecoder) {
        this.jwtEncoder = jwtEncoder;
        this.jwtDecoder = jwtDecoder;
    }

    public String generateToken(User user) {
        Instant now = Instant.now();
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer("auth-service")
                .subject(user.getUsername())
                .claim("roles", user.getRoles().stream().map(role -> role.getName()).toList())
                .claim("type", "access")
                .issuedAt(now)
                .expiresAt(now.plusSeconds(3600)) // 1 hour
                .build();

        return jwtEncoder.encode(JwtEncoderParameters.from(claims)).getTokenValue();
    }

    public String generateRefreshToken(User user) {
        Instant now = Instant.now();
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .issuer("auth-service")
                .subject(user.getUsername())
                .claim("type", "refresh")
                .issuedAt(now)
                .expiresAt(now.plusSeconds(86400)) // 24 hours
                .build();

        return jwtEncoder.encode(JwtEncoderParameters.from(claims)).getTokenValue();
    }

    public boolean validateToken(String token) {
        try {
            Jwt jwt = jwtDecoder.decode(token);
            return jwt.getExpiresAt() == null || jwt.getExpiresAt().isAfter(Instant.now());
        } catch (JwtException e) {
            return false;
        }
    }

    public String getUsernameFromToken(String token) {
        try {
            Jwt jwt = jwtDecoder.decode(token);
            return jwt.getSubject();
        } catch (JwtException e) {
            throw new RuntimeException("Недействительный токен", e);
        }
    }

    public String getTokenType(String token) {
        try {
            Jwt jwt = jwtDecoder.decode(token);
            return jwt.getClaim("type");
        } catch (JwtException e) {
            throw new RuntimeException("Недействительный токен", e);
        }
    }
}