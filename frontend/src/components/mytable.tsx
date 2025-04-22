import Navbar from "./navbar.tsx";

export default function Mytable ()  {
    return (
        <div className="bg-Dino-light ">
            <div className='flex flex-col justify-center items-center h-screen gap-2 text-2xl text-Aubergine'>
                <div>My Table</div>
            </div>
            <Navbar></Navbar>
        </div>
    )
}