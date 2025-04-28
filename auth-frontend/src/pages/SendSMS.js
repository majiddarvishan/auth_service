import React, { useState } from "react";
import axios from "axios";

const SendSMS = () => {
  const [form, setForm] = useState({ recipient: "", message: "" });

  const handleSendSMS = async (e) => {
    e.preventDefault();
    try {
      await axios.post("https://localhost:8443/sms/sendsms", form, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      alert("SMS sent!");
    } catch (error) {
      alert("Failed to send SMS");
    }
  };

  return (
    <div className="container mt-5">
      <h2>Send SMS</h2>
      <form onSubmit={handleSendSMS}>
        <input
          type="text"
          placeholder="Recipient"
          value={form.recipient}
          onChange={(e) => setForm({ ...form, recipient: e.target.value })}
        />
        <textarea
          placeholder="Message"
          value={form.message}
          onChange={(e) => setForm({ ...form, message: e.target.value })}
        />
        <button type="submit" className="btn btn-warning">Send SMS</button>
      </form>
    </div>
  );
};

export default SendSMS;
