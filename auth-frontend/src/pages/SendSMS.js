import React, { useState } from "react";
import api from '../services/api';

const SendSMS = () => {
  const [form, setForm] = useState({
    senders: "",
    receivers: "",
    text: "",
    track_ids: ""
  });

  const countLines = (value) =>
    value
      .split('\n')
      .map(v => v.trim())
      .filter(Boolean).length;

  const handleSendSMS = async (e) => {
    e.preventDefault();

    const payload = {
      senders: form.senders.split('\n').map(s => s.trim()).filter(Boolean),
      receivers: form.receivers.split('\n').map(r => r.trim()).filter(Boolean),
      text: form.text,
      track_ids: form.track_ids.split('\n').map(t => t.trim()).filter(Boolean),
    };

    try {
      await api.post("/sms/sendsms", payload, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      alert("SMS sent!");
    } catch (error) {
      alert("Failed to send SMS");
    }
  };

  return (
    <div className="container mt-5">
      <div className="card shadow p-4">
        <h3 className="mb-4 text-center">📨 Send Bulk SMS</h3>
        <form onSubmit={handleSendSMS}>
          <div className="mb-3">
            <label className="form-label fw-bold">Senders (one per line)</label>
            <textarea
              className="form-control"
              rows={3}
              placeholder="e.g.\nsender1\nsender2"
              value={form.senders}
              onChange={(e) => setForm({ ...form, senders: e.target.value })}
            />
            <small className="text-muted">{countLines(form.senders)} sender(s)</small>
          </div>

          <div className="mb-3">
            <label className="form-label fw-bold">Receivers (one per line)</label>
            <textarea
              className="form-control"
              rows={3}
              placeholder="e.g.\n+123456789\n+987654321"
              value={form.receivers}
              onChange={(e) => setForm({ ...form, receivers: e.target.value })}
            />
            <small className="text-muted">{countLines(form.receivers)} receiver(s)</small>
          </div>

          <div className="mb-3">
            <label className="form-label fw-bold">Track IDs (one per line)</label>
            <textarea
              className="form-control"
              rows={3}
              placeholder="e.g.\nid123\nid456"
              value={form.track_ids}
              onChange={(e) => setForm({ ...form, track_ids: e.target.value })}
            />
            <small className="text-muted">{countLines(form.track_ids)} track ID(s)</small>
          </div>

          <div className="mb-3">
            <label className="form-label fw-bold">Text</label>
            <textarea
              className="form-control"
              rows={3}
              placeholder="Type your message here..."
              value={form.text}
              onChange={(e) => setForm({ ...form, text: e.target.value })}
            />
          </div>

          <div className="d-grid">
            <button type="submit" className="btn btn-primary btn-lg">🚀 Send SMS</button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default SendSMS;
