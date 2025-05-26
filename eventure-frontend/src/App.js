import { useState } from "react";
import RegisterForm from "./components/RegisterForm";
import LoginForm from "./components/LoginForm";
import TokenDisplay from "./components/TokenDisplay";
import ProtectedActions from "./components/ProtectedActions";
import "./App.css";

function App() {
  const [accessToken, setAccessToken] = useState("");
  const [refreshToken, setRefreshToken] = useState("");

  return (
    <div className="container">
      <h1>Eventure</h1>

      <section style={{ marginBottom: "2rem" }}>
        <RegisterForm />
      </section>

      <section style={{ marginBottom: "2rem" }}>
        <LoginForm
          onLogin={({ accessToken = "", refreshToken = "" } = {}) => {
            setAccessToken(accessToken);
            setRefreshToken(refreshToken);
          }}
        />
      </section>

      <section>
        <TokenDisplay accessToken={accessToken} refreshToken={refreshToken} />
      </section>

      <section>
        <ProtectedActions accessToken={accessToken} />
      </section>
    </div>
  );
}

export default App;
