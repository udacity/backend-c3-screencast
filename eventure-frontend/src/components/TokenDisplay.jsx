import { useEffect } from "react";

export default function TokenDisplay({ accessToken, refreshToken }) {
  useEffect(() => {
    if (accessToken || refreshToken) {
      console.log("Access Token:", accessToken);
      console.log("Refresh Token:", refreshToken);
    }
  }, [accessToken, refreshToken]);

  if (!accessToken && !refreshToken) {
    return null;
  }

  return (
    <div>
      <h2>Tokens</h2>
      <div>
        <strong>Access Token:</strong>
        <pre>{accessToken}</pre>
      </div>
      <div>
        <strong>Refresh Token:</strong>
        <pre>{refreshToken}</pre>
      </div>
    </div>
  );
}
