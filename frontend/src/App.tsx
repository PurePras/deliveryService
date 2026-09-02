import { useEffect, useState } from "react";
import { getHealth } from "./api/client";

function App() {
  const [apiStatus, setApiStatus] = useState("Checking API...");

  useEffect(() => {
    getHealth()
      .then((data) => {
        setApiStatus(`${data.service}: ${data.status}`);
      })
      .catch(() => {
        setApiStatus("API unavailable");
      });
  }, []);

  return (
    <div>
      <header>
        <h1>Shri Ram Sabji Delivery</h1>
      </header>

      <main>
        <h2>ताज़ी सब्ज़ियाँ, अब आपके घर तक</h2>

        <p>
          Fresh vegetables delivered to your doorstep.
        </p>

        <button>
          Order Fresh Vegetables
        </button>

        <p>Backend: {apiStatus}</p>
      </main>
    </div>
  );
}

export default App;