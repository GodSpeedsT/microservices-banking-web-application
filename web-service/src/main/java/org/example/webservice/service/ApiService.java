package org.example.webservice.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class ApiService {

    private final RestTemplate restTemplate;

    @Value("${services.auth.url:http://auth-service:8081}")
    private String authServiceUrl;

    @Value("${services.deposit.url:http://deposit-service:8082}")
    private String depositServiceUrl;

    public ApiService(RestTemplate restTemplate) {
        this.restTemplate = restTemplate;
    }

    public <T> T post(String url, Object request, Class<T> responseType, String token) {
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        if (token != null) {
            headers.setBearerAuth(token);
        }

        HttpEntity<Object> entity = new HttpEntity<>(request, headers);
        ResponseEntity<T> response = restTemplate.exchange(url, HttpMethod.POST, entity, responseType);
        return response.getBody();
    }

    public <T> T get(String url, Class<T> responseType, String token) {
        HttpHeaders headers = new HttpHeaders();
        if (token != null) {
            headers.setBearerAuth(token);
        }

        HttpEntity<String> entity = new HttpEntity<>(headers);
        ResponseEntity<T> response = restTemplate.exchange(url, HttpMethod.GET, entity, responseType);
        return response.getBody();
    }

    public Object registerUser(Object userRequest, String token) {
        return post(authServiceUrl + "/auth/auth/register", userRequest, Object.class, token);
    }

    public Object loginUser(Object loginRequest) {
        return post(authServiceUrl + "/auth/auth/login", loginRequest, Object.class, null);
    }
}
