import { useEffect, useState, useCallback } from "react";
import { Link, useParams } from "react-router";
import { apiJson, apiFetch } from "../utils/api";

type Group = {
  ID: string;
  Name: string;
};

type User = {
  ID: string;
  Name: string;
  Email: string;
  Status: string;
};

type UserResponse = {
  groups?: Group[];
  user?: User;
};

export default function UserEdit() {
  const { id } = useParams();

  const [user, setUser] = useState<User | null>(null);
  const [groups, setGroups] = useState<Group[]>([]);

  const [name, setName] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [originalName, setOriginalName] = useState<string>("");

  const [saving, setSaving] = useState(false);
  const [statusSaving, setStatusSaving] = useState<boolean>(false);
  const [showPasswordModal, setShowPasswordModal] = useState<boolean>(false);
  const [newPassword, setNewPassword] = useState<string>("");
  const [confirmPassword, setConfirmPassword] = useState<string>("");
  const [passwordSaving, setPasswordSaving] = useState<boolean>(false);
  const [showRemoveModal, setShowRemoveModal] = useState<boolean>(false);
  const [removeSelected, setRemoveSelected] = useState<string[]>([]);
  const [removeSaving, setRemoveSaving] = useState<boolean>(false);

  const [showAddModal, setShowAddModal] = useState<boolean>(false);
  const [availableGroups, setAvailableGroups] = useState<Group[]>([]);
  const [addSelected, setAddSelected] = useState<string[]>([]);
  const [addSaving, setAddSaving] = useState<boolean>(false);
  const [loadingAvailable, setLoadingAvailable] = useState<boolean>(false);

  const loadUser = useCallback(async () => {
    if (!id) return;
    try {
      const { res, data } = await apiJson(`/users/${id}`);
      if (!res.ok) {
        console.error("Failed to fetch user", res.status);
        return;
      }

      if (data) {
        const payload = data as UserResponse;
        setUser(payload.user || null);
        setGroups(payload.groups || []);

        const u = payload.user;
        setName(u?.Name ?? "");
        setEmail(u?.Email ?? "");
        setOriginalName(u?.Name ?? "");
      }
    } catch (err) {
      console.error(err);
    }
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const t = setTimeout(() => {
      void loadUser();
    }, 0);
    return () => clearTimeout(t);
  }, [id, loadUser]);

  const changed = name !== originalName;

  const handleSave = async () => {
    if (!id || !changed) return;
    setSaving(true);
    try {
      const body = JSON.stringify({ Name: name });
      const res = await apiFetch(`/users/${id}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body,
      });

      if (!res.ok) {
        console.error("Save failed", res.status);
        // TODO: show error to user
      } else {
        // update originals
        setOriginalName(name);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleToggleStatus = async () => {
    if (!id || !user) return;
    const newStatus = user.Status === "active" ? "inactive" : "active";
    setStatusSaving(true);
    try {
      const res = await apiFetch(`/users/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ Status: newStatus }),
      });

      if (!res.ok) {
        console.error("Status update failed", res.status);
      } else {
        setUser({ ...user, Status: newStatus });
      }
    } catch (err) {
      console.error(err);
    } finally {
      setStatusSaving(false);
    }
  };

  const handleConfirmPassword = async () => {
    if (!id) return;
    if (!newPassword || newPassword !== confirmPassword) return;
    setPasswordSaving(true);
    try {
      const res = await apiFetch(`/users/${id}/password`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ Password: newPassword }),
      });

      if (!res.ok) {
        console.error("Password change failed", res.status);
      } else {
        setShowPasswordModal(false);
        setNewPassword("");
        setConfirmPassword("");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setPasswordSaving(false);
    }
  };

  return (
    <div className="flex gap-8">
      <div className="flex-1 flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <h2 className="text-xl">Edit User Profile</h2>
          <h3>Name</h3>
          <input
            className="py-1 px-2 w-full bg-white border"
            value={name}
            onChange={(e) => setName(e.target.value)}
          ></input>
          <h3>Email</h3>
          <input
            className="py-1 px-2 w-full bg-gray-200 border"
            value={email}
            disabled
          ></input>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={handleSave}
            disabled={!changed || saving}
          >
            {saving ? "Saving..." : "Save"}
          </button>
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-xl">Password</h2>
          <p>
            Warning: changing your password will log you out of your current
            session.
          </p>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={() => setShowPasswordModal(true)}
          >
            Change
          </button>
        </div>

        <div className="flex flex-col gap-2">
          <h2 className="text-xl">Activation</h2>
          <p>User is currently {user?.Status || "unknown"}</p>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={handleToggleStatus}
            disabled={statusSaving || !user}
          >
            {statusSaving
              ? user?.Status === "active"
                ? "Deactivating..."
                : "Activating..."
              : (user?.Status == "active" && "Deactivate") || "Activate"}
          </button>
        </div>
      </div>

      <div className="flex flex-col flex-1 gap-4">
        <h2 className="text-xl">User's Groups</h2>
        <div className="flex-1 overflow-y-auto">
          {groups.map((g) => (
            <Link to={`/groups/${g.ID}`} className="py-1">
              {g.Name}
            </Link>
          ))}
        </div>
        <div className="flex">
          <button
            className="p-1 w-24 bg-white border"
            onClick={async () => {
              setShowAddModal(true);
              setLoadingAvailable(true);
              try {
                const { res, data } = await apiJson(
                  `/users/${id}/available-groups`,
                );
                if (res.ok) {
                  setAvailableGroups((data as Group[]) || []);
                } else {
                  console.error("Failed to fetch available groups", res.status);
                }
              } catch (err) {
                console.error(err);
              } finally {
                setLoadingAvailable(false);
              }
            }}
          >
            + Add
          </button>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={() => {
              setRemoveSelected([]);
              setShowRemoveModal(true);
            }}
          >
            - Remove
          </button>
        </div>
      </div>
      {showPasswordModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Change Password</h3>
            <div className="flex flex-col gap-2">
              <input
                type="password"
                className="py-1 px-2 w-full border"
                placeholder="New password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
              />
              <input
                type="password"
                className="py-1 px-2 w-full border"
                placeholder="Confirm password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => {
                  setShowPasswordModal(false);
                  setNewPassword("");
                  setConfirmPassword("");
                }}
                disabled={passwordSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={handleConfirmPassword}
                disabled={
                  passwordSaving ||
                  !newPassword ||
                  newPassword !== confirmPassword
                }
              >
                {passwordSaving ? "Saving..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}

      {showRemoveModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Remove Groups</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {groups.map((g) => (
                <label key={g.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={removeSelected.includes(g.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setRemoveSelected((prev) =>
                        checked
                          ? [...prev, g.ID]
                          : prev.filter((id) => id !== g.ID),
                      );
                    }}
                  />
                  <span>{g.Name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => setShowRemoveModal(false)}
                disabled={removeSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={async () => {
                  if (!id || removeSelected.length === 0) return;
                  setRemoveSaving(true);
                  try {
                    const res = await apiFetch(`/users/${id}/groups`, {
                      method: "DELETE",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(removeSelected),
                    });
                    if (!res.ok) {
                      console.error("Remove groups failed", res.status);
                    } else {
                      await loadUser();
                      setShowRemoveModal(false);
                      setRemoveSelected([]);
                    }
                  } catch (err) {
                    console.error(err);
                  } finally {
                    setRemoveSaving(false);
                  }
                }}
                disabled={removeSaving || removeSelected.length === 0}
              >
                {removeSaving ? "Removing..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}

      {showAddModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Add Groups</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {loadingAvailable && <div>Loading...</div>}
              {!loadingAvailable && availableGroups.length === 0 && (
                <div>No available groups</div>
              )}
              {availableGroups.map((g) => (
                <label key={g.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={addSelected.includes(g.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setAddSelected((prev) =>
                        checked
                          ? [...prev, g.ID]
                          : prev.filter((id) => id !== g.ID),
                      );
                    }}
                  />
                  <span>{g.Name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => {
                  setShowAddModal(false);
                  setAddSelected([]);
                }}
                disabled={addSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={async () => {
                  if (!id || addSelected.length === 0) return;
                  setAddSaving(true);
                  try {
                    const res = await apiFetch(`/users/${id}/groups`, {
                      method: "POST",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(addSelected),
                    });
                    if (!res.ok) {
                      console.error("Add groups failed", res.status);
                    } else {
                      await loadUser();
                      setShowAddModal(false);
                      setAddSelected([]);
                    }
                  } catch (err) {
                    console.error(err);
                  } finally {
                    setAddSaving(false);
                  }
                }}
                disabled={addSaving || addSelected.length === 0}
              >
                {addSaving ? "Adding..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
