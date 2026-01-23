package org.work.depositservice.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.math.BigDecimal;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class DepositRequest {
    private String accountNumber;
    private Long depositTypeId;
    private BigDecimal amount;
}