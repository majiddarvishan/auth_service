import React from "react";
import { Navigate } from "react-router-dom";
import { jwtDecode } from "jwt-decode";

const ProtectedRoute = ({ allowedRole, children }) => {
  const token = localStorage.getItem("token");
  if (!token) return <Navigate to="/login" />;

  const decodedToken = jwtDecode(token);
  // For admin route, allow role "admin"; for non-admin routes, you could allow any role except "admin", or you could adjust the logic.
  if (allowedRole === "admin" && decodedToken.role !== "admin") {
    return <Navigate to="/sms" />;
  }

  // You can adjust more logic here if needed.
  return children;
};

export default ProtectedRoute;
