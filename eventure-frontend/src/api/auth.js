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
      tokens: {
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
      },
    };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || "Login failed",
    };
  }
}

export async function fetchProfile(accessToken) {
  try {
    const response = await axios.get("http://localhost:8080/auth/profile", {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return { success: true, data: response.data };
  } catch (err) {
    return {
      success: false,
      error: err.response?.data?.error || "Error fetching profile",
    };
  }
}

export async function fetchAdminUsers(accessToken) {
  try {
    const response = await axios.get("http://localhost:8080/admin/users", {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });
    return { success: true, data: response.data };
  } catch (err) {
    return {
      success: false,
      error: err.response?.data?.error || "Error fetching users",
    };
  }
}
