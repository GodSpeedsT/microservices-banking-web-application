package org.example.apigateway.controllers;

import org.example.apigateway.dto.UserDto;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Mono;

@RestController
public class UserController {

@GetMapping("/api/users/me")
public Mono<UserDto> getUser(@AuthenticationPrincipal OidcUser oidcUser) {

    if (oidcUser == null) {
        return Mono.empty();
    }
    return Mono.just(new UserDto(
            oidcUser.getPreferredUsername(),
            oidcUser.getGivenName(),
            oidcUser.getFamilyName(),
            oidcUser.getEmail()
    ));
}
}
