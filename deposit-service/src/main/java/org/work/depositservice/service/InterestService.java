package org.work.depositservice.service;

import org.springframework.stereotype.Service;
import org.work.depositservice.entity.Deposit;
import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDateTime;
import java.time.temporal.ChronoUnit;

@Service
public class InterestService {

    public BigDecimal calculateInterest(Deposit deposit) {
        BigDecimal principal = deposit.getAmount();
        BigDecimal annualRate = deposit.getDepositType().getInterestRate();
        long days = ChronoUnit.DAYS.between(deposit.getStartDate(), LocalDateTime.now());

        BigDecimal dailyRate = annualRate.divide(BigDecimal.valueOf(365), 10, RoundingMode.HALF_UP);
        BigDecimal interest = principal.multiply(dailyRate)
                .multiply(BigDecimal.valueOf(days))
                .setScale(2, RoundingMode.HALF_UP);

        return interest;
    }

    public BigDecimal calculateMaturityAmount(Deposit deposit) {
        BigDecimal interest = calculateInterest(deposit);
        return deposit.getAmount().add(interest);
    }
}