import { useState } from "react";
import toast from "react-hot-toast";
import "../styles/contact-modal.css";

interface ContactModalProps {
  onClose: () => void;
}

interface FormData {
  name: string;
  email: string;
  phone: string;
  subject: string;
  message: string;
}

const initialFormData: FormData = {
  name: "",
  email: "",
  phone: "",
  subject: "Introduction",
  message: "",
};

export default function ContactModal({ onClose }: ContactModalProps) {
  const [formData, setFormData] = useState<FormData>(initialFormData);
  const [isSending, setIsSending] = useState(false);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleClear = () => {
    setFormData(initialFormData);
    toast.success("Form cleared");
  };

  const handleUseAI = () => {
    // Placeholder for AI feature
    toast.success("AI feature coming soon", {
      icon: "🤖",
    });
  };

  const handleSend = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!formData.name || !formData.email || !formData.message) {
      toast.error("Please fill in all required fields");
      return;
    }

    setIsSending(true);
    const sendingToast = toast.loading("Sending your message...");

    try {
      const response = await fetch("/email", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || "Failed to send email");
      }

      toast.success("Message sent successfully!", { id: sendingToast });
      setTimeout(() => {
        handleClear();
        onClose();
      }, 1500);
    } catch (error) {
      const errorMsg =
        error instanceof Error ? error.message : "Failed to send email";
      toast.error(errorMsg, { id: sendingToast });
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="contact-modal" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onClose}>
          ✕
        </button>

        <h2>Get In Touch</h2>

        <form onSubmit={handleSend}>
          <div className="form-group">
            <label htmlFor="name">Name *</label>
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Your name"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="email">Email *</label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="your@email.com"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="phone">Phone Number</label>
            <input
              type="tel"
              id="phone"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              placeholder="(123) 456-7890"
            />
          </div>

          <div className="form-group">
            <label htmlFor="subject">Subject *</label>
            <select
              id="subject"
              name="subject"
              value={formData.subject}
              onChange={handleChange}
              required
            >
              <option value="Introduction">Introduction</option>
              <option value="Hire Me">Hire Me</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="message">Message *</label>
            <textarea
              id="message"
              name="message"
              value={formData.message}
              onChange={handleChange}
              placeholder="Your message..."
              rows={6}
              required
            />
          </div>

          <div className="form-buttons">
            <button type="button" className="btn-clear" onClick={handleClear}>
              Clear
            </button>
            <button type="button" className="btn-ai" onClick={handleUseAI}>
              Use AI
            </button>
            <button type="submit" className="btn-send" disabled={isSending}>
              {isSending ? "Sending..." : "Send"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
