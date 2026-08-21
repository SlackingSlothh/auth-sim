import { Link } from "react-router";

export default function HomeIndex() {
  return (
    <div className="flex flex-col gap-8">
      <h1 className="text-4xl text-center">Welcome to Admin Control Panel!</h1>
      <div className="flex gap-4">
        <div className="flex flex-col flex-1 gap-2 bg-white p-4 border">
          <h2 className="text-2xl">Users</h2>
          <p>
            Register user, update user info, activate / deactivate user, and add
            / remove user from groups
          </p>
          <Link to="/users" className="mt-auto ml-auto bg-gray-100 p-2 border">
            See users
          </Link>
        </div>
        <div className="flex flex-col flex-1 gap-2 bg-white p-4 border">
          <h2 className="text-2xl">Groups</h2>
          <p>
            Add and edit groups, control group memberships, control group policy
          </p>
          <Link to="/groups" className="mt-auto ml-auto bg-gray-100 p-2 border">
            See groups
          </Link>
        </div>
        <div className="flex flex-col flex-1 gap-2 bg-white p-4 border">
          <h2 className="text-2xl">Apps</h2>
          <p>Register client apps, control app policy</p>
          <Link to="/apps" className="mt-auto ml-auto bg-gray-100 p-2 border">
            See apps
          </Link>
        </div>
      </div>
    </div>
  );
}
