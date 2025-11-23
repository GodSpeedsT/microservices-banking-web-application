package org.work.depositservice.dto;

import java.math.BigDecimal;
import java.time.LocalDateTime;

public class DepositResponse {
    private Long id;
    private String accountNumber;
    private String depositTypeName;
    private BigDecimal amount;
    private BigDecimal interestRate;
    private LocalDateTime startDate;
    private LocalDateTime endDate;
    private String status;
    private BigDecimal earnedInterest;

    // Конструкторы, геттеры, сеттеры


    public Long getId() {
        return id;
    }

    public String getAccountNumber() {
        return accountNumber;
    }

    public String getDepositTypeName() {
        return depositTypeName;
    }

    public BigDecimal getInterestRate() {
        return interestRate;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public LocalDateTime getStartDate() {
        return startDate;
    }

    public LocalDateTime getEndDate() {
        return endDate;
    }

    public String getStatus() {
        return status;
    }

    public BigDecimal getEarnedInterest() {
        return earnedInterest;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public void setAccountNumber(String accountNumber) {
        this.accountNumber = accountNumber;
    }

    public void setDepositTypeName(String depositTypeName) {
        this.depositTypeName = depositTypeName;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    public void setInterestRate(BigDecimal interestRate) {
        this.interestRate = interestRate;
    }

    public void setStartDate(LocalDateTime startDate) {
        this.startDate = startDate;
    }

    public void setEndDate(LocalDateTime endDate) {
        this.endDate = endDate;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public void setEarnedInterest(BigDecimal earnedInterest) {
        this.earnedInterest = earnedInterest;
    }
}