import { ApiResponse, CaptchaData, LoginRequest, LoginResponse } from '../types';
import { httpClient } from './http';

export const authService = {
  /**
   * Fetch captcha image and ID
   * @returns Captcha data including id, base64 image, and length
   */
  async fetchCaptcha(): Promise<ApiResponse<CaptchaData>> {
    return httpClient.get<CaptchaData>('/v1/system/captcha');
  },

  /**
   * Login with username, password, and captcha
   * @param credentials - Login credentials including username, password, captcha, and captcha_id
   * @returns Login response with token and username
   */
  async login(credentials: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return httpClient.post<LoginResponse>('/v1/system/login', credentials);
  },

  /**
   * Logout from the system
   * @returns Success response
   */
  async logout(): Promise<ApiResponse<void>> {
    return httpClient.get<void>('/v1/system/logout');
  },
};
