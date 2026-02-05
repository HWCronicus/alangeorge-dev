import { useState, useEffect } from "react";
import ThreeScene from "../components/ThreeScene";
import Terminal from "../components/Terminal";
import ContactModal from "../components/ContactModal";
import "../styles/home.css";

export default function Home() {
  const [isTerminalOpen, setIsTerminalOpen] = useState(false);
  const [isContactOpen, setIsContactOpen] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText("ssh alangeorge.dev");
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setIsTerminalOpen(false);
        setIsContactOpen(false);
      }
    };
    window.addEventListener("keydown", handleEscape);
    return () => window.removeEventListener("keydown", handleEscape);
  }, []);

  return (
    <>
      <ThreeScene />
      <div className="about">
        <p>
          This is just a landing page for my personal website,{" "}
          <strong>AlanGeorge.Dev</strong>. I'm Alan George, a software developer
          specializing in Go and web development.
        </p>
        <p>
          The meat and potatoes of this project is an interactive SSH resume
          built with Go using BubbleTea and Wish. Available at:{" "}
        </p>
        <code className="codeBlock">
          ssh alangeorge.dev
          <button
            className={`copyButton ${copied ? "copied" : ""}`}
            onClick={handleCopy}
            // title={copied ? "Copied!" : "Copy to clipboard"}
            data-tooltip={copied ? "Copied!" : "Copy to clipboard"}
          >
            {copied ? "✓" : "📋"}
          </button>
        </code>
        <p>
          If you don't want to SSH into the terminal yourself, you can click the
          "Terminal Example" button below to try it out.
        </p>
      </div>
      <div className="navButtons">
        <button
          className="navButton"
          data-tooltip="Check out my projects"
          onClick={() => window.open("https://github.com/HWCronicus", "_blank")}
        >
          GitHub
        </button>
        <button
          className="navButton"
          data-tooltip="Get in touch with me"
          onClick={() => setIsContactOpen(true)}
        >
          Email
        </button>
        <button
          className="navButton"
          data-tooltip="Try my interactive SSH resume built with Go"
          onClick={() => setIsTerminalOpen(true)}
        >
          Terminal Example
        </button>
      </div>

      {isTerminalOpen && (
        <div className="modal-overlay" onClick={() => setIsTerminalOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button
              className="modal-close"
              onClick={() => setIsTerminalOpen(false)}
            >
              ✕
            </button>
            <Terminal />
          </div>
        </div>
      )}

      {isContactOpen && (
        <ContactModal onClose={() => setIsContactOpen(false)} />
      )}
    </>
  );
}
