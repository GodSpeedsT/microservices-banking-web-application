package org.work.depositservice.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class AccountDto {
    private Long id;
    private String accountNumber;
    private String clientId;
    private BigDecimal balance;
    private String currency;
    private LocalDateTime createdAt;
}