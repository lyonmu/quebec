// API Response Types
export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

export interface CaptchaData {
  id: string;
  pictures: string; // base64 image
  length: number;
}
