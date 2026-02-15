import axios from 'axios';

const api = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
            if (!window.location.pathname.startsWith('/login')) {
                window.location.href = '/oauth2/authorization/gateway';
            }
        }
        return Promise.reject(error);
    }
);

export default api;