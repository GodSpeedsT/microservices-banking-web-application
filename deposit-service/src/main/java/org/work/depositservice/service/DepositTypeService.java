package org.work.depositservice.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;
import org.work.depositservice.entity.DepositType;
import org.work.depositservice.repository.DepositTypeRepository;

import java.util.List;
import java.util.Optional;

@Service
public class DepositTypeService {

    @Autowired
    private DepositTypeRepository depositTypeRepository;

    @Cacheable("depositTypes")
    public List<DepositType> getAllActiveDepositTypes() {
        return depositTypeRepository.findByIsActiveTrue();
    }

    public Optional<DepositType> getActiveDepositType(Long id) {
        return depositTypeRepository.findByIdAndIsActiveTrue(id);
    }

    public List<DepositType> getAllDepositTypes() {
        return depositTypeRepository.findAll();
    }

    public DepositType createDepositType(DepositType depositType) {
        if (depositType.getInterestRate().compareTo(java.math.BigDecimal.ZERO) <= 0) {
            throw new IllegalArgumentException("Процентная ставка должна быть положительной");
        }

        if (depositType.getTermMonths() <= 0) {
            throw new IllegalArgumentException("Срок депозита должен быть положительным");
        }

        return depositTypeRepository.save(depositType);
    }

    public DepositType updateDepositType(Long id, DepositType updatedDepositType) {
        DepositType existingDepositType = depositTypeRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Тип депозита не найден"));

        existingDepositType.setName(updatedDepositType.getName());
        existingDepositType.setInterestRate(updatedDepositType.getInterestRate());
        existingDepositType.setTermMonths(updatedDepositType.getTermMonths());
        existingDepositType.setDescription(updatedDepositType.getDescription());
        existingDepositType.setIsActive(updatedDepositType.getIsActive());

        return depositTypeRepository.save(existingDepositType);
    }

    public void deactivateDepositType(Long id) {
        DepositType depositType = depositTypeRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Тип депозита не найден"));

        depositType.setIsActive(false);
        depositTypeRepository.save(depositType);
    }

    public void activateDepositType(Long id) {
        DepositType depositType = depositTypeRepository.findById(id)
                .orElseThrow(() -> new RuntimeException("Тип депозита не найден"));

        depositType.setIsActive(true);
        depositTypeRepository.save(depositType);
    }

    public boolean isDepositTypeAvailable(Long id) {
        return depositTypeRepository.findByIdAndIsActiveTrue(id).isPresent();
    }
}