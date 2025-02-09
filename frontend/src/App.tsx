import './App.css'
import Home from "@/pages/Home.tsx";
import {Toaster} from "@/components/ui/toaster.tsx";

function App() {

    return (
        <>
            <main className="min-h-screen bg-background"><Home/></main>
            <Toaster/>
        </>
    )
}

export default App
