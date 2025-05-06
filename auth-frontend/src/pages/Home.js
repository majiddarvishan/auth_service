import React from "react";
import { Link } from "react-router-dom";

const Home = () => {
  return (
    <div className="container mt-5 d-flex justify-content-center">
      <div className="card shadow-lg p-5 w-100" style={{ maxWidth: "600px" }}>
        <h1 className="text-center mb-3">🔐 Welcome</h1>
        <p className="text-center text-muted fs-5">
          Manage users and send SMS securely with ease.
        </p>

        <div className="d-grid gap-3 mt-4">
          <Link to="/login" className="btn btn-primary btn-lg">
            🔑 Login
          </Link>
          <Link to="/register" className="btn btn-success btn-lg">
            🧾 Register
          </Link>
          <Link to="/sms" className="btn btn-warning btn-lg text-white">
            ✉️ Send SMS
          </Link>
          <Link to="/admin" className="btn btn-danger btn-lg">
            🛠️ Admin Dashboard
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
