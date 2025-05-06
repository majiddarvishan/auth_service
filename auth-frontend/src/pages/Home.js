import React from "react";
import { Link } from "react-router-dom";

const Home = () => {
  return (
    <div className="container mt-5 text-center">
      <h1>Welcome to the Authentication & Accounting System</h1>
      <p>Manage users, track balances, and send SMS securely.</p>

      <div className="mt-4">
        <Link to="/login" className="btn btn-primary mx-2">Login</Link>
        <Link to="/register" className="btn btn-success mx-2">Register</Link>
        <Link to="/sms" className="btn btn-warning mx-2">Send SMS</Link>
        <Link to="/dashboard" className="btn btn-danger mx-2">Admin Dashboard</Link>
      </div>
    </div>
  );
};

export default Home;
