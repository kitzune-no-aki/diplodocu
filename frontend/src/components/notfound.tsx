import { Link } from "react-router-dom";

export default function NotFound() {
    return (
        <div className="flex min-h-screen items-center justify-center bg-Dino-dark px-4">
        <div className="max-w-lg text-center">
        <h1 className="mb-4 text-6xl font-bold text-Sand">404</h1>
            <p className="mb-8 text-2xl text-Sand">
        Oops! The page you&apos;re looking for has vanished into the void.
    </p>
    <Link
    to="/"
    className="inline-block rounded-lg bg-Dino-light px-8 py-3 text-lg border border-Dino-light font-semibold text-Sand transition-colors hover:border hover:border-Sand"
        >
        Return to Safety
    </Link>
    </div>
    </div>
);
}