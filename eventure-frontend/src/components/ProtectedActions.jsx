import { useState } from "react";
import { fetchProfile, fetchAdminUsers } from "../api/auth";

export default function ProtectedActions({ accessToken }) {
  const [profileResult, setProfileResult] = useState(null);
  const [adminUsersResult, setAdminUsersResult] = useState(null);

  const handleFetchProfile = async () => {
    const result = await fetchProfile(accessToken);
    setProfileResult(result);
  };

  const handleFetchAdminUsers = async () => {
    const result = await fetchAdminUsers(accessToken);
    setAdminUsersResult(result);
  };

  return (
    <div>
      <h2>Protected Actions</h2>

      <div className="action-block">
        <button onClick={handleFetchProfile}>Fetch Profile</button>
        {profileResult && (
          <pre>
            {profileResult.success
              ? JSON.stringify(profileResult.data, null, 2)
              : `❌ ${profileResult.error}`}
          </pre>
        )}
      </div>

      <div className="action-block">
        <button onClick={handleFetchAdminUsers}>Fetch Admin Users</button>
        {adminUsersResult && (
          <pre>
            {adminUsersResult.success
              ? JSON.stringify(adminUsersResult.data, null, 2)
              : `❌ ${adminUsersResult.error}`}
          </pre>
        )}
      </div>
    </div>
  );
}
