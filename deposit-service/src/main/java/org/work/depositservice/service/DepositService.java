package org.work.depositservice.service;

import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.depositservice.dto.DepositRequest;
import org.work.depositservice.dto.DepositResponse;
import org.work.depositservice.entity.Deposit;
import org.work.depositservice.entity.DepositStatus;
import org.work.depositservice.entity.DepositType;
import org.work.depositservice.entity.Account;
import org.work.depositservice.handler.InsufficientFundsException;
import org.work.depositservice.repository.DepositRepository;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
public class DepositService {

    private final DepositRepository depositRepository;
    private final AccountService accountService;
    private final DepositTypeService depositTypeService;
    private final InterestService interestService;

    @Transactional
    public DepositResponse createDeposit(DepositRequest request) {
        Account account = accountService.getAccountByNumber(request.getAccountNumber())
                .orElseThrow(() -> new EntityNotFoundException("Счет не найден"));

        DepositType depositType = depositTypeService.getCurrentDepositType(request.getDepositTypeId())
                .orElseThrow(() -> new IllegalStateException("Тип депозита недоступен"));

        if (account.getBalance().compareTo(request.getAmount()) < 0) {
            throw new InsufficientFundsException("Недостаточно средств для открытия вклада");
        }

        account.withdraw(request.getAmount());

        Deposit deposit = new Deposit();
        deposit.setAccount(account);
        deposit.setDepositType(depositType);
        deposit.setAmount(request.getAmount());
        deposit.setEndDate(LocalDateTime.now().plusMonths(depositType.getTermMonths()));

        return convertToResponse(depositRepository.save(deposit));
    }

    @Transactional
    public void closeDeposit(Long depositId, String clientId) {
        Deposit deposit = depositRepository.findByIdAndAccount_ClientId(depositId, clientId)
                .orElseThrow(() -> new EntityNotFoundException("Депозит не найден или принадлежит не вам"));

        if (DepositStatus.ACTIVE != deposit.getStatus()) {
            throw new IllegalStateException("Депозит уже закрыт");
        }

        completeClosing(deposit);
    }

    @Transactional
    public void closeDepositByAdmin(Long depositId) {
        Deposit deposit = depositRepository.findById(depositId)
                .orElseThrow(() -> new EntityNotFoundException("Депозит не найден"));
        completeClosing(deposit);
    }

    private void completeClosing(Deposit deposit) {
        BigDecimal interest = interestService.calculateInterest(deposit);
        deposit.setEarnedInterest(interest);
        deposit.setStatus(DepositStatus.CLOSED);

        deposit.getAccount().deposit(deposit.getAmount().add(interest));
        depositRepository.save(deposit);
    }

    public List<DepositResponse> getDepositsByClient(String currentUserId) {
        return depositRepository.findByAccount_ClientId(currentUserId)
                .stream().map(
                        this::convertToResponse
                ).collect(Collectors.toList());
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