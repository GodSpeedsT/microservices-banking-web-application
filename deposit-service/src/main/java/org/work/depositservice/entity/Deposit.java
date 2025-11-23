package org.work.depositservice.entity;

import jakarta.persistence.*;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "deposits")
public class Deposit {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "account_id", nullable = false)
    private Account account;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "deposit_type_id", nullable = false)
    private DepositType depositType;

    @Column(nullable = false)
    private BigDecimal amount;

    @Column(nullable = false)
    private LocalDateTime startDate;

    @Column(nullable = false)
    private LocalDateTime endDate;

    @Column(nullable = false)
    private String status; // ACTIVE, CLOSED, MATURED

    private BigDecimal earnedInterest;

    // Конструкторы, геттеры, сеттеры
    public Deposit() {
        this.startDate = LocalDateTime.now();
        this.status = "ACTIVE";
        this.earnedInterest = BigDecimal.ZERO;
    }

    public Long getId() {
        return id;
    }

    public Account getAccount() {
        return account;
    }

    public DepositType getDepositType() {
        return depositType;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public LocalDateTime getEndDate() {
        return endDate;
    }

    public LocalDateTime getStartDate() {
        return startDate;
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

    public void setAccount(Account account) {
        this.account = account;
    }

    public void setDepositType(DepositType depositType) {
        this.depositType = depositType;
    }

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
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