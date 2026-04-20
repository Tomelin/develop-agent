import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios';

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
<<<<<<< feat/phase-02-frontend-17026144929788359576
  withCredentials: true, // Importante para enviar os cookies (refresh token)
=======
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
>>>>>>> main
});

let isRefreshing = false;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let failedQueue: { resolve: (value?: unknown) => void; reject: (reason?: any) => void }[] = [];

const processQueue = (error: AxiosError | null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });

  failedQueue = [];
};

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Apenas rodar no lado do cliente
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('@agency:token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise(function (resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            if (originalRequest.headers) {
               originalRequest.headers.Authorization = 'Bearer ' + token;
            }
            return api(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
<<<<<<< feat/phase-02-frontend-17026144929788359576
        const { data } = await axios.post(
          `${api.defaults.baseURL}/auth/refresh`,
          {},
          { withCredentials: true } // envia o refresh token via cookie
        );

        const newAccessToken = data.access_token;
        localStorage.setItem('@agency:token', newAccessToken);

        if (originalRequest.headers) {
           originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
=======
        const refreshResponse = await api.post('/auth/refresh');
        const newAccessToken = refreshResponse.data?.access_token;

        if (newAccessToken && typeof window !== 'undefined') {
          localStorage.setItem('access_token', newAccessToken);
          originalRequest.headers = originalRequest.headers ?? {};
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return api(originalRequest);
>>>>>>> main
        }
        processQueue(null, newAccessToken);

        return api(originalRequest);
      } catch (err) {
        processQueue(err as AxiosError, null);
        localStorage.removeItem('@agency:token');

        // Redirecionar para login apenas no client-side
        if (typeof window !== 'undefined') {
<<<<<<< feat/phase-02-frontend-17026144929788359576
            window.location.href = '/login';
=======
          localStorage.removeItem('access_token');
          window.location.href = '/login';
>>>>>>> main
        }
        return Promise.reject(err);
      } finally {
        isRefreshing = false;
      }

      if (typeof window !== 'undefined') {
        localStorage.removeItem('access_token');
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);
