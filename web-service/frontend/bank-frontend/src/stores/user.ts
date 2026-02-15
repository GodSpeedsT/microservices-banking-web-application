import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import api from '@/services/api';
import type { UserProfile } from '@/types';

export const useUserStore = defineStore('user', () => {
    const user = ref<UserProfile | null>(null);
    const loading = ref(false);

    const isAuthenticated = computed(() => !!user.value);
    const username = computed(() => user.value?.username || 'Клиент');
    const isAdmin = computed(() => user.value?.roles.includes('ROLE_ADMIN') || false);

    async function fetchUser() {
        if (user.value) return;

        loading.value = true;
        try {

            const response = await api.get('/user');
            user.value = {
                username: response.data.preferred_username || response.data.name || 'User',
                email: response.data.email,
                roles: response.data.roles || []
            };
        } catch (error) {
            user.value = null
           // console.error('Failed to fetch user profile', error);
        } finally {
            loading.value = false;
        }
    }

    function clearUser() {
        user.value = null;
    }

    return { user, loading, isAuthenticated, username, isAdmin, fetchUser, clearUser };
});