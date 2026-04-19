import React, { useState, useEffect } from "react";
import { jwtDecode } from "jwt-decode";
import { AuthContext } from "./AuthContext";

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const isTokenExpired = (token) => {
    try {
      const decodedToken = jwtDecode(token);
      const currentTime = Math.floor(Date.now() / 1000);

      const isExpired = decodedToken.exp <= currentTime;
      return isExpired;
    } catch (error) {
      console.error("Error checking token expiration:", error);
      return true;
    }
  };

  const getTimeUntilExpiration = (token) => {
    try {
      const decodedToken = jwtDecode(token);
      const currentTime = Date.now() / 1000;
      const timeLeft = decodedToken.exp - currentTime;

      return Math.max(0, timeLeft * 1000);
    } catch (error) {
      console.error("Error getting time until expiration:", error);
      return 0;
    }
  };

  const handleExpiredToken = () => {
    console.log("Token expired, logging out user");
    localStorage.removeItem("accessToken");
    setUser(null);
  };

  useEffect(() => {
    try {
      const token = localStorage.getItem("accessToken");

      if (token) {
        if (isTokenExpired(token)) {
          handleExpiredToken();
        } else {
          const decodedToken = jwtDecode(token);
          setUser({
            id: decodedToken.uid,
            email: decodedToken.email,
            role: decodedToken.role,
          });

          const timeUntilExpiration = getTimeUntilExpiration(token);
          if (timeUntilExpiration > 0) {
            const timeoutId = setTimeout(() => {
              handleExpiredToken();
            }, timeUntilExpiration);

            return () => {
              clearTimeout(timeoutId);
            };
          }
        }
      } else {
        console.log("No token found");
      }
    } catch (error) {
      console.error("failed to process token", error);
      localStorage.removeItem("accessToken");
    } finally {
      setLoading(false);
    }
  }, []);

  const login = (token) => {
    if (isTokenExpired(token)) {
      console.error("cannot login with expired token");
      return false;
    }

    localStorage.setItem("accessToken", token);
    const decoded = jwtDecode(token);
    setUser({ id: decoded.uid, email: decoded.email, role: decoded.role });

    const timeUntilExpiration = getTimeUntilExpiration(token);
    if (timeUntilExpiration > 0) {
      setTimeout(() => {
        handleExpiredToken();
      }, timeUntilExpiration);
    }

    return true;
  };

  const logout = () => {
    localStorage.removeItem("accessToken");
    setUser(null);
  };

  const isAuthenticated = () => {
    const token = localStorage.getItem("accessToken");
    return token && !isTokenExpired(token);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated,
        isTokenExpired: () => {
          const token = localStorage.getItem("accessToken");
          return !token || isTokenExpired(token);
        },
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
