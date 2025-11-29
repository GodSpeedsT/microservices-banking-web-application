package org.example.webservice.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.ViewControllerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class MvcConfig implements WebMvcConfigurer {

    @Override
    public void addViewControllers(ViewControllerRegistry registry) {
        // Маппинг URL на HTML файлы
        registry.addViewController("/").setViewName("forward:/index.html");
        registry.addViewController("/auth/login").setViewName("forward:/login.html");
        registry.addViewController("/auth/register").setViewName("forward:/register.html");
    }
}