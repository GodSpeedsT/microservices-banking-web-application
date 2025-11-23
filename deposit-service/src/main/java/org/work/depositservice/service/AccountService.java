package org.work.depositservice.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.depositservice.dto.AccountDto;
import org.work.depositservice.entity.Account;
import org.work.depositservice.repository.AccountRepository;
import java.math.BigDecimal;
import java.util.Optional;

@Service
public class AccountService {

    @Autowired
    private AccountRepository accountRepository;

    @Transactional
    public Account createAccount(String clientId, String currency) {
        String accountNumber = generateAccountNumber();

        Account account = new Account();
        account.setAccountNumber(accountNumber);
        account.setClientId(clientId);
        account.setCurrency(currency);
        account.setBalance(BigDecimal.ZERO);

        return accountRepository.save(account);
    }

    public Optional<Account> getAccountByNumber(String accountNumber) {
        return accountRepository.findByAccountNumber(accountNumber);
    }

    public Optional<Account> getAccountByClientId(String clientId) {
        return accountRepository.findByClientId(clientId);
    }


    @Transactional
    public void updateBalance(String accountNumber, BigDecimal amount) {
        Account account = accountRepository.findByAccountNumber(accountNumber)
                .orElseThrow(() -> new RuntimeException("Счет не найден"));

        account.deposit(amount);
        accountRepository.save(account);
    }

    private String generateAccountNumber() {
        return "ACC" + System.currentTimeMillis();
    }

    public AccountDto convertToDto(Account account) {
        AccountDto dto = new AccountDto();
        dto.setId(account.getId());
        dto.setAccountNumber(account.getAccountNumber());
        dto.setClientId(account.getClientId());
        dto.setBalance(account.getBalance());
        dto.setCurrency(account.getCurrency());
        dto.setCreatedAt(account.getCreatedAt());
        return dto;
    }
}