import { useEffect, useRef, useState } from "react";
import Table from "../components/table";
import type { Column } from "../components/table";
import { apiJson } from "../utils/api";

type App = {
  id: string | number;
  name: string;
  clientId: string;
  status: string;
};

export default function AppsRoute() {
  const [data, setData] = useState<App[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showNew, setShowNew] = useState(false);
  const [name, setName] = useState("");
  const [logoutNotificationUrl, setLogoutNotificationUrl] = useState("");
  const [launchUrl, setLaunchUrl] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const mountedRef = useRef(true);

  async function fetchApps() {
    if (!mountedRef.current) return;
    setLoading(true);
    try {
      const { data: d } = await apiJson("/apps");
      if (!mountedRef.current) return;
      setData((Array.isArray(d) ? d : []) as App[]);
      setError(null);
    } catch (err) {
      if (!mountedRef.current) return;
      setError(String(err));
    } finally {
      if (mountedRef.current) setLoading(false);
    }
  }

  useEffect(() => {
    mountedRef.current = true;
    const t = setTimeout(() => {
      void fetchApps();
    }, 0);
    return () => {
      mountedRef.current = false;
      clearTimeout(t);
    };
  }, []);

  const columns: Column<App>[] = [
    { key: "name", header: "Name" },
    { key: "clientId", header: "Client ID" },
    { key: "status", header: "Status" },
  ];

  const resetForm = () => {
    setName("");
    setLogoutNotificationUrl("");
    setLaunchUrl("");
    setFormError(null);
  };

  const openNew = () => {
    resetForm();
    setShowNew(true);
  };

  const closeNew = () => {
    setShowNew(false);
    resetForm();
  };

  async function handleCreate() {
    setFormError(null);
    if (!name.trim()) return setFormError("Name is required");
    if (!logoutNotificationUrl.trim())
      return setFormError("Logout notification URL is required");

    setSubmitting(true);
    try {
      const payload = { name, logoutNotificationUrl, launchUrl };
      const { res, data: resp } = await apiJson("/apps", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!res.ok) {
        const msg =
          resp && typeof resp === "string" ? resp : JSON.stringify(resp || {});
        setFormError(`Create failed: ${msg}`);
      } else {
        closeNew();
        await fetchApps();
      }
    } catch (err) {
      setFormError(String(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex">
        <button
          onClick={openNew}
          className="ml-auto p-2 bg-gray-100 border cursor-pointer"
        >
          + New
        </button>
      </div>

      {loading && <div>Loading apps...</div>}
      {error && <div className="text-red-600">{error}</div>}
      {!loading && !error && (
        <Table data={data} columns={columns} rowKey="id" />
      )}

      {showNew && (
        <div className="fixed inset-0 flex items-center justify-center z-50">
          <div
            className="absolute inset-0 bg-black opacity-40"
            onClick={closeNew}
          />
          <div className="relative bg-white rounded shadow-lg p-6 w-full max-w-md z-10">
            <h2 className="text-xl mb-4">Create App</h2>
            {formError && <div className="text-red-600 mb-2">{formError}</div>}
            <div className="flex flex-col gap-2">
              <label className="flex flex-col">
                <span className="text-sm">Name</span>
                <input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="border p-2"
                />
              </label>
              <label className="flex flex-col">
                <span className="text-sm">Logout Notification URL</span>
                <input
                  value={logoutNotificationUrl}
                  onChange={(e) => setLogoutNotificationUrl(e.target.value)}
                  className="border p-2"
                />
              </label>
              <label className="flex flex-col">
                <span className="text-sm">Launch URL</span>
                <input
                  value={launchUrl}
                  onChange={(e) => setLaunchUrl(e.target.value)}
                  className="border p-2"
                />
              </label>
            </div>
            <div className="flex gap-2 mt-4 justify-end">
              <button
                onClick={closeNew}
                className="px-3 py-1 border bg-gray-100"
              >
                Cancel
              </button>
              <button
                onClick={handleCreate}
                disabled={submitting}
                className="px-3 py-1 bg-blue-600 text-white rounded"
              >
                {submitting ? "Creating..." : "Create"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
