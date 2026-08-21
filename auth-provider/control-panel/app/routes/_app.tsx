import { Outlet, Link } from "react-router";

// import type { Route } from "../../.react-router/types/app/routes/+types/_app";
// import { getCurrentUser } from "../auth/auth.server";

// export async function loader({ request }: Route.LoaderArgs) {
//   const user = await getCurrentUser(request);

//   if (!user) {
//     throw redirect("/login");
//   }

//   return { user };
// }

// export default function ProtectedLayout() {
//   return <Outlet />;
// }

export default function AppLayout() {
  return (
    <div>
      <header className="flex gap-4 fixed w-full top-0 left-0 p-4 bg-gray-400">
        <h1>Admin Control Panel</h1>
        <nav className="flex gap-4 ml-auto">
          <Link to="/">Home</Link>
          <Link to="/users">Users</Link>
          <Link to="/groups">Groups</Link>
          <Link to="/apps">Apps</Link>
        </nav>
      </header>

      <main className="m-8 mt-16">
        <Outlet />
      </main>
    </div>
  );
}
export {};
