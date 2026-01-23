package org.work.depositservice.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor
@NoArgsConstructor
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
}