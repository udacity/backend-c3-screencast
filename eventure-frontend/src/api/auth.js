import axios from "axios";

const API_BASE = "http://localhost:8080/auth";

export async function register(username, password) {
  try {
    const response = await axios.post(`${API_BASE}/register`, {
      username,
      password,
    });

    return {
      success: true,
      message: response.data.message,
    };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || "Registration failed",
    };
  }
}

export async function login(username, password) {
  try {
    const response = await axios.post(`${API_BASE}/login`, {
      username,
      password,
    });

    return {
      success: true,
      accessToken: response.data.access_token,
      refreshToken: response.data.refresh_token,
    };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || "Login failed",
    };
  }
}
