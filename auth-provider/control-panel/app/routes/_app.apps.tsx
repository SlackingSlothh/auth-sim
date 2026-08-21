import { Outlet } from "react-router";

export default function AppsLayout() {
  return (
    <div className="flex flex-col gap-4">
      <header className="flex">
        <h1 className="text-2xl">Apps Control Panel</h1>
      </header>
      <Outlet />
    </div>
  );
}
