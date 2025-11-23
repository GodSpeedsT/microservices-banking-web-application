package org.work.depositservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import org.work.depositservice.entity.Deposit;
import java.util.List;
import java.util.Optional;

@Repository
public interface DepositRepository extends JpaRepository<Deposit, Long> {
    List<Deposit> findByAccount_AccountNumber(String accountNumber);
    List<Deposit> findByAccount_ClientId(String clientId);
    List<Deposit> findByStatus(String status);
    Optional<Deposit> findByIdAndAccount_ClientId(Long id, String clientId);
}