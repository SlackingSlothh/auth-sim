import { useEffect, useState, useRef } from "react";
import { Link } from "react-router";
import Table from "../components/table";
import type { Column } from "../components/table";
import { apiJson } from "../utils/api";

type User = {
  id: string | number;
  email: string;
  name: string;
  status: string;
};

export default function UsersRoute() {
  const [data, setData] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showNew, setShowNew] = useState(false);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const mountedRef = useRef(true);

  async function fetchUsers() {
    if (!mountedRef.current) return;
    setLoading(true);
    try {
      const { data: d } = await apiJson("/users");
      if (!mountedRef.current) return;
      setData((Array.isArray(d) ? d : []) as User[]);
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
      void fetchUsers();
    }, 0);
    return () => {
      mountedRef.current = false;
      clearTimeout(t);
    };
  }, []);

  const columns: Column<User>[] = [
    { key: "name", header: "Name" },
    { key: "email", header: "Email" },
    { key: "status", header: "Status" },
    {
      key: "id",
      header: "",
      render: (_v, row) => (
        <Link
          to={`/users/${row.id}`}
          className="inline-block px-2 py-1 bg-blue-500 text-white rounded"
        >
          Edit
        </Link>
      ),
      width: 120,
    },
  ];

  const resetForm = () => {
    setName("");
    setEmail("");
    setPassword("");
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
    const emailRegex =
      /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
    if (!emailRegex.test(email)) {
      return setFormError("Valid email is required");
    }
    if (password.length < 8)
      return setFormError("Password must be at least 8 characters");

    setSubmitting(true);
    try {
      const payload = { name, email, password };
      const { res, data: resp } = await apiJson("/users", {
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
        await fetchUsers();
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

      {loading && <div>Loading users...</div>}
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
            <h2 className="text-xl mb-4">Create User</h2>
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
                <span className="text-sm">Email</span>
                <input
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="border p-2"
                />
              </label>
              <label className="flex flex-col">
                <span className="text-sm">Password</span>
                <input
                  value={password}
                  type="password"
                  onChange={(e) => setPassword(e.target.value)}
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
