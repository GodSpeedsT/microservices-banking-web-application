package org.work.depositservice.dto;

import java.math.BigDecimal;

public class DepositRequest {
    private String accountNumber;
    private Long depositTypeId;
    private BigDecimal amount;

    // Конструкторы, геттеры, сеттеры
    public DepositRequest() {}

    public DepositRequest(String accountNumber, Long depositTypeId, BigDecimal amount) {
        this.accountNumber = accountNumber;
        this.depositTypeId = depositTypeId;
        this.amount = amount;
    }

    public String getAccountNumber() {
        return accountNumber;
    }

    public Long getDepositTypeId() {
        return depositTypeId;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public void setAccountNumber(String accountNumber) {
        this.accountNumber = accountNumber;
    }

    public void setDepositTypeId(Long depositTypeId) {
        this.depositTypeId = depositTypeId;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }
}