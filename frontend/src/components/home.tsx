import Navbar from "./navbar.tsx";

export default function Home ()  {
    return (
        <div className="bg-Dino-light ">
            <div className='flex flex-col justify-center items-center h-screen gap-2 text-2xl text-Aubergine'>
                <div>Hallo</div>
                <div>你好</div>
                <div>Hello</div>
                <div>こんにちは</div>
                <div>Grüzi</div>
                <div>안녕하세요</div>
                <div>Salut</div>
            </div>
            <Navbar></Navbar>
        </div>
    )
}
