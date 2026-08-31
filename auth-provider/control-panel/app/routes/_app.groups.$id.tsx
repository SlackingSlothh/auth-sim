import { useEffect, useState, useCallback } from "react";
import { Link, useParams } from "react-router";
import { apiJson, apiFetch } from "../utils/api";

type Group = {
  ID: string;
  Name: string;
  Description: string;
};

type User = {
  ID: string;
  Name: string;
  Email: string;
};

type App = {
  ID: string;
  Name: string;
  Status: string;
};

type GroupResponse = {
  group?: Group;
  users?: User[];
  apps?: App[];
};

export default function GroupEdit() {
  const { id } = useParams();

  const [users, setUsers] = useState<User[]>([]);
  const [apps, setApps] = useState<App[]>([]);

  const [name, setName] = useState<string>("");
  const [description, setDescription] = useState<string>("");
  const [originalName, setOriginalName] = useState<string>("");
  const [originalDesc, setOriginalDesc] = useState<string>("");

  const [saving, setSaving] = useState(false);
  const [showRemoveModal, setShowRemoveModal] = useState<boolean>(false);
  const [removeSelected, setRemoveSelected] = useState<string[]>([]);
  const [removeSaving, setRemoveSaving] = useState<boolean>(false);

  const [showAddModal, setShowAddModal] = useState<boolean>(false);
  const [nonmembers, setNonmembers] = useState<User[]>([]);
  const [addSelected, setAddSelected] = useState<string[]>([]);
  const [addSaving, setAddSaving] = useState<boolean>(false);
  const [loadingAvailable, setLoadingAvailable] = useState<boolean>(false);

  const [showRemoveAppModal, setShowRemoveAppModal] = useState<boolean>(false);
  const [removeAppSelected, setRemoveAppSelected] = useState<string[]>([]);
  const [removeAppSaving, setRemoveAppSaving] = useState<boolean>(false);

  const [showAddAppModal, setShowAddAppModal] = useState<boolean>(false);
  const [denied, setDenied] = useState<App[]>([]);
  const [addAppSelected, setAddAppSelected] = useState<string[]>([]);
  const [addAppSaving, setAddAppSaving] = useState<boolean>(false);
  const [loadingDeniedApps, setLoadingDeniedApps] = useState<boolean>(false);

  const loadGroup = useCallback(async () => {
    if (!id) return;
    try {
      const { res, data } = await apiJson(`/groups/${id}`);
      if (!res.ok) {
        console.error("Failed to fetch group", res.status);
        return;
      }

      if (data) {
        const payload = data as GroupResponse;
        setUsers(payload.users || []);
        setApps(payload.apps ?? []);

        const u = payload.group;
        setName(u?.Name ?? "");
        setDescription(u?.Description ?? "");
        setOriginalName(u?.Name ?? "");
        setOriginalDesc(u?.Description ?? "");
      }
    } catch (err) {
      console.error(err);
    }
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const t = setTimeout(() => {
      void loadGroup();
    }, 0);
    return () => clearTimeout(t);
  }, [id, loadGroup]);

  const changed = name !== originalName || description !== originalDesc;

  const handleSave = async () => {
    if (!id || !changed) return;
    setSaving(true);
    try {
      const body = JSON.stringify({ Name: name, Description: description });
      const res = await apiFetch(`/groups/${id}`, {
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
        setOriginalDesc(description);
      }
    } catch (err) {
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="flex gap-8">
      <div className="flex-1 flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <h2 className="text-xl">Edit Group Profile</h2>
          <h3>Name</h3>
          <input
            className="py-1 px-2 w-full bg-white border"
            value={name}
            onChange={(e) => setName(e.target.value)}
          ></input>
          <h3>Description</h3>
          <textarea
            className="py-1 px-2 w-full resize-none bg-white border"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={6}
          ></textarea>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={handleSave}
            disabled={!changed || saving}
          >
            {saving ? "Saving..." : "Save"}
          </button>
        </div>
      </div>

      <div className="flex flex-col flex-1 gap-4">
        <h2 className="text-xl">Group Members</h2>
        <div className="flex flex-col flex-1 overflow-y-auto">
          {users.map((u) => (
            <Link to={`/users/${u.ID}`} className="py-1">
              {u.Name} ({u.Email})
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
                const { res, data } = await apiJson(`/groups/${id}/nonmembers`);
                if (res.ok) {
                  setNonmembers((data as User[]) || []);
                } else {
                  console.error("Failed to fetch nonmembers", res.status);
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

      <div className="flex flex-col flex-1 gap-4">
        <h2 className="text-xl">Allowed Apps</h2>
        <div className="flex-1 overflow-y-auto">
          {apps.map((u) => (
            <Link to={`/apps/${u.ID}`} className="flex py-1">
              <span
                className={`h-3 w-3 rounded-full bg-${u.Status !== "active" ? "red" : "green"}-500`}
              ></span>
              <span>{u.Name}</span>
            </Link>
          ))}
        </div>
        <div className="flex">
          <button
            className="p-1 w-24 bg-white border"
            onClick={async () => {
              setShowAddAppModal(true);
              setLoadingDeniedApps(true);
              try {
                const { res, data } = await apiJson(
                  `/groups/${id}/denied-apps`,
                );
                if (res.ok) {
                  setDenied((data as App[]) || []);
                } else {
                  console.error("Failed to fetch denied apps", res.status);
                }
              } catch (err) {
                console.error(err);
              } finally {
                setLoadingDeniedApps(false);
              }
            }}
          >
            + Add
          </button>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={() => {
              setRemoveAppSelected([]);
              setShowRemoveAppModal(true);
            }}
          >
            - Remove
          </button>
        </div>
      </div>

      {showRemoveModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Remove Members</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {users.map((u) => (
                <label key={u.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={removeSelected.includes(u.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setRemoveSelected((prev) =>
                        checked
                          ? [...prev, u.ID]
                          : prev.filter((id) => id !== u.ID),
                      );
                    }}
                  />
                  <span>
                    {u.Name} ({u.Email})
                  </span>
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
                    const res = await apiFetch(`/groups/${id}/members`, {
                      method: "DELETE",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(removeSelected),
                    });
                    if (!res.ok) {
                      console.error("Remove members failed", res.status);
                    } else {
                      await loadGroup();
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
            <h3 className="text-lg mb-2">Add Members</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {loadingAvailable && <div>Loading...</div>}
              {!loadingAvailable && nonmembers.length === 0 && (
                <div>No members outside this group</div>
              )}
              {nonmembers.map((u) => (
                <label key={u.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={addSelected.includes(u.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setAddSelected((prev) =>
                        checked
                          ? [...prev, u.ID]
                          : prev.filter((id) => id !== u.ID),
                      );
                    }}
                  />
                  <span>
                    {u.Name} ({u.Email})
                  </span>
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
                    const res = await apiFetch(`/groups/${id}/members`, {
                      method: "POST",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(addSelected),
                    });
                    if (!res.ok) {
                      console.error("Add members failed", res.status);
                    } else {
                      await loadGroup();
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

      {showRemoveAppModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Deny Apps</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {apps.map((u) => (
                <label key={u.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={removeAppSelected.includes(u.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setRemoveAppSelected((prev) =>
                        checked
                          ? [...prev, u.ID]
                          : prev.filter((id) => id !== u.ID),
                      );
                    }}
                  />
                  <span>{u.Name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => setShowRemoveAppModal(false)}
                disabled={removeAppSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={async () => {
                  if (!id || removeAppSelected.length === 0) return;
                  setRemoveAppSaving(true);
                  try {
                    const res = await apiFetch(`/groups/${id}/allowed-apps`, {
                      method: "DELETE",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(removeAppSelected),
                    });
                    if (!res.ok) {
                      console.error("Deny apps failed", res.status);
                    } else {
                      await loadGroup();
                      setShowRemoveAppModal(false);
                      setRemoveAppSelected([]);
                    }
                  } catch (err) {
                    console.error(err);
                  } finally {
                    setRemoveAppSaving(false);
                  }
                }}
                disabled={removeAppSaving || removeAppSelected.length === 0}
              >
                {removeAppSaving ? "Denying..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}

      {showAddAppModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Allow Apps</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {loadingDeniedApps && <div>Loading...</div>}
              {!loadingDeniedApps && denied.length === 0 && (
                <div>All apps allowed</div>
              )}
              {denied.map((u) => (
                <label key={u.ID} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={addAppSelected.includes(u.ID)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setAddAppSelected((prev) =>
                        checked
                          ? [...prev, u.ID]
                          : prev.filter((id) => id !== u.ID),
                      );
                    }}
                  />
                  <span>{u.Name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => {
                  setShowAddAppModal(false);
                  setAddAppSelected([]);
                }}
                disabled={addAppSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={async () => {
                  if (!id || addAppSelected.length === 0) return;
                  setAddAppSaving(true);
                  try {
                    const res = await apiFetch(`/groups/${id}/allowed-apps`, {
                      method: "POST",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(addAppSelected),
                    });
                    if (!res.ok) {
                      console.error("Allow apps failed", res.status);
                    } else {
                      await loadGroup();
                      setShowAddAppModal(false);
                      setAddAppSelected([]);
                    }
                  } catch (err) {
                    console.error(err);
                  } finally {
                    setAddAppSaving(false);
                  }
                }}
                disabled={addAppSaving || addAppSelected.length === 0}
              >
                {addAppSaving ? "Allowing..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
