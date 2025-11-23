package org.work.depositservice.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.depositservice.dto.AccountDto;
import org.work.depositservice.service.AccountService;
import java.util.Optional;

@RestController
@RequestMapping("/api/accounts")
public class AccountController {

    @Autowired
    private AccountService accountService;

    @PostMapping
    public ResponseEntity<AccountDto> createAccount(
            @RequestParam String clientId,
            @RequestParam String currency) {
        return ResponseEntity.ok(
                accountService.convertToDto(
                        accountService.createAccount(clientId, currency)
                )
        );
    }

    @GetMapping("/{accountNumber}")
    public ResponseEntity<AccountDto> getAccount(@PathVariable String accountNumber) {
        Optional<AccountDto> account = accountService.getAccountByNumber(accountNumber)
                .map(accountService::convertToDto);

        return account.map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @GetMapping("/client/{clientId}")
    public ResponseEntity<AccountDto> getAccountByClient(@PathVariable String clientId) {
        Optional<AccountDto> account = accountService.getAccountByClientId(clientId)
                .map(accountService::convertToDto);

        return account.map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }
}