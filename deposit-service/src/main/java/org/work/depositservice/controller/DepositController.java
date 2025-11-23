package org.work.depositservice.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.depositservice.dto.DepositRequest;
import org.work.depositservice.dto.DepositResponse;
import org.work.depositservice.service.DepositService;
import java.util.List;

@RestController
@RequestMapping("/api/deposits")
public class DepositController {

    @Autowired
    private DepositService depositService;

    @PostMapping
    public ResponseEntity<DepositResponse> createDeposit(@RequestBody DepositRequest request) {
        try {
            DepositResponse response = depositService.createDeposit(request);
            return ResponseEntity.ok(response);
        } catch (RuntimeException e) {
            return ResponseEntity.badRequest().build();
        }
    }

    @GetMapping("/client/{clientId}")
    public ResponseEntity<List<DepositResponse>> getClientDeposits(@PathVariable String clientId) {
        List<DepositResponse> deposits = depositService.getDepositsByClient(clientId);
        return ResponseEntity.ok(deposits);
    }

    @GetMapping("/account/{accountNumber}")
    public ResponseEntity<List<DepositResponse>> getAccountDeposits(@PathVariable String accountNumber) {
        List<DepositResponse> deposits = depositService.getDepositsByAccount(accountNumber);
        return ResponseEntity.ok(deposits);
    }

    @PostMapping("/{depositId}/close")
    public ResponseEntity<Void> closeDeposit(
            @PathVariable Long depositId,
            @RequestParam String clientId) {
        try {
            depositService.closeDeposit(depositId, clientId);
            return ResponseEntity.ok().build();
        } catch (RuntimeException e) {
            return ResponseEntity.badRequest().build();
        }
    }
}