package org.example.webservice.service;

import lombok.RequiredArgsConstructor;
import org.example.webservice.dto.AuthRequest;
import org.example.webservice.dto.AuthResponse;
import org.example.webservice.dto.UserResponse;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
@RequiredArgsConstructor
public class ApiService {

    private final RestTemplate restTemplate;

    private static final String AUTH_SERVICE_URL = "http://localhost:8081/auth";

    public UserResponse registerUser(AuthRequest authRequest) {
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        HttpEntity<AuthRequest> entity = new HttpEntity<>(authRequest, headers);

        return restTemplate.postForObject(
                AUTH_SERVICE_URL + "/register",
                entity,
                UserResponse.class
        );
    }

}
