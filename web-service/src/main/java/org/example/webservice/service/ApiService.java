package org.example.webservice.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.dto.AuthResponse;
import org.example.webservice.dto.UserResponse;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
@RequiredArgsConstructor
@Slf4j
public class ApiService {

    private final RestTemplate restTemplate;

    @Value("${services.auth.url:http://localhost:8081}")
    private String authServiceUrl;

    @Value("${services.deposit.url:http://localhost:8082}")
    private String depositServiceUrl;

    @Value("${services.transaction.url:http://localhost:8083}")
    private String transactionServiceUrl;

    public AuthResponse loginUser(AuthRequest loginRequest) {
        String url = authServiceUrl + "/auth/login";
        log.debug("Calling auth service: {}", url);

        return post(url, loginRequest, AuthResponse.class, null);
    }

    public UserResponse registerUser(AuthRequest registerRequest) {
        String url = authServiceUrl + "/auth/register";
        log.debug("Calling auth service for registration: {}", url);

        return post(url, registerRequest, UserResponse.class, null);
    }

    public <T> T post(String url, Object request, Class<T> responseType, String token) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            if (token != null) {
                headers.setBearerAuth(token);
            }

            HttpEntity<Object> entity = new HttpEntity<>(request, headers);
            ResponseEntity<T> response = restTemplate.exchange(url, HttpMethod.POST, entity, responseType);

            log.debug("POST {} - Status: {}", url, response.getStatusCode());
            return response.getBody();

        } catch (Exception e) {
            log.error("Error calling POST {}: {}", url, e.getMessage());
            throw new RuntimeException("Service unavailable: " + e.getMessage(), e);
        }
    }

    public <T> T get(String url, Class<T> responseType, String token) {
        try {
            HttpHeaders headers = new HttpHeaders();
            headers.setContentType(MediaType.APPLICATION_JSON);
            if (token != null) {
                headers.setBearerAuth(token);
            }

            HttpEntity<String> entity = new HttpEntity<>(headers);
            ResponseEntity<T> response = restTemplate.exchange(url, HttpMethod.GET, entity, responseType);

            log.debug("GET {} - Status: {}", url, response.getStatusCode());
            return response.getBody();

        } catch (Exception e) {
            log.error("Error calling GET {}: {}", url, e.getMessage());
            throw new RuntimeException("Service unavailable: " + e.getMessage(), e);
        }
    }

    // Методы для проверки доступности сервисов
    public boolean isAuthServiceAvailable() {
        try {
            String url = authServiceUrl + "/auth/login";
            restTemplate.headForHeaders(url);
            return true;
        } catch (Exception e) {
            return false;
        }
    }
}