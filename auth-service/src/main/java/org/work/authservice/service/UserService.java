package org.work.authservice.service;

import lombok.RequiredArgsConstructor;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.work.authservice.entity.Role;
import org.work.authservice.entity.User;
import org.work.authservice.repository.UserRepository;

import java.util.Optional;
import java.util.Set;

@Service
@RequiredArgsConstructor
public class UserService {
    private final UserRepository userRepository;
    private final PasswordEncoder passwordEncoder; // Используем PasswordEncoder напрямую
    private final RoleService roleService;
    private final ValidationService validationService;

    @Transactional
    public User registerUser(String username, String rawPassword, Set<String> roleNames) {
        validationService.validateUsername(username);
        validationService.validatePassword(rawPassword);

        if (userRepository.findByUsername(username).isPresent()) {
            throw new RuntimeException("Username is already in use");
        }

        Set<Role> roles = roleNames != null && !roleNames.isEmpty()
                ? roleService.getRolesByNames(roleNames)
                : roleService.getUserRoles();

        User user = new User();
        user.setUsername(username);
        user.setPassword(passwordEncoder.encode(rawPassword)); // Используем напрямую
        user.setRoles(roles);

        return userRepository.save(user);
    }

    @Transactional
    public User registerUser(String username, String rawPassword) {
        return registerUser(username, rawPassword, null);
    }

    public Optional<User> findByUsername(String username) {
        return userRepository.findByUsername(username);
    }

    // УДАЛЯЕМ этот метод, чтобы разорвать цикл
    // public PasswordEncoder getPasswordEncoder() {
    //     return passwordEncoder;
    // }
}