package org.example.webservice.service;

import jakarta.servlet.http.HttpSession;
import org.springframework.stereotype.Service;

@Service
public class AuthService {

    private static final String TOKEN_KEY = "auth_token";
    private static final String USERNAME_KEY = "username";

    public void loginUser(HttpSession session, String token, String username) {
        session.setAttribute(TOKEN_KEY, token);
        session.setAttribute(USERNAME_KEY, username);
    }

    public void logoutUser(HttpSession session) {
        session.removeAttribute(TOKEN_KEY);
        session.removeAttribute(USERNAME_KEY);
        session.invalidate();
    }

    public boolean isAuthenticated(HttpSession session) {
        return session.getAttribute(TOKEN_KEY) != null;
    }

    public String getToken(HttpSession session) {
        return (String) session.getAttribute(TOKEN_KEY);
    }

    public String getUsername(HttpSession session) {
        return (String) session.getAttribute(USERNAME_KEY);
    }
}