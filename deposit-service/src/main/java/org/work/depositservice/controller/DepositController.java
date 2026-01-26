package org.work.depositservice.controller;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;
import org.work.depositservice.dto.DepositRequest;
import org.work.depositservice.dto.DepositResponse;
import org.work.depositservice.service.DepositService;
import org.work.depositservice.service.SecurityContextService;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/deposits")
@RequiredArgsConstructor
public class DepositController {

    private final DepositService depositService;
    private final SecurityContextService securityContext;

    @GetMapping("/my")
    @PreAuthorize("hasAuthority('ROLE_USER')")
    public ResponseEntity<List<DepositResponse>> getMyDeposits() {
        // Берем ID прямо из токена, пользователь не может подсмотреть чужие вклады
        String currentUserId = securityContext.getCurrentUserId();
        return ResponseEntity.ok(depositService.getDepositsByClient(currentUserId));
    }

    @GetMapping("/debug-token")
    public Map<String, Object> debug(@AuthenticationPrincipal Jwt jwt) {
        return Map.of(
                "sub_claim", jwt.getSubject(), // То, что идет в clientId
                "all_claims", jwt.getClaims()   // Все данные из токена
        );
    }

    @PostMapping
    @PreAuthorize("hasAuthority('ROLE_USER')")
    public ResponseEntity<DepositResponse> createDeposit(@RequestBody DepositRequest request) {
        return ResponseEntity.ok(depositService.createDeposit(request));
    }

    @PostMapping("/{depositId}/close")
    @PreAuthorize("hasAuthority('ROLE_USER')")
    public ResponseEntity<Void> closeMyDeposit(@PathVariable Long depositId) {
        String currentUserId = securityContext.getCurrentUserId();
        depositService.closeDeposit(depositId, currentUserId);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{depositId}")
    @PreAuthorize("hasAuthority('ROLE_ADMIN')")
    public ResponseEntity<Void> forceCloseDeposit(@PathVariable Long depositId) {
        // Админ закрывает любой вклад (логику можно расширить в сервисе)
        depositService.closeDepositByAdmin(depositId);
        return ResponseEntity.ok().build();
    }
}