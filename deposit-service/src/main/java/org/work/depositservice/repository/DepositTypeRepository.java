package org.work.depositservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import org.work.depositservice.entity.DepositType;
import java.util.List;
import java.util.Optional;

@Repository
public interface DepositTypeRepository extends JpaRepository<DepositType, Long> {
    List<DepositType> findByIsActiveTrue();
    Optional<DepositType> findByIdAndIsActiveTrue(Long id);
}