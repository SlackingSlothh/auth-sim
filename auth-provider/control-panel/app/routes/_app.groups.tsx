import { Outlet } from "react-router";

export default function GroupsRoute() {
  return (
    <div className="flex flex-col gap-4">
      <header className="flex">
        <h1 className="text-2xl">Groups Control Panel</h1>
      </header>
      <Outlet />
    </div>
  );
}
