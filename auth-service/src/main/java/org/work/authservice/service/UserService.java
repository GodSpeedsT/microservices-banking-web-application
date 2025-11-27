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
    private final PasswordEncoder passwordEncoder;
    private final RoleService roleService;
    private final ValidationService validationService;

    /**
     * Регистрирует нового пользователя с возможностью указания ролей.
     */
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
        // Хеширование пароля перед сохранением
        user.setPassword(passwordEncoder.encode(rawPassword));
        user.setRoles(roles);

        return userRepository.save(user);
    }

    /**
     * Регистрирует нового пользователя с дефолтной ролью.
     */
    @Transactional
    public User registerUser(String username, String rawPassword) {
        return registerUser(username, rawPassword, null);
    }

    /**
     * Ищет пользователя по имени, возвращает Optional.
     * Используется в CustomUserDetailsService.
     */
    public Optional<User> findByUsername(String username) {
        return userRepository.findByUsername(username);
    }

    /**
     * Ищет пользователя по имени и возвращает объект User.
     * Выбрасывает RuntimeException, если пользователь не найден.
     * Удобен для использования в контроллерах (например, UserController), где токен уже проверен.
     */
    public User findByUsernameOrFail(String username) {
        return userRepository.findByUsername(username)
                .orElseThrow(() -> new RuntimeException("User not found: " + username));
    }
}