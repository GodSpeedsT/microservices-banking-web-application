package org.work.depositservice.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.work.depositservice.entity.DepositType;
import org.work.depositservice.service.DepositTypeService;
import java.util.List;

@RestController
@RequestMapping("/api/deposit-types")
public class DepositTypeController {

    @Autowired
    private DepositTypeService depositTypeService;

    @GetMapping("/active")
    public ResponseEntity<List<DepositType>> getActiveDepositTypes() {
        List<DepositType> depositTypes = depositTypeService.getAllActiveDepositTypes();
        return ResponseEntity.ok(depositTypes);
    }

    @GetMapping
    public ResponseEntity<List<DepositType>> getAllDepositTypes() {
        List<DepositType> depositTypes = depositTypeService.getAllDepositTypes();
        return ResponseEntity.ok(depositTypes);
    }

    @GetMapping("/{id}")
    public ResponseEntity<DepositType> getDepositType(@PathVariable Long id) {
        return depositTypeService.getActiveDepositType(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public ResponseEntity<DepositType> createDepositType(@RequestBody DepositType depositType) {
        try {
            DepositType createdDepositType = depositTypeService.createDepositType(depositType);
            return ResponseEntity.ok(createdDepositType);
        } catch (IllegalArgumentException e) {
            return ResponseEntity.badRequest().build();
        }
    }

    @PutMapping("/{id}")
    public ResponseEntity<DepositType> updateDepositType(
            @PathVariable Long id,
            @RequestBody DepositType depositType) {
        try {
            DepositType updatedDepositType = depositTypeService.updateDepositType(id, depositType);
            return ResponseEntity.ok(updatedDepositType);
        } catch (RuntimeException e) {
            return ResponseEntity.notFound().build();
        }
    }

    @PostMapping("/{id}/deactivate")
    public ResponseEntity<Void> deactivateDepositType(@PathVariable Long id) {
        try {
            depositTypeService.deactivateDepositType(id);
            return ResponseEntity.ok().build();
        } catch (RuntimeException e) {
            return ResponseEntity.notFound().build();
        }
    }

    @PostMapping("/{id}/activate")
    public ResponseEntity<Void> activateDepositType(@PathVariable Long id) {
        try {
            depositTypeService.activateDepositType(id);
            return ResponseEntity.ok().build();
        } catch (RuntimeException e) {
            return ResponseEntity.notFound().build();
        }
    }
}