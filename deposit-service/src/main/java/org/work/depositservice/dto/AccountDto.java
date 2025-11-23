package org.work.depositservice.dto;

import java.math.BigDecimal;
import java.time.LocalDateTime;

public class AccountDto {
    private Long id;
    private String accountNumber;
    private String clientId;
    private BigDecimal balance;
    private String currency;
    private LocalDateTime createdAt;

    public AccountDto(Long id, String accountNumber, String clientId, BigDecimal balance, String currency, LocalDateTime createdAt) {
        this.id = id;
        this.accountNumber = accountNumber;
        this.clientId = clientId;
        this.balance = balance;
        this.currency = currency;
        this.createdAt = createdAt;
    }

    public AccountDto() {
    }

    public Long getId() {
        return id;
    }

    public String getAccountNumber() {
        return accountNumber;
    }

    public BigDecimal getBalance() {
        return balance;
    }

    public String getClientId() {
        return clientId;
    }

    public String getCurrency() {
        return currency;
    }

    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public void setAccountNumber(String accountNumber) {
        this.accountNumber = accountNumber;
    }

    public void setClientId(String clientId) {
        this.clientId = clientId;
    }

    public void setBalance(BigDecimal balance) {
        this.balance = balance;
    }

    public void setCurrency(String currency) {
        this.currency = currency;
    }

    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }
}