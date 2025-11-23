package org.work.depositservice.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.depositservice.dto.DepositRequest;
import org.work.depositservice.dto.DepositResponse;
import org.work.depositservice.entity.Deposit;
import org.work.depositservice.entity.DepositType;
import org.work.depositservice.entity.Account;
import org.work.depositservice.repository.DepositRepository;
import org.work.depositservice.repository.DepositTypeRepository;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

@Service
public class DepositService {

    @Autowired
    private DepositRepository depositRepository;

    @Autowired
    private AccountService accountService;

    @Autowired
    private DepositTypeService depositTypeService;

    @Autowired
    private InterestService interestService;

    @Autowired
    private DepositTypeRepository depositTypeRepository;

    @Transactional
    public DepositResponse createDeposit(DepositRequest request) {

        Account account = accountService.getAccountByNumber(request.getAccountNumber())
                .orElseThrow(() -> new RuntimeException("Счет не найден"));

        DepositType depositType = depositTypeService.getActiveDepositType(request.getDepositTypeId())
                .orElseThrow(() -> new RuntimeException("Тип депозита не найден или не активен"));

        if (account.getBalance().compareTo(request.getAmount()) < 0) {
            throw new RuntimeException("Недостаточно средств на счете");
        }

        account.withdraw(request.getAmount());
        accountService.updateBalance(account.getAccountNumber(), account.getBalance());

        Deposit deposit = new Deposit();
        deposit.setAccount(account);
        deposit.setDepositType(depositType);
        deposit.setAmount(request.getAmount());
        deposit.setStartDate(LocalDateTime.now());
        deposit.setEndDate(LocalDateTime.now().plusMonths(depositType.getTermMonths()));

        Deposit savedDeposit = depositRepository.save(deposit);

        return convertToResponse(savedDeposit);
    }

    public List<DepositResponse> getDepositsByClient(String clientId) {
        return depositRepository.findByAccount_ClientId(clientId)
                .stream()
                .map(this::convertToResponse)
                .collect(Collectors.toList());
    }

    public List<DepositResponse> getDepositsByAccount(String accountNumber) {
        return depositRepository.findByAccount_AccountNumber(accountNumber)
                .stream()
                .map(this::convertToResponse)
                .collect(Collectors.toList());
    }

    @Transactional
    public void closeDeposit(Long depositId, String clientId) {
        Deposit deposit = depositRepository.findByIdAndAccount_ClientId(depositId, clientId)
                .orElseThrow(() -> new RuntimeException("Депозит не найден"));

        if (!"ACTIVE".equals(deposit.getStatus())) {
            throw new RuntimeException("Депозит уже закрыт");
        }

        BigDecimal totalInterest = interestService.calculateInterest(deposit);
        deposit.setEarnedInterest(totalInterest);

        BigDecimal totalAmount = deposit.getAmount().add(totalInterest);
        accountService.updateBalance(deposit.getAccount().getAccountNumber(), totalAmount);

        deposit.setStatus("CLOSED");
        depositRepository.save(deposit);
    }

    public Optional<DepositType> getActiveDepositType(Long depositTypeId) {
        return depositTypeRepository.findByIdAndIsActiveTrue(depositTypeId);
    }

    private DepositResponse convertToResponse(Deposit deposit) {
        DepositResponse response = new DepositResponse();
        response.setId(deposit.getId());
        response.setAccountNumber(deposit.getAccount().getAccountNumber());
        response.setDepositTypeName(deposit.getDepositType().getName());
        response.setAmount(deposit.getAmount());
        response.setInterestRate(deposit.getDepositType().getInterestRate());
        response.setStartDate(deposit.getStartDate());
        response.setEndDate(deposit.getEndDate());
        response.setStatus(deposit.getStatus());
        response.setEarnedInterest(deposit.getEarnedInterest());
        return response;
    }
}