import { Routes, Route } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import Home from "./pages/Home";

function App() {
  return (
    <>
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 4000,
          style: {
            background: "#1a1a1a",
            color: "#ff8800",
            border: "1px solid #ff8800",
            fontFamily: "'Fira Code', monospace",
          },
          success: {
            iconTheme: {
              primary: "#ff8800",
              secondary: "#1a1a1a",
            },
          },
          error: {
            iconTheme: {
              primary: "#ff3333",
              secondary: "#1a1a1a",
            },
            style: {
              border: "1px solid #ff3333",
            },
          },
        }}
      />
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </>
  );
}

export default App;
