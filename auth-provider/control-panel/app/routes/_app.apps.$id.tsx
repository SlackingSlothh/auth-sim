import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router";
import { apiFetch, apiJson } from "../utils/api";
import { PencilIcon, TrashIcon } from "@heroicons/react/24/outline";

type AppInfo = {
  id: string;
  name: string;
  client_id: string;
  status: string;
  launch_url?: string | null;
  logout_notification_url: string;
  redirect_uris?: string[];
  allowed_groups?: AppAllowedGroup[];
};

type AppAllowedGroup = {
  id: string;
  name: string;
};

type AppDetailResponse = {
  app?: AppInfo;
};

export default function AppEdit() {
  const { id } = useParams();

  const [app, setApp] = useState<AppInfo | null>(null);
  const [redirectURIs, setRedirectURIs] = useState<string[]>([]);
  const [allowedGroups, setAllowedGroups] = useState<AppAllowedGroup[]>([]);

  const [name, setName] = useState<string>("");
  const [launchURL, setLaunchURL] = useState<string>("");
  const [logoutNotificationURL, setLogoutNotificationURL] =
    useState<string>("");
  const [originalName, setOriginalName] = useState<string>("");
  const [originalLaunchURL, setOriginalLaunchURL] = useState<string>("");
  const [originalLogoutNotificationURL, setOriginalLogoutNotificationURL] =
    useState<string>("");

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [statusSaving, setStatusSaving] = useState<boolean>(false);

  const [selectedRedirectURI, setSelectedRedirectURI] = useState<string>("");
  const [showRedirectModal, setShowRedirectModal] = useState<boolean>(false);
  const [redirectMode, setRedirectMode] = useState<"add" | "edit">("add");
  const [redirectValue, setRedirectValue] = useState<string>("");
  const [redirectSaving, setRedirectSaving] = useState<boolean>(false);

  const [showRemoveGroupModal, setShowRemoveGroupModal] =
    useState<boolean>(false);
  const [removeSelected, setRemoveSelected] = useState<string[]>([]);
  const [removeSaving, setRemoveSaving] = useState<boolean>(false);

  const [showAddGroupModal, setShowAddGroupModal] = useState<boolean>(false);
  const [availableGroups, setAvailableGroups] = useState<AppAllowedGroup[]>([]);
  const [addSelected, setAddSelected] = useState<string[]>([]);
  const [addSaving, setAddSaving] = useState<boolean>(false);
  const [loadingAvailable, setLoadingAvailable] = useState<boolean>(false);

  const loadApp = useCallback(async () => {
    if (!id) return;

    try {
      const { res, data } = await apiJson(`/apps/${id}`);
      if (!res.ok) {
        console.error("Failed to fetch app", res.status);
        return;
      }

      const payload = (data as AppDetailResponse) || {};
      const loadedApp = payload.app ?? null;
      const loadedRedirectURIs = loadedApp?.redirect_uris ?? [];

      setApp(loadedApp);
      setRedirectURIs(loadedRedirectURIs);
      setAllowedGroups(loadedApp?.allowed_groups ?? []);

      setName(loadedApp?.name ?? "");
      setLaunchURL(loadedApp?.launch_url ?? "");
      setLogoutNotificationURL(loadedApp?.logout_notification_url ?? "");
      setOriginalName(loadedApp?.name ?? "");
      setOriginalLaunchURL(loadedApp?.launch_url ?? "");
      setOriginalLogoutNotificationURL(
        loadedApp?.logout_notification_url ?? "",
      );
      setSelectedRedirectURI((current) =>
        loadedRedirectURIs.includes(current) ? current : "",
      );
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    if (!id) return;

    const timer = setTimeout(() => {
      void loadApp();
    }, 0);

    return () => clearTimeout(timer);
  }, [id, loadApp]);

  const changed =
    name !== originalName ||
    launchURL !== originalLaunchURL ||
    logoutNotificationURL !== originalLogoutNotificationURL;

  const handleSave = async () => {
    if (!id || !changed) return;
    setSaving(true);
    try {
      const launch = launchURL.trim();
      const res = await apiFetch(`/apps/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          launch_url: launch ? launch : null,
          logout_notification_url: logoutNotificationURL,
        }),
      });

      if (!res.ok) {
        console.error("Save failed", res.status);
      } else {
        setOriginalName(name);
        setOriginalLaunchURL(launch);
        setOriginalLogoutNotificationURL(logoutNotificationURL);
        setApp((current) =>
          current
            ? {
                ...current,
                name,
                launch_url: launch ? launch : null,
                logout_notification_url: logoutNotificationURL,
              }
            : current,
        );
      }
    } catch (err) {
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleToggleStatus = async () => {
    if (!id || !app) return;
    const newStatus = app.status === "active" ? "inactive" : "active";
    setStatusSaving(true);
    try {
      const launch = launchURL.trim();
      const res = await apiFetch(`/apps/${id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name,
          launch_url: launch ? launch : null,
          logout_notification_url: logoutNotificationURL,
          status: newStatus,
        }),
      });

      if (!res.ok) {
        console.error("Status update failed", res.status);
      } else {
        setApp({ ...app, status: newStatus });
      }
    } catch (err) {
      console.error(err);
    } finally {
      setStatusSaving(false);
    }
  };

  const openAddRedirectModal = () => {
    setRedirectMode("add");
    setRedirectValue("");
    setShowRedirectModal(true);
  };

  const openEditRedirectModal = (uri: string) => {
    setSelectedRedirectURI(uri);
    setRedirectMode("edit");
    setRedirectValue(uri);
    setShowRedirectModal(true);
  };

  const handleConfirmRedirect = async () => {
    if (!id || !redirectValue.trim()) return;
    setRedirectSaving(true);
    try {
      const res = await apiFetch(`/apps/${id}/redirect-uris`, {
        method: redirectMode === "add" ? "POST" : "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(
          redirectMode === "add"
            ? { redirectUri: redirectValue }
            : {
                oldRedirectUri: selectedRedirectURI,
                newRedirectUri: redirectValue,
              },
        ),
      });

      if (!res.ok) {
        console.error("Redirect URI update failed", res.status);
      } else {
        await loadApp();
        setSelectedRedirectURI(redirectValue.trim());
        setShowRedirectModal(false);
        setRedirectValue("");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setRedirectSaving(false);
    }
  };

  const handleDeleteRedirect = async (uri: string) => {
    if (!id) return;
    setRedirectSaving(true);
    try {
      const res = await apiFetch(`/apps/${id}/redirect-uris`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ redirectUri: uri }),
      });

      if (!res.ok) {
        console.error("Delete redirect URI failed", res.status);
      } else {
        await loadApp();
        setSelectedRedirectURI("");
      }
    } catch (err) {
      console.error(err);
    } finally {
      setRedirectSaving(false);
    }
  };

  if (!id) {
    return <div>App not found.</div>;
  }

  if (loading) {
    return <div>Loading app...</div>;
  }

  return (
    <div className="flex gap-8">
      <div className="flex-1 flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <h2 className="text-xl">Edit App Profile</h2>
          <h3>Name</h3>
          <input
            className="py-1 px-2 w-full bg-white border"
            value={name}
            onChange={(e) => setName(e.target.value)}
          ></input>
          <h3>Client ID</h3>
          <input
            className="py-1 px-2 w-full bg-gray-200 border"
            value={app?.client_id ?? ""}
            disabled
          ></input>
          <h3>Launch URL</h3>
          <input
            className="py-1 px-2 w-full bg-white border"
            value={launchURL}
            onChange={(e) => setLaunchURL(e.target.value)}
          ></input>
          <h3>Logout Notification URL</h3>
          <input
            className="py-1 px-2 w-full bg-white border"
            value={logoutNotificationURL}
            onChange={(e) => setLogoutNotificationURL(e.target.value)}
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
          <h2 className="text-xl">Activation</h2>
          <p>App is currently {app?.status || "unknown"}</p>
          <button
            className="p-1 w-24 bg-white border ml-auto"
            onClick={handleToggleStatus}
            disabled={statusSaving || !app}
          >
            {statusSaving
              ? app?.status === "active"
                ? "Deactivating..."
                : "Activating..."
              : (app?.status == "active" && "Deactivate") || "Activate"}
          </button>
        </div>
      </div>

      <div className="flex flex-col flex-1 gap-4">
        <h2 className="text-xl">App Redirect URIs</h2>
        <div className="flex min-h-56 flex-1 flex-col border bg-white p-2">
          <div className="flex flex-1 flex-col gap-2 overflow-y-auto">
            {redirectURIs.at(0) == null && (
              <span>No registered redirect URI</span>
            )}

            {redirectURIs.map((uri) => (
              <div key={uri} className="flex items-center gap-2 border-b pb-2">
                <span className="flex-1 break-all">{uri}</span>
                <button
                  className="p-1 h-full bg-white border"
                  onClick={() => openEditRedirectModal(uri)}
                  disabled={redirectSaving}
                >
                  <PencilIcon className="h-full"></PencilIcon>
                </button>
                <button
                  className="p-1 h-full bg-white border"
                  onClick={() => void handleDeleteRedirect(uri)}
                  disabled={redirectSaving}
                >
                  <TrashIcon className="h-full"></TrashIcon>
                </button>
              </div>
            ))}
          </div>
        </div>
        <div className="mt-2 flex">
          <button
            className="p-1 w-24 bg-white border"
            onClick={openAddRedirectModal}
          >
            + Add
          </button>
        </div>
      </div>

      <div className="flex flex-col flex-1 gap-4">
        <h2 className="text-xl">Allowed Groups</h2>
        <div className="flex flex-col flex-1 overflow-y-auto">
          {allowedGroups.map((group) => (
            <Link key={group.id} to={`/groups/${group.id}`} className="py-1">
              {group.name}
            </Link>
          ))}
        </div>
        <div className="flex">
          <button
            className="p-1 w-24 bg-white border"
            onClick={async () => {
              setShowAddGroupModal(true);
              setLoadingAvailable(true);
              try {
                const { res, data } = await apiJson(
                  `/apps/${id}/available-groups`,
                );
                if (res.ok) {
                  setAvailableGroups((data as AppAllowedGroup[]) || []);
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
              setShowRemoveGroupModal(true);
            }}
          >
            - Remove
          </button>
        </div>
      </div>

      {showRedirectModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">
              {redirectMode === "add"
                ? "Add Redirect URI"
                : "Edit Redirect URI"}
            </h3>
            <input
              className="py-1 px-2 w-full border"
              value={redirectValue}
              onChange={(e) => setRedirectValue(e.target.value)}
              autoFocus
            />
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => {
                  setShowRedirectModal(false);
                  setRedirectValue("");
                }}
                disabled={redirectSaving}
              >
                Cancel
              </button>
              <button
                className="px-3 py-1 border bg-white ml-auto"
                onClick={handleConfirmRedirect}
                disabled={redirectSaving || !redirectValue.trim()}
              >
                {redirectSaving ? "Saving..." : "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}

      {showRemoveGroupModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Remove Groups</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {allowedGroups.map((group) => (
                <label key={group.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={removeSelected.includes(group.id)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setRemoveSelected((prev) =>
                        checked
                          ? [...prev, group.id]
                          : prev.filter((id) => id !== group.id),
                      );
                    }}
                  />
                  <span>{group.name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => setShowRemoveGroupModal(false)}
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
                    const res = await apiFetch(`/apps/${id}/groups`, {
                      method: "DELETE",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(removeSelected),
                    });
                    if (!res.ok) {
                      console.error("Remove groups failed", res.status);
                    } else {
                      await loadApp();
                      setShowRemoveGroupModal(false);
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

      {showAddGroupModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white p-4 rounded shadow w-full max-w-md">
            <h3 className="text-lg mb-2">Add Groups</h3>
            <div className="flex flex-col gap-2 max-h-64 overflow-y-auto">
              {loadingAvailable && <div>Loading...</div>}
              {!loadingAvailable && availableGroups.length === 0 && (
                <div>No available groups</div>
              )}
              {availableGroups.map((group) => (
                <label key={group.id} className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    checked={addSelected.includes(group.id)}
                    onChange={(e) => {
                      const checked = e.target.checked;
                      setAddSelected((prev) =>
                        checked
                          ? [...prev, group.id]
                          : prev.filter((id) => id !== group.id),
                      );
                    }}
                  />
                  <span>{group.name}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-2 mt-4">
              <button
                className="px-3 py-1 border bg-white"
                onClick={() => {
                  setShowAddGroupModal(false);
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
                    const res = await apiFetch(`/apps/${id}/groups`, {
                      method: "POST",
                      headers: { "Content-Type": "application/json" },
                      body: JSON.stringify(addSelected),
                    });
                    if (!res.ok) {
                      console.error("Add groups failed", res.status);
                    } else {
                      await loadApp();
                      setShowAddGroupModal(false);
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
