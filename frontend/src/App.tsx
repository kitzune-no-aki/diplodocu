import Home from "./components/home.tsx";
import Mytable from "./components/mytable.tsx";
import NotFound from "./components/notfound.tsx";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import ProtectedRoute from "./components/protectedroute.tsx";

const router = createBrowserRouter([
    {
        path: "/",
        element: <Home />,
        errorElement: <NotFound />,
    },
    {
        path: "/mytable",
        element: (
            <ProtectedRoute>
                <Mytable />
            </ProtectedRoute>
        ),
    },
    {
        path: "*", // Catch-all 404 route
        element: <NotFound />,
    },
]);

export default function App() {
    return <RouterProvider router={router} />;
}